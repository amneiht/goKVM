package server

import (
	"log"
	"runtime"

	"github.com/amneiht/goKVM/app"
	"github.com/amneiht/goKVM/connect"
	"github.com/amneiht/goKVM/connect/message/data"
	"github.com/amneiht/goKVM/device/clipboard"
	"github.com/amneiht/goKVM/device/meta"
	"github.com/go-vgo/robotgo"
	"google.golang.org/protobuf/proto"
)

func (t *serverContext) handleClipBroad(newClip []byte) {

	var mess = &data.Message{
		Type:    data.MessType_CLIPBROAD,
		Request: true,
		Payload: newClip,
	}
	buff, _ := proto.Marshal(mess)
	_, err := t.Write(buff)
	if err == nil {
		log.Default().Printf("Send %d byte to client", len(newClip))
	}
}
func sendRelease(t *serverContext, x int, y int) {
	mess := &data.Message{
		Type:    data.MessType_RELEASE,
		Request: true,
	}
	point := &data.Point{
		X: int32(x),
		Y: int32(y),
	}
	sbuf, _ := proto.Marshal(point)
	mess.Payload = sbuf
	buf, _ := proto.Marshal(mess)
	t.Write(buf)
}
func (t *serverContext) control(mess *data.Event) {
	if mess.Type == meta.EV_REL && mess.Code == meta.REL_X {
		x, y := robotgo.Location()
		if t.letfSwitch && x == 0 {

			robotgo.Move(app.DISTANCE, y)
			sendRelease(t, x, y)
		} else if !t.letfSwitch && x == t.sizeScreen.X-1 {
			robotgo.Move(x-app.DISTANCE, y)
			sendRelease(t, x, y)
		}
	}
}
func (t *serverContext) startSession(sock *connect.KVMSocket) {
	// khoa lai goroutin
	runtime.LockOSThread()
	t.sock = sock
	t.state = STATE_FREE
	buf := make([]byte, clipboard.MAXLENGTH+1024)
	if t.x11 {
		t.sizeScreen.X, t.sizeScreen.Y = robotgo.GetScreenSize()
	}
	logger := log.Default()
	logger.Println("Screen size is", t.sizeScreen)
	for {

		n, err := sock.Read(buf)
		if err != nil {
			break
		}
		var mess data.Message
		err = proto.Unmarshal(buf[:n], &mess)
		switch mess.Type {
		case data.MessType_EVENT:
			var mevent data.Event
			err = proto.Unmarshal(mess.Payload, &mevent)
			// fmt.Printf("Get event %d %d \n", mevent.Code, mevent.Type)
			t.emu.Handle(&mevent)
			if t.autoSwitch && t.x11 {
				t.control(&mevent)
			}
		case data.MessType_ENTER:
			if t.x11 && t.autoSwitch {
				var point data.Point
				proto.Unmarshal(mess.Payload, &point)
				// todo : check  for diffrent display resolution
				var y int = int(point.Y)
				if point.Y > int32(t.sizeScreen.Y) {
					y = t.sizeScreen.X
				}
				if t.letfSwitch {
					robotgo.Move(app.DISTANCE, y)
				} else {
					robotgo.Move(t.sizeScreen.X-app.DISTANCE, y)
				}

			}
		case data.MessType_RELEASE:
			t.emu.ClearKey()
		case data.MessType_CLIPBROAD:
			if mess.Request {
				t.clip.SetClipBoard(mess.Payload)
			}
		}
	}
}
