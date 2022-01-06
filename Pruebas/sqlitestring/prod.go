package main

import (
	"database/sql"
	"fmt"
	"os"
	"resource/db"
	"resource/utils"
	"strconv"

	"github.com/valyala/fasthttp"
)

type MyHandler struct {
	Dbs   *sql.DB `json:"Dbs"`
	Total int64   `json:"Total"`
}

func main() {

	i, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err == nil {

		sqlite, err := db.GetDbFiltroBytes("sFiltrodb0")
		s := "ATTACH DATABASE '/var/db/sFiltrodb1' as 'db1'"
		_, err = sqlite.Exec(s)

		//sqlite.SetMaxIdleConns(5)

		if err == nil {

			filtro := db.Filtro{}
			filtro.C = []db.Campos{db.Campos{T: 1, N: "Procesador", V: []string{"X15", "IntelC3", "Amd71", "X15", "IntelC3", "Amd71", "X15", "IntelC3", "Amd71", "X15", "IntelC3", "Amd71", "X15", "IntelC3", "Amd71", "X15", "IntelC3", "Amd71", "X15", "IntelC3", "Amd71", "X15", "IntelC3", "Amd71", "X15", "IntelC3", "Amd71", "X15", "IntelC3", "Amd71", "X15", "IntelC3", "Amd71"}}, db.Campos{T: 1, N: "Pantalla", V: []string{"4", "4.5", "5.5", "4.5", "5.5", "4.5", "5.5", "4.5", "5.5", "4.5", "5.5"}}, db.Campos{T: 1, N: "Memoria", V: []string{"2GB", "4GB", "8GB", "16GB", "32GB", "64GB", "128GB"}}, db.Campos{T: 1, N: "Marca", V: []string{"Samsung", "Motorola", "Nokia", "Samsung", "Motorola", "Nokia", "Samsung", "Motorola", "Nokia", "Samsung", "Motorola", "Nokia", "Samsung", "Motorola", "Nokia"}}}
			filtro.E = []db.Evals{db.Evals{T: 1, N: "Buena"}, db.Evals{T: 1, N: "Nelson"}, db.Evals{T: 1, N: "Hola"}, db.Evals{T: 1, N: "Mundo"}}
			//db.FiltroBytesInit(sqlite, filtro, i)

		}

		h := &MyHandler{Dbs: sqlite, Total: i}
		fasthttp.ListenAndServe(":80", h.HandleFastHTTP)

	}

}

func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {

	ctx.Response.Header.Set("Content-Type", "application/json")
	switch string(ctx.Path()) {
	case "/get0":
		ran := utils.Random(h.Total)
		content, err := db.GetFiltroStringContent(h.Dbs, ran)
		if err == nil {
			fmt.Println("FOUND", ran)
			fmt.Fprintf(ctx, content)
		} else {
			fmt.Println("NOT FOUND", ran)
			ctx.Error("Not Found", fasthttp.StatusNotFound)
		}
	default:
		ctx.Error("Not Found", fasthttp.StatusNotFound)
	}

}
