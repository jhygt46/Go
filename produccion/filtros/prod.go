package main

import (
	"os"
	"fmt"
	"time"
	"strconv"
	"io/ioutil"
	"path/filepath"
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

	d1 := []byte("{\"Id\":1,\"Data\":{\"C\":[{ \"T\": 1, \"N\": \"Nacionalidad\", \"V\": [\"Chilena\", \"Argentina\", \"Brasileña\", \"Uruguaya\"] }, { \"T\": 2, \"N\": \"Servicios\", \"V\": [\"Americana\", \"Rusa\", \"Bailarina\", \"Masaje\"] },{ \"T\": 3, \"N\": \"Edad\" }],\"E\": [{ \"T\": 1, \"N\": \"Rostro\" },{ \"T\": 1, \"N\": \"Senos\" },{ \"T\": 1, \"N\": \"Trasero\" }]}}")

	x := make([]int, 2)

	for j, _ := range x {

		v := 1000
		folder := getFolder(j)
		cant := uint64(v)

		newpath := filepath.Join("/var/Go/pruebas/utils/filtros", folder)
		err := os.MkdirAll(newpath, os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Carpeta %v creada...\n", folder)

		time1 := time.Now()
		for i := 0; i < v; i++ {
			err := os.WriteFile("/var/Go/pruebas/utils/filtros/"+folder+"/"+strconv.Itoa(i), d1, 0644)
			if err != nil {
				fmt.Println(err)
			}
		}
		elapsed1 := uint64(time.Since(time1) / time.Nanosecond) / cant


		time2 := time.Now()
		for i := 0; i < v; i++ {
			file, err := os.Open("/var/Go/pruebas/utils/filtros/"+folder+"/"+strconv.Itoa(i))
			if err != nil{
				fmt.Println(err)
			}
			file.Close()
			byteValue, _ := ioutil.ReadAll(file)
			read(byteValue)
		}
		elapsed2 := uint64(time.Since(time2) / time.Nanosecond) / cant
		cantidad := strconv.Itoa(v)
		fmt.Printf("DuracionEscritura c/u [%v] / DuracionLectura c/u [%v] / Cantidad [%v] \n", elapsed1, elapsed2, cantidad)

	}


	
	

	
	

	pass := &MyHandler {}
	fasthttp.ListenAndServe(":80", pass.HandleFastHTTP)
	
}
func read(x []byte){

}
func getFolder(num int) string {

	var c1 int = num / 1000000
	var c2 int = num / 10000
	var c3 int = num / 100

	fmt.Printf("num[%v] c1[%v] c2[%v]", num, c1, c2)
	return strconv.Itoa(c1)+"/"+strconv.Itoa(c2)+"/"+strconv.Itoa(c3)
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