package circuitbreaker

// The different states that a circuit breaker can be in.
type State int

const (
	CLOSED State = iota
	OPEN
	HALFOPEN
)
