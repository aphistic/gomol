package gomol

import (
	. "gopkg.in/check.v1"
)

func (s *GomolSuite) TestMemInitLogger(c *C) {
	ml := NewMemLogger()
	c.Check(ml.IsInitialized(), Equals, false)
	ml.InitLogger()
	c.Check(ml.IsInitialized(), Equals, true)
}

func (s *GomolSuite) TestMemInitLoggerFail(c *C) {
	ml := NewMemLoggerWithConfig(MemLoggerConfig{FailInit: true})
	c.Check(ml.IsInitialized(), Equals, false)
	err := ml.InitLogger()
	c.Check(ml.IsInitialized(), Equals, false)
	c.Check(err, NotNil)
	c.Check(err.Error(), Equals, "Init failed")
}

func (s *GomolSuite) TestMemShutdownLogger(c *C) {
	ml := NewMemLogger()
	c.Check(ml.isShutdown, Equals, false)
	ml.ShutdownLogger()
	c.Check(ml.isShutdown, Equals, true)
}

func (s *GomolSuite) TestMemShutdownLoggerFail(c *C) {
	ml := NewMemLoggerWithConfig(MemLoggerConfig{FailShutdown: true})
	c.Check(ml.isShutdown, Equals, false)
	err := ml.ShutdownLogger()
	c.Check(ml.isShutdown, Equals, false)
	c.Check(err, NotNil)
	c.Check(err.Error(), Equals, "Shutdown failed")
}

func (s *GomolSuite) TestMemClearMessages(c *C) {
	ml := NewMemLogger()
	c.Check(ml.Messages, HasLen, 0)
	ml.Dbg("test")
	c.Check(ml.Messages, HasLen, 1)
	ml.ClearMessages()
	c.Check(ml.Messages, HasLen, 0)
}

func (s *GomolSuite) TestMemDbg(c *C) {
	ml := NewMemLogger()
	ml.Dbg("test")
	c.Assert(ml.Messages, HasLen, 1)
	c.Check(ml.Messages[0].Level, Equals, LEVEL_DEBUG)
	c.Check(ml.Messages[0].Message, Equals, "test")
	c.Check(ml.Messages[0].Attrs, HasLen, 0)
}

func (s *GomolSuite) TestMemDbgf(c *C) {
	ml := NewMemLogger()
	ml.Dbgf("test %v", 1234)
	c.Assert(ml.Messages, HasLen, 1)
	c.Check(ml.Messages[0].Level, Equals, LEVEL_DEBUG)
	c.Check(ml.Messages[0].Message, Equals, "test 1234")
	c.Check(ml.Messages[0].Attrs, HasLen, 0)
}

func (s *GomolSuite) TestMemDbgm(c *C) {
	ml := NewMemLogger()
	ml.Dbgm(
		map[string]interface{}{
			"attr1": 4321,
		},
		"test %v", 1234)
	c.Assert(ml.Messages, HasLen, 1)
	c.Check(ml.Messages[0].Level, Equals, LEVEL_DEBUG)
	c.Check(ml.Messages[0].Message, Equals, "test 1234")
	c.Assert(ml.Messages[0].Attrs, HasLen, 1)
	c.Check(ml.Messages[0].Attrs["attr1"], Equals, 4321)
}

func (s *GomolSuite) TestMemInfo(c *C) {
	ml := NewMemLogger()
	ml.Info("test")
	c.Assert(ml.Messages, HasLen, 1)
	c.Check(ml.Messages[0].Level, Equals, LEVEL_INFO)
	c.Check(ml.Messages[0].Message, Equals, "test")
	c.Check(ml.Messages[0].Attrs, HasLen, 0)
}

func (s *GomolSuite) TestMemInfof(c *C) {
	ml := NewMemLogger()
	ml.Infof("test %v", 1234)
	c.Assert(ml.Messages, HasLen, 1)
	c.Check(ml.Messages[0].Level, Equals, LEVEL_INFO)
	c.Check(ml.Messages[0].Message, Equals, "test 1234")
	c.Check(ml.Messages[0].Attrs, HasLen, 0)
}

func (s *GomolSuite) TestMemInfom(c *C) {
	ml := NewMemLogger()
	ml.Infom(
		map[string]interface{}{
			"attr1": 4321,
		},
		"test %v", 1234)
	c.Assert(ml.Messages, HasLen, 1)
	c.Check(ml.Messages[0].Level, Equals, LEVEL_INFO)
	c.Check(ml.Messages[0].Message, Equals, "test 1234")
	c.Assert(ml.Messages[0].Attrs, HasLen, 1)
	c.Check(ml.Messages[0].Attrs["attr1"], Equals, 4321)
}

func (s *GomolSuite) TestMemWarn(c *C) {
	ml := NewMemLogger()
	ml.Warn("test")
	c.Assert(ml.Messages, HasLen, 1)
	c.Check(ml.Messages[0].Level, Equals, LEVEL_WARNING)
	c.Check(ml.Messages[0].Message, Equals, "test")
	c.Check(ml.Messages[0].Attrs, HasLen, 0)
}

func (s *GomolSuite) TestMemWarnf(c *C) {
	ml := NewMemLogger()
	ml.Warnf("test %v", 1234)
	c.Assert(ml.Messages, HasLen, 1)
	c.Check(ml.Messages[0].Level, Equals, LEVEL_WARNING)
	c.Check(ml.Messages[0].Message, Equals, "test 1234")
	c.Check(ml.Messages[0].Attrs, HasLen, 0)
}

func (s *GomolSuite) TestMemWarnm(c *C) {
	ml := NewMemLogger()
	ml.Warnm(
		map[string]interface{}{
			"attr1": 4321,
		},
		"test %v", 1234)
	c.Assert(ml.Messages, HasLen, 1)
	c.Check(ml.Messages[0].Level, Equals, LEVEL_WARNING)
	c.Check(ml.Messages[0].Message, Equals, "test 1234")
	c.Assert(ml.Messages[0].Attrs, HasLen, 1)
	c.Check(ml.Messages[0].Attrs["attr1"], Equals, 4321)
}

func (s *GomolSuite) TestMemErr(c *C) {
	ml := NewMemLogger()
	ml.Err("test")
	c.Assert(ml.Messages, HasLen, 1)
	c.Check(ml.Messages[0].Level, Equals, LEVEL_ERROR)
	c.Check(ml.Messages[0].Message, Equals, "test")
	c.Check(ml.Messages[0].Attrs, HasLen, 0)
}

func (s *GomolSuite) TestMemErrf(c *C) {
	ml := NewMemLogger()
	ml.Errf("test %v", 1234)
	c.Assert(ml.Messages, HasLen, 1)
	c.Check(ml.Messages[0].Level, Equals, LEVEL_ERROR)
	c.Check(ml.Messages[0].Message, Equals, "test 1234")
	c.Check(ml.Messages[0].Attrs, HasLen, 0)
}

func (s *GomolSuite) TestMemErrm(c *C) {
	ml := NewMemLogger()
	ml.Errm(
		map[string]interface{}{
			"attr1": 4321,
		},
		"test %v", 1234)
	c.Assert(ml.Messages, HasLen, 1)
	c.Check(ml.Messages[0].Level, Equals, LEVEL_ERROR)
	c.Check(ml.Messages[0].Message, Equals, "test 1234")
	c.Assert(ml.Messages[0].Attrs, HasLen, 1)
	c.Check(ml.Messages[0].Attrs["attr1"], Equals, 4321)
}

func (s *GomolSuite) TestMemFatal(c *C) {
	ml := NewMemLogger()
	ml.Fatal("test")
	c.Assert(ml.Messages, HasLen, 1)
	c.Check(ml.Messages[0].Level, Equals, LEVEL_FATAL)
	c.Check(ml.Messages[0].Message, Equals, "test")
	c.Check(ml.Messages[0].Attrs, HasLen, 0)
}

func (s *GomolSuite) TestMemFatalf(c *C) {
	ml := NewMemLogger()
	ml.Fatalf("test %v", 1234)
	c.Assert(ml.Messages, HasLen, 1)
	c.Check(ml.Messages[0].Level, Equals, LEVEL_FATAL)
	c.Check(ml.Messages[0].Message, Equals, "test 1234")
	c.Check(ml.Messages[0].Attrs, HasLen, 0)
}

func (s *GomolSuite) TestMemFatalm(c *C) {
	ml := NewMemLogger()
	ml.Fatalm(
		map[string]interface{}{
			"attr1": 4321,
		},
		"test %v", 1234)
	c.Assert(ml.Messages, HasLen, 1)
	c.Check(ml.Messages[0].Level, Equals, LEVEL_FATAL)
	c.Check(ml.Messages[0].Message, Equals, "test 1234")
	c.Assert(ml.Messages[0].Attrs, HasLen, 1)
	c.Check(ml.Messages[0].Attrs["attr1"], Equals, 4321)
}

func (s *GomolSuite) TestMemBaseAttrs(c *C) {
	b := newBase()
	b.SetAttr("attr1", 7890)
	b.SetAttr("attr2", "val2")

	ml := NewMemLogger()
	b.AddLogger(ml)
	ml.Dbgm(
		map[string]interface{}{
			"attr1": 4321,
			"attr3": "val3",
		},
		"test %v", 1234)
	c.Assert(ml.Messages, HasLen, 1)
	c.Check(ml.Messages[0].Level, Equals, LEVEL_DEBUG)
	c.Check(ml.Messages[0].Message, Equals, "test 1234")
	c.Assert(ml.Messages[0].Attrs, HasLen, 3)
	c.Check(ml.Messages[0].Attrs["attr1"], Equals, 4321)
	c.Check(ml.Messages[0].Attrs["attr2"], Equals, "val2")
	c.Check(ml.Messages[0].Attrs["attr3"], Equals, "val3")
}
