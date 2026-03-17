package sharecb

import (
	"context"
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
func (t *Watcher) SetClipBoard(buf []byte) {

	t.mu.Lock()
	defer t.mu.Unlock()

	t.old = make([]byte, len(buf))
	copy(t.old, buf)

	clipboard.Write(clipboard.FmtText, buf)
}
func Init() {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
}
func (t *Watcher) Check() {

	ch := clipboard.Watch(t.ctx, clipboard.FmtText)
	for data := range ch {
		t.mu.Lock()
		if t.OnChange != nil && !equal(data, t.old) {
			// save old clipbroad
			t.old = make([]byte, len(data))
			copy(t.old, data)
			t.OnChange(data)
		}
		t.mu.Unlock()
	}
}
