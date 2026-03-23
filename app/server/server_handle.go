package server

import (
	"log"

	"github.com/amneiht/goKVM/connect"
	"github.com/amneiht/goKVM/connect/message/data"
	"github.com/amneiht/goKVM/device/clipboard"
	"github.com/go-vgo/robotgo"
	"github.com/holoplot/go-evdev"
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
func (t *serverContext) control(mess *data.Event) {
	if mess.Type == evdev.EV_REL && mess.Code == evdev.REL_X {
		x, y := robotgo.Location()

		mess := &data.Message{
			Type:    data.MessType_RELEASE,
			Request: true,
		}
		buf, _ := proto.Marshal(mess)
		if t.letfSwitch && x == 0 {
			robotgo.Move(1, y)
			t.Write(buf)
		} else if !t.letfSwitch && x == t.sizeScreen-1 {
			robotgo.Move(x-1, y)
			t.Write(buf)
		}
	}
}
func (t *serverContext) startSession(sock *connect.KVMSocket) {
	t.sock = sock
	t.state = STATE_FREE
	buf := make([]byte, clipboard.MAXLENGTH+1024)
	t.sizeScreen, _ = robotgo.GetScreenSize()
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
			if t.autoSwitch {
				t.control(&mevent)
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
