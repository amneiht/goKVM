package main

import (
	"github.com/go-vgo/robotgo"
)

func main() {

	robotgo.KeyDown(robotgo.Shift)
	robotgo.KeyDown("[")
	robotgo.KeyUp(robotgo.Shift)
	robotgo.KeyUp("[")
}
