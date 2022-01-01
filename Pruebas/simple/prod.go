package main

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

type MyHandler struct {
}

func main() {

	h := &MyHandler{}
	fasthttp.ListenAndServe(":80", h.HandleFastHTTP)

}
func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "OK")
}
