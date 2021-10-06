package main

import (
    "os"
    "testing"
	"fmt"
    "strconv"
)

func main() {
    fmt.Println("Hello World")
}


func BenchmarkCalculateA(b *testing.B) {

    d1 := []byte{115, 111, 109, 101, 10}
    for i := 0; i < b.N; i++ {
        err := os.WriteFile("/var/Go/pruebas/utils/filtros/1/"+strconv.Itoa(i), d1, 0644)
        if err != nil {
            fmt.Println(err)
        }
    }

}






/*
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