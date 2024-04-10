package protocol

type State int

const (
	StateHandshaking State = iota
	StateStatus
	StateLogin
)
