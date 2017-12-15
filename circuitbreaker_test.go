package circuitbreaker_test

import (
	"time"

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
			Expect(cb.State()).To(Equal(Closed))
			Expect(cb.Calls()).To(Equal(0))
			Expect(cb.Fails()).To(Equal(0))
		})

		Context("when provided an actual function", func() {
			BeforeEach(func() {
				cb = New(&MockCaller{})
			})

			It("returns a circuit breaker in closed state", func() {
				Expect(cb.State()).To(Equal(Closed))
				Expect(cb.Calls()).To(Equal(0))
				Expect(cb.Fails()).To(Equal(0))
			})
		})
	})

	Context("circuitBreaker", func() {

		// Test the CircuitBreaker implementation.
		Context("CircuitBreaker", func() {
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

			Context("Countdown()", func() {
				It("returns false", func() {
					Expect(cb.Countdown()).To(Equal(false))
				})
			})

			Context("Open()", func() {
				Context("when in closed state", func() {
					It("changes to open state", func() {
						cb.Open()
						Expect(cb.State()).To(Equal(Open))
					})

					It("starts a waiting countdown", func() {
						cb.Open()
						Expect(cb.Countdown()).To(Equal(true))
					})

					Context("after the waiting period is over", func() {
						var calls int

						BeforeEach(func() {
							cb = NewWithOptions(
								&MockCaller{},
								&Options{
									MaxFails: 2,
									Timeout:  time.Millisecond,
								},
							)

							for cb.State() != Open {
								cb.Call()
								calls = cb.Calls()
							}

							<-time.NewTimer(time.Millisecond).C
						})

						It("tries to call the operation", func() {
							for calls >= cb.Calls() {
								<-time.NewTimer(time.Millisecond).C
							}
						})

						It("handles error and success correctly", func() {
							if cb.Fails() == MAXFAILS {

								// Circuit breaker is still in error path, and
								// should be reseting state to open to wait.
								for cb.State() != Open {
									<-time.NewTimer(time.Millisecond).C
								}

								Expect(cb.Countdown()).To(Equal(true))
							} else {

								// Circuit breaker should be closed or
								// closing anytime now.
								for cb.State() != Closed {
									<-time.NewTimer(time.Millisecond).C
								}

								Expect(cb.Fails()).To(Equal(0))
							}
						})
					})
				})

				Context("when in open state", func() {
					It("doesn't do anything", func() {
						cb.Open()
						cb.Open()
						Expect(cb.State()).To(Equal(Open))
						Expect(cb.Countdown()).To(Equal(true))
					})
				})
			})

			Context("Call()", func() {
				Context("with a nil caller", func() {
					It("returns nil", func() {
						r, err := cb.Call()
						Expect(r).To(BeNil())
						Expect(err).To(BeNil())
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
							val, err := cb.Call()

							if err == nil {
								Expect(val).NotTo(BeNil())
								succeeded = true
								break
							}
						}
					})

					Context("when there is a failure", func() {
						BeforeEach(func() {
							failed := false

							for !failed {
								_, err := cb.Call()

								if err != nil {
									failed = true
								}
							}
						})

						It("increments the fail count", func() {
							Expect(cb.Fails()).To(Equal(1))
						})

						Context("then a success before fail threshold", func() {
							It("sets fail count to zero", func() {

								// Our MockCaller will succeed if we tell it to.
								_, err := cb.Call(true)
								Expect(err).To(BeNil())
								Expect(cb.Fails()).To(Equal(0))
							})
						})

						Context("and fail threshold is reached", func() {
							BeforeEach(func() {
								for cb.Fails() < MAXFAILS {
									cb.Call()
								}
							})

							It("changes state to open", func() {
								Expect(cb.State()).To(Equal(Open))
							})
						})
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

						It("all calls should return tripped error", func() {
							_, err := cb.Call()
							Expect(err).NotTo(BeNil())
							Expect(err).To(Equal(TrippedError))
						})
					})
				})
			})
		})
	})
})
