package gomol

import (
	. "gopkg.in/check.v1"
)

func (s *GomolSuite) TestDefaultInitLogger(c *C) {
	curDefault = NewBase()
	c.Check(IsInitialized(), Equals, false)
	AddLogger(newDefaultMemLogger())
	defLogger := curDefault.loggers[0].(*memLogger)
	c.Check(defLogger.IsInitialized(), Equals, false)
	InitLoggers()
	c.Check(IsInitialized(), Equals, true)
	c.Check(defLogger.IsInitialized(), Equals, true)
	ShutdownLoggers()
}

func (s *GomolSuite) TestDefaultShutdownLogger(c *C) {
	curDefault = NewBase()
	c.Check(IsInitialized(), Equals, false)
	AddLogger(newDefaultMemLogger())
	InitLoggers()
	c.Check(IsInitialized(), Equals, true)
	defLogger := curDefault.loggers[0].(*memLogger)
	c.Check(defLogger.isShutdown, Equals, false)
	ShutdownLoggers()
	c.Check(defLogger.isShutdown, Equals, true)
	c.Check(IsInitialized(), Equals, false)
}

func (s *GomolSuite) TestDefaultAddLogger(c *C) {
	curDefault = NewBase()
	c.Check(curDefault.loggers, HasLen, 0)
	AddLogger(newDefaultMemLogger())
	c.Check(curDefault.loggers, HasLen, 1)
}

func (s *GomolSuite) TestDefaultRemoveLogger(c *C) {
	curDefault = NewBase()

	ml1 := newDefaultMemLogger()
	ml2 := newDefaultMemLogger()
	ml3 := newDefaultMemLogger()
	AddLogger(ml1)
	AddLogger(ml2)
	AddLogger(ml3)

	InitLoggers()

	c.Check(ml1.IsInitialized(), Equals, true)
	c.Check(ml2.IsInitialized(), Equals, true)
	c.Check(ml3.IsInitialized(), Equals, true)
	c.Check(curDefault.loggers, HasLen, 3)

	err := RemoveLogger(ml2)
	c.Assert(err, IsNil)
	c.Check(ml1.IsInitialized(), Equals, true)
	c.Check(ml2.IsInitialized(), Equals, false)
	c.Check(ml3.IsInitialized(), Equals, true)
	c.Check(curDefault.loggers, HasLen, 2)
}

func (s *GomolSuite) TestDefaultClearLoggers(c *C) {
	curDefault = NewBase()

	ml1 := newDefaultMemLogger()
	ml2 := newDefaultMemLogger()
	ml3 := newDefaultMemLogger()
	AddLogger(ml1)
	AddLogger(ml2)
	AddLogger(ml3)

	InitLoggers()

	c.Check(ml1.IsInitialized(), Equals, true)
	c.Check(ml2.IsInitialized(), Equals, true)
	c.Check(ml3.IsInitialized(), Equals, true)
	c.Check(curDefault.loggers, HasLen, 3)

	err := ClearLoggers()
	c.Assert(err, IsNil)
	c.Check(ml1.IsInitialized(), Equals, false)
	c.Check(ml2.IsInitialized(), Equals, false)
	c.Check(ml3.IsInitialized(), Equals, false)
	c.Check(curDefault.loggers, HasLen, 0)
}

func (s *GomolSuite) TestDefaultSetLogLevel(c *C) {
	curDefault = NewBase()
	InitLoggers()
	ml := newDefaultMemLogger()
	AddLogger(ml)

	SetLogLevel(LEVEL_WARNING)
	Dbg("test")
	Info("test")
	Warn("test")
	Err("test")
	Fatal("test")
	ShutdownLoggers()
	c.Check(ml.Messages, HasLen, 3)
}

func (s *GomolSuite) TestDefaultSetAttr(c *C) {
	curDefault = NewBase()
	c.Check(curDefault.BaseAttrs, HasLen, 0)
	SetAttr("attr", 1234)
	c.Check(curDefault.BaseAttrs, HasLen, 1)
	c.Check(curDefault.BaseAttrs["attr"], Equals, 1234)
}

func (s *GomolSuite) TestDefaultRemoveAttr(c *C) {
	curDefault = NewBase()
	c.Check(curDefault.BaseAttrs, HasLen, 0)
	SetAttr("attr", 1234)
	c.Check(curDefault.BaseAttrs, HasLen, 1)
	c.Check(curDefault.BaseAttrs["attr"], Equals, 1234)
	RemoveAttr("attr")
	c.Check(curDefault.BaseAttrs, HasLen, 0)
}

func (s *GomolSuite) TestDefaultClearAttrs(c *C) {
	curDefault = NewBase()
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

func (s *GomolSuite) TestDefaultNewLogAdapter(c *C) {
	la := NewLogAdapter(map[string]interface{}{"foo": "bar"})
	defLogger := curDefault.loggers[0].(*memLogger)

	la.Dbgm(map[string]interface{}{"attr": "val"}, "test")

	ShutdownLoggers()

	c.Assert(len(defLogger.Messages), Equals, 1)
	c.Check(defLogger.Messages[0], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LEVEL_DEBUG,
		Message:   "test",
		Attrs: map[string]interface{}{
			"foo":  "bar",
			"attr": "val",
		},
	})
}

func (s *GomolSuite) TestDefaultDbg(c *C) {
	Dbg("test")
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*memLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LEVEL_DEBUG,
		Message:   "test",
		Attrs:     map[string]interface{}{},
	})
}

func (s *GomolSuite) TestDefaultDbgf(c *C) {
	Dbgf("test %v", 1234)
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*memLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LEVEL_DEBUG,
		Message:   "test 1234",
		Attrs:     map[string]interface{}{},
	})
}

func (s *GomolSuite) TestDefaultDbgm(c *C) {
	Dbgm(
		map[string]interface{}{
			"attr1": 4321,
		},
		"test %v", 1234)
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*memLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LEVEL_DEBUG,
		Message:   "test 1234",
		Attrs: map[string]interface{}{
			"attr1": 4321,
		},
	})
}

func (s *GomolSuite) TestDefaultInfo(c *C) {
	Info("test")
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*memLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LEVEL_INFO,
		Message:   "test",
		Attrs:     map[string]interface{}{},
	})
}

func (s *GomolSuite) TestDefaultInfof(c *C) {
	Infof("test %v", 1234)
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*memLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LEVEL_INFO,
		Message:   "test 1234",
		Attrs:     map[string]interface{}{},
	})
}

func (s *GomolSuite) TestDefaultInfom(c *C) {
	Infom(
		map[string]interface{}{
			"attr1": 4321,
		},
		"test %v", 1234)
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*memLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LEVEL_INFO,
		Message:   "test 1234",
		Attrs: map[string]interface{}{
			"attr1": 4321,
		},
	})
}

func (s *GomolSuite) TestDefaultWarn(c *C) {
	Warn("test")
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*memLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LEVEL_WARNING,
		Message:   "test",
		Attrs:     map[string]interface{}{},
	})
}

func (s *GomolSuite) TestDefaultWarnf(c *C) {
	Warnf("test %v", 1234)
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*memLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LEVEL_WARNING,
		Message:   "test 1234",
		Attrs:     map[string]interface{}{},
	})
}

func (s *GomolSuite) TestDefaultWarnm(c *C) {
	Warnm(
		map[string]interface{}{
			"attr1": 4321,
		},
		"test %v", 1234)
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*memLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LEVEL_WARNING,
		Message:   "test 1234",
		Attrs: map[string]interface{}{
			"attr1": 4321,
		},
	})
}

func (s *GomolSuite) TestDefaultErr(c *C) {
	Err("test")
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*memLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LEVEL_ERROR,
		Message:   "test",
		Attrs:     map[string]interface{}{},
	})
}

func (s *GomolSuite) TestDefaultErrf(c *C) {
	Errf("test %v", 1234)
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*memLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LEVEL_ERROR,
		Message:   "test 1234",
		Attrs:     map[string]interface{}{},
	})
}

func (s *GomolSuite) TestDefaultErrm(c *C) {
	Errm(
		map[string]interface{}{
			"attr1": 4321,
		},
		"test %v", 1234)
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*memLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LEVEL_ERROR,
		Message:   "test 1234",
		Attrs: map[string]interface{}{
			"attr1": 4321,
		},
	})
}

func (s *GomolSuite) TestDefaultFatal(c *C) {
	Fatal("test")
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*memLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LEVEL_FATAL,
		Message:   "test",
		Attrs:     map[string]interface{}{},
	})
}

func (s *GomolSuite) TestDefaultFatalf(c *C) {
	Fatalf("test %v", 1234)
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*memLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LEVEL_FATAL,
		Message:   "test 1234",
		Attrs:     map[string]interface{}{},
	})
}

func (s *GomolSuite) TestDefaultFatalm(c *C) {
	Fatalm(
		map[string]interface{}{
			"attr1": 4321,
		},
		"test %v", 1234)
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*memLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LEVEL_FATAL,
		Message:   "test 1234",
		Attrs: map[string]interface{}{
			"attr1": 4321,
		},
	})
}

func (s *GomolSuite) TestDefaultDie(c *C) {
	Die(1234, "test")
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*memLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LEVEL_FATAL,
		Message:   "test",
		Attrs:     map[string]interface{}{},
	})

	c.Check(curDefault.isInitialized, Equals, false)
	c.Check(curTestExiter.exited, Equals, true)
	c.Check(curTestExiter.code, Equals, 1234)
}

func (s *GomolSuite) TestDefaultDief(c *C) {
	Dief(1234, "test %v", 1234)
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*memLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LEVEL_FATAL,
		Message:   "test 1234",
		Attrs:     map[string]interface{}{},
	})
	c.Check(defLogger.Messages[0].Level, Equals, LEVEL_FATAL)
	c.Check(defLogger.Messages[0].Message, Equals, "test 1234")
	c.Check(defLogger.Messages[0].Attrs, HasLen, 0)

	c.Check(curDefault.isInitialized, Equals, false)
	c.Check(curTestExiter.exited, Equals, true)
	c.Check(curTestExiter.code, Equals, 1234)
}

func (s *GomolSuite) TestDefaultDiem(c *C) {
	Diem(
		1234,
		map[string]interface{}{
			"attr1": 4321,
		},
		"test %v", 1234)
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*memLogger)
	c.Assert(defLogger.Messages, HasLen, 1)
	c.Check(defLogger.Messages[0], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LEVEL_FATAL,
		Message:   "test 1234",
		Attrs: map[string]interface{}{
			"attr1": 4321,
		},
	})

	c.Check(curDefault.isInitialized, Equals, false)
	c.Check(curTestExiter.exited, Equals, true)
	c.Check(curTestExiter.code, Equals, 1234)
}
