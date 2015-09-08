package gomol

import (
	. "gopkg.in/check.v1"
)

func (s *GomolSuite) TestSetDefault(c *C) {
	SetDefault(nil)
	c.Check(curDefault, NotNil)

	b := newBase()
	b.AddLogger(NewMemLogger())
	SetDefault(b)
	c.Check(curDefault, NotNil)
	c.Check(curDefault, Equals, b)
	c.Check(curDefault.loggers, HasLen, 1)
}

func (s *GomolSuite) TestDefaultInitLogger(c *C) {
	defLogger := curDefault.loggers[0].(*MemLogger)
	c.Check(defLogger.IsInitialized, Equals, false)
	InitLoggers()
	c.Check(defLogger.IsInitialized, Equals, true)
}

func (s *GomolSuite) TestDefaultShutdownLogger(c *C) {
	defLogger := curDefault.loggers[0].(*MemLogger)
	c.Check(defLogger.IsShutdown, Equals, false)
	ShutdownLoggers()
	c.Check(defLogger.IsShutdown, Equals, true)
}

func (s *GomolSuite) TestDefaultDbg(c *C) {
	Dbg("test")
	defLogger := curDefault.loggers[0].(*MemLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0].Level, Equals, levelDbg)
	c.Check(defLogger.Messages[0].Message, Equals, "test")
	c.Check(defLogger.Messages[0].Attrs, HasLen, 0)
}

func (s *GomolSuite) TestDefaultDbgf(c *C) {
	Dbgf("test %v", 1234)
	defLogger := curDefault.loggers[0].(*MemLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0].Level, Equals, levelDbg)
	c.Check(defLogger.Messages[0].Message, Equals, "test 1234")
	c.Check(defLogger.Messages[0].Attrs, HasLen, 0)
}

func (s *GomolSuite) TestDefaultDbgm(c *C) {
	Dbgm(
		map[string]interface{}{
			"attr1": 4321,
		},
		"test %v", 1234)
	defLogger := curDefault.loggers[0].(*MemLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0].Level, Equals, levelDbg)
	c.Check(defLogger.Messages[0].Message, Equals, "test 1234")
	c.Assert(defLogger.Messages[0].Attrs, HasLen, 1)
	c.Check(defLogger.Messages[0].Attrs["attr1"], Equals, 4321)
}

func (s *GomolSuite) TestDefaultInfo(c *C) {
	Info("test")
	defLogger := curDefault.loggers[0].(*MemLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0].Level, Equals, levelInfo)
	c.Check(defLogger.Messages[0].Message, Equals, "test")
	c.Check(defLogger.Messages[0].Attrs, HasLen, 0)
}

func (s *GomolSuite) TestDefaultInfof(c *C) {
	Infof("test %v", 1234)
	defLogger := curDefault.loggers[0].(*MemLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0].Level, Equals, levelInfo)
	c.Check(defLogger.Messages[0].Message, Equals, "test 1234")
	c.Check(defLogger.Messages[0].Attrs, HasLen, 0)
}

func (s *GomolSuite) TestDefaultInfom(c *C) {
	Infom(
		map[string]interface{}{
			"attr1": 4321,
		},
		"test %v", 1234)
	defLogger := curDefault.loggers[0].(*MemLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0].Level, Equals, levelInfo)
	c.Check(defLogger.Messages[0].Message, Equals, "test 1234")
	c.Assert(defLogger.Messages[0].Attrs, HasLen, 1)
	c.Check(defLogger.Messages[0].Attrs["attr1"], Equals, 4321)
}

func (s *GomolSuite) TestDefaultWarn(c *C) {
	Warn("test")
	defLogger := curDefault.loggers[0].(*MemLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0].Level, Equals, levelWarn)
	c.Check(defLogger.Messages[0].Message, Equals, "test")
	c.Check(defLogger.Messages[0].Attrs, HasLen, 0)
}

func (s *GomolSuite) TestDefaultWarnf(c *C) {
	Warnf("test %v", 1234)
	defLogger := curDefault.loggers[0].(*MemLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0].Level, Equals, levelWarn)
	c.Check(defLogger.Messages[0].Message, Equals, "test 1234")
	c.Check(defLogger.Messages[0].Attrs, HasLen, 0)
}

func (s *GomolSuite) TestDefaultWarnm(c *C) {
	Warnm(
		map[string]interface{}{
			"attr1": 4321,
		},
		"test %v", 1234)
	defLogger := curDefault.loggers[0].(*MemLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0].Level, Equals, levelWarn)
	c.Check(defLogger.Messages[0].Message, Equals, "test 1234")
	c.Assert(defLogger.Messages[0].Attrs, HasLen, 1)
	c.Check(defLogger.Messages[0].Attrs["attr1"], Equals, 4321)
}

func (s *GomolSuite) TestDefaultErr(c *C) {
	Err("test")
	defLogger := curDefault.loggers[0].(*MemLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0].Level, Equals, levelError)
	c.Check(defLogger.Messages[0].Message, Equals, "test")
	c.Check(defLogger.Messages[0].Attrs, HasLen, 0)
}

func (s *GomolSuite) TestDefaultErrf(c *C) {
	Errf("test %v", 1234)
	defLogger := curDefault.loggers[0].(*MemLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0].Level, Equals, levelError)
	c.Check(defLogger.Messages[0].Message, Equals, "test 1234")
	c.Check(defLogger.Messages[0].Attrs, HasLen, 0)
}

func (s *GomolSuite) TestDefaultErrm(c *C) {
	Errm(
		map[string]interface{}{
			"attr1": 4321,
		},
		"test %v", 1234)
	defLogger := curDefault.loggers[0].(*MemLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0].Level, Equals, levelError)
	c.Check(defLogger.Messages[0].Message, Equals, "test 1234")
	c.Assert(defLogger.Messages[0].Attrs, HasLen, 1)
	c.Check(defLogger.Messages[0].Attrs["attr1"], Equals, 4321)
}

func (s *GomolSuite) TestDefaultFatal(c *C) {
	Fatal("test")
	defLogger := curDefault.loggers[0].(*MemLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0].Level, Equals, levelFatal)
	c.Check(defLogger.Messages[0].Message, Equals, "test")
	c.Check(defLogger.Messages[0].Attrs, HasLen, 0)
}

func (s *GomolSuite) TestDefaultFatalf(c *C) {
	Fatalf("test %v", 1234)
	defLogger := curDefault.loggers[0].(*MemLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0].Level, Equals, levelFatal)
	c.Check(defLogger.Messages[0].Message, Equals, "test 1234")
	c.Check(defLogger.Messages[0].Attrs, HasLen, 0)
}

func (s *GomolSuite) TestDefaultFatalm(c *C) {
	Fatalm(
		map[string]interface{}{
			"attr1": 4321,
		},
		"test %v", 1234)
	defLogger := curDefault.loggers[0].(*MemLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0].Level, Equals, levelFatal)
	c.Check(defLogger.Messages[0].Message, Equals, "test 1234")
	c.Assert(defLogger.Messages[0].Attrs, HasLen, 1)
	c.Check(defLogger.Messages[0].Attrs["attr1"], Equals, 4321)
}
