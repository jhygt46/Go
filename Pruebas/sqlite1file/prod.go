package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/valyala/fasthttp"
)

type Config struct {
}
type Filtro struct {
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
	Dbs    *sql.DB `json:"Dbs"`
	Config Config  `json:"Config"`
	Total  int64   `json:"Total"`
}

func main() {

	total := 1000000
	db, err := getsqlite(0)
	if err == nil {
		add_db(db, total)
		h := &MyHandler{Dbs: db, Total: int64(total)}
		fasthttp.ListenAndServe(":80", h.HandleFastHTTP)
	}

}

func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {

	ctx.Response.Header.Set("Content-Type", "application/json")

	switch string(ctx.Path()) {
	case "/get":

		content, err := get_content(h.Dbs, random(h.Total))
		if err == nil {
			fmt.Fprintf(ctx, content)
		} else {
			ctx.Error("Not Found", fasthttp.StatusNotFound)
		}

	default:
		ctx.Error("Not Found", fasthttp.StatusNotFound)
	}

}
func add_db(db *sql.DB, total int) {

	filtro := Filtro{}
	filtro.C = []Campos{Campos{T: 1, N: "Procesador", V: []string{"X15", "IntelC3", "Amd71", "X15", "IntelC3", "Amd71", "X15", "IntelC3", "Amd71", "X15", "IntelC3", "Amd71", "X15", "IntelC3", "Amd71", "X15", "IntelC3", "Amd71", "X15", "IntelC3", "Amd71", "X15", "IntelC3", "Amd71", "X15", "IntelC3", "Amd71", "X15", "IntelC3", "Amd71", "X15", "IntelC3", "Amd71"}}, Campos{T: 1, N: "Pantalla", V: []string{"4", "4.5", "5.5", "4.5", "5.5", "4.5", "5.5", "4.5", "5.5", "4.5", "5.5"}}, Campos{T: 1, N: "Memoria", V: []string{"2GB", "4GB", "8GB", "16GB", "32GB", "64GB", "128GB"}}, Campos{T: 1, N: "Marca", V: []string{"Samsung", "Motorola", "Nokia", "Samsung", "Motorola", "Nokia", "Samsung", "Motorola", "Nokia", "Samsung", "Motorola", "Nokia", "Samsung", "Motorola", "Nokia"}}}
	filtro.E = []Evals{Evals{T: 1, N: "Buena"}, Evals{T: 1, N: "Nelson"}, Evals{T: 1, N: "Hola"}, Evals{T: 1, N: "Mundo"}}

	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
	}
	defer tx.Rollback() // The rollback will be ignored if the tx has been committed later in the function.
	stmt, err := tx.Prepare("INSERT INTO contents(content) VALUES(?)")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close() // Prepared statements take up server resources and should be closed after use.
	now := time.Now()
	for i := 0; i < total; i++ {

		filtro.Id = int32(i)
		u, err := json.Marshal(filtro)
		if err == nil {
			if _, err := stmt.Exec(string(u)); err != nil {
				fmt.Println(err)
			}
		}

	}
	elapsed := time.Since(now)
	fmt.Printf("WRITES FILES %v [%s] c/u total %v\n", total, time_cu(elapsed, total), elapsed)
	if err := tx.Commit(); err != nil {
		fmt.Println(err)
	}

}
func getsqlite(i int) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./filtros"+strconv.Itoa(i)+".db")
	if err == nil {
		stmt, err := db.Prepare(`create table if not exists contents (id integer not null primary key autoincrement,content text)`)
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
func get_content(db *sql.DB, id int64) (string, error) {
	rows, err := db.Query("SELECT content FROM contents WHERE id=?", id)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	var content string
	for rows.Next() {
		err := rows.Scan(&content)
		if err != nil {
			return "", err
		}
	}
	return content, nil
}
func divmod(numerator, denominator int64) (quotient, remainder int64) {
	quotient = numerator / denominator
	remainder = numerator % denominator
	return
}
func random(max int64) int64 {
	n, _ := rand.Int(rand.Reader, big.NewInt(max))
	return n.Int64()
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
