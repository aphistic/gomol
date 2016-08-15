package gomol

import . "gopkg.in/check.v1"

func (s *GomolSuite) TestNewConfig(c *C) {
	cfg := NewConfig()
	c.Check(cfg.FilenameAttr, Equals, "")
	c.Check(cfg.LineNumberAttr, Equals, "")
}
