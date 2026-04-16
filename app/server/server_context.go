package server

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/amneiht/goKVM/app"
	"github.com/amneiht/goKVM/connect"
	"github.com/amneiht/goKVM/device"
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

	clip      *clipboard.CBService
	shareClip bool
	emu       emulator.Device
	// auto switch
	autoSwitch bool
	letfSwitch bool
	state      remoteState
	// mointor size
	sizeScreen device.Vsize
	// runing control
	run    bool
	x11    bool
	listen *connect.KVMListener
	sock   *connect.KVMSocket
	// robo   device.Robo
}

// for window
func readInput(str string) ([]byte, error) {

	data, err := os.ReadFile(str)
	if err != nil {
		return nil, err
	}
	k := 0
	for i, v := range data {
		if v > ' ' && v < '~' {
			k = i
			break
		}
	}
	ret := data[k:]
	return ret, nil
}
func newServerContext(str string) *serverContext {
	cfg, err := ini.Load(str)
	if err != nil {
		panic(err)
	}
	svctx := new(serverContext)
	gb := cfg.Section(app.DEFAULTSESSION)

	file := gb.Key(app.LOG).String()
	if len(file) > 0 {
		svctx.log = &lumberjack.Logger{
			Filename:   file,
			MaxSize:    10, // 10MB
			MaxBackups: 1,
			MaxAge:     28,
		}
		log.SetOutput(svctx.log)
	}
	inf := gb.Key(app.HOST).String()
	port, _ := gb.Key(app.PORT).Int()
	psk := gb.Key(app.PSK).String()
	if len(psk) == 0 {
		psk = "abc@12345"
	}
	if len(inf) == 0 {
		inf = "0.0.0.0"
	}
	log.Default().Println("psk = ", psk)
	if port == 0 {
		port = 1357
	}

	svctx.listen = connect.NewListener(inf, port, psk)
	svctx.clip = clipboard.NewClipBroadService()
	svctx.emu = emulator.CreateVirtualDevice()
	svctx.shareClip, err = gb.Key(app.CLIPBROAD).Bool()
	if err != nil {
		svctx.shareClip = false
	}

	svctx.x11 = true
	sw := gb.Key(app.SWITCH).String()
	if len(sw) > 0 {
		svctx.autoSwitch = true
		if strings.Compare(sw, app.MODELEFT) == 0 {
			svctx.letfSwitch = true
		}
	}
	// svctx.robo = device.CreateWarrper()

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
	ret := ctx.clip.Init()

	if ret == false {
		log.Default().Println("X11 system is not avaible")
		ctx.x11 = false
	}
	if ctx.shareClip {
		log.Default().Println("Enable clipbroad sharing")
		ctx.clip.OnChange = ctx.handleClipBroad
		go ctx.clip.StartService()
	}
	defer ctx.clip.Close()
	ctx.listen.Start(ctx.startSession)
}
