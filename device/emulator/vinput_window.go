//go:build windows

package emulator

import (
	"log"
	"syscall"

	"github.com/amneiht/goKVM/connect/message/data"
	"github.com/amneiht/goKVM/device/meta"
	"github.com/go-vgo/robotgo"
)

var (
	user32         = syscall.NewLazyDLL("user32.dll")
	procKeybdEvent = user32.NewProc("keybd_event")
)

type winDevice struct {
	input map[int32]int32
}

/*
create virtual device
*/
func CreateVirtualDevice() Device {

	var dev = new(winDevice)
	maps := make(map[int32]int32)
	dev.input = maps
	return dev
}

/*
clear all key event trigger by app
*/

func keyPress(code int32) {

	str, ok := meta.GetKey(code)
	if ok {
		robotgo.KeyDown(str)
	}

}
func keyUp(code int32) {
	str, ok := meta.GetKey(code)
	if ok {
		robotgo.KeyUp(str)
	}

}
func (t *winDevice) ClearKey() {
	for key := range t.input {
		switch key {
		case meta.BTN_LEFT:
			robotgo.MouseUp(robotgo.Mleft)
		case meta.BTN_RIGHT:
			robotgo.MouseUp(robotgo.Mright)
		case meta.BTN_MIDDLE:
			robotgo.MouseUp(robotgo.Center)
		default:
			keyUp(key)
		}
	}
	log.Default().Println("Clear all key")
	clear(t.input)
}
func (t *winDevice) Close() {

}
func (t *winDevice) Handle(mevt *data.Event) {

	switch mevt.Type {
	case meta.EV_REL:
		// mouse handle
		switch mevt.Code {
		case meta.REL_Y:
			/*
				if mevt.Value > 0 {
					t.mouse.MoveDown(mevt.Value)
				} else {
					t.mouse.MoveUp(-mevt.Value)
				}
			*/
			robotgo.MoveRelative(0, int(mevt.Value))
		case meta.REL_X:
			/*
				if mevt.Value > 0 {
					t.mouse.MoveRight(mevt.Value)
				} else {
					t.mouse.MoveLeft(-mevt.Value)
				} */
			robotgo.MoveRelative(int(mevt.Value), 0)
		case meta.REL_WHEEL:
			// t.mouse.Wheel(false, mevt.Value)
			// robotgo.WheelLeft
			if mevt.Value > 0 {
				robotgo.ScrollDir(int(mevt.Value), "up")
			} else {
				robotgo.ScrollDir(int(-mevt.Value), "down")
			}
		}
	case meta.EV_KEY:
		if mevt.Value > 0 {
			t.input[mevt.Code] = mevt.Value
		} else {
			delete(t.input, mevt.Code)
		}
		switch mevt.Code {
		case meta.BTN_LEFT:
			if mevt.Value > 0 {
				robotgo.MouseDown(robotgo.Mleft)
			} else {
				robotgo.MouseUp(robotgo.Mleft)
			}
		case meta.BTN_RIGHT:
			if mevt.Value > 0 {
				robotgo.MouseDown(robotgo.Mright)
			} else {
				robotgo.MouseUp(robotgo.Mright)
			}
		case meta.BTN_MIDDLE:
			if mevt.Value > 0 {
				robotgo.MouseDown(robotgo.Center)
			} else {
				robotgo.MouseUp(robotgo.Center)
			}
		default:
			if mevt.Value > 0 {
				keyPress(mevt.Code)
			} else {
				keyUp(mevt.Code)
			}
		}
	}
}
