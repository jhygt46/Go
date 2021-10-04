package main

import (
    "testing"
	"fmt"
    "encoding/binary"
)

func main() {
    fmt.Println("Hello World")
}

func BenchmarkCalculateA(b *testing.B) {

    byteNumber := []byte{0, 0, 0, 0, 0, 0, 0, 54, 53, 54, 54, 54, 54, 54, 54, 54, 54, 54}
    for i := 0; i < b.N; i++ {
        data := binary.LittleEndian.Uint16(byteNumber)
        data++
    }

}

func BenchmarkCalculateB(b *testing.B) {

    byteNumber := []byte{0, 0, 0, 0, 0, 0, 0, 244}
    for i := 0; i < b.N; i++ {
        data := binary.BigEndian.Uint64(byteNumber)
        data++
    }

}

func Calculate(x int) (result int) {
    return x + 2
}