package event

import (
	"fmt"
	"sync"

	evdev "github.com/gvalkov/golang-evdev"
)

type Handle struct {
	// context
	isRun bool
	wg    sync.WaitGroup

	// capture
	key   *evdev.InputDevice
	mouse *evdev.InputDevice
	// control function
	ctlset   bool
	shiftSet bool
	gap      bool

	OnTrigerEvent func(uint16, uint16, int32)
	OnGapChange   func(bool)
}

type Context struct {
}

func (t *Handle) Run() bool {
	return t.isRun
}
func (t *Handle) Stop() {
	t.isRun = false
}
func (t *Handle) emitEvent(dtype uint16, code uint16, value int32) {

	// loc cac su kien
	if dtype != evdev.EV_KEY && dtype != evdev.EV_REL {
		return
	}
	// handle event

	if dtype == evdev.EV_KEY {
		switch code {
		case evdev.KEY_RIGHTCTRL:
			t.ctlset = value > 0
		case evdev.KEY_RIGHTSHIFT:
			t.shiftSet = value > 0
		}
		if t.ctlset && t.shiftSet {
			// change gap mode
			t.gap = !t.gap
			if t.OnGapChange != nil {
				t.OnGapChange(t.gap)
			}
		}
	}
	// fmt.Printf("Got event %d %d %d\n", dtype, code, value)
	if t.gap {
		t.OnTrigerEvent(dtype, code, value)
	}
}

func (t *Handle) Close() {
	t.isRun = false
	t.wg.Wait()

	// read last event
	t.key.Read()
	t.mouse.Read()
	// release all
	t.key.Release()
	t.mouse.Release()

	fmt.Println("Close program")
}

func NewHandle() *Handle {
	var p Handle
	p.isRun = true
	p.ctlset = false
	p.shiftSet = false

	p.gap = false
	// config dev
	p.key, _ = evdev.Open(keyBroadInput)
	p.mouse, _ = evdev.Open(mouseInput)

	return &p
}

func (t *Handle) Start(triger func(uint16, uint16, int32)) {

	t.wg.Add(2)
	t.OnTrigerEvent = triger
	go captureDevice(t, t.mouse)
	go captureDevice(t, t.key)

	t.wg.Wait()
}
