package circuitbreaker

import (
    "time"
)

type Options struct {
    MaxFails    int
    Timeout     time.Duration
}
