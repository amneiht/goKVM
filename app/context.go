package app

import (
	"crypto/sha256"
	"log"

	"github.com/amneiht/goKVM/connect/message/data"
	"github.com/amneiht/goKVM/util"
	"gopkg.in/ini.v1"
	"gopkg.in/natefinch/lumberjack.v2"
)

type State int32

const (
	UNKNOW State = iota
	UNAUTH
	AUTH
	DISCONNECT
)

type AppCtx struct {
	Log                *lumberjack.Logger
	IsClipBroadSupport bool
	status             State
	psk                string
	port               int
}

func (t *AppCtx) Close() {
	if t.Log != nil {
		t.Log.Close()
	}
}

func (t *AppCtx) hashData(mess *data.Auth) []byte {
	str := mess.User + ":" + mess.Nonce + ":" + t.psk
	res := sha256.Sum256([]byte(str))
	return res[:]
}
func (t *AppCtx) checkUser(mess *data.Auth) bool {

	respone := t.hashData(mess)
	return util.Equal(respone[:], mess.Result)
}
func (t *AppCtx) CreateResult(mess *data.Auth) []byte {
	return t.hashData(mess)
}
func CreateContext(cfile *string) *AppCtx {

	ctx := new(AppCtx)
	cfg, err := ini.Load(*cfile)
	if err != nil {
		log.Fatal(err)
	}

	ctx.port, err = cfg.Section("global").Key("port").Int()
	if err != nil {
		ctx.port = 1357
	}
	file := cfg.Section("global").Key("log").String()
	if len(file) > 0 {
		ctx.Log = &lumberjack.Logger{
			Filename:   file,
			MaxSize:    10, // 10MB
			MaxBackups: 1,
			MaxAge:     28,
		}
		log.SetOutput(ctx.Log)
	}
	ctx.psk = cfg.Section("global").Key("psk").String()
	// log.Default().Printf("psk %s", ctx.psk)
	return ctx
}
