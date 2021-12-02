package main

import (
	"fmt"
    "context"
    "github.com/go-redis/redis/v8"
	"github.com/valyala/fasthttp"
)

type MyHandler struct {
	redis *redis.Client
}

var ctx = context.Background()

func main() {

	rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })

	err := rdb.Set(ctx, "buena", "Nelson", 0).Err()
    if err != nil {
        panic(err)
    }

	h := &MyHandler{ redis: rdb }
	fasthttp.ListenAndServe(":80", h.HandleFastHTTP)	

}

func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {

	ctx.Response.Header.Set("Content-Type", "application/json")
	switch string(ctx.Path()) {
	case "/get0":

		val, err := h.redis.Get(ctx, "buena").Result()
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(ctx, val)

	case "/get1":

	default:
		ctx.Error("Not Found", fasthttp.StatusNotFound)
	}
	
}

func ExampleClient() {

    rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })

    err := rdb.Set(ctx, "key", "value", 0).Err()
    if err != nil {
        panic(err)
    }

    val, err := rdb.Get(ctx, "key").Result()
    if err != nil {
        panic(err)
    }
    fmt.Println("key", val)

    val2, err := rdb.Get(ctx, "key2").Result()
    if err == redis.Nil {
        fmt.Println("key2 does not exist")
    } else if err != nil {
        panic(err)
    } else {
        fmt.Println("key2", val2)
    }
    // Output: key value
    // key2 does not exist
}