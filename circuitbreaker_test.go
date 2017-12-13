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
		cb = New(nil)
	})

	Context("New()", func() {
		It("returns a circuit breaker in closed state", func() {
			Expect(cb.State()).To(Equal(CLOSED))
		})

		Context("when provided an actual function", func() {
			BeforeEach(func() {
				cb = New(&MockCaller{})
			})

			It("returns a circuit breaker in closed state", func() {
				Expect(cb.State()).To(Equal(CLOSED))
			})
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
							if i%2 == 0 {
								cb.Close()
								continue
							}

							cb.Open()
						}
					}()

					wg.Wait()
				})
			})

			Context("Calls()", func() {
				It("always returns zero", func() {
					Expect(cb.Calls()).To(Equal(0))
					cb.Call()
					Expect(cb.Calls()).To(Equal(0))
					cb.Call()
					Expect(cb.Calls()).To(Equal(0))
				})

				Context("with a valid caller that made a call", func() {
					BeforeEach(func() {
						cb = New(&MockCaller{})
						cb.Call()
					})

					It("returns an int greater than zero", func() {
						Expect(cb.Calls()).To(Equal(1))
					})
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

			Context("Call()", func() {
				Context("with a nil caller", func() {
					It("returns nil", func() {
						Expect(cb.Call()).To(BeNil())
					})
				})

				Context("with a valid caller", func() {
					BeforeEach(func() {
						cb = New(&MockCaller{})
					})

					It("increments the call count", func() {
						before := cb.Calls()
						cb.Call()
						Expect(cb.Calls()).To(Equal(before + 1))
					})

					// We know what output to expect because we have
					// MockCaller, the circuit breaker implementation
					// shouldn't care about what the output of an
					// arbitrary operation is, simply return it.
					It("returns the output", func() {
						succeeded := false

						for !succeeded {
							val := cb.Call()

							if val != nil {
								succeeded = true
								break
							}
						}

						Expect(succeeded).To(Equal(true))
					})

					Context("while in an open state", func() {
						BeforeEach(func() {
							cb.Open()
						})

						It("should not increment the call count", func() {
							before := cb.Calls()
							cb.Call()
							Expect(cb.Calls()).To(Equal(before))
						})

						It("all calls should return nil", func() {
							Expect(cb.Call()).To(BeNil())
						})
					})
				})
			})
		})
	})
})
