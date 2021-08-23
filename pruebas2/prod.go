package main

import (
	//"fmt"
	"strconv"
	"encoding/json"
    "github.com/valyala/fasthttp"
    //"github.com/dgraph-io/ristretto"
)

type Data struct {
	C int64 `json:"C"`
	F int64 `json:"F"`
	E int64 `json:"E"`
}
type MyHandler struct {
	minicache map[int]*Data
}

func main() {

	pass := &MyHandler{ minicache: make(map[int]*Data) }
	for n := 0; n <= 3000000; n++ {
		pass.minicache[n] = &Data{ int64(n), 1844674407370955161, 1844674407370955161 }
	}
    fasthttp.ListenAndServe(":81", pass.HandleFastHTTP)

}

func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {

	id, _ := strconv.Atoi(string(ctx.QueryArgs().Peek("id")))
	json.NewEncoder(ctx).Encode(h.minicache[id])

}