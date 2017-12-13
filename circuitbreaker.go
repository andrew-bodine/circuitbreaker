package circuitbreaker

func New() CircuitBreaker {
    c := &circuitBreaker{
        state:  make(chan State, 1),
    }

    c.state <- CLOSED

    return c
}

type CircuitBreaker interface {

    // State returns the state of a circuit breaker instance.
    State() State

    // Forces the circuit breaker into open state.
    Open()

    // Forces the circuit breaker into closed state.
    Close()
}

type circuitBreaker struct {
    state   chan State
}

// Implement the CircuitBreaker interface.
func (c *circuitBreaker) State() State {
    s := <- c.state
    c.state <- s

    return s
}

// Implement the CircuitBreaker interface.
func (c *circuitBreaker) Open() {
    <- c.state
    c.state <- OPEN
}

// Implement the CircuitBreaker interface.
func (c *circuitBreaker) Close() {
    <- c.state
    c.state <- CLOSED
}
