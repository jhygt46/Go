package main

import (
	"os"
	"fmt"
	"time"
	"strconv"
	"math/big"
	"io/ioutil"
	"crypto/rand"
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

	escribirArchivos()
	leerArchivos()

	pass := &MyHandler {}
	fasthttp.ListenAndServe(":80", pass.HandleFastHTTP)
	
}
func read(x []byte){

}

func escribirArchivos(){

	d1 := []byte("{\"Id\":1,\"Data\":{\"C\":[{ \"T\": 1, \"N\": \"Nacionalidad\", \"V\": [\"Chilena\", \"Argentina\", \"Brasileña\", \"Uruguaya\"] }, { \"T\": 2, \"N\": \"Servicios\", \"V\": [\"Americana\", \"Rusa\", \"Bailarina\", \"Masaje\"] },{ \"T\": 3, \"N\": \"Edad\" }],\"E\": [{ \"T\": 1, \"N\": \"Rostro\" },{ \"T\": 1, \"N\": \"Senos\" },{ \"T\": 1, \"N\": \"Trasero\" }]}}")

	x := make([]int, 10000)
	c := 0
	time1 := time.Now()

	for j, _ := range x {

		//j = j + 1863100
		v := 100
		folder := getFolder(j)
		//cant := uint64(v)

		newpath := filepath.Join("/var/tmp/utils/filtros", folder)
		err := os.MkdirAll(newpath, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			fmt.Println("FOLDER ERROR: ", err)
		}

		//time1 := time.Now()
		for i := 0; i < v; i++ {
			err := os.WriteFile("/var/tmp/utils/filtros/"+folder+"/"+strconv.Itoa(i), d1, 0644)
			if err != nil {
				fmt.Println(err)
			}
		}
		//elapsed1 := uint64(time.Since(time1) / time.Nanosecond) / cant
		//fmt.Printf("utils/filtros/%v [%v] [%v]\n", folder, j, elapsed1)
		c++
	}
	elapsed1 := time.Since(time1)
	fmt.Printf("Cantidad %v / Tiempo: [%v]\n", c, elapsed1)

}

func leerArchivos(){
	
	time1 := time.Now()
	i := 0
	for i < 2000 {

		n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
		folder := getFolder(int(n.Int64()))
		file, err := os.Open("/home/admin/Go/pruebas/utils/filtros/"+folder)
		if err != nil{
			fmt.Println(err)
		}
		file.Close()
		byteValue, _ := ioutil.ReadAll(file)
		read(byteValue)
		i++

	}

	elapsed1 := uint64(time.Since(time1) / time.Nanosecond) / 2000
	fmt.Printf("DuracionLectura [%v]", elapsed1)

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
func read_int32(data []byte) uint32 {
    var x uint32
    for _, c := range data {
        x = x * 10 + uint32(c - '0')
    }
    return x
}
func printelaped(start time.Time, str string){
	elapsed := time.Since(start)
	fmt.Printf("%s / Tiempo [%v]\n", str, elapsed)
}
func getFolder(num int) string {

	var c1 int = num / 1000000
	var n1 int = num - c1 * 1000000

	var c2 int = n1 / 10000
	n1 = n1 - c2 * 10000

	var c3 int = n1 / 100
	var c4 int = n1 % 100

	//fmt.Printf("num[%v] c1[%v] c2[%v]", num, c1, c2)
	return strconv.Itoa(c1)+"/"+strconv.Itoa(c2)+"/"+strconv.Itoa(c3)+"/"+strconv.Itoa(c4)
}