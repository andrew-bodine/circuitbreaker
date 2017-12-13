package circuitbreaker

// A Caller enables a circuit breaker to handle any kind of function,
// along with any inputs and outputs. This interface decouples the
// circuit breaker logic from the actual operation being wrapped.
type Caller interface {

	// Implement Call so the circuit breaker knows what function
	// you actually want to call, and so it can return values
	// properly.
	Call(...interface{}) interface{}

	// Implement OnOpen if you want to be notified when the circuit
	// breaker state changes to open.
	OnOpen()

	// Implement OnClose if you want to be notified when the circuit
	// breaker state changes to closed.
	OnClose()
}
