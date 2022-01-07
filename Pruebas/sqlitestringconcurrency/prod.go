package main

import (
	"fmt"
	"resource/utils"

	"github.com/bvinc/go-sqlite-lite/sqlite3"
	"github.com/valyala/fasthttp"
)

type MyHandler struct {
	Total int64         `json:"Total"`
	Dbs   *sqlite3.Conn `json:"Dbs"`
	Stmt  *sqlite3.Stmt `json:"Stmt"`
}

func main() {

	conn, err := sqlite3.Open("file:/var/db/sFiltrodb0?cache=shared&mode=ro")
	if err != nil {
		fmt.Println("Err1", err)
	}

	stmt, err := conn.Prepare("SELECT filtro FROM filtros WHERE id=?")
	if err != nil {
		fmt.Println("Err2", err)
	}

	defer stmt.Close()
	defer conn.Close()

	h := &MyHandler{Stmt: stmt, Total: 1000000}
	fasthttp.ListenAndServe(":80", h.HandleFastHTTP)

}

func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {

	err := h.Stmt.Bind(utils.Random(h.Total))
	_, err = h.Stmt.Step()
	check(err)
	var filtro string
	err = h.Stmt.Scan(&filtro)
	check(err)
	err = h.Stmt.Reset()
	check(err)
	fmt.Fprintf(ctx, filtro)

	/*
		ctx.Response.Header.Set("Content-Type", "application/json")
		switch string(ctx.Path()) {
		case "/get0":

				stmt, err := h.Dbs.Prepare(`SELECT filtro FROM filtros WHERE id = ?`, 18)
				check(err)
				defer stmt.Close()

				for {
					hasRow, err := stmt.Step()
					check(err)
					if !hasRow {
						break
					}

					var filtro string
					err = stmt.Scan(&filtro)
					check(err)
					fmt.Fprintf(ctx, filtro)
				}

		default:
			ctx.Error("Not Found", fasthttp.StatusNotFound)
		}
	*/

}

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

/*
package main

import (
	"fmt"
	"log"
	"net/http"
	"resource/utils"

	"zombiezen.com/go/sqlite/sqlitex"
)

var dbpool *sqlitex.Pool
var Total int64 = 1000000

// Using a Pool to execute SQL in a concurrent HTTP handler.
func main() {
	var err error
	dbpool, err = sqlitex.Open("/var/db/sFiltrodb0", 0, 10)
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
	stmt.SetInt64("$id", utils.Random(Total))
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
*/
