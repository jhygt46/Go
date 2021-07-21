package main

import (
    "fmt"
    "os"
    "math"
    "time"
    "log"
    "bufio"
	"strconv"
	//"reflect"
	"io/ioutil"
	"path/filepath"
    "encoding/json"
    "github.com/valyala/fasthttp"
    "github.com/dgraph-io/ristretto"
)

type Data struct {
	T int `json:"T"`
	I int `json:"I"`
	P string `json:"P"`
}
type SingleData []Data
type Palabras struct {
	Id string `json:"Id"`
	Data SingleData `json:"Data"`
}
var MultipleDatas []SingleData

type MyHandler struct {
	cache *ristretto.Cache
}

var leng = [27]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "ñ", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}

func main() {

    cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil { panic(err) }

    // READ AUTOCOMPLETE FOLDER
	var files []string
	errs := filepath.Walk("autocomplete", func(path string, info os.FileInfo, err error) error {
        files = append(files, path)
        return nil
    })
	if errs != nil { panic(err) }
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
				var m Palabras
				err := dec.Decode(&m)
				if err != nil {
					log.Fatal(err)
				}
				i++
                MultipleDatas = append(MultipleDatas, m.Data)
				cache.Set(m.Id, &MultipleDatas[len(MultipleDatas) - 1], 1)
			}
			dec.Token()

			elapsed := time.Since(start)
			fmt.Printf("Cantidad [%v] Peso [%s] Tiempo [%v] .\n", i, FileSize(fi.Size()), elapsed)

		}
    }
    
    pass := &MyHandler{ cache: cache, }
    fasthttp.ListenAndServe(":80", pass.HandleFastHTTP)

}

func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	
	ctx.Response.Header.Set("Content-Type", "application/json")
    id := string(ctx.QueryArgs().Peek("id"))
	num := string(ctx.QueryArgs().Peek("num"))

	if num == "" {

		value, found := h.cache.Get(id)
		if !found {
			jsonFile, err := os.Open("json/"+id+".json")
			if err == nil{
				byteValue, _ := ioutil.ReadAll(jsonFile)
				fmt.Fprintf(ctx, string(byteValue))
			}else{
				fmt.Fprintf(ctx, "[]");
			}
		}else{
			json.NewEncoder(ctx).Encode(value)
		}

	}else{

		d, _ := strconv.Atoi(num)
		cant := len(id) - d
		fmt.Fprintf(ctx, "[");
		for n := 1; n <= d; n++ {
			nid := id[0:cant+n]
			value, found := h.cache.Get(nid)
			if !found {
				jsonFile, err := os.Open("json/"+nid+".json")
				if err == nil{
					byteValue, _ := ioutil.ReadAll(jsonFile)
					fmt.Fprintf(ctx, string(byteValue))
				}else{
					fmt.Fprintf(ctx, "[]");
				}
			}else{
				//f := *value.(*SingleData)
				//fmt.Printf("type  %T, valor %v\n", f, f)
				json.NewEncoder(ctx).Encode(value)
			}
			if n < d {
				fmt.Fprintf(ctx, ",");
			}
		}
		fmt.Fprintf(ctx, "]");
		
	}
	
}

func ParseAuto(s string) string {
	le := len(s) / 2
	bs := make([]byte, le)
	bl := 0
	for n := 0; n < le; n++ {
		d, _ := strconv.Atoi(s[n*2:n*2+2])
		j := d - 10
		bl += copy(bs[bl:], leng[j])
	}
	return string(bs)
}
func FileExists(name string) bool {
    if fi, err := os.Stat(name); err == nil {
        if fi.Mode().IsRegular() {
            return true
        }
    }
    return false
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