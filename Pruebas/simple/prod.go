package main

import (
	"encoding/json"
	"fmt"
	"resource/utils"

	"github.com/valyala/fasthttp"
)

type MyHandler struct {
	cache []*[]uint8
}
type Data struct {
	C  []Campos `json:"C"`
	E  []Evals  `json:"E"`
	Id int32    `json:"Id"`
}
type Campos struct {
	T int      `json:"T"`
	N string   `json:"N"`
	V []string `json:"V"`
}
type Evals struct {
	T int    `json:"T"`
	N string `json:"N"`
}

func main() {

	h := &MyHandler{cache: make([]*[]uint8, 0)}

	aux := make([][]uint8, 0)

	data := Data{}
	data.C = []Campos{Campos{T: 1, N: "Procesador", V: []string{"X15", "IntelC3", "", "Amd71"}}, Campos{T: 1, N: "Pantalla", V: []string{"4", "4.5", "5.5"}}, Campos{T: 1, N: "Memoria", V: []string{"2GB", "4GB", "8GB"}}, Campos{T: 1, N: "Marca", V: []string{"Samsung", "Motorola", "Nokia"}}}
	data.E = []Evals{Evals{T: 1, N: "Buena"}, Evals{T: 1, N: "Nelson"}, Evals{T: 1, N: "Hola"}, Evals{T: 1, N: "Mundo"}}

	for i := 0; i < 1000; i++ {

		data.Id = int32(i)

		u, err := json.Marshal(data)
		if err == nil {
			aux = append(aux, u)
			h.cache = append(h.cache, &aux[len(aux)-1])
		}

	}

	fasthttp.ListenAndServe(":80", h.HandleFastHTTP)

}
func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {

	id := utils.Read_uint32(ctx.QueryArgs().Peek("id"))
	fmt.Fprintf(ctx, string(*h.cache[id]))

}
