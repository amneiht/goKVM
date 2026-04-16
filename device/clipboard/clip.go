package clipboard

import (
	"context"
	"log"
	"sync"

	"github.com/amneiht/goKVM/util"
	"golang.design/x/clipboard"
)

const (
	MAXLENGTH = 10000
)

type CBService struct {
	// call back
	OnChange func(data []byte)

	//
	init bool
	run  bool
	// control
	wg  sync.WaitGroup
	mu  sync.Mutex
	ctx context.Context
	old []byte
}

func NewClipBroadService() *CBService {

	sv := new(CBService)
	sv.run = true
	sv.init = false

	sv.ctx = context.Background()
	return sv
}
func (t *CBService) Close() {
	t.run = false
	t.ctx.Done()

	t.wg.Done()
}
func TrimStr(input []byte) []byte {
	start := 0
	end := len(input) - 1
	for input[start] <= ' ' && start < end {
		start++
	}

	for input[end] <= ' ' && end > start {
		end--
	}
	return input[start : end+1]

}
func (t *CBService) SetClipBoard(input []byte) {

	if !t.init {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	buf := TrimStr(input)
	if len(buf) == 0 {
		log.Default().Println("Reject empty string")
		return
	}
	t.old = make([]byte, len(buf))
	copy(t.old, buf)

	clipboard.Write(clipboard.FmtText, buf)
}
func (t *CBService) Init() bool {

	err := clipboard.Init()
	t.init = err == nil
	if err != nil {
		log.Default().Println("Cannot init clipbroad ", err)
	}
	return t.init
}
func (t *CBService) StartService() {

	t.wg.Add(1)
	defer t.wg.Done()
	loger := log.Default()
	// loop for init service

	ch := clipboard.Watch(t.ctx, clipboard.FmtText)
	for data := range ch {
		if len(data) < 2 {
			continue
		}
		t.mu.Lock()
		if t.OnChange != nil && !util.Equal(data, t.old) {
			// save old clipbroad
			if len(data) < MAXLENGTH {
				t.old = make([]byte, len(data))
				copy(t.old, data)
				t.OnChange(data)
			} else {
				loger.Printf("Cannot send data over %d\n", MAXLENGTH)
			}
		}
		t.mu.Unlock()
	}
}
