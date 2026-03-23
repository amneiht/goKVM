package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/amneiht/goKVM/connect"
)

func startServer(psk string) {
	fmt.Println("Start server with psk", psk)
	listener := connect.NewListener("0.0.0.0", 1257, psk)
	listener.Start(func(t *connect.KVMSocket) {
		buff := make([]byte, 4096)
		for {
			n, _ := t.Read(buff)
			fmt.Println(buff[:n])
		}
	})
}

func startClient(host string, psk string) {
	fmt.Println("Start Client with psk", psk)
	sock := connect.CreateSocket(psk, host, 1257)
	if sock.Connect() {
		for {
			_, err := sock.Write([]byte("Hello"))
			if err != nil {
				break
			}
			time.Sleep(1 * time.Second)
		}
	} else {
		fmt.Println("Connect false")
	}
}
func main() {

	server := flag.Bool("s", false, "Server mode")
	host := flag.String("c", "127.0.0.1", "Server ip adderess")
	psk := flag.String("psk", "Amneiht@1232323", "Pre share key")

	flag.Parse()

	fmt.Println("Connect mode server:", *server)
	if *server {
		startServer(*psk)
	} else {
		startClient(*host, *psk)
	}

}
