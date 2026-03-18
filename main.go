package main

import (
	"flag"
	"fmt"

	"github.com/amneiht/goKVM/app"
)

func print_support() {
	fmt.Println("# Config example")
	fmt.Println("[global]\n")
	fmt.Println("port = 1357")
	fmt.Println("psk = Amneiht@12345")
	fmt.Println("log = /var/log/goKVM.log")
}
func main() {

	var iServer = flag.Bool("s", false, "run as server")
	var host = flag.String("c", "192.168.10.1", "server ip to connect")
	var support = flag.Bool("print", false, "print sample config")
	var cfile = flag.String("f", "", "Config file ")
	flag.Parse()

	if *support {
		print_support()
		return
	}

	if len(*cfile) == 0 {
		// fmt.Println("We need config file")
		panic("We need config file")
	}
	var ctx *app.AppCtx

	ctx = app.CreateContext(cfile)
	defer ctx.Close()
	if *iServer {
		app.StartServer(ctx)
	} else {
		// fmt.Println(host)?
		app.ClientConnect(ctx, host)
	}

}
