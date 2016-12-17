package gomol

import (
	"testing"
	"time"

	. "github.com/onsi/gomega"
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

func (s *GomolSuite) TestTestClockNow(t *testing.T) {
	realNow := time.Now().AddDate(0, 0, 1)

	setClock(newTestClock(realNow))

	Expect(clock().Now()).To(Equal(realNow))
}

func (s *GomolSuite) TestRealClockNow(t *testing.T) {
	// This test is completely pointless because it's not something that can really
	// be tested but I was sick of seeing the red for a lack of a unit test.  So I created
	// this one and figure even on slow systems the two lines should be executed within
	// one second of each other. :P
	setClock(&realClock{})

	timeNow := time.Now()
	clockNow := clock().Now()

	diff := clockNow.Sub(timeNow)
	Expect(diff < time.Second).To(Equal(true))
}
