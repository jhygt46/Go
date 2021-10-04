package main

import (
	"fmt"
    "github.com/valyala/fasthttp"
)

type MyHandler struct {
	
}

func main() {

	pass := &MyHandler {}

	//go func() {
		fasthttp.ListenAndServe(":80", pass.HandleFastHTTP)
	//}()

}

func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "ERROR DDos");
}