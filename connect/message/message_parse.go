package message

import (
	"encoding/binary"

	"github.com/amneiht/goKVM/connect/message/data"
	"github.com/amneiht/goKVM/util"
	"google.golang.org/protobuf/proto"
)

type Parser struct {
	blen        int
	tmpb        []byte
	OnPasreDone func(mess *data.Message)
}

func NewParser(bufferLeng int) *Parser {

	p := new(Parser)
	p.blen = bufferLeng
	p.tmpb = make([]byte, 0, bufferLeng+1024)
	p.tmpb = p.tmpb[:0]
	return p
}

func (t *Parser) Append(buf []byte) {
	end := len(buf)
	start := 0
	mess := &data.Message{}
	var temp []byte
	var n uint32

	var vadd int
	for start < end {
		nadd := end
		if nadd+len(t.tmpb) > t.blen {
			nadd = t.blen - len(t.tmpb)
		}
		t.tmpb = append(t.tmpb, buf[start:nadd]...)
		temp = t.tmpb
		// vadd = 0
		for {
			leng := len(temp)
			if leng < 4 {
				break
			}
			n = binary.BigEndian.Uint32(temp)
			// fmt.Println("n is ", n, "lengt is", leng)
			if n == 0 || n > uint32(leng-4) {
				// fmt.Println("break")
				break
			}
			ptemp := temp[4 : n+4]
			err := proto.Unmarshal(ptemp, mess)
			if err == nil {
				t.OnPasreDone(mess)
			}
			temp = temp[n+4:]
			vadd = vadd + 4 + int(n)
		}
		start = nadd
		// fmt.Println("Leng is ", len(t.tmpb))
		remain := len(t.tmpb) - vadd
		// fmt.Println("Clear size with  ", len(t.tmpb), "and leng ", remain)
		t.tmpb = util.ClearSlice(t.tmpb, remain)
		buf = buf[nadd:]

		if len(buf) < 5 || len(t.tmpb) < 5 {
			break
		}
	}
}
