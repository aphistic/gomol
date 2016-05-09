package gomol

import (
	. "gopkg.in/check.v1"
)

func newDefaultMemLogger() *memLogger {
	cfg := newMemLoggerConfig()
	l, _ := newMemLogger(cfg)
	return l
}

func (s *GomolSuite) TestMemInitLogger(c *C) {
	ml := newDefaultMemLogger()
	c.Check(ml.IsInitialized(), Equals, false)
	ml.InitLogger()
	c.Check(ml.IsInitialized(), Equals, true)
}

func (s *GomolSuite) TestMemInitLoggerFail(c *C) {
	mlCfg := newMemLoggerConfig()
	mlCfg.FailInit = true
	ml, err := newMemLogger(mlCfg)
	c.Check(err, IsNil)
	c.Check(ml.IsInitialized(), Equals, false)
	err = ml.InitLogger()
	c.Check(ml.IsInitialized(), Equals, false)
	c.Check(err, NotNil)
	c.Check(err.Error(), Equals, "Init failed")
}

func (s *GomolSuite) TestMemShutdownLogger(c *C) {
	ml := newDefaultMemLogger()
	c.Check(ml.isShutdown, Equals, false)
	ml.ShutdownLogger()
	c.Check(ml.isShutdown, Equals, true)
}

func (s *GomolSuite) TestMemShutdownLoggerFail(c *C) {
	mlCfg := newMemLoggerConfig()
	mlCfg.FailShutdown = true
	ml, err := newMemLogger(mlCfg)
	c.Check(err, IsNil)
	c.Check(ml.isShutdown, Equals, false)
	err = ml.ShutdownLogger()
	c.Check(ml.isShutdown, Equals, false)
	c.Check(err, NotNil)
	c.Check(err.Error(), Equals, "Shutdown failed")
}

func (s *GomolSuite) TestMemClearMessages(c *C) {
	ml := newDefaultMemLogger()
	c.Check(ml.Messages, HasLen, 0)
	ml.Logm(LEVEL_DEBUG, nil, "test")
	c.Check(ml.Messages, HasLen, 1)
	ml.ClearMessages()
	c.Check(ml.Messages, HasLen, 0)
}

func (s *GomolSuite) TestMemLogmNoAttrs(c *C) {
	ml := newDefaultMemLogger()
	ml.Logm(LEVEL_DEBUG, nil, "test")
	c.Assert(ml.Messages, HasLen, 1)
	c.Check(ml.Messages[0].Level, Equals, LEVEL_DEBUG)
	c.Check(ml.Messages[0].Message, Equals, "test")
	c.Check(ml.Messages[0].Attrs, HasLen, 0)
}

func (s *GomolSuite) TestMemLogmAttrs(c *C) {
	ml := newDefaultMemLogger()
	ml.Logm(
		LEVEL_DEBUG,
		map[string]interface{}{
			"attr1": 4321,
		},
		"test 1234")
	c.Assert(ml.Messages, HasLen, 1)
	c.Check(ml.Messages[0].Level, Equals, LEVEL_DEBUG)
	c.Check(ml.Messages[0].Message, Equals, "test 1234")
	c.Assert(ml.Messages[0].Attrs, HasLen, 1)
	c.Check(ml.Messages[0].Attrs["attr1"], Equals, 4321)
}

func (s *GomolSuite) TestMemBaseAttrs(c *C) {
	b := newBase()
	b.SetAttr("attr1", 7890)
	b.SetAttr("attr2", "val2")

	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	ml.Logm(
		LEVEL_DEBUG,
		map[string]interface{}{
			"attr1": 4321,
			"attr3": "val3",
		},
		"test 1234")
	c.Assert(ml.Messages, HasLen, 1)
	c.Check(ml.Messages[0].Level, Equals, LEVEL_DEBUG)
	c.Check(ml.Messages[0].Message, Equals, "test 1234")
	c.Assert(ml.Messages[0].Attrs, HasLen, 3)
	c.Check(ml.Messages[0].Attrs["attr1"], Equals, 4321)
	c.Check(ml.Messages[0].Attrs["attr2"], Equals, "val2")
	c.Check(ml.Messages[0].Attrs["attr3"], Equals, "val3")
}
