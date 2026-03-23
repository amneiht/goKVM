package main

import (
	"fmt"
	"time"

	"github.com/go-vgo/robotgo"
)

func main() {

	for {
		x, y := robotgo.Location()
		fmt.Println("Mouse location:", x, y)
		time.Sleep(50 * time.Millisecond)

	}
}
