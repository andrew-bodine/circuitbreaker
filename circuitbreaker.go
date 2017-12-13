package circuitbreaker

import (
	"time"
)

const (
	MAXFAILS = 5
	TIMEOUT = time.Second
)

type CircuitBreaker interface {

	// State returns the state of a circuit breaker instance.
	State() State

	// Calls returns the number of times the wrapped function has
	// been called.
	Calls() int

	// Fails returns the number of times the wrapped function has
	// failed since the last successful call.
	Fails() int

	// Forces the circuit breaker into open state.
	Open()

	// Calls the wrapped function when in a closed state and
	// returns the result.
	Call(...interface{}) (interface{}, error)
}

func New(c Caller) CircuitBreaker {
	cb := &circuitBreaker{
		state:  make(chan State, 1),
		calls:	make(chan int, 1),
		fails:	make(chan int, 1),
		timer:	make(chan *time.Timer, 1),
		caller: make(chan Caller, 1),
	}

	cb.state <- Closed

	cb.calls <- 0
	cb.fails <- 0

	cb.timer <- nil

	cb.caller <- c

	return cb
}

type circuitBreaker struct {
	state	chan State

	calls	chan int
	fails	chan int

	timer	chan *time.Timer

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
func (cb *circuitBreaker) Fails() int {
	n := <- cb.fails
	cb.fails <- n

	return n
}

// Implement the CircuitBreaker interface.
func (cb *circuitBreaker) Open() {
	s := <-cb.state
	cb.state <- Open

	if s != Open {
		<- cb.timer
		cb.timer <- time.NewTimer(TIMEOUT)
	}
}

// Implement the CircuitBreaker interface.
func (cb *circuitBreaker) Call(args ...interface{}) (interface{}, error) {
	if cb.State() == Open {
		return nil, TrippedError
	}

	c := <-cb.caller
	cb.caller <- c
	if c == nil {
		return nil, nil
	}

	// Increment call count.
	cs := <- cb.calls
	cb.calls <- cs + 1

	r, err := c.Call(args...)
	if err != nil {
		n := <- cb.fails
		n += 1
		cb.fails <- n

		if n >= MAXFAILS {
			cb.Open()
		}
	} else {

		// If there was a successful call, then set fail
		// count to zero.
		<- cb.fails
		cb.fails <- 0
	}

	return r, err
}
