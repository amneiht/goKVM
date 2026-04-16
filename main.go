package main

import (
	"flag"
	"fmt"
	"runtime"

	"github.com/amneiht/goKVM/app"
	"github.com/amneiht/goKVM/app/client"
	"github.com/amneiht/goKVM/app/server"
)

func printClientConfig() {
	fmt.Printf("# Config example\n")
	fmt.Printf("[global]\n")
	fmt.Printf("%s = /var/log/goKVM.log\n", app.LOG)
	fmt.Printf("#%s = left | right\n", app.SWITCH)
	fmt.Printf("%s = left     # move mose to left to switch change\n", app.SWITCH)
	fmt.Printf("%s = no    # share clibroad  default no\n", app.CLIPBROAD)
	fmt.Printf("\n[laptop]\n")
	fmt.Printf("%s = 1            # using for switch view\n", app.ID)
	fmt.Printf("%s = 1357\n", app.PORT)
	fmt.Printf("%s = Amneiht@12345\n", app.PSK)
	fmt.Printf("%s = 192.168.1.1\n", app.HOST)
	fmt.Printf("\n[laptop]\n")
	fmt.Printf("%s = 1            # using for switch view\n", app.ID)
	fmt.Printf("%s = 1357\n", app.PORT)
	fmt.Printf("%s = Amneiht@12345\n", app.PSK)
	fmt.Printf("%s = 192.168.1.w\n", app.HOST)

}
func printServerConfig() {
	fmt.Printf("# Config example\n")
	fmt.Printf("[global]\n")
	fmt.Printf("%s = /var/log/goKVM.log\n", app.LOG)
	fmt.Printf("%s = right    # move mose to top right to switch change\n", app.SWITCH)
	fmt.Printf("%s = 1357\n", app.PORT)
	fmt.Printf("%s = Amneiht@12345\n", app.PSK)
	fmt.Printf("%s = no    # share clibroad  default no\n", app.CLIPBROAD)
	fmt.Printf("%s = 0.0.0.0  # listen on all interface\n", app.LISTEN)
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
		// fmt.Printf("We need config file")
		panic("We need config file")
	}
	fmt.Println("using  config file ", *cfile)
	if *iServer {
		server.StartServer(*cfile)
	} else {
		if runtime.GOOS == "linux" {
			fmt.Println("Start client")
			client.StartClient(*cfile)
		} else {
			fmt.Println("Donot suppport on ", runtime.GOOS)
		}
	}

}
