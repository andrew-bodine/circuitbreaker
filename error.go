package circuitbreaker

import (
	"errors"
)

var TrippedError error = errors.New("Circuit breaker is tripped.")
