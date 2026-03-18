package sharecb

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/amneiht/goKVM/util"
	"golang.design/x/clipboard"
)

type Watcher struct {
	ctx      context.Context
	OnChange func(newClip []byte)
	old      []byte
	mu       sync.Mutex
}

func CreateWatcher() *Watcher {
	watch := new(Watcher)
	watch.OnChange = nil
	watch.ctx = context.Background()
	return watch
}

func (t *Watcher) Close() {
	t.ctx.Done()
}

func TrimStr(input []byte) []byte {
	// ret := strings.TrimFunc(string(input), func(r rune) bool {
	// 	return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	// })
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
func (t *Watcher) SetClipBoard(input []byte) {

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
func Init() bool {
	err := clipboard.Init()
	if err != nil {
		fmt.Println(err.Error())
	}
	return err == nil
}
func (t *Watcher) Check() {

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
				log.Default().Printf("Cannot send data over %d\n", MAXLENGTH)
			}
		}
		t.mu.Unlock()
	}
}
