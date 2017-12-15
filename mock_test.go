package circuitbreaker_test

import (
	"errors"
	"math/rand"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/andrew-bodine/circuitbreaker"
)

// Tests for mocks.
var _ = Describe("circuitbreaker", func() {
	Context("faulty()", func() {
		Context("when told to succeed", func() {
			It("returns success", func() {
				r, err := faulty(true)
				Expect(err).To(BeNil())
				Expect(r).NotTo(BeNil())
			})
		})

		Context("over many calls", func() {
			It("will succeed and fail at least once", func() {
				succeeded := false
				failed := false

				for !succeeded || !failed {
					_, err := faulty()

					if err != nil {
						failed = true
						continue
					}

					succeeded = true
				}

				Expect(succeeded).To(Equal(true))
				Expect(failed).To(Equal(true))
			})
		})
	})
})

// This function is meant to work most of the time, but it will randomly
// return a nil after a random timeout. It tries to simulate a problem someone
// might have doing something like a network bound operation.
func faulty(args ...interface{}) (interface{}, error) {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	chaos := r.Intn(10)

	succeedRegardless := false
	if len(args) > 0 {
		val, ok := args[0].(bool)

		if ok {
			succeedRegardless = val
		}
	}

	if chaos < 5 && !succeedRegardless {
		t := time.NewTimer(time.Millisecond)
		<-t.C

		return nil, errors.New("Weird error happened while trying to faulty()")
	}

	return chaos, nil
}

// Mock is a made-up type that serves as an example type someone is
// trying to define, that has one or more methods which need to be
// protected with circuit breakers.
type Mock struct {

	// Still using the above faulty() implementation, but we'll
	// wrap it with a circuit breaker.
	faulty circuitbreaker.CircuitBreaker
}

// This demonstrates how existing logic will simply utilize the
// circuit breaker in place of the wrapped method.
func (m *Mock) eventuallyCallsFaulty() interface{} {
	result, err := m.faulty.Call()

	// Circuit breaker is tripped, if this is the first trip then we just
	// started failing fast, otherwise this was fail fast.
	if err != nil {
		return err
	}

	return result
}

// MockCaller is an example caller implementation to demonstrate
// how to wrap problematic operations.
type MockCaller struct{}

// Implement the circuitbreaker.Caller interface.
func (m *MockCaller) Call(args ...interface{}) (interface{}, error) {

	// Still run the same faulty() operation as before, but it will be
	// wrapped inside this Call() implementation.
	return faulty(args...)
}
