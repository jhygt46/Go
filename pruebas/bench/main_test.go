package main

import (
    //"os"
    "testing"
	"fmt"
    "time"
    //"strconv"
    //"crypto/rand"
    //"math/big"
)

//https://gist.github.com/arsham/bbc93990d8e5c9b54128a3d88901ab90#file-go_cpu_memory_profiling_benchmarks-sh
//go test -bench=. -benchmem

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

type MyHandler struct {
	Minicache map[uint64]*Data
    Config Config
}

type Config struct {
	Fecha time.Time `json:"Fecha"`
	AutoCache bool `json:"Cachetipo"`
	TotalCache int32 `json:"TotalCache"`
	CountCache int32 `json:"CountCache"`
	MetricCount int32 `json:"MetricCount"`
	MetricCache int32 `json:"MetricCache"`
	MetricFile int32 `json:"MetricFile"`
}

func main() {
    fmt.Println("Hello World")
}

var letras []int32 = []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

func BenchmarkFoo1(b *testing.B) {
    for i := 0; i < b.N; i++ {
        IndexOf(letras, 5)
    }
}

func IndexOf(arr []int32, candidate int32) int {
    for index, c := range arr {
        if c == candidate {
            return index
        }
    }
    return -1
}

/*
func BenchmarkFoo1(b *testing.B) {
    h := &MyHandler{}
    genericBenchmarkFoo1(b, h)
}
func BenchmarkFoo2(b *testing.B) {
    h := &MyHandler{}
    genericBenchmarkFoo2(b, h)
}
func BenchmarkFoo3(b *testing.B) {
    h := &MyHandler{}
    genericBenchmarkFoo3(b, h)
}
func BenchmarkFoo4(b *testing.B) {
    h := &MyHandler{}
    genericBenchmarkFoo4(b, h)
}
*/
/*
func genericBenchmarkFoo1(b *testing.B, h *MyHandler) {
    for i := 0; i < b.N; i++ {
        id := uint64(i)
        if res, found := h.Minicache[id]; found {
			json.NewEncoder(ctx).Encode(res)
			h.Config.MetricCache++
		}else{
            jsonFiltro, err := os.Open("/var/filtros/"+getFolder64(id))
			if err == nil {
				byteValueFiltro, _ := ioutil.ReadAll(jsonFiltro)
				if h.Config.AutoCache {
					data := Data{}
					if err := json.Unmarshal(byteValueFiltro, &data); err == nil {
						h.Minicache[uint64(id)] = &data
						h.Config.CountCache++
						if h.Config.CountCache >= h.Config.TotalCache {
							h.Config.AutoCache = false
						}
					}
				}
				h.Config.MetricFile++
				fmt.Fprintf(string(byteValueFiltro))
			}
			defer jsonFiltro.Close()
        }
    }
}
func genericBenchmarkFoo2(b *testing.B, n *MyHandler) {
}
func genericBenchmarkFoo3(b *testing.B, n *MyHandler) {
}
func genericBenchmarkFoo4(b *testing.B, n *MyHandler) {  
}
*/

/*
func BenchmarkCalculateA(b *testing.B) {
    for i := 0; i < b.N; i++ {
        n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
        n.Int64()
    }
}

func BenchmarkCalculateB(b *testing.B) {
    d1 := []byte{115, 111, 109, 101, 10}
    for i := 0; i < b.N; i++ {
        err := os.WriteFile("/var/Go/pruebas/utils/filtros/1/"+strconv.Itoa(i), d1, 0644)
        if err != nil {
            fmt.Println(err)
        }
    }
}
func BenchmarkCalculateA(b *testing.B) {
    byteNumber := []byte{49, 50, 53, 52, 49, 51, 52}
    for i := 0; i < b.N; i++ {
        read_int32(byteNumber)
    }
}
func BenchmarkCalculateB(b *testing.B) {
    byteNumber := []byte{49, 50, 53, 52, 49, 51, 52}
    for i := 0; i < b.N; i++ {
        read_int32b(byteNumber)
    }
}
func read_int32(data []byte) int32 {
    var x int32
    for _, c := range data {
        x = x * 10 + int32(c - '0')
    }
    return x
}
func read_int32b(data []byte) int {
    x, _ := strconv.Atoi(string(data))
    return x
}
*/