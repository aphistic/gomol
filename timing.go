package gomol

import (
	"time"
)

var curClock gomolClock = &realClock{}

type gomolClock interface {
	Now() time.Time
}

func setClock(clock gomolClock) {
	curClock = clock
}
func clock() gomolClock {
	return curClock
}

type realClock struct{}

func (*realClock) Now() time.Time {
	return time.Now()
}
