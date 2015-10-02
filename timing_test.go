package gomol

import (
	. "gopkg.in/check.v1"
	"time"
)

type TestClock struct {
	curTime time.Time
}

func NewTestClock(curTime time.Time) *TestClock {
	return &TestClock{curTime: curTime}
}

func (c *TestClock) Now() time.Time {
	return c.curTime
}

func (s *GomolSuite) TestTestClockNow(c *C) {
	realNow := time.Now().AddDate(0, 0, 1)

	clock = NewTestClock(realNow)

	c.Check(clock.Now(), Equals, realNow)
}
