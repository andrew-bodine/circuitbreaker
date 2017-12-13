package circuitbreaker_test

import (
    "sync"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"

    . "github.com/andrew-bodine/circuitbreaker"
)

var _ = Describe("circuitbreaker", func() {
    var cb CircuitBreaker

    BeforeEach(func() {
        cb = New()
    })

    Context("New()", func() {
        It("returns a circuit breaker in closed state", func() {
            Expect(cb.State()).To(Equal(CLOSED))
        })
    })

    Context("circuitBreaker", func() {

        // Test the CircuitBreaker implementation.
        Context("CircuitBreaker", func() {
            Context("State()", func() {
                It("is concurrently safe", func() {
                    var wg sync.WaitGroup

                    wg.Add(1)
                    go func() {
                        defer wg.Done()

                        for i := 0; i < 50; i++ {
                            cb.State()
                        }
                    }()

                    wg.Add(1)
                    go func() {
                        defer wg.Done()

                        for i := 0; i < 50; i++ {
                            if i % 2 == 0 {
                                cb.Close()
                                continue
                            }

                            cb.Open()
                        }
                    }()

                    wg.Wait()
                })
            })

            Context("Open()", func() {
                Context("when in closed state", func() {
                    It("changes to open state", func() {
                        cb.Open()
                        Expect(cb.State()).To(Equal(OPEN))
                    })
                })

                Context("when in open state", func() {})

                Context("when in half open state", func() {})
            })

            Context("Close()", func() {
                Context("when in open state", func() {
                    BeforeEach(func() {
                        cb.Open()
                    })

                    It("changes to closed state", func() {
                        cb.Close()
                        Expect(cb.State()).To(Equal(CLOSED))
                    })
                })

                Context("when in closed state", func() {})

                Context("when in half open state", func() {})
            })
        })
    })
})
