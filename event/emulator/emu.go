package emulator

import (
	"log"

	"github.com/amneiht/goKVM/connect/message/data"
	"github.com/bendahl/uinput"
	evdev "github.com/gvalkov/golang-evdev"
)

type Device struct {
	// add mutext heare
	mouse uinput.Mouse
	keyb  uinput.Keyboard
	input map[int32]int32
}

/*
create virtual device
*/
func CreateVirtualDevice() *Device {
	mouse, err := uinput.CreateMouse("/dev/uinput", []byte("go-mouse"))
	if err != nil {
		panic("Cannot create mouse input")
	}
	key, _ := uinput.CreateKeyboard("/dev/uinput", []byte("go-key"))

	var dev Device
	dev.mouse = mouse
	dev.keyb = key
	maps := make(map[int32]int32)
	dev.input = maps
	return &dev
}

/*
clear all key event trigger by app
*/
func (t *Device) ClearKey() {
	for key := range t.input {
		switch key {
		case evdev.BTN_LEFT:
			t.mouse.LeftRelease()
		case evdev.BTN_RIGHT:
			t.mouse.RightRelease()
		case evdev.BTN_MIDDLE:
			t.mouse.MiddleRelease()

		}
	}
	log.Default().Println("Clear all key")
	clear(t.input)
}
func (t *Device) Close() {
	t.mouse.Close()
	t.keyb.Close()
}
func (t *Device) Handle(mevt *data.Event) {

	switch mevt.Type {
	case evdev.EV_REL:
		// mouse handle
		switch mevt.Code {
		case evdev.REL_Y:
			if mevt.Value > 0 {
				t.mouse.MoveDown(mevt.Value)
			} else {
				t.mouse.MoveUp(-mevt.Value)
			}
		case evdev.REL_X:
			// t.mouse.MoveLeft(mevt.Value)
			if mevt.Value > 0 {
				t.mouse.MoveRight(mevt.Value)
			} else {
				t.mouse.MoveLeft(-mevt.Value)
			}
		case evdev.REL_WHEEL:
			t.mouse.Wheel(false, mevt.Value)
		}
	case evdev.EV_KEY:
		if mevt.Value > 0 {
			t.input[mevt.Code] = mevt.Value
		} else {
			delete(t.input, mevt.Code)
		}
		switch mevt.Code {
		case evdev.BTN_LEFT:
			if mevt.Value > 0 {
				t.mouse.LeftPress()
			} else {
				t.mouse.LeftRelease()
			}
		case evdev.BTN_RIGHT:
			if mevt.Value > 0 {
				t.mouse.RightPress()
			} else {
				t.mouse.RightRelease()
			}
		case evdev.BTN_MIDDLE:
			if mevt.Value > 0 {
				t.mouse.MiddlePress()
			} else {
				t.mouse.MiddleRelease()
			}
		default:
			if mevt.Value > 0 {
				t.keyb.KeyDown(int(mevt.Code))
			} else {
				t.keyb.KeyUp(int(mevt.Code))
			}
		}
	}
}
