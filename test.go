package main

import (
	"fmt"

	"github.com/amneiht/goKVM/event/sharecb"
)

func main() {

	input := []byte("   fucking wow shit   d   ")
	input2 := []byte(" ")
	fmt.Println(string(sharecb.TrimStr(input)))
	fmt.Println(string(sharecb.TrimStr(input2)))
}
