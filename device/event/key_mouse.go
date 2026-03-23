package event

import (
	"strings"

	evdev "github.com/holoplot/go-evdev"
)

func keycheck(code []evdev.EvCode, keyA evdev.EvCode, keyB evdev.EvCode) bool {
	var checkA, CheckB bool
	for _, key := range code {

		switch key {
		case keyA:
			checkA = true
		case keyB:
			CheckB = true
		}
		if checkA && CheckB {
			break
		}
	}
	return checkA == true && CheckB == true
}
func findDevice() (pkey string, pmouse string) {
	list, _ := evdev.ListDevicePaths()
	fmouse := false
	fkey := false
	for _, path := range list {
		// logger.Println(path)
		if strings.Contains(path.Name, "go-") {
			// reject virtual device
			continue
		}
		d, err := evdev.Open(path.Path)
		if err != nil {
			continue
		}
		defer d.Close()

		if !fkey {
			ecode := d.CapableEvents(evdev.EV_KEY)

			if keycheck(ecode[:], evdev.KEY_RIGHTCTRL, evdev.KEY_RIGHTSHIFT) {
				fkey = true
				pkey = d.Path()
			}

		}
		if !fmouse {
			ecode := d.CapableEvents(evdev.EV_REL)
			if keycheck(ecode[:], evdev.REL_X, evdev.REL_Y) {
				fmouse = true
				pmouse = d.Path()
			}
		}
		if fkey && fmouse {
			break
		}
	}

	return

}
