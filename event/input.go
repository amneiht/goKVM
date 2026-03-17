package event

import (
	"log"

	evdev "github.com/gvalkov/golang-evdev"
)

func captureDevice(handle *Handle, device *evdev.InputDevice) {
	// mouse event

	isGap := handle.gap
	defer handle.wg.Done()

	for handle.isRun {
		events, err := device.Read()
		if err != nil {
			log.Fatal(err)
		}

		for _, ev := range events {
			handle.emitEvent(ev.Type, ev.Code, ev.Value)
		}

		if handle.gap != isGap {
			if handle.gap {
				device.Grab()
			} else {
				device.Release()
			}
			isGap = handle.gap
		}

	}
}
