package gomol

import (
	"time"
)

var clock Clock = &RealClock{}

type Clock interface {
	Now() time.Time
}

type RealClock struct{}

func (*RealClock) Now() time.Time {
	return time.Now()
}
