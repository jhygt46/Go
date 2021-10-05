package main

import (
    "testing"
	"fmt"
    "encoding/binary"
    "bytes"
)

func main() {
    fmt.Println("Hello World")
}

func read_int32(data []byte) (ret uint32) {
    buf := bytes.NewBuffer(data)
    binary.Read(buf, binary.LittleEndian, &ret)
    return
}

func BenchmarkCalculateA(b *testing.B) {

    byteNumber := []byte{0, 0, 0, 244, 0, 0, 0, 244}
    for i := 0; i < b.N; i++ {
        fmt.Println(read_int32(byteNumber))
    }

}

func BenchmarkCalculateB(b *testing.B) {

    byteNumber := []byte{0, 0, 0, 244, 0, 0, 0, 244}
    for i := 0; i < b.N; i++ {
        if len(byteNumber) > 7 {
            data := binary.BigEndian.Uint64(byteNumber)
            data++
        }
    }

}

func Calculate(x int) (result int) {
    return x + 2
}