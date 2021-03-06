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
type Config struct {

}
type MyHandler struct {
	Config Config `json:"Config"`
	MDBS []*sql.DB `json:"MDBS"`
	Total int64 `json:"Total"`
}

func main() {

	total := 1000000
	subtotal := 100000
	dbs := make([]*sql.DB, 0)
	len := 10
	for i:=0; i<len; i++ {
		db, err := getsqlite(i)
		if err == nil {
			add_db(db, subtotal)
			dbs = append(dbs, db)
		}
	}

	h := &MyHandler{ MDBS: dbs, Total: int64(total) }
	fasthttp.ListenAndServe(":80", h.HandleFastHTTP)
	
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

func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {

	ctx.Response.Header.Set("Content-Type", "application/json")
	
	switch string(ctx.Path()) {
	case "/get":
		
		db, id := getdbid(random(h.Total), 0)
		content, err := get_content(h.MDBS[db], id)
		if err == nil{
			fmt.Fprintf(ctx, content)
		}else{
			ctx.Error("Not Found", fasthttp.StatusNotFound)
		}

	default:
		ctx.Error("Not Found", fasthttp.StatusNotFound)
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
	}else{
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
func getdbid(num, base int64) (db, id int64) {
	c, n := divmod(num-base, 100000)
	return c, n
}
func random(max int64) int64 {
	n, _ := rand.Int(rand.Reader, big.NewInt(max))
	return n.Int64()
}