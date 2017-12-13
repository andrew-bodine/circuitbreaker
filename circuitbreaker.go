package circuitbreaker

func New(c Caller) CircuitBreaker {
    cb := &circuitBreaker{
        state:  make(chan State, 1),
        caller: make(chan Caller, 1),
    }

    cb.state <- CLOSED
    cb.caller <- c

    return cb
}

type CircuitBreaker interface {

    // State returns the state of a circuit breaker instance.
    State() State

    // Forces the circuit breaker into open state.
    Open()

    // Forces the circuit breaker into closed state.
    Close()

    // Calls the wrapped function and returns the result.
    Call(...interface{}) interface{}
}

type circuitBreaker struct {
    state   chan State
    caller  chan Caller
}

// Implement the CircuitBreaker interface.
func (cb *circuitBreaker) State() State {
    s := <- cb.state
    cb.state <- s

    return s
}

// Implement the CircuitBreaker interface.
func (cb *circuitBreaker) Open() {
    <- cb.state
    cb.state <- OPEN
}

// Implement the CircuitBreaker interface.
func (cb *circuitBreaker) Close() {
    <- cb.state
    cb.state <- CLOSED
}

// Implement the CircuitBreaker interface.
func (cb *circuitBreaker) Call(args ...interface{}) interface{} {
    c := <- cb.caller
    cb.caller <- c

    if c == nil {
        return nil
    }

    return c.Call(args)
}
