package main

import (
	"fmt"
	"time"

	"github.com/amneiht/goKVM/device/clipboard"
	"github.com/amneiht/goKVM/device/event"
	"github.com/holoplot/go-evdev"
)

func testClib() {
	service := clipboard.NewClipBroadService()
	service.OnChange = func(data []byte) {
		fmt.Println(string(data))
	}
	go service.StartService()
	time.Sleep(30 * time.Second)
	service.Close()

}

func test_key() {
	cap := event.NewCapture()

	cap.OnGrapChange = func(b bool) {
		fmt.Println("Enter to grab mod:", b)
	}
	gkey := []evdev.EvCode{evdev.KEY_RIGHTSHIFT, evdev.KEY_RIGHTCTRL}
	cap.SetKey(gkey)
	// add grap mode
	go cap.Start()
	time.Sleep(30 * time.Second)
	cap.Close()
}
func main() {
	testClib()
}
