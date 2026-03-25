package main

import (
	"flag"
	"fmt"

	"github.com/amneiht/goKVM/app/client"
	"github.com/amneiht/goKVM/app/server"
)

func printClientConfig() {
	fmt.Println("# Config example")
	fmt.Println("[global]")
	fmt.Println("log = /var/log/goKVM.log")
	fmt.Println("# switch = left | right")
	fmt.Println("switch = left     # move mose to left to switch change")
	fmt.Println("clipboard = no    # share clibroad  default no")
	fmt.Println("\n[laptop]")
	fmt.Println("id = 1            # using for switch view")
	fmt.Println("port = 1357")
	fmt.Println("psk = Amneiht@12345")
	fmt.Println("host = 192.168.1.1")
	fmt.Println("\n[pc1]")
	fmt.Println("id = 2  # using for switch view")
	fmt.Println("port = 1357")
	fmt.Println("psk = Amneiht@12345")
	fmt.Println("host = 192.168.1.2")

}
func printServerConfig() {
	fmt.Println("# Config example")
	fmt.Println("[global]")
	fmt.Println("log = /var/log/goKVM.log")
	fmt.Println("switch = right    # move mose to top right to switch change")
	fmt.Println("port = 1357")
	fmt.Println("psk = Amneiht@12345")
	fmt.Println("clipboard = no    # share clibroad  default no")
	fmt.Println("listen = 0.0.0.0  # listen on all interface")
}
func main() {

	var iServer = flag.Bool("s", false, "run as server")
	var supportc = flag.Bool("print-client", false, "print client sample config")
	var supports = flag.Bool("print-server", false, "print server sample config")
	var cfile = flag.String("f", "", "Config file ")

	flag.Parse()

	if *supportc {
		printClientConfig()
		return
	}

	if *supports {
		printServerConfig()
		return
	}

	if len(*cfile) == 0 {
		// fmt.Println("We need config file")
		panic("We need config file")
	}

	if *iServer {
		server.StartServer(*cfile)
	} else {
		fmt.Println("start client")
		client.StartClient(*cfile)
	}

}
