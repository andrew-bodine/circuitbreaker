package circuitbreaker

// A Caller enables a circuit breaker to handle any kind of function,
// along with any inputs and outputs. This interface decouples the
// circuit breaker logic from the actual operation being wrapped.
type Caller interface {

	// Implement Call so the circuit breaker has an actual operation
	// to run, this is also how you will receive any return values
	// from said operation.
	Call(...interface{}) (interface{}, error)

	// Implement OnOpen if you want to be notified when the circuit
	// breaker state changes to open.
	OnOpen()

	// Implement OnClose if you want to be notified when the circuit
	// breaker state changes to closed.
	OnClose()
}
