package app

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/amneiht/goKVM/connect/message/data"
	"github.com/amneiht/goKVM/event"
	"github.com/amneiht/goKVM/event/sharecb"
	"github.com/amneiht/goKVM/util"
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
		conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
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

	logger := log.Default()
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
		logger.Printf("Send %d byte to server", len(newClip))
	}

	// run check session
	go watch.Check()

	handle.OnGapChange = func(gap bool) {
		logger.Printf("Capture mode = %t\n", gap)
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
		conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
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
			logger.Println("get data from server")
			switch mess.Type {
			case data.MessType_CLIPBROAD:
				logger.Printf("Buffer from server %d\n", len(mess.Payload))
				watch.SetClipBoard(mess.Payload)
			}
		}
	}
	wg.Wait()
}

func authentic(ctx *AppCtx, conn net.Conn) bool {
	user, _ := util.RandomString(10)
	maxtry := 3
	ctx.status = UNAUTH

	mess := &data.Message{
		Request: false,
		Type:    data.MessType_AUTH,
	}

	authm := &data.Auth{
		User:   user,
		Method: "sha256",
	}

	readbuff := make([]byte, 2048)
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	_, err := conn.Read(readbuff)
	if err != nil {
		return false
	}
	var rmess data.Message
	proto.Unmarshal(readbuff, &rmess)

	logger := log.Default()
	logger.Println("Read from server")

	for ctx.status != AUTH && maxtry > 0 {

		if !rmess.Request && rmess.Type != data.MessType_AUTH {
			maxtry = maxtry - 1
			continue
		}
		var mauth data.Auth
		proto.Unmarshal(rmess.Payload, &mauth)

		authm.Nonce = mauth.Nonce
		authm.Result = ctx.CreateResult(authm)

		p1, _ := proto.Marshal(authm)
		mess.Payload = p1
		p2, _ := proto.Marshal(mess)
		conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
		_, err = conn.Write(p2)
		if err != nil {
			break
		}
		logger.Println("Write authention message to server")
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		_, err = conn.Read(readbuff)
		if err != nil {
			break
		}
		logger.Println("Read new mess from server")
		proto.Unmarshal(readbuff, &rmess)
		if rmess.Type == data.MessType_REGISTER {
			ctx.status = AUTH
			logger.Println("Register complease")
			break
		}
		maxtry--

	}
	return ctx.status == AUTH

}
func ClientConnect(ctx *AppCtx, s *string) {

	conf := &tls.Config{
		InsecureSkipVerify: true,
	}
	connect_str := fmt.Sprintf("%s:%d", *s, ctx.port)

	conn, err := tls.Dial("tcp", connect_str, conf)
	log.Default().Printf(" Connect to : %s\n", connect_str)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	if !authentic(ctx, conn) {
		return
	}
	runSession(conn)

}
