package circuitbreaker

type State int

const (
	CLOSED State = iota
	OPEN
	HALFOPEN
)
