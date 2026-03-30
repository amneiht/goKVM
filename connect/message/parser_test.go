package message

import (
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/amneiht/goKVM/connect/message/data"
	"google.golang.org/protobuf/proto"
)

func TestParser(t *testing.T) {

	parser := NewParser(10240)
	mess := new(data.Message)

	buf := make([]byte, 10240)
	buf = buf[:0]
	// tao buff

	mess.Type = data.MessType_AUTH
	mess.Request = true
	parser.OnPasreDone = func(mess *data.Message) {
		fmt.Println("got mess")
	}
	// fmt.Println(len(buf))
	num := make([]byte, 4)
	for range 10 {
		mess.Type = data.MessType_CLIPBROAD
		sbuf, _ := proto.Marshal(mess)
		ul := uint32(len(sbuf))
		// fmt.Println(ul)
		binary.BigEndian.PutUint32(num, ul)
		buf = append(buf, num...)
		buf = append(buf, sbuf...)
	}
	// fmt.Println(buf)
	parser.Append(buf[:19])
	// fmt.Println("Parser: ", parser)
	parser.Append(buf[19:])
}
