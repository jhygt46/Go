package main

import (
    "testing"
	"fmt"
)

func main() {
    fmt.Println("Hello World")
}

func BenchmarkCalculate(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Calculate(2)
    }
}

func Calculate(x int) (result int) {
    return x + 2
}