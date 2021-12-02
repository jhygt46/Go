package main

import (
	"fmt"
	"time"
	"strconv"
	"math/big"
	"crypto/rand"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/valyala/fasthttp"
)

type MyHandler struct {
	Dbs1 *sql.DB `json:"Dbs1"`
	Dbs2 *sql.DB `json:"Dbs2"`
	Total int64 `json:"Total"`
}

func main() {

	totalfile := 3500000
	totalmem := 350000

	db1, err1 := getsqliteDbfile(0)
	if err1 == nil {

	}
	db2, err2 := getsqliteDbmem()
	if err2 == nil {
		
	}

	add_db(db1, totalfile)
	add_db(db2, totalmem)

	h := &MyHandler{ Dbs1: db1, Dbs2: db2, Total: int64(totalfile) }
	fasthttp.ListenAndServe(":80", h.HandleFastHTTP)
	
}

func getsqliteDbfile(i int) (*sql.DB, error) {
	db1, err := sql.Open("sqlite3", "./filtros"+strconv.Itoa(i)+".db")
	if err != nil {
		fmt.Printf("cannot open an SQLite memory database: %v", err)
		return nil, err
	}
	_, err = db1.Exec("CREATE TABLE contents (id integer not null primary key autoincrement,content text)")
	if err != nil {
		fmt.Printf("cannot create schema: %v", err)
		return nil, err
	}
	return db1, nil
}
func getsqliteDbmem() (*sql.DB, error) {
	db1, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		fmt.Printf("cannot open an SQLite memory database: %v", err)
		return nil, err
	}
	//defer db1.Close()
	_, err = db1.Exec("CREATE TABLE contents (id integer not null primary key autoincrement,content text)")
	if err != nil {
		fmt.Printf("cannot create schema: %v", err)
		return nil, err
	}
	return db1, nil
}



func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {

	ctx.Response.Header.Set("Content-Type", "application/json")
	total := random(h.Total)
	switch string(ctx.Path()) {
	case "/get0":
		content1, err := get_content(h.Dbs1, total)
		if err == nil{
			fmt.Fprintf(ctx, content1)
		}else{
			content2, err := get_content(h.Dbs2, total)
			if err == nil{
				fmt.Fprintf(ctx, content2)
			}else{
				ctx.Error("Not Found", fasthttp.StatusNotFound)
			}
		}
	case "/get1":

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
func random(max int64) int64 {
	n, _ := rand.Int(rand.Reader, big.NewInt(max - 1))
	return n.Int64() + 1
}