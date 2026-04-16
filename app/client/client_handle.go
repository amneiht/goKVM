//go:build linux

package client

import (
	"log"

	"github.com/amneiht/goKVM/app"
	"github.com/amneiht/goKVM/connect"
	"github.com/amneiht/goKVM/connect/message/data"
	"github.com/amneiht/goKVM/device/clipboard"
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
			x, y := t.robo.Location()
			if t.letfSwitch && x == 0 {
				t.robo.Move(app.DISTANCE, y)
				t.cap.GrabChange(true)
			} else if !t.letfSwitch && x == t.sizeS.X-1 {
				t.robo.Move(x-app.DISTANCE, y)
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
		// using x11
		x, y := t.robo.Location()
		point := &data.Point{
			X: int32(x),
			Y: int32(y),
		}
		mess.Payload, _ = proto.Marshal(point)
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
	for t.run && t.activeSock.Status() == connect.AUTH {
		n, err := t.Read(buf)
		if err == nil {
			err = proto.Unmarshal(buf[:n], &mess)
			logger.Println("get data from server")
			switch mess.Type {
			case data.MessType_CLIPBROAD:
				logger.Printf("Buffer from server %d\n", len(mess.Payload))
				t.cp.SetClipBoard(mess.Payload)
			case data.MessType_RELEASE:
				t.cap.GrabChange(false)
				if t.autoSwitch {
					var point = new(data.Point)
					proto.Unmarshal(mess.Payload, point)
					var y int = int(point.Y)
					if y > t.sizeS.Y {
						y = t.sizeS.Y
					}
					if t.letfSwitch {
						t.robo.Move(app.DISTANCE, y)
					} else {
						t.robo.Move(t.sizeS.X-app.DISTANCE, y)
					}
				}
			}
		}
	}
}
