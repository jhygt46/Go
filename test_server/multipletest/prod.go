package main

import (
	"fmt"
	"time"
	"strconv"
	"math/big"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"github.com/valyala/fasthttp"
	_ "github.com/mattn/go-sqlite3"
)
type Config struct {

}
type Data struct {
	C [] Campos `json:"C"`
	E [] Evals `json:"E"`
	N string `json:"N"`
}
type Campos struct {
	T int `json:"T"`
	N string `json:"N"`
	V [] string `json:"V"`
}
type Evals struct {
	T int `json:"T"`
	N string `json:"N"`
}
type MyHandler struct {
	Dbs *sql.DB `json:"Dbs"`
	Config Config `json:"Config"`
	Minicache *Minicache `json:"Minicache"`
	Total int64 `json:"Total"`
}
type Minicache struct {
	Cache map[int64]Data `json:"Cache"`
}

func main() {

	totalcache := 350000
	total := 1000000

	db, err := getsqlite(0)
	if err == nil {

		h := &MyHandler{ Dbs: db, Minicache: &Minicache{ Cache: make(map[int64]Data, totalcache) }, Total: int64(total) }
		add_db(db, total)
		h.db_to_cache(db)
		fasthttp.ListenAndServe(":80", h.HandleFastHTTP)	

	}
	
}

func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {

	ctx.Response.Header.Set("Content-Type", "application/json")
	total := random(h.Total)
	switch string(ctx.Path()) {
	case "/get":
		if res, found := h.Minicache.Cache[total]; found {
			json.NewEncoder(ctx).Encode(res)
		}else{
			content, err := get_content(h.Dbs, total)
			if err == nil{
				fmt.Fprintf(ctx, content)
			}else{
				ctx.Error("Not Found", fasthttp.StatusNotFound)
			}
		}
	case "/get1":
		content, err := get_content(h.Dbs, total)
		if err == nil{
			fmt.Fprintf(ctx, content)
		}else{
			ctx.Error("Not Found", fasthttp.StatusNotFound)
		}
	case "/get2":
		db, err := getsqlite2(0)
		if err == nil {
			content, err := get_content(db, total)
			if err == nil{
				fmt.Fprintf(ctx, content)
			}else{
				ctx.Error("Not Found", fasthttp.StatusNotFound)
			}
		}
	default:
		ctx.Error("Not Found", fasthttp.StatusNotFound)
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
func (h *MyHandler) db_to_cache(db *sql.DB) {

	now := time.Now()
	rows, err := db.Query("SELECT id, content FROM contents LIMIT 350000")
	if err != nil { 
		fmt.Println(err)
	}
	defer rows.Close()
	var content string
	var id int64
	data := Data{}
	c := 0
	for rows.Next() {
		err := rows.Scan(&id, &content)
		if err != nil { 
			fmt.Println(err)
		}
		if err := json.Unmarshal([]byte(content), &data); err == nil {
			h.Minicache.Cache[id] = data
			c++
		}
	}
	elapsed := time.Since(now)
	fmt.Printf("WRITES FILES %v [%s] c/u total %v\n", c, time_cu(elapsed, c), elapsed)

}
func getsqlite(i int) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./filtros"+strconv.Itoa(i)+".db")
	defer db.Close()
	if err == nil {
		stmt, err := db.Prepare(`create table if not exists contents (id integer not null primary key autoincrement,content text)`)
		if err != nil {
			fmt.Println("err1")
			fmt.Println(err)
			return db, err
		}
		stmt.Exec()
		return db, nil
	}else{
		fmt.Println("err2")
		fmt.Println(err)
		return db, err
	}
}
func getsqlite2(i int) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./filtros"+strconv.Itoa(i)+".db")
	defer db.Close()
	if err == nil {
		return db, nil
	}
	return db, err
}
func add_db(db *sql.DB, total int){

	str1 := []byte("{\"C\":[{ \"T\": 1, \"N\": \"Nacionalidad\", \"V\": [\"Chilena\", \"Argentina\", \"Brasileña\", \"Uruguaya\"] }, { \"T\": 2, \"N\": \"Servicios\", \"V\": [\"Americana\", \"Rusa\", \"Bailarina\", \"Masaje\"] },{ \"T\": 3, \"N\": \"Edad\" }],\"E\": [{ \"T\": 1, \"N\": \"Rostro\" },{ \"T\": 1, \"N\": \"Senos\" },{ \"T\": 1, \"N\": \"Trasero\" }]}")
	str := string(str1)
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
	for i:=0; i<total; i++ {
		if _, err := stmt.Exec(str); err != nil {
			fmt.Println(err)
		}
	}
	elapsed := time.Since(now)
	fmt.Printf("WRITES FILES %v [%s] c/u total %v\n", total, time_cu(elapsed, total), elapsed)
	if err := tx.Commit(); err != nil {
		fmt.Println(err)
	}

}
func time_cu(t time.Duration, c int) string {
	ms := float64(t / time.Nanosecond)
	res := ms / float64(c)
	var s string
	if res < 1000 {
		s = fmt.Sprintf("%.2f NanoSec", res)
	} else if res >= 1000 && res < 1000000{
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
	return strconv.FormatInt(c1, 10)+"/"+strconv.FormatInt(c2, 10)+"/"+strconv.FormatInt(c3, 10)

}
func getFolderFile64(num int64) string {

	c1, n1 := divmod(num, 1000000)
	c2, n2 := divmod(n1, 10000)
	c3, c4 := divmod(n2, 100)
	return strconv.FormatInt(c1, 10)+"/"+strconv.FormatInt(c2, 10)+"/"+strconv.FormatInt(c3, 10)+"/"+strconv.FormatInt(c4, 10)

}
func random(max int64) int64 {
	n, _ := rand.Int(rand.Reader, big.NewInt(max - 1))
	return n.Int64() + 1
}