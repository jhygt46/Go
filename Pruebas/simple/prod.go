package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"

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
type Cuadrantes struct {
	Minlat             float32 `json:"Minlat"`
	Maxlat             float32 `json:"Minlat"`
	Minlng             float32 `json:"Minlat"`
	Maxlng             float32 `json:"Minlat"`
	DimensionCuadrante float32 `json:"Minlat"`
}

type Runes struct {
}

func main() {

	h := &MyHandler{cache: make([]*[]uint8, 0)}

	aux := make([][]uint8, 0)

	data := Data{}
	data.C = []Campos{Campos{T: 1, N: "Procesador", V: []string{"X15", "IntelC3", "", "Amd71"}}, Campos{T: 1, N: "Pantalla", V: []string{"4", "4.5", "5.5"}}, Campos{T: 1, N: "Memoria", V: []string{"2GB", "4GB", "8GB"}}, Campos{T: 1, N: "Marca", V: []string{"Samsung", "Motorola", "Nokia"}}}
	data.E = []Evals{Evals{T: 1, N: "Buena"}, Evals{T: 1, N: "Nelson"}, Evals{T: 1, N: "Hola"}, Evals{T: 1, N: "Mundo"}}

	for i := 0; i < 10; i++ {
		data.Id = int32(i)
		u, err := json.Marshal(data)
		if err == nil {
			aux = append(aux, u)
			h.cache = append(h.cache, &aux[len(aux)-1])
		}
	}
	now := time.Now()
	x, y, w := getCuads(31.75, 75.50, 0.25)
	lapsed := time.Since(now)

	/*
		x := Get_Cuad(31.75, 75.50)
		fmt.Println(x)
		y := Get_Cuad(31.75, 75.25)
		fmt.Println(y)
	*/

	fmt.Println("x:", x)
	fmt.Println("y:", y)
	fmt.Println("w:", w)
	fmt.Println("Lapsed:", lapsed)

	fasthttp.ListenAndServe(":80", h.HandleFastHTTP)

}
func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {

	x := ctx.QueryArgs().Peek("id")
	id := Read_uint32bytes(x)

	fmt.Printf("%v %T\n", id, id)

	//fmt.Fprintf(ctx, string(*h.cache[id]))
	fmt.Fprintf(ctx, string(id))

}

func Read_uint32bytes(data []byte) []byte {
	var x uint32
	for _, c := range data {
		x = x*10 + uint32(c-'0')
	}
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, x)
	return b
}

func Get_Cuad(lat, lng float32) uint16 {

	var cantcuadlat float32 = 8
	var cuaddim float32 = 0.25
	var minlat float32 = 31.0
	var minlng float32 = 75.0
	x := (lat - minlat) / cuaddim
	y := ((lng - minlng) / cuaddim) * cantcuadlat
	return uint16(x + y)

}

func getCuads(lat, lng, distancia float32) (uint16, uint16, uint16) {

	var cantcuadlat uint16 = 8

	x := Get_Cuad(lat, lng)
	y := Get_Cuad(lat, lng+distancia)
	w := (y - x) / cantcuadlat

	return x, y, w

}
