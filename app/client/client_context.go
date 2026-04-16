//go:build linux

package client

import (
	"log"
	"strings"
	"time"

	"github.com/amneiht/goKVM/app"
	"github.com/amneiht/goKVM/connect"
	"github.com/amneiht/goKVM/device"
	"github.com/amneiht/goKVM/device/clipboard"
	"github.com/amneiht/goKVM/device/event"
	"github.com/holoplot/go-evdev"
	"gopkg.in/ini.v1"
	"gopkg.in/natefinch/lumberjack.v2"
)

type clientContext struct {
	log *lumberjack.Logger

	cp             *clipboard.CBService
	shareClipboard bool
	cap            *event.Capture
	// auto switch
	autoSwitch bool
	socks      [10]*connect.KVMSocket
	activeSock *connect.KVMSocket
	keySwitch  map[evdev.EvCode]bool
	acitveId   int
	letfSwitch bool
	// mointor size
	sizeS device.Vsize
	// runing control
	run  bool
	robo device.Robo
}

func NewClient(config string) *clientContext {
	client := new(clientContext)
	cfg, err := ini.Load(config)
	if err != nil {
		panic(err)
	}
	client.cp = clipboard.NewClipBroadService()
	client.cap = event.NewCapture()
	client.run = true
	gkey := []evdev.EvCode{evdev.KEY_RIGHTSHIFT, evdev.KEY_RIGHTCTRL}
	client.cap.SetKey(gkey)
	gb := cfg.Section(app.DEFAULTSESSION)
	logfile := gb.Key(app.LOG).String()
	client.robo = device.CreateWarrper()
	if len(logfile) > 0 {
		client.log = &lumberjack.Logger{
			Filename:   logfile,
			MaxSize:    10, // 10MB
			MaxBackups: 1,
			MaxAge:     28,
		}
		log.SetOutput(client.log)
	}
	// TODO : the only opotion avaiable is home screen size
	client.sizeS.X, client.sizeS.Y = client.robo.GetScreenSize()
	sw := gb.Key(app.SWITCH).String()
	if len(sw) > 0 {
		client.autoSwitch = true
		if strings.Compare(sw, app.MODELEFT) == 0 {
			client.letfSwitch = true
		}
	}
	list := cfg.Sections()

	logger := log.Default()
	for _, section := range list {
		if section.Name() == app.DEFAULTSESSION {
			continue
		}
		id, _ := section.Key(app.ID).Int()
		if id == 0 || id > 9 {
			if id > 9 {
				logger.Println("max support for 10 devices")
			}
			continue
		}
		psk := section.Key(app.PSK).String()
		port, _ := section.Key(app.PORT).Int()
		host := section.Key(app.HOST).String()

		if client.socks[id] != nil {
			logger.Println("id is duplicate on section", section.Name())
			continue
		}
		conn := connect.CreateSocket(psk, host, port)
		client.socks[id] = conn
	}
	client.shareClipboard, _ = gb.Key(app.CLIPBROAD).Bool()
	return client
}

func (t *clientContext) Write(data []byte) (int, error) {
	return t.activeSock.Write(data)
}
func (t *clientContext) Read(data []byte) (int, error) {
	return t.activeSock.Read(data)
}

func (t *clientContext) Close() {
	// TODO Control gorountie

}
func StartClient(config string) {
	ctx := NewClient(config)
	logger := log.Default()
	logger.Println("Create new context")
	for i := range ctx.socks {
		if ctx.socks[i] != nil {
			ctx.activeSock = ctx.socks[i]
			ctx.acitveId = i
			break
		}
	}
	log.Default().Println("Active socket on ", ctx.acitveId)
	if ctx.activeSock == nil {
		log.Default().Println("Config is no avaiable config")
		return
	}
	// logger := log.Default()
	ctx.cap.OnGrapChange = ctx.handleGrap
	// setting capppture control
	ctx.cap.OnEventChange = ctx.handleEvent

	ctx.cp.Init()
	if ctx.shareClipboard {
		// clipbroad send
		ctx.cp.OnChange = ctx.handleClipBroad
		go ctx.cp.StartService()
	}

	logger.Println("Start client")
	for {
		sock := ctx.activeSock
		if !sock.Connect() {
			time.Sleep(5 * time.Second)
			continue
		}
		// start service
		logger.Println("Connect to server")
		go ctx.cap.Start()
		ctx.run = true

		ctx.handleMessage()
		logger.Println("Disconnect server")
		ctx.cap.Stop()
		// clear socket
		ctx.activeSock.Disconnect()
		// change socket
		ctx.activeSock = ctx.socks[ctx.acitveId]
	}

}
