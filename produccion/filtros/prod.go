package main

import (
	"fmt"
	"time"
    "github.com/valyala/fasthttp"
)

type MyHandler struct {
	
}

func main() {

	pass := &MyHandler {}
	fasthttp.ListenAndServe(":80", pass.HandleFastHTTP)
	
}
func read_int32(data []byte) int32 {
    var x int32
    for _, c := range data {
        x = x * 10 + int32(c - '0')
    }
    return x
}
func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {

	time := time.Now()
	id := read_int32(ctx.QueryArgs().Peek("id"))
	fmt.Fprintf(ctx, "HOLA");
	printelaped(time, "HTTP")
	fmt.Println(id)

}

func printelaped(start time.Time, str string){
	elapsed := time.Since(start)
	fmt.Printf("%s / Tiempo [%v]\n", str, elapsed)
}