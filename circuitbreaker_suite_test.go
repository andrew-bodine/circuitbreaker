package circuitbreaker_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCircuitbreaker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Circuitbreaker Suite")
}
