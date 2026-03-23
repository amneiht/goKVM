package client

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/amneiht/goKVM/connect"
	"github.com/amneiht/goKVM/device/clipboard"
	"github.com/amneiht/goKVM/device/event"
	"github.com/go-vgo/robotgo"
	"github.com/holoplot/go-evdev"
	"gopkg.in/ini.v1"
	"gopkg.in/natefinch/lumberjack.v2"
)

type clientContext struct {
	log *lumberjack.Logger

	cp  *clipboard.CBService
	cap *event.Capture
	// auto switch
	autoSwitch bool
	socks      [10]*connect.KVMSocket
	activeSock *connect.KVMSocket
	keySwitch  map[evdev.EvCode]bool
	acitveId   int
	letfSwitch bool
	// mointor size
	sizeScreen int
	// runing control
	run bool
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
	gb := cfg.Section("global")
	logfile := gb.Key("log").String()
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
	client.sizeScreen, _ = robotgo.GetScreenSize()
	sw := gb.Key("switch").String()
	if len(sw) > 0 {
		client.autoSwitch = true
		if strings.Compare(sw, "left") == 0 {
			client.letfSwitch = true
		}
	}
	list := cfg.Sections()

	logger := log.Default()
	for _, section := range list {
		if section.Name() == "global" {
			continue
		}
		id, _ := section.Key("id").Int()
		if id == 0 || id > 9 {
			if id > 9 {
				logger.Println("max support for 10 devices")
			}
			continue
		}
		psk := section.Key("psk").String()
		port, _ := section.Key("port").Int()
		host := section.Key("host").String()

		if client.socks[id] != nil {
			logger.Println("id is duplicate on section", section.Name())
			continue
		}
		conn := connect.CreateSocket(psk, host, port)
		client.socks[id] = conn
	}
	return client
}

func (t *clientContext) Write(data []byte) (int, error) {
	if t.activeSock != nil {
		return t.activeSock.Write(data)
	} else {
		return 0, errors.New("No Available Sock")
	}
}
func (t *clientContext) Read(data []byte) (int, error) {

	if t.activeSock != nil {
		return t.activeSock.Read(data)
	} else {
		return 0, errors.New("No Available Sock")
	}
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
	// clipbroad send
	ctx.cp.OnChange = ctx.handleClipBroad

	// start service
	go ctx.cap.Start()
	go ctx.cp.StartService()

	log.Default().Println("Start client")
	for {
		sock := ctx.activeSock
		if !sock.Connect() {
			time.Sleep(5 * time.Second)
		}
		ctx.run = true

		ctx.handleMessage()
		// clear socket
		ctx.activeSock.Disconnect()
		// change socket
		ctx.activeSock = ctx.socks[ctx.acitveId]
	}

}
