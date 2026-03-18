package app

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"log"
	"math/big"
	"net"
	"time"

	"github.com/amneiht/goKVM/connect/message/data"
	"github.com/amneiht/goKVM/event/emulator"
	"github.com/amneiht/goKVM/event/sharecb"
	"google.golang.org/protobuf/proto"
)

func generateTLSCert() tls.Certificate {

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

var numConnect int = 0

func handle(conn net.Conn, dev *emulator.Device, ctx *AppCtx) {

	defer conn.Close()

	logger := log.Default()
	var watch *sharecb.Watcher = nil
	if ctx.IsClipBroadSupport {
		watch = sharecb.CreateWatcher()
		defer watch.Close()
		watch.OnChange = func(newClip []byte) {
			var mess = &data.Message{
				Type:    data.MessType_CLIPBROAD,
				Request: true,
				Payload: newClip,
			}
			buff, _ := proto.Marshal(mess)
			conn.SetReadDeadline(time.Now().Add(5 * time.Second))
			conn.Write(buff)
			logger.Printf("Send %d to client\n", len(newClip))
		}
		go watch.Check()
	}
	// run check session

	buf := make([]byte, sharecb.MAXLENGTH+1000)

	for {
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		n, err := conn.Read(buf)
		if err != nil {
			break
		}
		var mess data.Message
		err = proto.Unmarshal(buf[:n], &mess)
		switch mess.Type {
		case data.MessType_EVENT:
			var mevent data.Event
			err = proto.Unmarshal(mess.Payload, &mevent)
			// fmt.Printf("Get event %d %d \n", mevent.Code, mevent.Type)
			dev.Handle(&mevent)
		case data.MessType_RELEASE:
			dev.ClearKey()
		case data.MessType_CLIPBROAD:
			if mess.Request {
				if watch != nil {
					watch.SetClipBoard(mess.Payload)
				} else {
					logger.Println("UnSupport Clibroad")
				}
			}
		}
	}
	logger.Println("close connect")
	numConnect = numConnect - 1

}

func reject(conn net.Conn) {
	defer conn.Close()
	buf := []byte("Connect refuse")
	conn.Write(buf)
}

func StartServer(ctx *AppCtx) {
	cert := generateTLSCert()
	ctx.IsClipBroadSupport = sharecb.Init()

	logger := log.Default()
	logger.Println("Create virtual Output")
	dev := emulator.CreateVirtualDevice()
	defer dev.Close()

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	ln, err := tls.Listen("tcp", ":1597", config)
	if err != nil {
		log.Fatal(err)
	}

	logger.Println("TLS server listening :1597")

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		if numConnect == 0 {
			numConnect = 1
			go handle(conn, dev, ctx)
		} else {
			go reject(conn)
		}
	}
}
