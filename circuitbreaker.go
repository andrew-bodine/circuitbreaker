package circuitbreaker

type CircuitBreaker interface {

	// State returns the state of a circuit breaker instance.
	State() State

	// Calls returns the number of times the wrapped function has
	// been called.
	Calls() int

	// Forces the circuit breaker into open state.
	Open()

	// Forces the circuit breaker into closed state.
	Close()

	// Calls the wrapped function when in a closed state and
	// returns the result.
	Call(...interface{}) interface{}
}

func New(c Caller) CircuitBreaker {
	cb := &circuitBreaker{
		state:  make(chan State, 1),
		calls:	make(chan int, 1),
		caller: make(chan Caller, 1),
	}

	cb.state <- CLOSED

	cb.calls <- 0
	cb.caller <- c

	return cb
}

type circuitBreaker struct {
	state	chan State

	calls	chan int
	caller	chan Caller
}

// Implement the CircuitBreaker interface.
func (cb *circuitBreaker) State() State {
	s := <-cb.state
	cb.state <- s

	return s
}

// Implement the CircuitBreaker interface.
func (cb *circuitBreaker) Calls() int {
	c := <- cb.calls
	cb.calls <- c

	return c
}

// Implement the CircuitBreaker interface.
func (cb *circuitBreaker) Open() {
	<-cb.state
	cb.state <- OPEN
}

// Implement the CircuitBreaker interface.
func (cb *circuitBreaker) Close() {
	<-cb.state
	cb.state <- CLOSED
}

// Implement the CircuitBreaker interface.
func (cb *circuitBreaker) Call(args ...interface{}) interface{} {
	if cb.State() == OPEN {
		return nil
	}

	c := <-cb.caller
	cb.caller <- c
	if c == nil {
		return nil
	}

	cs := <- cb.calls
	cb.calls <- cs + 1

	return c.Call(args)
}
