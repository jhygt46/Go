package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/valyala/fasthttp"
)

type MyHandler struct {
	Prods   map[uint32][]uint8 `json:"Prods"`
	Empresa map[uint32][]uint8 `json:"Empresa"`
}

func main() {

	h := &MyHandler{Prods: make(map[uint32][]uint8), Empresa: make(map[uint32][]uint8)}

	total, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err == nil {

		PrintMemUsage()
		h.Prods = make(map[uint32][]uint8, total)
		h.Empresa = make(map[uint32][]uint8, total)
		for i := uint32(0); i < uint32(total); i++ {
			h.Prods[i] = make([]byte, 40)
			h.Empresa[i] = make([]byte, 40)
			rand.Read(h.Prods[i])
			rand.Read(h.Empresa[i])
		}
		PrintMemUsage()
		fmt.Printf("TOTAL: (%v)\n", total)

	}

	fasthttp.ListenAndServe(":80", h.HandleFastHTTP)

}

func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {

	x := ctx.QueryArgs().Peek("id")
	id := Read_uint32bytes(x)

	fmt.Println("ID", id)
	fmt.Println("PRO", h.Prods[id])
	fmt.Println("EMP", h.Empresa[id])
	fmt.Fprintf(ctx, "HOLA MUNDO")

}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
func Read_uint32bytes(data []byte) uint32 {
	var x uint32
	for _, c := range data {
		x = x*10 + uint32(c-'0')
	}
	return x
}
