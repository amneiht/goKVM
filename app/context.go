package app

import "gopkg.in/natefinch/lumberjack.v2"

type AppCtx struct {
	Log                *lumberjack.Logger
	IsClipBroadSupport bool
}

func (t *AppCtx) Close() {
	if t.Log != nil {
		t.Log.Close()
	}
}

func CreateContext1(s *string) *AppCtx {
	ctx := new(AppCtx)
	ctx.Log = &lumberjack.Logger{
		Filename:   *s,
		MaxSize:    10, // 10MB
		MaxBackups: 1,
		MaxAge:     28,
	}
	return ctx
}

func CreateContext() *AppCtx {
	ctx := new(AppCtx)
	return ctx
}
