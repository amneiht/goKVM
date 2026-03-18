package sharecb

import (
	"context"
	"fmt"
	"sync"

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
func equal(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
func trimStr(input []byte) []byte {
	// ret := strings.TrimFunc(string(input), func(r rune) bool {
	// 	return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	// })
	start := 0
	end := len(input)
	for input[start] <= ' ' && start < end {
		start++
	}
	end = end - 1
	for input[end] <= ' ' && end > start {
		end--
	}
	return input[start : end+1]

}
func (t *Watcher) SetClipBoard(input []byte) {

	t.mu.Lock()
	defer t.mu.Unlock()
	buf := trimStr(input)
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
		t.mu.Lock()
		if t.OnChange != nil && !equal(data, t.old) {
			// save old clipbroad
			if len(data) < MAXLENGTH {
				t.old = make([]byte, len(data))
				copy(t.old, data)
				t.OnChange(data)
			} else {
				fmt.Printf("Cannot send data over %d\n", MAXLENGTH)
			}
		}
		t.mu.Unlock()
	}
}
