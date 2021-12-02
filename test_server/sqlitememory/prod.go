package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/mattn/go-sqlite3"
	"github.com/valyala/fasthttp"
)

type MyHandler struct {
	Dbs *sql.DB `json:"Dbs"`
}

func main() {

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatalf("cannot open an SQLite memory database: %v", err)
	}
	defer db.Close()

	// sqlite> select strftime('%J', '2015-04-13T19:22:19.773Z'), strftime('%J', '2015-04-13T19:22:19');
	_, err = db.Exec("CREATE TABLE unix_time (time text); INSERT INTO unix_time (time) VALUES ('BuenaNelson')")
	if err != nil {
		log.Fatalf("cannot create schema: %v", err)
	}

	h := &MyHandler{ Dbs: db }
	fasthttp.ListenAndServe(":80", h.HandleFastHTTP)
	
}

func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {

	ctx.Response.Header.Set("Content-Type", "application/json")
	switch string(ctx.Path()) {
	case "/get0":
		row := h.Dbs.QueryRow("SELECT time FROM unix_time")
		var t string
		err := row.Scan(&t)
		if err != nil {
			log.Fatalf("cannot scan addition: %v", err)
		}
		fmt.Fprintf(ctx, t)
	case "/get1":

	default:
		ctx.Error("Not Found", fasthttp.StatusNotFound)
	}

}