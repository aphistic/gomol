package gomol

import (
	. "gopkg.in/check.v1"
)

func (s *GomolSuite) TestGelfNew(c *C) {
	l := NewGelfLogger("host", 1234)
	c.Check(l, NotNil)
	c.Check(l.Hostname, Equals, "host")
	c.Check(l.Port, Equals, 1234)
}
