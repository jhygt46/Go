package main

import (
	"os"
	"log"
	"fmt"
	"flag"
	"time"
	"math"
	"bufio"
	"strconv"
	//"reflect"
	//"runtime"
	//"io/ioutil"
	"path/filepath"
	"encoding/json"
    "github.com/valyala/fasthttp"
    "github.com/dgraph-io/ristretto"
)

type Filtros struct {
	Id int `json:"Id"`
	Data Data `json:"Data"`
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
type Data struct {
	C [] Campos `json:"C"`
	E [] Evals `json:"E"`
}
type MyHandler struct {
	cache *ristretto.Cache
	minicache *map[int]*Data
	//data *[]Data
}

func main() {

	var files []string

	file := flag.String("file", "filtros_go.json", "")
	isfile := flag.Bool("isfile", false, "")
	folder := flag.String("folder", "../utils/files/", "")

	if !*isfile  {
		i:=0
		errs := filepath.Walk(*folder, func(path string, info os.FileInfo, err error) error {
			if i > 0{
				files = append(files, path)
			}
			i++
			return nil
		})
		if errs != nil { panic(errs) }
	}else{
		files = append(files, *folder+*file)
	}

	var Datas []Data
	var minicache = make(map[int]*Data)

	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil { panic(err) }

	// READ AUTOCOMPLETE FOLDER
    for _, file := range files {
		if FileExists(file) {

			start := time.Now()

			f, err := os.Open(file)
			if err != nil { log.Fatalf("Error to read [file=%v]: %v", file, err.Error()) }
			fi, err := f.Stat()
			if err != nil { log.Fatalf("Could not obtain stat, handle error: %v", err.Error()) }

			r := bufio.NewReader(f)
			dec := json.NewDecoder(r)
			i := 0

			dec.Token()
			for dec.More() {
				var m Filtros
				err := dec.Decode(&m)
				if err != nil {
					log.Fatal(err)
				}
				i++
                Datas = append(Datas, m.Data)
				minicache[m.Id] = &Datas[len(Datas) - 1]
			}
			dec.Token()

			elapsed := time.Since(start)
			fmt.Printf("Cantidad [%v] Peso [%s] Tiempo [%v] .\n", i, FileSize(fi.Size()), elapsed)

		}
    }

	pass := &MyHandler{ cache: cache, minicache: &minicache }
    fasthttp.ListenAndServe(":80", pass.HandleFastHTTP)

}

func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	
	id, _ := strconv.Atoi(string(ctx.QueryArgs().Peek("id")))
	//val := *h.minicache
	//fmt.Println(val[id])
	//fmt.Println(id)
    fmt.Fprintf(ctx, id);
	
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
func FileExists(name string) bool {
    if fi, err := os.Stat(name); err == nil {
        if fi.Mode().IsRegular() {
            return true
        }
    }
    return false
}
/*

	jsonFile, err := os.Open("daemon.json")
    if err != nil {
        fmt.Println(err)
    }
    defer jsonFile.Close()
    byteValue, _ := ioutil.ReadAll(jsonFile)

    var result map[string]interface{}
    json.Unmarshal([]byte(byteValue), &result)

	fmt.Println(reflect.TypeOf(result))

	for key, val := range result {
		if key == "Servicios" {
			for _, arrv := range val.([]interface{}) {
				for key2, val2 := range arrv.(map[string]interface{}) {
					if key2 == "Servers" {
						for _, arrv2 := range val2.([]interface{}) {
							for key3, _ := range arrv2.(map[string]interface{}) {
								if key3 == "Ip"{
									//fmt.Println("IP: ",val3)
								}
								if key3 == "Nombre"{
									//fmt.Println("Nombre: ",val3)
								}
							}
						}
					}
					if key2 == "Alertas" {
						//fmt.Println(val2)
					}
					if key2 == "Info" {
						//fmt.Println(val2)
					}
				}
			}
		}
	}

	//typ := reflect.TypeOf(fo1)
	//fmt.Printf("Struct is %d bytes long\n", typ.Size())
	//fmt.Println(unsafe.Sizeof(fo1))

	*/