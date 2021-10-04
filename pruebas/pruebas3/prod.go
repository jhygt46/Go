package main

import (
	//"flag"
	"fmt"
	//"math"
	//"time"
	"strconv"
	//"unsafe"
	//"encoding/json"
	//"github.com/DmitriyVTitov/size"
    "github.com/valyala/fasthttp"
    //"github.com/dgraph-io/ristretto"
)


type Config struct {
	Tipo int8 `json:"Tipo"`
}
type Data struct {
	C int64 `json:"C"`
	F int64 `json:"F"`
	E int64 `json:"E"`
}
type MyHandler struct {
	minicache map[int]*Data
	config Config
}

func main() {

	pass := &MyHandler{}

	/*
	config := Config{ Tipo: 1 }
	start := time.Now()
	numbPtr := flag.Int("numb", 3000000, "an int")
	flag.Parse()
	pass := &MyHandler{ minicache: make(map[int]*Data, *numbPtr), config: config }
	for n := 0; n < *numbPtr; n++ {
		pass.minicache[n] = &Data{ int64(n), 1844674407370955161, 1844674407370955161 }
	}
	fmt.Println(FileSize(int64(size.Of(pass))))
	printelaped(start, "MyHandler listo")
	fmt.Println("Memory size of Data", unsafe.Sizeof(Data{}))
	fmt.Println("Se crearon: ", *numbPtr)
	*/
	
    fasthttp.ListenAndServe(":80", pass.HandleFastHTTP)

}

func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {

	id, err := strconv.Atoi(string(ctx.QueryArgs().Peek("id")))
	if err == nil {
		fmt.Fprintf(ctx, "%d", id)
	}

	/*
	start := time.Now()
	id, _ := strconv.Atoi(string(ctx.QueryArgs().Peek("id")))
	fmt.Println("Tipo: ", h.config.Tipo)
	json.NewEncoder(ctx).Encode(h.minicache[id])
	printelaped(start, "Visita")
	*/
	
}
/*
func printelaped(start time.Time, str string){
	elapsed := time.Since(start)
	fmt.Printf("%s / Tiempo [%v]\n", str, elapsed)
}

func logn(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}
func humanateBytes(s uint64, base float64, sizes []string) string {
	if s < 10 {
		return fmt.Sprintf("%dB", s)
	}
	e := math.Floor(logn(float64(s), base))
	suffix := sizes[int(e)]
	val := float64(s) / math.Pow(base, math.Floor(e))
	f := "%.0f"
	if val < 10 {
		f = "%.1f"
	}
	return fmt.Sprintf(f+"%s", val, suffix)
}
func FileSize(s int64) string {
	sizes := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	return humanateBytes(uint64(s), 1024, sizes)
}
*/