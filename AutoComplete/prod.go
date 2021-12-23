package main

import (
	"fmt"
	"utils"
)

type MyHandler struct {

}

var letras []int32 = []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36}
var lens float64 = float64(len(letras))

func main() {

	var post []int32 = []int32{36, 36, 36, 36, 36}
	num, err := utils.GetNum(post, letras, lens)
	if err == nil {
		fmt.Println(num)
	}else{
		fmt.Println(err, num)
	}

}