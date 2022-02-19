package bytes

import (
	"encoding/binary"
	"fmt"
	"testing"
)

func m(x uint64, l int) {
	//fmt.Printf("%v/", x)
}
func m2(x uint8, l int) {
	//fmt.Printf("%v/", x)
}
func m3(x uint64) {
	fmt.Printf("%v|", x)
}
func Benchmark_GetMultipleBytes1(b *testing.B) {
	var bytes = []uint8{250, 255, 11, 1, 10}
	for i := 0; i < b.N; i++ {
		x, l := GetMultipleBytes(bytes[0:3], 0)
		m(x, l)
	}
}
func Benchmark_GetMultipleBytes2(b *testing.B) {
	var bytes = []uint8{250, 255, 11, 1, 10}
	for i := 0; i < b.N; i++ {
		x, l := GetMultipleBytes2(bytes[0:3])
		m(x, l)
	}
}
func Benchmark_GetMultipleBytes3(b *testing.B) {
	var bytes = []uint8{233, 255, 11, 1, 10}
	for i := 0; i < b.N; i++ {
		x := GetMultipleBytes3(bytes[0:3])
		m(x, 2)
	}
}
func Benchmark_GetMultipleBytes4(b *testing.B) {
	var bytes = []uint8{255, 254, 11, 1, 10}
	for i := 0; i < b.N; i++ {
		x := uint64(bytes[0:1][0])
		if x == 255 {
			x = Bytes2toInt64(bytes[0:2])
			if x == 65535 {
				x = Bytes3toInt64(bytes[0:3])
			}
		}
		m(x, 2)
	}
}

func Benchmark_2ByteInt(b *testing.B) {
	var bytes = []uint8{255, 255, 255, 255, 10}
	for i := 0; i < b.N; i++ {
		x := Bytes2toInt64(bytes[0:2])
		m(x, 2)
	}
}
func Benchmark_1ByteInt(b *testing.B) {
	var bytes = []uint8{255, 255, 255, 255, 10}
	for i := 0; i < b.N; i++ {
		x := bytes[0:1][0]
		m2(x, 2)
	}
}
func GetMultipleBytes2(bytes []uint8) (uint64, int) {

	switch bytes[1:2][0] {
	case 255:
		return Bytes3toInt64(bytes[0:3]), 3
	case 0:
		return uint64(bytes[0:1][0]), 1
	default:
		return Bytes2toInt64(bytes[0:2]), 2
	}
}
func GetMultipleBytes(bytes []uint8, j int) (uint64, int) {

	var res uint64
	res = Bytes2toInt64(bytes[j : j+2])
	if res < 65535 {
		return res + 1, 2
	} else {
		res = Bytes3toInt64(bytes[j : j+3])
		if res < 16777215 {
			return res + 2, 3
		} else {
			res = Bytes4toInt64(bytes[j : j+4])
			if res < 4294967295 {
				return res + 3, 4
			} else {
				res = Bytes5toInt64(bytes[j : j+5])
				if res < 4294967295 {
					return res + 4, 5
				}
			}
		}
	}

	return res, 0
}
func GetMultipleBytes3(bytes []uint8) uint64 {
	x := uint64(bytes[0:1][0])
	if x == 255 {
		x = Bytes2toInt64(bytes[0:2])
		if x == 65535 {
			x = Bytes3toInt64(bytes[0:3])
		}
	}
	return x
}

func Bytes2toInt64(b []uint8) uint64 {
	bytes := make([]byte, 6, 8)
	bytes = append(bytes, b...)
	return binary.BigEndian.Uint64(bytes)
}
func Bytes3toInt64(b []uint8) uint64 {
	bytes := make([]byte, 5, 8)
	bytes = append(bytes, b...)
	return binary.BigEndian.Uint64(bytes)
}
func Bytes4toInt64(b []uint8) uint64 {
	bytes := make([]byte, 4, 8)
	bytes = append(bytes, b...)
	return binary.BigEndian.Uint64(bytes)
}
func Bytes5toInt64(b []uint8) uint64 {
	bytes := make([]byte, 3, 8)
	bytes = append(bytes, b...)
	return binary.BigEndian.Uint64(bytes)
}

//go test -benchmem -bench=.
