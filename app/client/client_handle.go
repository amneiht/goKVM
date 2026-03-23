package client

import (
	"log"

	"github.com/amneiht/goKVM/connect/message/data"
	"github.com/amneiht/goKVM/device/clipboard"
	"github.com/go-vgo/robotgo"
	"github.com/holoplot/go-evdev"
	"google.golang.org/protobuf/proto"
)

func (t *clientContext) handleEvent(evtype uint16, evcode uint16, value int32) {

	if t.cap.IsGrab() {
		var mevent = &data.Event{
			Type: int32(evtype), Code: int32(evcode), Value: value,
		}
		buf, _ := proto.Marshal(mevent)

		var mess = &data.Message{
			Request: true,
			Type:    data.MessType_EVENT,
			Payload: buf,
		}
		sendbuff, _ := proto.Marshal(mess)
		t.Write(sendbuff)
	} else if t.autoSwitch && evtype == evdev.EV_REL {
		// handle mouse vent if not in grab mod
		if evcode == evdev.REL_X {
			x, y := robotgo.Location()
			if t.letfSwitch && x == 0 {
				robotgo.Move(1, y)
				t.cap.GrabChange(true)
			} else if !t.letfSwitch && x == t.sizeScreen-1 {
				robotgo.Move(x-1, y)
				t.cap.GrabChange(true)

			}
		}
	}

}

func (t *clientContext) handleGrap(grab bool) {
	log.Default().Printf("Capture mode = %t\n", grab)
	var mess = &data.Message{
		Request: true,
		Type:    data.MessType_RELEASE}
	if grab {
		mess.Type = data.MessType_ENTER
	}
	sendbuff, _ := proto.Marshal(mess)
	t.Write(sendbuff)

}

func (t *clientContext) handleClipBroad(newClip []byte) {

	var mess = &data.Message{
		Type:    data.MessType_CLIPBROAD,
		Request: true,
		Payload: newClip,
	}
	buff, _ := proto.Marshal(mess)
	t.Write(buff)
	log.Default().Printf("Send %d byte to server", len(newClip))
}

func (t *clientContext) handleMessage() {

	buf := make([]byte, clipboard.MAXLENGTH+1024)
	logger := log.Default()
	var mess data.Message
	for t.run {
		n, err := t.Read(buf)
		if err == nil {
			proto.Unmarshal(buf[:n], &mess)
			logger.Println("get data from server")
			switch mess.Type {
			case data.MessType_CLIPBROAD:
				logger.Printf("Buffer from server %d\n", len(mess.Payload))
				t.cp.SetClipBoard(mess.Payload)
			case data.MessType_RELEASE:
				t.cap.GrabChange(false)
			}
		}
	}
}
