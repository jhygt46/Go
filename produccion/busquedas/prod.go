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

	go func() {
		fasthttp.ListenAndServe(":80", pass.HandleFastHTTP)
	}()

}

func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "ERROR DDos");
}

func printelaped(start time.Time, str string){
	elapsed := time.Since(start)
	fmt.Printf("%s / Tiempo [%v]\n", str, elapsed)
}