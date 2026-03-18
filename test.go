package main

import (
	"fmt"

	"github.com/amneiht/goKVM/event/sharecb"
)

func main() {

	input := []byte("   fucking wow shit      ")
	input2 := []byte("   fucking wow shit     ss")
	fmt.Println(string(sharecb.TrimStr(input)))
	fmt.Println(string(sharecb.TrimStr(input2)))
}
