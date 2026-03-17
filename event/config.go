package event

import (
	"fmt"
	"strings"

	evdev "github.com/gvalkov/golang-evdev"
)

var mouseInput string = "/dev/input/event2"
var keyBroadInput string = "/dev/input/event3"

func Config(mouse string, keyBroad string) {
	mouseInput = mouse
	keyBroadInput = keyBroad
}

func Init() {
	// only for linux
	devices, _ := evdev.ListInputDevices()

	var setMouse = false
	var setKey = false
	for _, dev := range devices {
		// fmt.Printf("%s \n", dev.Name)
		if !setMouse {
			if strings.Contains(dev.Name, "Mouse") {
				setMouse = true
				mouseInput = dev.Fn
			}
		}
		if !setKey {
			if strings.Contains(dev.Name, "Keyboard") {
				setKey = true
				keyBroadInput = dev.Fn
			}
		}
		if setKey && setMouse {
			break
		}
	}
	fmt.Printf("Mouse in put is %s \n", mouseInput)
	fmt.Printf("KeyBroad put is %s \n", keyBroadInput)
}
