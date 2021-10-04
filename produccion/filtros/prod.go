package main

import (
	"time"
	"fmt"
    "github.com/valyala/fasthttp"
)

type MyHandler struct {
	
}

func main() {

	pass := &MyHandler {}
	fasthttp.ListenAndServe(":80", pass.HandleFastHTTP)
	
}

func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	time := time.Now()
	fmt.Fprintf(ctx, "ERROR DDos");
	printelaped(time, "HTTP")
}

func printelaped(start time.Time, str string){
	elapsed := time.Since(start)
	fmt.Printf("%s / Tiempo [%v]\n", str, elapsed)
}