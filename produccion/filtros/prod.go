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

	id := binary.LittleEndian.Uint16(ctx.QueryArgs().Peek("id"))
	time := time.Now()
	fmt.Fprintf(ctx, id);
	printelaped(time, "HTTP")

}

func printelaped(start time.Time, str string){
	elapsed := time.Since(start)
	fmt.Printf("%s / Tiempo [%v]\n", str, elapsed)
}