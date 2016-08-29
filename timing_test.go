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

func (s *GomolSuite) TestRealClockNow(c *C) {
	// This test is completely pointless because it's not something that can really
	// be tested but I was sick of seeing the red for a lack of a unit test.  So I created
	// this one and figure even on slow systems the two lines should be executed within
	// one second of each other. :P
	setClock(&realClock{})

	timeNow := time.Now()
	clockNow := clock().Now()

	diff := clockNow.Sub(timeNow)
	c.Check(diff < time.Second, Equals, true)
}
