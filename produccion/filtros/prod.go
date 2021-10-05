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
	fmt.Println(id)

	if len(id) > 7 {
		ids := binary.BigEndian.Uint64(id)
		fmt.Println(ids)
	}else{
		ids := uint64(6)
		fmt.Println(ids)
	}
	
	time := time.Now()
	fmt.Fprintf(ctx, "HOLA");
	printelaped(time, "HTTP")
	

}

func printelaped(start time.Time, str string){
	elapsed := time.Since(start)
	fmt.Printf("%s / Tiempo [%v]\n", str, elapsed)
}