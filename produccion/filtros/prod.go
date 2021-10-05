package main

import (
	"fmt"
	"time"
	"encoding/binary"
    "github.com/valyala/fasthttp"
)

type MyHandler struct {
	
}

func main() {

	pass := &MyHandler {}
	fasthttp.ListenAndServe(":80", pass.HandleFastHTTP)
	
}

func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {

	id := ctx.QueryArgs().Peek("id")
	ids := uint64(6)
	if len(id) > 7 {
		ids := binary.BigEndian.Uint64(id)
	}
	
	time := time.Now()
	fmt.Fprintf(ctx, "HOLA");
	printelaped(time, "HTTP")
	fmt.Println(ids)

}

func printelaped(start time.Time, str string){
	elapsed := time.Since(start)
	fmt.Printf("%s / Tiempo [%v]\n", str, elapsed)
}