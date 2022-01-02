package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/valyala/fasthttp"
)

type MyHandler struct {
}

func main() {

	i, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err == nil {
		fmt.Printf("%v %T", i, i)
	}

	h := &MyHandler{}
	fasthttp.ListenAndServe(":80", h.HandleFastHTTP)

}
func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "OK")
}
