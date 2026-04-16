package meta

import (
	"fmt"
	"testing"

	"github.com/go-vgo/robotgo"
)

func TestKey(t *testing.T) {

	for _, value := range keymap {
		fmt.Println("Type: ", value)
		robotgo.KeyPress(value)
		robotgo.KeyUp(value)
	}
}
