package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"resource/utils"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/valyala/fasthttp"
)

type Filtro struct {
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

type MyHandler struct {
	redis *redis.Client
	Total int64 `json:"Total"`
}

var ctxs = context.Background()

func main() {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	i, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err == nil {

		filtro := Filtro{}
		filtro.C = []Campos{Campos{T: 1, N: "Procesador", V: []string{"X15", "IntelC3", "Amd71", "X15", "IntelC3", "Amd71", "X15", "IntelC3", "Amd71", "X15", "IntelC3", "Amd71", "X15", "IntelC3", "Amd71", "X15", "IntelC3", "Amd71", "X15", "IntelC3", "Amd71", "X15", "IntelC3", "Amd71", "X15", "IntelC3", "Amd71", "X15", "IntelC3", "Amd71", "X15", "IntelC3", "Amd71"}}, Campos{T: 1, N: "Pantalla", V: []string{"4", "4.5", "5.5", "4.5", "5.5", "4.5", "5.5", "4.5", "5.5", "4.5", "5.5"}}, Campos{T: 1, N: "Memoria", V: []string{"2GB", "4GB", "8GB", "16GB", "32GB", "64GB", "128GB"}}, Campos{T: 1, N: "Marca", V: []string{"Samsung", "Motorola", "Nokia", "Samsung", "Motorola", "Nokia", "Samsung", "Motorola", "Nokia", "Samsung", "Motorola", "Nokia", "Samsung", "Motorola", "Nokia"}}}
		filtro.E = []Evals{Evals{T: 1, N: "Buena"}, Evals{T: 1, N: "Nelson"}, Evals{T: 1, N: "Hola"}, Evals{T: 1, N: "Mundo"}}

		for x := 1; x <= int(i); x++ {
			filtro.Id = int32(x)

			u, err := json.Marshal(filtro)
			if err == nil {
				err := rdb.Set(ctxs, strconv.Itoa(x), u, 0).Err()
				if err != nil {
					panic(err)
				}
			}

		}

		h := &MyHandler{redis: rdb, Total: i}
		fasthttp.ListenAndServe(":80", h.HandleFastHTTP)

	}

}

func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {

	ctx.Response.Header.Set("Content-Type", "application/json")
	switch string(ctx.Path()) {
	case "/get0":

		val, err := h.redis.Get(ctxs, strconv.FormatInt(utils.Random(h.Total), 10)).Result()
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(ctx, string(val))

	case "/get1":

	default:
		ctx.Error("Not Found", fasthttp.StatusNotFound)
	}

}

/*
func ExampleClient() {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Set(ctxs, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctxs, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctxs, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
}
*/
