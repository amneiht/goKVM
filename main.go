package main

import (
	"flag"
	"log"

	"github.com/amneiht/goKVM/app"
)

func main() {

	var iServer = flag.Bool("s", false, "run as server")
	var host = flag.String("c", "192.168.10.1", "server ip to connect")
	var logdir = flag.String("l", "/var/log/goKVM.log", "log file")
	var nolog = flag.Bool("no-log", false, "disable log to file")
	flag.Parse()

	var ctx *app.AppCtx

	if !*nolog {
		ctx = app.CreateContext()
	} else {
		ctx = app.CreateContext1(logdir)
		log.SetOutput(ctx.Log)
	}

	defer ctx.Close()
	if *iServer {
		app.StartServer(ctx)
	} else {
		// fmt.Println(host)?
		app.ClientConnect(host)
	}

}
