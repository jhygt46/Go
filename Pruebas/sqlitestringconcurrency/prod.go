package main

import (
	"fmt"
	"log"
	"net/http"

	"zombiezen.com/go/sqlite/sqlitex"
)

var dbpool *sqlitex.Pool

// Using a Pool to execute SQL in a concurrent HTTP handler.
func main() {
	var err error
	dbpool, err = sqlitex.Open("C:/Allin/db/sFiltrodb0", 0, 10)
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServe(":80", nil))
}

func handle(w http.ResponseWriter, r *http.Request) {
	conn := dbpool.Get(r.Context())
	if conn == nil {
		return
	}
	defer dbpool.Put(conn)
	stmt := conn.Prep("SELECT filtro FROM filtros WHERE id = $id")
	stmt.SetText("$id", "1")
	for {
		if hasRow, err := stmt.Step(); err != nil {
			// ... handle error
		} else if !hasRow {
			break
		}
		foo := stmt.GetText("filtro")
		// ... use foo
		fmt.Fprintln(w, foo)
	}
}
