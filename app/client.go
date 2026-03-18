package app

import (
	"crypto/tls"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/amneiht/goKVM/connect/message/data"
	"github.com/amneiht/goKVM/event"
	"github.com/amneiht/goKVM/event/sharecb"
	"google.golang.org/protobuf/proto"
)

func keepAlive(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()

	var mess = &data.Message{
		Request: true,
		Type:    data.MessType_KEEPALIVE,
	}
	buff, _ := proto.Marshal(mess)
	for {
		_, err := conn.Write(buff)
		if err != nil {
			break
		}
		time.Sleep(2 * time.Second)
	}
}
func runSession(conn net.Conn) {
	var wg sync.WaitGroup
	wg.Add(1)
	go keepAlive(conn, &wg)
	event.Init()

	var handle = event.NewHandle()
	defer handle.Close()

	watch := sharecb.CreateWatcher()
	defer watch.Close()

	watch.OnChange = func(newClip []byte) {

		var mess = &data.Message{
			Type:    data.MessType_CLIPBROAD,
			Request: true,
			Payload: newClip,
		}
		buff, _ := proto.Marshal(mess)
		conn.Write(buff)
		fmt.Println("Send data to server")
	}

	// run check session
	go watch.Check()

	handle.OnGapChange = func(gap bool) {
		fmt.Printf("Capture mode = %t\n", gap)
		if !gap {
			// sen release event
			var mess = &data.Message{
				Request: true,
				Type:    data.MessType_RELEASE}
			sendbuff, _ := proto.Marshal(mess)
			conn.Write(sendbuff)
		}
	}

	go handle.Start(func(dtype uint16, dcode uint16, dvalue int32) {
		var mevent = &data.Event{
			Type: int32(dtype), Code: int32(dcode), Value: dvalue,
		}
		buf, _ := proto.Marshal(mevent)

		var mess = &data.Message{
			Request: true,
			Type:    data.MessType_EVENT,
			Payload: buf,
		}
		sendbuff, _ := proto.Marshal(mess)
		_, err := conn.Write(sendbuff)
		if err != nil {
			handle.Stop()
		}
		// fmt.Printf("buffWrite :%d \n", len(sendbuff))
	})

	buf := make([]byte, sharecb.MAXLENGTH+1000)
	for handle.Run() {
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		n, err := conn.Read(buf)
		if err == nil {
			var mess data.Message
			proto.Unmarshal(buf[:n], &mess)

			switch mess.Type {
			case data.MessType_CLIPBROAD:
				fmt.Printf("Buffer from server %d\n", len(mess.Payload))
				watch.SetClipBoard(mess.Payload)
			}
		}
	}
	wg.Wait()
}
func ClientConnect(s *string) {

	conf := &tls.Config{
		InsecureSkipVerify: true,
	}
	var connect_str = *s + ":1597"

	conn, err := tls.Dial("tcp", connect_str, conf)
	fmt.Printf(" Connect to : %s\n", connect_str)
	if err != nil {
		panic(err)
	}
	runSession(conn)
	defer conn.Close()
}
