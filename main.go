package main

import (
	"flag"

	"github.com/amneiht/goKVM/app"
)

func main() {

	var iServer = flag.Bool("s", false, "run as server")
	var host = flag.String("c", "192.168.10.1", "server ip to connect")
	flag.Parse()
	if *iServer {
		app.StartServer()
	} else {
		// fmt.Println(host)?
		app.ClientConnect(host)
	}

}
