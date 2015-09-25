package gomol

import (
	. "gopkg.in/check.v1"
)

func (s *GomolSuite) TestDefaultInitLogger(c *C) {
	curDefault = newBase()
	curDefault.AddLogger(NewMemLogger())
	defLogger := curDefault.loggers[0].(*MemLogger)
	c.Check(defLogger.IsInitialized(), Equals, false)
	InitLoggers()
	c.Check(defLogger.IsInitialized(), Equals, true)
	curDefault.ShutdownLoggers()
}

func (s *GomolSuite) TestDefaultShutdownLogger(c *C) {
	curDefault = newBase()
	curDefault.AddLogger(NewMemLogger())
	curDefault.InitLoggers()
	defLogger := curDefault.loggers[0].(*MemLogger)
	c.Check(defLogger.isShutdown, Equals, false)
	ShutdownLoggers()
	c.Check(defLogger.isShutdown, Equals, true)
}

func (s *GomolSuite) TestDefaultAddLogger(c *C) {
	curDefault = newBase()
	c.Check(curDefault.loggers, HasLen, 0)
	AddLogger(NewMemLogger())
	c.Check(curDefault.loggers, HasLen, 1)
}

func (s *GomolSuite) TestDefaultSetAttr(c *C) {
	curDefault = newBase()
	c.Check(curDefault.BaseAttrs, HasLen, 0)
	SetAttr("attr", 1234)
	c.Check(curDefault.BaseAttrs, HasLen, 1)
	c.Check(curDefault.BaseAttrs["attr"], Equals, 1234)
}

func (s *GomolSuite) TestDefaultRemoveAttr(c *C) {
	curDefault = newBase()
	c.Check(curDefault.BaseAttrs, HasLen, 0)
	SetAttr("attr", 1234)
	c.Check(curDefault.BaseAttrs, HasLen, 1)
	c.Check(curDefault.BaseAttrs["attr"], Equals, 1234)
	RemoveAttr("attr")
	c.Check(curDefault.BaseAttrs, HasLen, 0)
}

func (s *GomolSuite) TestDefaultClearAttrs(c *C) {
	curDefault = newBase()
	c.Check(curDefault.BaseAttrs, HasLen, 0)
	SetAttr("attr", 1234)
	c.Check(curDefault.BaseAttrs, HasLen, 1)
	c.Check(curDefault.BaseAttrs["attr"], Equals, 1234)
	SetAttr("attr2", 1234)
	c.Check(curDefault.BaseAttrs, HasLen, 2)
	c.Check(curDefault.BaseAttrs["attr2"], Equals, 1234)
	ClearAttrs()
	c.Check(curDefault.BaseAttrs, HasLen, 0)
}

func (s *GomolSuite) TestDefaultDbg(c *C) {
	Dbg("test")
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*MemLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0].Level, Equals, LEVEL_DEBUG)
	c.Check(defLogger.Messages[0].Message, Equals, "test")
	c.Check(defLogger.Messages[0].Attrs, HasLen, 0)
}

func (s *GomolSuite) TestDefaultDbgf(c *C) {
	Dbgf("test %v", 1234)
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*MemLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0].Level, Equals, LEVEL_DEBUG)
	c.Check(defLogger.Messages[0].Message, Equals, "test 1234")
	c.Check(defLogger.Messages[0].Attrs, HasLen, 0)
}

func (s *GomolSuite) TestDefaultDbgm(c *C) {
	Dbgm(
		map[string]interface{}{
			"attr1": 4321,
		},
		"test %v", 1234)
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*MemLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0].Level, Equals, LEVEL_DEBUG)
	c.Check(defLogger.Messages[0].Message, Equals, "test 1234")
	c.Assert(defLogger.Messages[0].Attrs, HasLen, 1)
	c.Check(defLogger.Messages[0].Attrs["attr1"], Equals, 4321)
}

func (s *GomolSuite) TestDefaultInfo(c *C) {
	Info("test")
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*MemLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0].Level, Equals, LEVEL_INFO)
	c.Check(defLogger.Messages[0].Message, Equals, "test")
	c.Check(defLogger.Messages[0].Attrs, HasLen, 0)
}

func (s *GomolSuite) TestDefaultInfof(c *C) {
	Infof("test %v", 1234)
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*MemLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0].Level, Equals, LEVEL_INFO)
	c.Check(defLogger.Messages[0].Message, Equals, "test 1234")
	c.Check(defLogger.Messages[0].Attrs, HasLen, 0)
}

func (s *GomolSuite) TestDefaultInfom(c *C) {
	Infom(
		map[string]interface{}{
			"attr1": 4321,
		},
		"test %v", 1234)
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*MemLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0].Level, Equals, LEVEL_INFO)
	c.Check(defLogger.Messages[0].Message, Equals, "test 1234")
	c.Assert(defLogger.Messages[0].Attrs, HasLen, 1)
	c.Check(defLogger.Messages[0].Attrs["attr1"], Equals, 4321)
}

func (s *GomolSuite) TestDefaultWarn(c *C) {
	Warn("test")
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*MemLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0].Level, Equals, LEVEL_WARNING)
	c.Check(defLogger.Messages[0].Message, Equals, "test")
	c.Check(defLogger.Messages[0].Attrs, HasLen, 0)
}

func (s *GomolSuite) TestDefaultWarnf(c *C) {
	Warnf("test %v", 1234)
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*MemLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0].Level, Equals, LEVEL_WARNING)
	c.Check(defLogger.Messages[0].Message, Equals, "test 1234")
	c.Check(defLogger.Messages[0].Attrs, HasLen, 0)
}

func (s *GomolSuite) TestDefaultWarnm(c *C) {
	Warnm(
		map[string]interface{}{
			"attr1": 4321,
		},
		"test %v", 1234)
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*MemLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0].Level, Equals, LEVEL_WARNING)
	c.Check(defLogger.Messages[0].Message, Equals, "test 1234")
	c.Assert(defLogger.Messages[0].Attrs, HasLen, 1)
	c.Check(defLogger.Messages[0].Attrs["attr1"], Equals, 4321)
}

func (s *GomolSuite) TestDefaultErr(c *C) {
	Err("test")
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*MemLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0].Level, Equals, LEVEL_ERROR)
	c.Check(defLogger.Messages[0].Message, Equals, "test")
	c.Check(defLogger.Messages[0].Attrs, HasLen, 0)
}

func (s *GomolSuite) TestDefaultErrf(c *C) {
	Errf("test %v", 1234)
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*MemLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0].Level, Equals, LEVEL_ERROR)
	c.Check(defLogger.Messages[0].Message, Equals, "test 1234")
	c.Check(defLogger.Messages[0].Attrs, HasLen, 0)
}

func (s *GomolSuite) TestDefaultErrm(c *C) {
	Errm(
		map[string]interface{}{
			"attr1": 4321,
		},
		"test %v", 1234)
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*MemLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0].Level, Equals, LEVEL_ERROR)
	c.Check(defLogger.Messages[0].Message, Equals, "test 1234")
	c.Assert(defLogger.Messages[0].Attrs, HasLen, 1)
	c.Check(defLogger.Messages[0].Attrs["attr1"], Equals, 4321)
}

func (s *GomolSuite) TestDefaultFatal(c *C) {
	Fatal("test")
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*MemLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0].Level, Equals, LEVEL_FATAL)
	c.Check(defLogger.Messages[0].Message, Equals, "test")
	c.Check(defLogger.Messages[0].Attrs, HasLen, 0)
}

func (s *GomolSuite) TestDefaultFatalf(c *C) {
	Fatalf("test %v", 1234)
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*MemLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0].Level, Equals, LEVEL_FATAL)
	c.Check(defLogger.Messages[0].Message, Equals, "test 1234")
	c.Check(defLogger.Messages[0].Attrs, HasLen, 0)
}

func (s *GomolSuite) TestDefaultFatalm(c *C) {
	Fatalm(
		map[string]interface{}{
			"attr1": 4321,
		},
		"test %v", 1234)
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*MemLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0].Level, Equals, LEVEL_FATAL)
	c.Check(defLogger.Messages[0].Message, Equals, "test 1234")
	c.Assert(defLogger.Messages[0].Attrs, HasLen, 1)
	c.Check(defLogger.Messages[0].Attrs["attr1"], Equals, 4321)
}
