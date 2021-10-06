package main

import (
	"os"
	"fmt"
	"time"
	"strconv"
	"io/ioutil"
	//"path/filepath"
	"encoding/json"
    "github.com/valyala/fasthttp"
)

type MyHandler struct {
	minicache map[uint32]*Data
}
type Data struct {
	C [] Campos `json:"C"`
	E [] Evals `json:"E"`
}
type Campos struct {
	T int `json:"T"`
	N string `json:"N"`
	V [] string `json:"V"`
}
type Evals struct {
	T int `json:"T"`
	N string `json:"N"`
}


func main() {

	/*
	newpath := filepath.Join("/var/Go/pruebas/utils/filtros", "1")
	err := os.MkdirAll(newpath, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
	*/

	d1 := []byte("{\"Id\":1,\"Data\":{\"C\":[{ \"T\": 1, \"N\": \"Nacionalidad\", \"V\": [\"Chilena\", \"Argentina\", \"Brasileña\", \"Uruguaya\"] }, { \"T\": 2, \"N\": \"Servicios\", \"V\": [\"Americana\", \"Rusa\", \"Bailarina\", \"Masaje\"] },{ \"T\": 3, \"N\": \"Edad\" }],\"E\": [{ \"T\": 1, \"N\": \"Rostro\" },{ \"T\": 1, \"N\": \"Senos\" },{ \"T\": 1, \"N\": \"Trasero\" }]}}")

	x := []int{10, 100}
	for j, v := range x {

		folder := strconv.Itoa(j)
		cantidad := strconv.Itoa(v)
		fmt.Println("CARPETA: "+folder+" // CANTIDAD: "+cantidad)

		time1 := time.Now()
		for i := 1; i <= v; i++ {
			err := os.WriteFile("/var/Go/pruebas/utils/filtros/"+folder+"/"+strconv.Itoa(i), d1, 0644)
			if err != nil {
				fmt.Println(err)
			}
		}
		printelaped(time1, "WRITE")

		time2 := time.Now()
		for i := 1; i <= v; i++ {
			jsonFile, err := os.Open("/var/Go/pruebas/utils/filtros/"+folder+"/"+strconv.Itoa(i))
			if err == nil{
				byteValue, _ := ioutil.ReadAll(jsonFile)
				read(byteValue)
			}
			defer jsonFile.Close()
		}
		printelaped(time2, "READ")

	}


	
	

	
	

	pass := &MyHandler {}
	fasthttp.ListenAndServe(":80", pass.HandleFastHTTP)
	
}
func read(x []byte){

}
func read_int32(data []byte) uint32 {
    var x uint32
    for _, c := range data {
        x = x * 10 + uint32(c - '0')
    }
    return x
}
func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {

	time := time.Now()
	
	switch string(ctx.Path()) {
	case "/filtro":
		ctx.Response.Header.Set("Content-Type", "application/json")
		id := read_int32(ctx.QueryArgs().Peek("id"))
		if res, found := h.minicache[id]; found {
			json.NewEncoder(ctx).Encode(res)
		}else{

		}
	default:
		ctx.Error("Not Found", fasthttp.StatusNotFound)
	}

	printelaped(time, "HTTP")
	

}

func printelaped(start time.Time, str string){
	elapsed := time.Since(start)
	fmt.Printf("%s / Tiempo [%v]\n", str, elapsed)
}