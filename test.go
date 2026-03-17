package main

import (
	"fmt"
	"time"

	"github.com/amneiht/goKVM/event/sharecb"
)

func main() {

	sharecb.Init()
	watch := sharecb.CreateWatcher()
	watch.OnChange = func(newClip []byte) {
		// watch.SetClipBoard(newClip)
		fmt.Printf("have data %s \n", newClip)
	}
	go watch.Check()
	buff := []byte("This is new clipboard data")
	for {
		time.Sleep(2 * time.Second)
		watch.SetClipBoard(buff)
	}
}
