package gomol

import (
	. "gopkg.in/check.v1"
)

func (s *GomolSuite) TestGelfNew(c *C) {
	cfg := NewGelfLoggerConfig()
	cfg.Hostname = "host"
	cfg.Port = 1234

	l := NewGelfLogger(cfg)
	c.Check(l, NotNil)
	c.Check(l.config.Hostname, Equals, "host")
	c.Check(l.config.Port, Equals, 1234)
}
