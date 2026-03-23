package server

import (
	"errors"
	"log"
	"strings"

	"github.com/amneiht/goKVM/connect"
	"github.com/amneiht/goKVM/device/clipboard"
	"github.com/amneiht/goKVM/device/emulator"
	"gopkg.in/ini.v1"
	"gopkg.in/natefinch/lumberjack.v2"
)

type remoteState int

const (
	STATE_CONTROL remoteState = iota
	STATE_FREE
)

type serverContext struct {
	log *lumberjack.Logger

	clip *clipboard.CBService
	emu  *emulator.Device
	// auto switch
	autoSwitch bool
	letfSwitch bool
	state      remoteState
	// mointor size
	sizeScreen int
	// runing control
	run    bool
	listen *connect.KVMListener
	sock   *connect.KVMSocket
}

func newServerContext(str string) *serverContext {
	cfg, err := ini.Load(str)
	if err != nil {
		panic(err)
	}
	svctx := new(serverContext)
	gb := cfg.Section("global")
	file := gb.Key("log").String()
	if len(file) > 0 {
		svctx.log = &lumberjack.Logger{
			Filename:   file,
			MaxSize:    10, // 10MB
			MaxBackups: 1,
			MaxAge:     28,
		}
		log.SetOutput(svctx.log)
	}
	inf := gb.Key("listen").String()
	port, _ := gb.Key("port").Int()
	psk := gb.Key("psk").String()
	if len(inf) == 0 {
		inf = "0.0.0.0"
	}
	if port == 0 {
		port = 1357
	}

	svctx.listen = connect.NewListener(inf, port, psk)
	svctx.clip = clipboard.NewClipBroadService()
	svctx.emu = emulator.CreateVirtualDevice()
	sw := gb.Key("switch").String()
	if len(sw) > 0 {
		svctx.autoSwitch = true
		if strings.Compare(sw, "left") == 0 {
			svctx.letfSwitch = true
		}
	}
	return svctx
}

func (t *serverContext) Write(data []byte) (int, error) {
	if t.sock != nil {
		return t.sock.Write(data)
	} else {
		return 0, errors.New("No Available Sock")
	}
}
func (t *serverContext) Read(data []byte) (int, error) {

	if t.sock != nil {
		return t.sock.Read(data)
	} else {
		return 0, errors.New("No Available Sock")
	}
}
func StartServer(s string) {
	ctx := newServerContext(s)
	ctx.clip.OnChange = ctx.handleClipBroad

	go ctx.clip.StartService()
	defer ctx.clip.Close()
	ctx.listen.Start(ctx.startSession)
}
