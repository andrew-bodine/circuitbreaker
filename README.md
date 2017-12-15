# circuitbreaker
[![Build Status](https://travis-ci.org/andrew-bodine/circuitbreaker.svg?branch=master)](https://travis-ci.org/andrew-bodine/circuitbreaker)

A golang circuit breaker package.

The implementation provided by this package follows the well known circuit breaker pattern.

<img alt="Circuit Breaker Pattern"
     width="600px"
     src="https://docs.google.com/drawings/d/e/2PACX-1vTz1nf8TKay0Uc1YkmLpUT70xl4dTkyXjuRs5W_Sq3FoftdSRx1j4_gO32ulxla2vg8efrSOOM3rOE7/pub?w=960&h=720"/>

## Example
```go
package circuitbreaker_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/andrew-bodine/circuitbreaker"
)

type HttpCaller struct{}

func (h *HttpCaller) Call(args ...interface{}) (interface{}, error) {
	method := args[0].(string)
	resource := args[1].(string)

	switch method {
	default:
		return http.Get(resource)
	}
}

// Validate basic usage.
var _ = Describe("circuitbreaker", func() {
	var cb circuitbreaker.CircuitBreaker

	BeforeEach(func() {
		cb = circuitbreaker.New(&HttpCaller{})
	})

	Context("basic usage", func() {
		Context("with an endpoint that is never down", func() {
			It("happily fetches the content", func() {
				resp, err := cb.Call(http.MethodGet, "https://www.google.com")
				Expect(err).To(BeNil())
				Expect(resp).NotTo(BeNil())
			})
		})
	})
})
```

> NOTE: An example of integrating a circuit breaker into a struct is covered in more detail in [mock_test.go](./mock_test.go)

## Documentation
All documentation can be found [here](https://godoc.org/github.com/andrew-bodine/circuitbreaker)
