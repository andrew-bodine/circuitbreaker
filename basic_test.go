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
