package gomol

import (
	"time"

	. "gopkg.in/check.v1"
)

type testClock struct {
	curTime time.Time
}

func newTestClock(curTime time.Time) *testClock {
	return &testClock{curTime: curTime}
}

func (c *testClock) Now() time.Time {
	return c.curTime
}

func (s *GomolSuite) TestTestClockNow(c *C) {
	realNow := time.Now().AddDate(0, 0, 1)

	setClock(newTestClock(realNow))

	c.Check(clock().Now(), Equals, realNow)
}
