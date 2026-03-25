package app

type ConnectState int

const (
	UNKNOW ConnectState = iota
	CONNECTING
	UNAUTH
	AUTH
	DISCONNECT
)

const (
	DISTANCE = 30
)

// View size
