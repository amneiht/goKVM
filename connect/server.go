package connect

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"log"
	"math/big"
	"net"
	"time"

	"github.com/amneiht/goKVM/connect/message/data"
	"github.com/amneiht/goKVM/util"
	"google.golang.org/protobuf/proto"
)

const (
	MaxConnect = 1
)

type KVMListener struct {
	port       int
	psk        string
	host       string
	ln         net.Listener
	numConnect int
}

func generateTLSCert() tls.Certificate {

	//	copy from AI
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),

		Subject: pkix.Name{
			Organization: []string{"Go TLS"},
		},

		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(24 * time.Hour),

		KeyUsage: x509.KeyUsageKeyEncipherment |
			x509.KeyUsageDigitalSignature,

		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
		},
	}

	der, err := x509.CreateCertificate(
		rand.Reader,
		&template,
		&template,
		&priv.PublicKey,
		priv,
	)

	if err != nil {
		panic(err)
	}

	return tls.Certificate{
		Certificate: [][]byte{der},
		PrivateKey:  priv,
	}
}
func NewListener(host string, port int, psk string) *KVMListener {
	list := new(KVMListener)
	list.host = host
	list.port = port
	list.psk = psk

	logger := log.Default()
	cert := generateTLSCert()

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	sport := fmt.Sprintf("%s:%d", list.host, list.port)
	ln, err := tls.Listen("tcp", sport, config)
	if err != nil {
		logger.Fatal(err)
	}
	list.ln = ln
	list.numConnect = 0

	return list
}

func reject(conn net.Conn) {
	defer conn.Close()
	buf := []byte("Connect refuse")
	conn.Write(buf)
}

func authenticationUser(t *KVMSocket) bool {
	nonce, _ := util.RandomString(32)
	maxtry := 3
	mess := &data.Message{
		Request: true,
		Type:    data.MessType_AUTH,
	}

	authm := &data.Auth{
		Nonce:  nonce,
		Method: "sha256",
	}
	mess.Payload, _ = proto.Marshal(authm)
	sendbuff, _ := proto.Marshal(mess)
	logger := log.Default()
	readbuff := make([]byte, 2048)

	for t.state != AUTH && maxtry > 0 {

		_, err := t.Write(sendbuff)
		if err != nil {
			break
		}
		logger.Println("Request authentic from server")
		_, err = t.Read(readbuff)
		if err != nil {
			break
		}
		logger.Println("Got message")
		var rmess data.Message
		proto.Unmarshal(readbuff, &rmess)
		if !rmess.Request && rmess.Type != data.MessType_AUTH {
			maxtry = maxtry - 1
			continue
		}
		var mauth data.Auth
		proto.Unmarshal(rmess.Payload, &mauth)
		respone := t.hashData(&mauth)
		if util.Equal(respone, mauth.Result) {
			t.state = AUTH
			rmess := &data.Message{
				Type:    data.MessType_REGISTER,
				Request: true,
			}
			cbuff, _ := proto.Marshal(rmess)
			t.Write(cbuff)
			logger.Println("Register complease")
		}
		maxtry--

	}
	return t.state == AUTH
}

func (t *KVMListener) handle(process func(t *KVMSocket), conn net.Conn) {

	sock := new(KVMSocket)
	sock.conn = conn
	sock.state = UNAUTH
	sock.psk = t.psk
	defer sock.Disconnect()

	if authenticationUser(sock) {
		process(sock)
	}

	t.numConnect = 0
}
func (t *KVMListener) Start(process func(t *KVMSocket)) {

	loger := log.Default()
	for {
		conn, err := t.ln.Accept()
		if err != nil {
			continue
		}

		if t.numConnect == 0 {
			t.numConnect = 1
			go t.handle(process, conn)
		} else {
			loger.Println("Reject connect")
			go reject(conn)
		}
	}
}
