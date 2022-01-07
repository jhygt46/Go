package main

import (
	"fmt"
	"sync"
)

type CacheObj struct {
	C string
}

func main() {
	pool := sync.Pool{
		New: func() interface{} {
			return &CacheObj{
				C: "1",
			}
		},
	}

	testObject := pool.Get().(*CacheObj)
	fmt.Println(testObject.C) // print 1
	testObject.C = "2"
	pool.Put(testObject)
}
