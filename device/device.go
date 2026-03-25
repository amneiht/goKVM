package device

import (
	"runtime"
	"sync"

	"github.com/go-vgo/robotgo"
)

type Vsize struct {
	X int
	Y int
}

// wrapper call for robot go
type Robo interface {
	Location() (int, int)
	Move(x, y int)
	GetScreenSize() (int, int)
	Close()
}

type rb struct {
	point     chan Vsize
	rlocation chan struct{}
	rsize     chan struct{}
	size      chan Vsize
	// conntrol
	stop chan struct{}
	wg   sync.WaitGroup
}

func (t *rb) Location() (x int, y int) {
	// t.mu.Lock()
	// defer t.mu.Unlock()
	var point Vsize
	t.rlocation <- struct{}{}
	// doc data
	point = <-t.size
	x = point.X
	y = point.Y
	return
}
func (t *rb) GetScreenSize() (x int, y int) {
	// t.mu.Lock()
	// defer t.mu.Unlock()
	// x, y = robotgo.GetScreenSize()
	var point Vsize
	t.rsize <- struct{}{}
	// get data
	point = <-t.size
	x = point.X
	y = point.Y
	return
}
func (t *rb) Move(x, y int) {
	// robotgo.Move(x, y)
	p := Vsize{X: x, Y: y}
	t.point <- p
}

func (t *rb) Close() {
	close(t.stop)
	t.wg.Wait()
	close(t.point)
	close(t.rlocation)
	close(t.size)
	close(t.rsize)
}
func (t *rb) moveInOne() {
	runtime.LockOSThread()
	defer t.wg.Done()
	run := true
	for run {
		select {
		case data := <-t.point:
			robotgo.Move(data.X, data.Y)
		case <-t.rlocation:
			x, y := robotgo.Location()
			var point = Vsize{X: x, Y: y}
			t.size <- point
		case <-t.rsize:
			x, y := robotgo.GetScreenSize()
			var point = Vsize{X: x, Y: y}
			t.size <- point
		case <-t.stop:
			run = false
		}

	}
}

func CreateWarrper() Robo {
	rbs := new(rb)
	rbs.point = make(chan Vsize)
	rbs.rlocation = make(chan struct{})
	rbs.rsize = make(chan struct{})
	rbs.size = make(chan Vsize)
	rbs.wg.Add(1)
	rbs.stop = make(chan struct{})
	go rbs.moveInOne()
	return rbs
}
