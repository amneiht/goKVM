package emulator

import "github.com/amneiht/goKVM/connect/message/data"

type Device interface {
	ClearKey()
	Close()
	Handle(mevt *data.Event)
}
