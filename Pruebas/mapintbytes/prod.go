package main

import (
	"database/sql"
	"fmt"
	"math/big"
	"os"
	"resource/utils"
	"strconv"
	"time"

	"crypto/rand"
	"encoding/json"

	_ "github.com/mattn/go-sqlite3"
	"github.com/valyala/fasthttp"
)

type Data struct {
	C  []Campos `json:"C"`
	E  []Evals  `json:"E"`
	Id int32    `json:"Id"`
}
type Campos struct {
	T int      `json:"T"`
	N string   `json:"N"`
	V []string `json:"V"`
}
type Evals struct {
	T int    `json:"T"`
	N string `json:"N"`
}
type MyHandler struct {
	Cache map[int64][]uint8 `json:"Cache"`
	Total int64             `json:"Total"`
}

func main() {

	i, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err == nil {

		files := []string{"filtrodb0"}
		CreateDb(files)

		h := &MyHandler{Cache: make(map[int64][]uint8, i), Total: i}
		h.AddCache(files[0], i)

		fasthttp.ListenAndServe(":80", h.HandleFastHTTP)

	}

}

func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {

	ctx.Response.Header.Set("Content-Type", "application/json")

	switch string(ctx.Path()) {
	case "/get0":
		if res, found := h.Cache[random(h.Total)]; found {
			fmt.Fprintf(ctx, string(res))
		} else {
			ctx.Error("Not Found", fasthttp.StatusNotFound)
		}
	case "/get1":
		fmt.Fprintf(ctx, "OK")
	default:
		ctx.Error("Not Found", fasthttp.StatusNotFound)
	}
}
func CreateDb(files []string) {

	total := 1000000
	for _, v := range files {
		if !utils.FileExists("/var/db/" + v) {
			db, err := getsqlite(v)
			if err == nil {
				now := time.Now()
				add_db(db, total)
				elapsed := time.Since(now)
				fmt.Printf("CREATE DB %s TOTAL %v [%s] c/u total %v\n", v, total, time_cu(elapsed, total), elapsed)
			}
		}
	}
}
func add_db(db *sql.DB, total int) {

	data := Data{}
	data.C = []Campos{Campos{T: 1, N: "Procesador", V: []string{"X15", "IntelC3", "", "Amd71"}}, Campos{T: 1, N: "Pantalla", V: []string{"4", "4.5", "5.5"}}, Campos{T: 1, N: "Memoria", V: []string{"2GB", "4GB", "8GB"}}, Campos{T: 1, N: "Marca", V: []string{"Samsung", "Motorola", "Nokia"}}}
	data.E = []Evals{Evals{T: 1, N: "Buena"}, Evals{T: 1, N: "Nelson"}, Evals{T: 1, N: "Hola"}, Evals{T: 1, N: "Mundo"}}

	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
	}
	defer tx.Rollback() // The rollback will be ignored if the tx has been committed later in the function.
	stmt, err := tx.Prepare("INSERT INTO filtros (filtro) VALUES(?)")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close() // Prepared statements take up server resources and should be closed after use.
	for i := 0; i < total; i++ {
		data.Id = int32(i)
		u, err := json.Marshal(data)
		if err == nil {
			if _, err := stmt.Exec(string(u)); err != nil {
				fmt.Println(err)
			}
		}
	}
	if err := tx.Commit(); err != nil {
		fmt.Println(err)
	}
}
func getsqlite(dbname string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "/var/db/"+dbname)
	if err == nil {
		stmt, err := db.Prepare(`create table if not exists filtros (id integer not null primary key autoincrement,filtro text, cache integer)`)
		if err != nil {
			fmt.Println("err1")
			fmt.Println(err)
			return db, err
		}
		stmt.Exec()
		return db, nil
	} else {
		fmt.Println("err2")
		fmt.Println(err)
		return db, err
	}
}
func (h *MyHandler) AddCache(file string, cant int64) {

	now := time.Now()
	db, err := sql.Open("sqlite3", "/var/db/"+file)

	if err == nil {
		rows, err := db.Query("SELECT id, filtro FROM filtros LIMIT ?", cant)
		if err == nil {
			defer rows.Close()
			var id int64
			var filtro string
			for rows.Next() {
				err := rows.Scan(&id, &filtro)
				if err == nil {
					/*
						data := Data{}
						if err := json.Unmarshal([]byte(filtro), &data); err == nil {
							h.Cache[id] = data
						}
					*/
					h.Cache[id] = []byte(filtro)
				} else {
					fmt.Print("ERR SCAN:")
					fmt.Println(err)
				}
			}
		} else {
			fmt.Print("ERR SELECT TABLE FILTROS:")
			fmt.Println(err)
		}
	} else {
		fmt.Print("ERR CONNECT DB:", file)
		fmt.Println(err)
	}
	elapsed := time.Since(now)
	fmt.Printf("ADD CACHE DB %s TOTAL %v [%s] c/u total %v\n", file, cant, time_cu(elapsed, int(cant)), elapsed)
}
func time_cu(t time.Duration, c int) string {
	ms := float64(t / time.Nanosecond)
	res := ms / float64(c)
	var s string
	if res < 1000 {
		s = fmt.Sprintf("%.2f NanoSec", res)
	} else if res >= 1000 && res < 1000000 {
		s = fmt.Sprintf("%.2f MicroSec", res/1000)
	} else {
		s = fmt.Sprintf("%.2f MilliSec", res/1000000)
	}
	return s
}
func divmod(numerator, denominator int64) (quotient, remainder int64) {
	quotient = numerator / denominator
	remainder = numerator % denominator
	return
}
func getFolder64(num int64) string {

	c1, n1 := divmod(num, 1000000)
	c2, n2 := divmod(n1, 10000)
	c3, _ := divmod(n2, 100)
	return strconv.FormatInt(c1, 10) + "/" + strconv.FormatInt(c2, 10) + "/" + strconv.FormatInt(c3, 10)
}
func getFolderFile64(num int64) string {

	c1, n1 := divmod(num, 1000000)
	c2, n2 := divmod(n1, 10000)
	c3, c4 := divmod(n2, 100)
	return strconv.FormatInt(c1, 10) + "/" + strconv.FormatInt(c2, 10) + "/" + strconv.FormatInt(c3, 10) + "/" + strconv.FormatInt(c4, 10)
}
func random(max int64) int64 {
	n, _ := rand.Int(rand.Reader, big.NewInt(max-1))
	return n.Int64() + 1
}
