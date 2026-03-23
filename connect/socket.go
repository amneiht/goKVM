package connect

import (
	"crypto/sha256"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/amneiht/goKVM/connect/message/data"
	"github.com/amneiht/goKVM/util"
	"google.golang.org/protobuf/proto"
)

type ConnectState int

const (
	UNKNOW ConnectState = iota
	CONNECTING
	UNAUTH
	AUTH
	DISCONNECT
)

type KVMSocket struct {
	psk   string
	host  string
	port  int
	state ConnectState
	conn  net.Conn
	wg    sync.WaitGroup
}

func CreateSocket(psk string, host string, port int) *KVMSocket {
	sock := new(KVMSocket)
	sock.psk = psk
	sock.host = host
	sock.port = port
	sock.state = UNKNOW
	return sock
}
func (t *KVMSocket) hashData(mess *data.Auth) []byte {
	str := mess.User + ":" + mess.Nonce + ":" + t.psk
	res := sha256.Sum256([]byte(str))
	return res[:]
}
func (t *KVMSocket) authentic() bool {
	user, _ := util.RandomString(10)
	maxtry := 3
	t.state = UNAUTH
	conn := t.conn
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

	for t.state != AUTH && maxtry > 0 {

		if !rmess.Request && rmess.Type != data.MessType_AUTH {
			maxtry = maxtry - 1
			continue
		}
		var mauth data.Auth
		proto.Unmarshal(rmess.Payload, &mauth)

		authm.Nonce = mauth.Nonce
		authm.Result = t.hashData(authm)

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
			t.state = AUTH
			logger.Println("Register complease")
			break
		}
		maxtry--

	}
	return t.state == AUTH

}
func (t *KVMSocket) keepAlive() {
	defer t.wg.Done()
	var mess = &data.Message{
		Request: true,
		Type:    data.MessType_KEEPALIVE,
	}
	buff, _ := proto.Marshal(mess)
	conn := t.conn
	for t.state == AUTH {
		conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
		_, err := conn.Write(buff)
		if err != nil {
			break
		}
		time.Sleep(2 * time.Second)
	}
}
func (t *KVMSocket) Disconnect() {
	t.state = DISCONNECT
	t.wg.Wait()
	t.conn.Close()
	t.conn = nil
}
func (t *KVMSocket) Write(data []byte) (int, error) {
	if t.conn != nil {
		t.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
		return t.conn.Write(data)
	} else {

		return 0, errors.New("Sock not avaiable")
	}
}
func (t *KVMSocket) Read(data []byte) (int, error) {
	if t.conn != nil {
		t.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		return t.conn.Read(data)
	} else {

		return 0, errors.New("Sock not avaiable")
	}
}

func (t *KVMSocket) Connect() bool {
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}
	connect_str := fmt.Sprintf("%s:%d", t.host, t.port)
	log.Default().Printf(" Connect to : %s\n", connect_str)
	var conn net.Conn
	var err error
	maxtry := 5
	for maxtry > 0 {
		conn, err = tls.Dial("tcp", connect_str, conf)
		if err != nil {
			time.Sleep(5 * time.Second)
			maxtry--
			continue
		}
		break
	}
	if maxtry <= 0 {
		return false
	}
	t.conn = conn
	// xac thuc
	ret := t.authentic()
	if ret {
		t.wg.Add(1)
		go t.keepAlive()
	}

	return ret
}
