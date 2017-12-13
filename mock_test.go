package circuitbreaker_test

import (
    "math/rand"
    "time"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"

    "github.com/andrew-bodine/circuitbreaker"
)

// This function is meant to work most of the time, but it will randomly
// return a nil after a random timeout to simulate a problem one might
// have doing something like a network operation.
func faulty() interface{} {
    s := rand.NewSource(time.Now().UnixNano())
    r := rand.New(s)

    chaos := r.Intn(1024)
    if chaos < 1 {
        t := time.NewTimer(time.Second)
        <- t.C

        return nil
    }

    return chaos
}

// MockCaller is an example caller implementation to demonstrate
// how to wrap problematic operations.
type MockCaller struct {}

// Implement the circuitbreaker.Caller interface.
func (m *MockCaller) Call(args ...interface{}) interface{} {
    return faulty()
}

// Implement the circuitbreaker.Caller interface.
func (m *MockCaller) OnOpen() {}

// Implement the circuitbreaker.Caller interface.
func (m *MockCaller) OnClose() {}

// Mock is a made-up type that serves as an example type someone is
// trying to define, that has one or more methods which need to be
// protected with circuit breakers.
type Mock struct {

    // Still using the above faulty() implementation, but we'll
    // wrap it with a circuit breaker.
    faulty  circuitbreaker.CircuitBreaker
}

// This demonstrates how existing logic will simply utilize the
// circuit breaker in place of the wrapped method.
func (m *Mock) eventuallyCallsFaulty() interface{} {
    return m.faulty.Call()
}

var _ = Describe("circuitbreaker", func() {
    Context("faulty()", func() {
        Context("over many calls", func() {
            It("will succeed and fail at least once", func () {
                succeeded := false
                failed := false

                for !succeeded || !failed {
                    val := faulty()

                    if val != nil {
                        succeeded = true
                        continue
                    }

                    failed = true
                }

                Expect(succeeded).To(Equal(true))
                Expect(failed).To(Equal(true))
            })
        })
    })
})
