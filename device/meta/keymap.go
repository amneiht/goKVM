package meta

import "github.com/go-vgo/robotgo"

var keymap map[int32]string

func init() {
	// add 104 key borad
	keymap = make(map[int32]string)
	keymap[KEY_ESC] = robotgo.Escape
	keymap[KEY_0] = robotgo.Key0
	keymap[KEY_1] = robotgo.Key1
	keymap[KEY_2] = robotgo.Key2
	keymap[KEY_3] = robotgo.Key3
	keymap[KEY_4] = robotgo.Key4
	keymap[KEY_5] = robotgo.Key5
	keymap[KEY_6] = robotgo.Key6
	keymap[KEY_7] = robotgo.Key7
	keymap[KEY_8] = robotgo.Key8
	keymap[KEY_9] = robotgo.Key9
	keymap[KEY_MINUS] = "-"
	keymap[KEY_EQUAL] = "="
	keymap[KEY_BACKSPACE] = robotgo.Backspace
	keymap[KEY_TAB] = robotgo.Tab
	keymap[KEY_Q] = robotgo.KeyQ
	keymap[KEY_W] = robotgo.KeyW
	keymap[KEY_E] = robotgo.KeyE
	keymap[KEY_R] = robotgo.KeyR
	keymap[KEY_T] = robotgo.KeyT
	keymap[KEY_Y] = robotgo.KeyY
	keymap[KEY_U] = robotgo.KeyU
	keymap[KEY_I] = robotgo.KeyI
	keymap[KEY_O] = robotgo.KeyO
	keymap[KEY_P] = robotgo.KeyP
	keymap[KEY_LEFTBRACE] = "["
	keymap[KEY_RIGHTBRACE] = "]"
	keymap[KEY_ENTER] = robotgo.Enter
	keymap[KEY_LEFTCTRL] = robotgo.Lctrl
	keymap[KEY_CAPSLOCK] = robotgo.Capslock
	keymap[KEY_A] = robotgo.KeyA
	keymap[KEY_S] = robotgo.KeyS
	keymap[KEY_D] = robotgo.KeyD
	keymap[KEY_F] = robotgo.KeyF
	keymap[KEY_G] = robotgo.KeyG
	keymap[KEY_H] = robotgo.KeyH
	keymap[KEY_J] = robotgo.KeyJ
	keymap[KEY_K] = robotgo.KeyK
	keymap[KEY_L] = robotgo.KeyL
	keymap[KEY_SEMICOLON] = ";"
	keymap[KEY_APOSTROPHE] = "'"
	keymap[KEY_GRAVE] = "`"
	keymap[KEY_BACKSLASH] = "\\"
	keymap[KEY_Z] = robotgo.KeyZ
	keymap[KEY_X] = robotgo.KeyX
	keymap[KEY_C] = robotgo.KeyC
	keymap[KEY_V] = robotgo.KeyV
	keymap[KEY_B] = robotgo.KeyB
	keymap[KEY_N] = robotgo.KeyN
	keymap[KEY_M] = robotgo.KeyM
	keymap[KEY_COMMA] = ","
	keymap[KEY_DOT] = "."
	keymap[KEY_SLASH] = "/"
	keymap[KEY_RIGHTSHIFT] = robotgo.Rshift
	keymap[KEY_LEFTALT] = robotgo.Lalt
	keymap[KEY_LEFTSHIFT] = robotgo.Lshift
	keymap[KEY_SPACE] = robotgo.Space
	keymap[KEY_F1] = robotgo.F1
	keymap[KEY_F2] = robotgo.F2
	keymap[KEY_F3] = robotgo.F3
	keymap[KEY_F4] = robotgo.F4
	keymap[KEY_F5] = robotgo.F5
	keymap[KEY_F6] = robotgo.F6
	keymap[KEY_F7] = robotgo.F7
	keymap[KEY_F8] = robotgo.F8
	keymap[KEY_F9] = robotgo.F9
	keymap[KEY_NUMLOCK] = robotgo.NumLock
	// keymap[KEY_SCROLLLOCK] = robotgo.S
	keymap[KEY_F10] = robotgo.F10
	keymap[KEY_F11] = robotgo.F11
	keymap[KEY_F12] = robotgo.F12
	keymap[KEY_RIGHTCTRL] = robotgo.Rctrl
	keymap[KEY_RIGHTALT] = robotgo.Ralt

	keymap[KEY_HOME] = robotgo.Home
	keymap[KEY_UP] = robotgo.Up
	keymap[KEY_PAGEUP] = robotgo.Pageup
	keymap[KEY_LEFT] = robotgo.Left
	keymap[KEY_RIGHT] = robotgo.Right
	keymap[KEY_END] = robotgo.End
	keymap[KEY_DOWN] = robotgo.Down
	keymap[KEY_PAGEDOWN] = robotgo.Pagedown
	keymap[KEY_INSERT] = robotgo.Insert
	keymap[KEY_DELETE] = robotgo.Delete
	keymap[KEY_SYSRQ] = robotgo.Printscreen

	keymap[KEY_KP0] = robotgo.Num0
	keymap[KEY_KP1] = robotgo.Num1
	keymap[KEY_KP2] = robotgo.Num2
	keymap[KEY_KP3] = robotgo.Num3
	keymap[KEY_KP4] = robotgo.Num4
	keymap[KEY_KP5] = robotgo.Num5
	keymap[KEY_KP6] = robotgo.Num6
	keymap[KEY_KP7] = robotgo.Num7
	keymap[KEY_KP8] = robotgo.Num8
	keymap[KEY_KP9] = robotgo.Num9

	keymap[KEY_KPDOT] = robotgo.NumDecimal
	keymap[KEY_KPENTER] = robotgo.NumEnter
	keymap[KEY_KPSLASH] = robotgo.NumDiv
	keymap[KEY_KPMINUS] = robotgo.NumMinus
	keymap[KEY_KPASTERISK] = robotgo.NumMul
	keymap[KEY_KPPLUS] = robotgo.NumPlus

	keymap[KEY_COMPOSE] = robotgo.Menu
	keymap[KEY_LEFTMETA] = "lcmd"
	keymap[KEY_RIGHTMETA] = "rcmd"
}

func GetKey(code int32) (string, bool) {

	str, ok := keymap[code]
	return str, ok
}
