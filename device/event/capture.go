package event

import (
	"log"
	"sync"
	"time"

	"github.com/bendahl/uinput"
	evdev "github.com/holoplot/go-evdev"
)

type Capture struct {
	run  bool // capture runiing
	cap  bool // capture loop
	grab bool
	// goroutin control
	wg sync.WaitGroup
	mu sync.Mutex
	// calback funtion
	OnGrapChange  func(bool)
	OnEventChange func(evtype uint16, evcode uint16, value int32)
	// key map
	grabMap map[evdev.EvCode]bool
	keyMap  map[evdev.EvCode]bool
	// device
	kinput uinput.Keyboard
	minput uinput.Mouse

	signal chan struct{}
}

func (t *Capture) Close() {

	t.cap = false
	if t.run {
		t.run = false
		<-t.signal
	}

	close(t.signal)

	t.kinput.Close()
	t.minput.Close()
}
func (t *Capture) Runing() bool {
	return t.run
}

func (t *Capture) IsGrab() bool {
	return t.grab
}

func (t *Capture) ClearInput() {
	for key := range t.keyMap {
		log.Default().Println("Clear key:", evdev.KEYNames[key])
		switch key {
		case evdev.BTN_LEFT:
			t.minput.LeftRelease()
		case evdev.BTN_RIGHT:
			t.minput.RightRelease()
		case evdev.BTN_MIDDLE:
			t.minput.MiddleRelease()
		default:
			t.kinput.KeyPress(int(key))
			// t.kinput.KeyUp(int(key))
		}
	}
	clear(t.keyMap)

}
func NewCapture() *Capture {
	cap := new(Capture)
	mouse, err := uinput.CreateMouse("/dev/uinput", []byte("go-mouse"))
	if err != nil {
		panic("Cannot create mouse input")
	}
	key, _ := uinput.CreateKeyboard("/dev/uinput", []byte("go-key"))
	cap.kinput = key
	cap.minput = mouse
	cap.run = true
	cap.signal = make(chan struct{})
	cap.keyMap = make(map[evdev.EvCode]bool)
	cap.grabMap = make(map[evdev.EvCode]bool)
	return cap
}
func check_grab(keymap map[evdev.EvCode]bool) bool {
	res := true
	for _, value := range keymap {
		res = res && value
	}
	return res
}
func (t *Capture) SetKey(mode []evdev.EvCode) {
	clear(t.grabMap)
	for _, key := range mode {
		t.grabMap[key] = false
	}
}
func (t *Capture) handle(ie *evdev.InputEvent) {

	// check grab mode
	if ie.Type == evdev.EV_KEY {
		// save state if not in grab mode
		_, ok := t.grabMap[ie.Code]
		if ok {
			t.grabMap[ie.Code] = ie.Value > 0
			check := check_grab(t.grabMap)
			if check == true {
				check = !t.grab
				t.GrabChange(check)
			}
		}
	}
	if t.OnEventChange != nil {
		t.OnEventChange(uint16(ie.Type), uint16(ie.Code), ie.Value)
	}
}
func (t *Capture) GrabChange(b bool) {
	t.grab = b
	if t.OnGrapChange != nil {
		t.OnGrapChange(t.grab)
	}

}
func captureMouse(t *Capture, dev string) {
	defer func() {
		if r := recover(); r != nil {
			t.cap = false
			log.Default().Println("Mouse capture Panic recover")
		}
	}()
	defer t.wg.Done()
	edev, _ := evdev.Open(dev)
	defer edev.Close()

	// edev.NonBlock()
	grab := false
	for t.cap {
		ie, err := edev.ReadOne()
		if err != nil {
			// fmt.Println(err)
			t.cap = false
			log.Default().Println("Error:", err)
			break
		}
		if ie.Type != evdev.EV_REL && ie.Type != evdev.EV_KEY {
			// loc su kien
			continue
		}
		t.handle(ie)

		if t.grab != grab {
			grab = t.grab
			if grab {
				log.Default().Println("Grab mouse")
				edev.Grab()
			} else {
				edev.Ungrab()
			}
		}
	}
}

func captureKeyBroad(t *Capture, dev string) {
	defer func() {
		if r := recover(); r != nil {
			t.cap = false
			log.Default().Println("Keyboard capture Panic recover")
		}
	}()

	defer t.wg.Done()
	edev, _ := evdev.Open(dev)
	defer edev.Close()
	// edev.NonBlock()
	grab := false
	// we must grap event first because readone is block until we press , so
	// first key is duplicate on two machine
	edev.Grab()

	for t.Runing() {
		ie, err := edev.ReadOne()
		if err != nil {
			log.Default().Println("Error:", err)
			t.cap = false
			break
		}
		if ie.Type != evdev.EV_KEY {
			continue
		}
		t.handle(ie)
		if t.grab != grab {
			grab = t.grab
			if grab {
				log.Default().Println("Grab keybroad")
				t.ClearInput()
			}
		}
		if !grab {
			if ie.Value > 0 {
				t.keyMap[ie.Code] = true
				// log.Default().Println("press key:", evdev.KEYNames[ie.Code])
				t.kinput.KeyDown(int(ie.Code))
			} else {
				// xoa key
				delete(t.keyMap, ie.Code)
				t.kinput.KeyUp(int(ie.Code))
			}
		}
	}
}

func (t *Capture) Start() {
	logger := log.Default()
	t.grab = false
	t.run = true
	for t.run {
		logger.Println("Find new Devive")
		dkey, dmouse := findDevice()
		if len(dmouse) == 0 && len(dkey) == 0 {
			time.Sleep(3 * time.Second)
			continue
		}
		logger.Println("Mouse is", dmouse)
		logger.Println("Keybroad is", dkey)
		t.wg.Add(2)
		t.cap = true
		go captureMouse(t, dmouse)
		go captureKeyBroad(t, dkey)

		t.wg.Wait()
	}
	t.signal <- struct{}{}
}

func (t *Capture) Stop() {
	t.cap = false
	if t.run {
		t.run = false
		<-t.signal
	}
}
