package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	setWindowsHookEx = user32.NewProc("SetWindowsHookExW")
	callNextHookEx   = user32.NewProc("CallNextHookEx")
	getMessage       = user32.NewProc("GetMessageW")
)

const (
	WH_KEYBOARD_LL = 13
	WM_KEYDOWN     = 0x0100
)

type KBDLLHOOKSTRUCT struct {
	VkCode    uint32
	ScanCode  uint32
	Flags     uint32
	Time      uint32
	ExtraInfo uintptr
}

func hookCallback(nCode int, wParam uintptr, lParam uintptr) uintptr {
	if nCode >= 0 {
		if wParam == WM_KEYDOWN {
			kbd := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))
			fmt.Println("Key:", kbd.VkCode)
		}
	}
	ret, _, _ := callNextHookEx.Call(0, uintptr(nCode), wParam, lParam)
	return ret
}

func main() {
	callback := syscall.NewCallback(hookCallback)

	hook, _, _ := setWindowsHookEx.Call(
		uintptr(WH_KEYBOARD_LL),
		callback,
		0,
		0,
	)

	var msg struct{}
	getMessage.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)

}
