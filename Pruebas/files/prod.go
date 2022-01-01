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
	"github.com/valyala/fasthttp"
)

type MyHandler struct {}

func main() {

	escribir_file("/var/db1_test", 350000)

	h := &MyHandler{}
	fasthttp.ListenAndServe(":80", h.HandleFastHTTP)
	
}

func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {

	ctx.Response.Header.Set("Content-Type", "application/json")
	switch string(ctx.Path()) {
	case "/get":
		
		folderfile := getFolderFile64(random(350000))
		file, err := os.Open("/var/db1_test/"+folderfile)
		if err != nil{
			fmt.Println(err)
		}
		byteValue, err := ioutil.ReadAll(file)
		if err != nil{
			fmt.Println(err)
		}
		file.Close()
		fmt.Fprintf(ctx, string(byteValue))

	default:
		ctx.Error("Not Found", fasthttp.StatusNotFound)
	}
	
}

func escribir_file(path string, numb int){

	d1 := []byte("{\"C\":[{ \"T\": 1, \"N\": \"Nacionalidad\", \"V\": [\"Chilena\", \"Argentina\", \"Brasileña\", \"Uruguaya\"] }, { \"T\": 2, \"N\": \"Servicios\", \"V\": [\"Americana\", \"Rusa\", \"Bailarina\", \"Masaje\"] },{ \"T\": 3, \"N\": \"Edad\" }],\"E\": [{ \"T\": 1, \"N\": \"Rostro\" },{ \"T\": 1, \"N\": \"Senos\" },{ \"T\": 1, \"N\": \"Trasero\" }]}")
	c := 0

	aux := numb / 100

	now := time.Now()
	for n := 0; n < aux; n++ {
		folder := getFolder64(int64(n*100))
		newpath := filepath.Join(path, folder)
		err := os.MkdirAll(newpath, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			fmt.Println("FOLDER ERROR: ", err)
		}
		for i := 0; i < 100; i++ {
			err := os.WriteFile(path+"/"+folder+"/"+strconv.Itoa(i), d1, 0644)
			if err != nil {
				fmt.Println(err)
			}
			c++
		}
	}
	elapsed := time.Since(now)
	fmt.Printf("WRITES FILES %v [%s] c/u total %v\n", c, time_cu(elapsed, c), elapsed)

}
func time_cu(t time.Duration, c int) string {
	ms := float64(t / time.Nanosecond)
	res := ms / float64(c)
	var s string
	if res < 1000 {
		s = fmt.Sprintf("%.2f NanoSec", res)
	} else if res >= 1000 && res < 1000000{
		s = fmt.Sprintf("%.2f MicroSec", res/1000)
	} else {
		s = fmt.Sprintf("%.2f MilliSec", res/1000000)
	}
	return s
}
func divmod(numerator, denominator int64) (quotient, remainder int64) {
	quotient = numerator / denominator
	remainder = numerator % denominator
	return
}
func getFolder64(num int64) string {

	c1, n1 := divmod(num, 1000000)
	c2, n2 := divmod(n1, 10000)
	c3, _ := divmod(n2, 100)
	return strconv.FormatInt(c1, 10)+"/"+strconv.FormatInt(c2, 10)+"/"+strconv.FormatInt(c3, 10)

}
func getFolderFile64(num int64) string {

	c1, n1 := divmod(num, 1000000)
	c2, n2 := divmod(n1, 10000)
	c3, c4 := divmod(n2, 100)
	return strconv.FormatInt(c1, 10)+"/"+strconv.FormatInt(c2, 10)+"/"+strconv.FormatInt(c3, 10)+"/"+strconv.FormatInt(c4, 10)

}
func random(max int64) int64 {
	n, _ := rand.Int(rand.Reader, big.NewInt(max))
	return n.Int64()
}