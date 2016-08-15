package gomol

import (
	"testing"
	"time"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type GomolSuite struct{}

var _ = Suite(&GomolSuite{})

var testBase *Base

type testExiter struct {
	exited bool
	code   int
}

func (exiter *testExiter) Exit(code int) {
	exiter.code = code
	exiter.exited = true
}

var curTestExiter *testExiter

func (s *GomolSuite) SetUpTest(c *C) {
	setClock(newTestClock(time.Now()))

	curTestExiter = &testExiter{}
	setExiter(curTestExiter)

	testBase = NewBase()
	testBase.AddLogger(newDefaultMemLogger())
	testBase.InitLoggers()

	curDefault = NewBase()
	curDefault.AddLogger(newDefaultMemLogger())
	curDefault.InitLoggers()
}

func (s *GomolSuite) TearDownTest(c *C) {
	curDefault.ShutdownLoggers()

	testBase.ShutdownLoggers()
}

func (s *GomolSuite) TestShouldLog(c *C) {
	b := NewBase()
	b.SetLogLevel(LevelInfo)
	c.Check(b.shouldLog(LevelUnknown), Equals, false)
	c.Check(b.shouldLog(LevelDebug), Equals, false)
	c.Check(b.shouldLog(LevelInfo), Equals, true)
	c.Check(b.shouldLog(LevelWarning), Equals, true)
	c.Check(b.shouldLog(LevelError), Equals, true)
	c.Check(b.shouldLog(LevelFatal), Equals, true)

	b.SetLogLevel(LevelFatal)
	c.Check(b.shouldLog(LevelUnknown), Equals, false)
	c.Check(b.shouldLog(LevelDebug), Equals, false)
	c.Check(b.shouldLog(LevelInfo), Equals, false)
	c.Check(b.shouldLog(LevelWarning), Equals, false)
	c.Check(b.shouldLog(LevelError), Equals, false)
	c.Check(b.shouldLog(LevelFatal), Equals, true)

	b.SetLogLevel(LevelNone)
	c.Check(b.shouldLog(LevelUnknown), Equals, false)
	c.Check(b.shouldLog(LevelDebug), Equals, false)
	c.Check(b.shouldLog(LevelInfo), Equals, false)
	c.Check(b.shouldLog(LevelWarning), Equals, false)
	c.Check(b.shouldLog(LevelError), Equals, false)
	c.Check(b.shouldLog(LevelFatal), Equals, false)
}

func (s *GomolSuite) TestNewBase(c *C) {
	b := NewBase()
	c.Check(b.isInitialized, Equals, false)
	c.Assert(b.config, NotNil)
	c.Check(b.config.FilenameAttr, Equals, "")
	c.Check(b.config.LineNumberAttr, Equals, "")
	c.Assert(b.queue, NotNil)
	c.Check(b.logLevel, Equals, LevelDebug)
	c.Check(b.loggers, HasLen, 0)
	c.Check(b.BaseAttrs, HasLen, 0)
}

func (s *GomolSuite) TestSetConfig(c *C) {
	b := NewBase()

	c.Assert(b.config, NotNil)
	c.Check(b.config.FilenameAttr, Equals, "")
	c.Check(b.config.LineNumberAttr, Equals, "")

	cfg := NewConfig()
	cfg.FilenameAttr = "filename"
	cfg.LineNumberAttr = "line_number"

	b.SetConfig(cfg)
	c.Assert(b.config, NotNil)
	c.Check(b.config.FilenameAttr, Equals, "filename")
	c.Check(b.config.LineNumberAttr, Equals, "line_number")
}

func (s *GomolSuite) TestSetLogLevel(c *C) {
	b := NewBase()
	b.InitLoggers()
	ml := newDefaultMemLogger()
	b.AddLogger(ml)

	b.SetLogLevel(LevelWarning)
	b.Dbg("test")
	b.Info("test")
	b.Warn("test")
	b.Err("test")
	b.Fatal("test")
	b.ShutdownLoggers()
	c.Check(ml.Messages, HasLen, 3)
}

func (s *GomolSuite) TestAddLogger(c *C) {
	b := NewBase()
	b.InitLoggers()
	c.Check(b.loggers, HasLen, 0)

	ml := newDefaultMemLogger()
	c.Check(ml.IsInitialized(), Equals, false)
	c.Check(ml.base, IsNil)

	b.AddLogger(ml)
	c.Check(b.IsInitialized(), Equals, true)
	c.Assert(b.loggers, HasLen, 1)
	c.Check(b.loggers[0].IsInitialized(), Equals, true)
	c.Check(ml.base, Equals, b)
}

func (s *GomolSuite) TestAddLoggerAfterInit(c *C) {
	b := NewBase()
	b.InitLoggers()

	ml := newDefaultMemLogger()
	c.Check(ml.IsInitialized(), Equals, false)

	ret := b.AddLogger(ml)
	c.Check(ret, IsNil)
	c.Check(ml.IsInitialized(), Equals, true)
}

func (s *GomolSuite) TestAddLoggerAfterShutdown(c *C) {
	b := NewBase()

	ml := newDefaultMemLogger()
	c.Check(ml.IsInitialized(), Equals, false)
	ml.InitLogger()
	c.Check(ml.IsInitialized(), Equals, true)

	ret := b.AddLogger(ml)
	c.Check(ret, IsNil)
	c.Check(ml.IsInitialized(), Equals, false)
}

func (s *GomolSuite) TestAddLoggerAfterInitFail(c *C) {
	b := NewBase()
	b.InitLoggers()

	mlCfg := newMemLoggerConfig()
	mlCfg.FailInit = true
	ml, err := newMemLogger(mlCfg)
	c.Check(err, IsNil)
	c.Check(ml.IsInitialized(), Equals, false)

	ret := b.AddLogger(ml)
	c.Check(ret, NotNil)
	c.Check(ret.Error(), Equals, "Init failed")
	c.Check(ml.IsInitialized(), Equals, false)
	c.Check(b.loggers, HasLen, 0)
}

func (s *GomolSuite) TestAddLoggerAfterShutdownFail(c *C) {
	b := NewBase()

	mlCfg := newMemLoggerConfig()
	mlCfg.FailShutdown = true
	ml, err := newMemLogger(mlCfg)
	c.Check(err, IsNil)
	c.Check(ml.IsInitialized(), Equals, false)
	ml.InitLogger()
	c.Check(ml.IsInitialized(), Equals, true)

	ret := b.AddLogger(ml)
	c.Check(ret, NotNil)
	c.Check(ret.Error(), Equals, "Shutdown failed")
	c.Check(ml.IsInitialized(), Equals, true)
	c.Check(b.loggers, HasLen, 0)
}

func (s *GomolSuite) TestBaseRemoveLogger(c *C) {
	b := NewBase()

	ml1 := newDefaultMemLogger()
	ml2 := newDefaultMemLogger()
	ml3 := newDefaultMemLogger()
	b.AddLogger(ml1)
	b.AddLogger(ml2)
	b.AddLogger(ml3)

	b.InitLoggers()

	c.Check(ml1.IsInitialized(), Equals, true)
	c.Check(ml2.IsInitialized(), Equals, true)
	c.Check(ml3.IsInitialized(), Equals, true)
	c.Check(b.loggers, HasLen, 3)

	err := b.RemoveLogger(ml2)
	c.Assert(err, IsNil)
	c.Check(ml1.IsInitialized(), Equals, true)
	c.Check(ml2.IsInitialized(), Equals, false)
	c.Check(ml3.IsInitialized(), Equals, true)
	c.Check(b.loggers, HasLen, 2)
}

func (s *GomolSuite) TestBaseRemoveLoggerNonExistent(c *C) {
	b := NewBase()

	ml1 := newDefaultMemLogger()
	ml2 := newDefaultMemLogger()
	b.AddLogger(ml1)

	b.InitLoggers()

	c.Check(ml1.IsInitialized(), Equals, true)
	c.Check(b.loggers, HasLen, 1)

	err := b.RemoveLogger(ml2)
	c.Assert(err, IsNil)
}

func (s *GomolSuite) TestBaseClearLoggers(c *C) {
	b := NewBase()

	ml1 := newDefaultMemLogger()
	ml2 := newDefaultMemLogger()
	ml3 := newDefaultMemLogger()
	b.AddLogger(ml1)
	b.AddLogger(ml2)
	b.AddLogger(ml3)

	b.InitLoggers()

	c.Check(ml1.IsInitialized(), Equals, true)
	c.Check(ml2.IsInitialized(), Equals, true)
	c.Check(ml3.IsInitialized(), Equals, true)
	c.Check(b.loggers, HasLen, 3)

	err := b.ClearLoggers()
	c.Assert(err, IsNil)
	c.Check(ml1.IsInitialized(), Equals, false)
	c.Check(ml2.IsInitialized(), Equals, false)
	c.Check(ml3.IsInitialized(), Equals, false)
	c.Check(b.loggers, HasLen, 0)
}

func (s *GomolSuite) TestInitLoggers(c *C) {
	b := NewBase()
	c.Check(b.IsInitialized(), Equals, false)

	ml1 := newDefaultMemLogger()
	ml2 := newDefaultMemLogger()

	b.AddLogger(ml1)
	b.AddLogger(ml2)

	b.InitLoggers()

	c.Check(b.IsInitialized(), Equals, true)
	c.Check(ml1.IsInitialized(), Equals, true)
	c.Check(ml2.IsInitialized(), Equals, true)
}

func (s *GomolSuite) TestInitLoggersTwice(c *C) {
	b := NewBase()
	c.Check(b.IsInitialized(), Equals, false)

	ml1 := newDefaultMemLogger()
	ml2 := newDefaultMemLogger()

	b.AddLogger(ml1)
	b.AddLogger(ml2)

	b.InitLoggers()
	b.InitLoggers()

	c.Check(b.IsInitialized(), Equals, true)
	c.Check(ml1.IsInitialized(), Equals, true)
	c.Check(ml2.IsInitialized(), Equals, true)
}

func (s *GomolSuite) TestInitLoggersFail(c *C) {
	b := NewBase()

	mlCfg := newMemLoggerConfig()
	mlCfg.FailInit = true
	ml1, err := newMemLogger(mlCfg)
	c.Check(err, IsNil)
	ml2, err := newMemLogger(mlCfg)
	c.Check(err, IsNil)

	b.AddLogger(ml1)
	b.AddLogger(ml2)

	err = b.InitLoggers()
	c.Check(err, NotNil)
	c.Check(err.Error(), Equals, "Init failed")

	c.Check(b.IsInitialized(), Equals, false)
	c.Check(ml1.IsInitialized(), Equals, false)
	c.Check(ml2.IsInitialized(), Equals, false)
}

func (s *GomolSuite) TestShutdownLoggers(c *C) {
	b := NewBase()

	ml1 := newDefaultMemLogger()
	ml2 := newDefaultMemLogger()

	b.AddLogger(ml1)
	b.AddLogger(ml2)

	b.InitLoggers()
	b.ShutdownLoggers()

	c.Check(ml1.isShutdown, Equals, true)
	c.Check(ml2.isShutdown, Equals, true)
}

func (s *GomolSuite) TestShutdownLoggersFail(c *C) {
	b := NewBase()

	mlCfg := newMemLoggerConfig()
	mlCfg.FailShutdown = true
	ml1, err := newMemLogger(mlCfg)
	c.Check(err, IsNil)
	ml2, err := newMemLogger(mlCfg)
	c.Check(err, IsNil)

	b.AddLogger(ml1)
	b.AddLogger(ml2)

	b.InitLoggers()
	err = b.ShutdownLoggers()
	c.Check(err, NotNil)
	c.Check(err.Error(), Equals, "Shutdown failed")

	c.Check(ml1.isShutdown, Equals, false)
	c.Check(ml2.isShutdown, Equals, false)
}

func (s *GomolSuite) TestShutdownLoggersTwice(c *C) {
	b := NewBase()

	ml1 := newDefaultMemLogger()
	ml2 := newDefaultMemLogger()

	b.AddLogger(ml1)
	b.AddLogger(ml2)

	b.InitLoggers()
	b.ShutdownLoggers()
	b.ShutdownLoggers()

	c.Check(ml1.isShutdown, Equals, true)
	c.Check(ml2.isShutdown, Equals, true)
}

func (s *GomolSuite) TestSetAttr(c *C) {
	b := NewBase()

	b.SetAttr("attr1", 1)
	c.Check(b.BaseAttrs, HasLen, 1)
	c.Check(b.BaseAttrs["attr1"], Equals, 1)
	b.SetAttr("attr2", "val2")
	c.Check(b.BaseAttrs, HasLen, 2)
	c.Check(b.BaseAttrs["attr2"], Equals, "val2")
}

func (s *GomolSuite) TestGetAttr(c *C) {
	b := NewBase()

	b.SetAttr("attr1", 1)
	b.SetAttr("attr2", "val2")

	c.Check(b.GetAttr("attr2"), Equals, "val2")
	c.Check(b.GetAttr("notakey"), IsNil)
}

func (s *GomolSuite) TestRemoveAttr(c *C) {
	b := NewBase()

	b.SetAttr("attr1", 1)
	c.Check(b.BaseAttrs, HasLen, 1)
	c.Check(b.BaseAttrs["attr1"], Equals, 1)

	b.RemoveAttr("attr1")
	c.Check(b.BaseAttrs, HasLen, 0)
}

func (s *GomolSuite) TestClearAttrs(c *C) {
	b := NewBase()

	b.SetAttr("attr1", 1)
	b.SetAttr("attr2", "val2")
	c.Check(b.BaseAttrs, HasLen, 2)

	b.ClearAttrs()
	c.Check(b.BaseAttrs, HasLen, 0)
}

// Base func tests

func (s *GomolSuite) BenchmarkBasicLogInsertion(c *C) {
	base := NewBase()
	base.InitLoggers()
	for i := 0; i < c.N; i++ {
		base.log(LevelDebug, map[string]interface{}{
			"attr1": "val1",
			"attr2": "val2",
			"attr3": 3,
			"attr4": 4,
		}, "msg %v %v", "string", 1234)
	}
	base.ShutdownLoggers()
}

func (s *GomolSuite) TestBaseDbg(c *C) {
	b := NewBase()

	l1 := newDefaultMemLogger()
	l2 := newDefaultMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Dbg("test")
	b.ShutdownLoggers()

	c.Assert(l1.Messages, HasLen, 1)
	c.Check(l1.Messages[0].Message, Equals, "test")
	c.Check(l1.Messages[0].Attrs, HasLen, 0)
	c.Check(l1.Messages[0].Level, Equals, LevelDebug)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test")
	c.Check(l2.Messages[0].Attrs, HasLen, 0)
	c.Check(l2.Messages[0].Level, Equals, LevelDebug)
}

func (s *GomolSuite) TestBaseDbgf(c *C) {
	b := NewBase()

	l1 := newDefaultMemLogger()
	l2 := newDefaultMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Dbgf("test %v", 1234)
	b.ShutdownLoggers()

	c.Assert(l1.Messages, HasLen, 1)
	c.Check(l1.Messages[0].Message, Equals, "test 1234")
	c.Check(l1.Messages[0].Attrs, HasLen, 0)
	c.Check(l1.Messages[0].Level, Equals, LevelDebug)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test 1234")
	c.Check(l2.Messages[0].Attrs, HasLen, 0)
	c.Check(l2.Messages[0].Level, Equals, LevelDebug)
}

func (s *GomolSuite) TestBaseDbgm(c *C) {
	b := NewBase()
	b.SetAttr("attr1", 1234)

	l1 := newDefaultMemLogger()
	l2 := newDefaultMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Dbgm(
		map[string]interface{}{
			"attr2": 4321,
			"attr3": "val3",
		},
		"test %v",
		1234)
	b.ShutdownLoggers()

	c.Assert(l1.Messages, HasLen, 1)
	c.Check(l1.Messages[0].Message, Equals, "test 1234")
	c.Assert(l1.Messages[0].Attrs, HasLen, 3)
	c.Check(l1.Messages[0].Attrs["attr1"], Equals, 1234)
	c.Check(l1.Messages[0].Attrs["attr2"], Equals, 4321)
	c.Check(l1.Messages[0].Attrs["attr3"], Equals, "val3")
	c.Check(l1.Messages[0].Level, Equals, LevelDebug)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test 1234")
	c.Assert(l2.Messages[0].Attrs, HasLen, 3)
	c.Check(l2.Messages[0].Attrs["attr1"], Equals, 1234)
	c.Check(l2.Messages[0].Attrs["attr2"], Equals, 4321)
	c.Check(l2.Messages[0].Attrs["attr3"], Equals, "val3")
	c.Check(l2.Messages[0].Level, Equals, LevelDebug)
}

func (s *GomolSuite) TestBaseInfo(c *C) {
	b := NewBase()

	l1 := newDefaultMemLogger()
	l2 := newDefaultMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Info("test")
	b.ShutdownLoggers()

	c.Assert(l1.Messages, HasLen, 1)
	c.Check(l1.Messages[0].Message, Equals, "test")
	c.Check(l1.Messages[0].Attrs, HasLen, 0)
	c.Check(l1.Messages[0].Level, Equals, LevelInfo)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test")
	c.Check(l2.Messages[0].Attrs, HasLen, 0)
	c.Check(l2.Messages[0].Level, Equals, LevelInfo)
}

func (s *GomolSuite) TestBaseInfof(c *C) {
	b := NewBase()

	l1 := newDefaultMemLogger()
	l2 := newDefaultMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Infof("test %v", 1234)
	b.ShutdownLoggers()

	c.Assert(l1.Messages, HasLen, 1)
	c.Check(l1.Messages[0].Message, Equals, "test 1234")
	c.Check(l1.Messages[0].Attrs, HasLen, 0)
	c.Check(l1.Messages[0].Level, Equals, LevelInfo)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test 1234")
	c.Check(l2.Messages[0].Attrs, HasLen, 0)
	c.Check(l2.Messages[0].Level, Equals, LevelInfo)
}

func (s *GomolSuite) TestBaseInfom(c *C) {
	b := NewBase()
	b.SetAttr("attr1", 1234)

	l1 := newDefaultMemLogger()
	l2 := newDefaultMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Infom(
		map[string]interface{}{
			"attr2": 4321,
			"attr3": "val3",
		},
		"test %v",
		1234)
	b.ShutdownLoggers()

	c.Assert(l1.Messages, HasLen, 1)
	c.Check(l1.Messages[0].Message, Equals, "test 1234")
	c.Assert(l1.Messages[0].Attrs, HasLen, 3)
	c.Check(l1.Messages[0].Attrs["attr1"], Equals, 1234)
	c.Check(l1.Messages[0].Attrs["attr2"], Equals, 4321)
	c.Check(l1.Messages[0].Attrs["attr3"], Equals, "val3")
	c.Check(l1.Messages[0].Level, Equals, LevelInfo)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test 1234")
	c.Assert(l2.Messages[0].Attrs, HasLen, 3)
	c.Check(l2.Messages[0].Attrs["attr1"], Equals, 1234)
	c.Check(l2.Messages[0].Attrs["attr2"], Equals, 4321)
	c.Check(l2.Messages[0].Attrs["attr3"], Equals, "val3")
	c.Check(l2.Messages[0].Level, Equals, LevelInfo)
}

func (s *GomolSuite) TestBaseWarn(c *C) {
	b := NewBase()

	l1 := newDefaultMemLogger()
	l2 := newDefaultMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Warn("test")
	b.ShutdownLoggers()

	c.Assert(l1.Messages, HasLen, 1)
	c.Check(l1.Messages[0].Message, Equals, "test")
	c.Check(l1.Messages[0].Attrs, HasLen, 0)
	c.Check(l1.Messages[0].Level, Equals, LevelWarning)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test")
	c.Check(l2.Messages[0].Attrs, HasLen, 0)
	c.Check(l2.Messages[0].Level, Equals, LevelWarning)
}

func (s *GomolSuite) TestBaseWarnf(c *C) {
	b := NewBase()

	l1 := newDefaultMemLogger()
	l2 := newDefaultMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Warnf("test %v", 1234)
	b.ShutdownLoggers()

	c.Assert(l1.Messages, HasLen, 1)
	c.Check(l1.Messages[0].Message, Equals, "test 1234")
	c.Check(l1.Messages[0].Attrs, HasLen, 0)
	c.Check(l1.Messages[0].Level, Equals, LevelWarning)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test 1234")
	c.Check(l2.Messages[0].Attrs, HasLen, 0)
	c.Check(l2.Messages[0].Level, Equals, LevelWarning)
}

func (s *GomolSuite) TestBaseWarnm(c *C) {
	b := NewBase()
	b.SetAttr("attr1", 1234)

	l1 := newDefaultMemLogger()
	l2 := newDefaultMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Warnm(
		map[string]interface{}{
			"attr2": 4321,
			"attr3": "val3",
		},
		"test %v",
		1234)
	b.ShutdownLoggers()

	c.Assert(l1.Messages, HasLen, 1)
	c.Check(l1.Messages[0].Message, Equals, "test 1234")
	c.Assert(l1.Messages[0].Attrs, HasLen, 3)
	c.Check(l1.Messages[0].Attrs["attr1"], Equals, 1234)
	c.Check(l1.Messages[0].Attrs["attr2"], Equals, 4321)
	c.Check(l1.Messages[0].Attrs["attr3"], Equals, "val3")
	c.Check(l1.Messages[0].Level, Equals, LevelWarning)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test 1234")
	c.Assert(l2.Messages[0].Attrs, HasLen, 3)
	c.Check(l2.Messages[0].Attrs["attr1"], Equals, 1234)
	c.Check(l2.Messages[0].Attrs["attr2"], Equals, 4321)
	c.Check(l2.Messages[0].Attrs["attr3"], Equals, "val3")
	c.Check(l2.Messages[0].Level, Equals, LevelWarning)
}

func (s *GomolSuite) TestBaseErr(c *C) {
	b := NewBase()

	l1 := newDefaultMemLogger()
	l2 := newDefaultMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Err("test")
	b.ShutdownLoggers()

	c.Assert(l1.Messages, HasLen, 1)
	c.Check(l1.Messages[0].Message, Equals, "test")
	c.Check(l1.Messages[0].Attrs, HasLen, 0)
	c.Check(l1.Messages[0].Level, Equals, LevelError)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test")
	c.Check(l2.Messages[0].Attrs, HasLen, 0)
	c.Check(l2.Messages[0].Level, Equals, LevelError)
}

func (s *GomolSuite) TestBaseErrf(c *C) {
	b := NewBase()

	l1 := newDefaultMemLogger()
	l2 := newDefaultMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Errf("test %v", 1234)
	b.ShutdownLoggers()

	c.Assert(l1.Messages, HasLen, 1)
	c.Check(l1.Messages[0].Message, Equals, "test 1234")
	c.Check(l1.Messages[0].Attrs, HasLen, 0)
	c.Check(l1.Messages[0].Level, Equals, LevelError)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test 1234")
	c.Check(l2.Messages[0].Attrs, HasLen, 0)
	c.Check(l2.Messages[0].Level, Equals, LevelError)
}

func (s *GomolSuite) TestBaseErrm(c *C) {
	b := NewBase()
	b.SetAttr("attr1", 1234)

	l1 := newDefaultMemLogger()
	l2 := newDefaultMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Errm(
		map[string]interface{}{
			"attr2": 4321,
			"attr3": "val3",
		},
		"test %v",
		1234)
	b.ShutdownLoggers()

	c.Assert(l1.Messages, HasLen, 1)
	c.Check(l1.Messages[0].Message, Equals, "test 1234")
	c.Assert(l1.Messages[0].Attrs, HasLen, 3)
	c.Check(l1.Messages[0].Attrs["attr1"], Equals, 1234)
	c.Check(l1.Messages[0].Attrs["attr2"], Equals, 4321)
	c.Check(l1.Messages[0].Attrs["attr3"], Equals, "val3")
	c.Check(l1.Messages[0].Level, Equals, LevelError)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test 1234")
	c.Assert(l2.Messages[0].Attrs, HasLen, 3)
	c.Check(l2.Messages[0].Attrs["attr1"], Equals, 1234)
	c.Check(l2.Messages[0].Attrs["attr2"], Equals, 4321)
	c.Check(l2.Messages[0].Attrs["attr3"], Equals, "val3")
	c.Check(l2.Messages[0].Level, Equals, LevelError)
}

func (s *GomolSuite) TestBaseFatal(c *C) {
	b := NewBase()

	l1 := newDefaultMemLogger()
	l2 := newDefaultMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Fatal("test")
	b.ShutdownLoggers()

	c.Assert(l1.Messages, HasLen, 1)
	c.Check(l1.Messages[0].Message, Equals, "test")
	c.Check(l1.Messages[0].Attrs, HasLen, 0)
	c.Check(l1.Messages[0].Level, Equals, LevelFatal)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test")
	c.Check(l2.Messages[0].Attrs, HasLen, 0)
	c.Check(l2.Messages[0].Level, Equals, LevelFatal)
}

func (s *GomolSuite) TestBaseFatalf(c *C) {
	b := NewBase()

	l1 := newDefaultMemLogger()
	l2 := newDefaultMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Fatalf("test %v", 1234)
	b.ShutdownLoggers()

	c.Assert(l1.Messages, HasLen, 1)
	c.Check(l1.Messages[0].Message, Equals, "test 1234")
	c.Check(l1.Messages[0].Attrs, HasLen, 0)
	c.Check(l1.Messages[0].Level, Equals, LevelFatal)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test 1234")
	c.Check(l2.Messages[0].Attrs, HasLen, 0)
	c.Check(l2.Messages[0].Level, Equals, LevelFatal)
}

func (s *GomolSuite) TestBaseFatalm(c *C) {
	b := NewBase()
	b.SetAttr("attr1", 1234)

	l1 := newDefaultMemLogger()
	l2 := newDefaultMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Fatalm(
		map[string]interface{}{
			"attr2": 4321,
			"attr3": "val3",
		},
		"test %v",
		1234)
	b.ShutdownLoggers()

	c.Assert(l1.Messages, HasLen, 1)
	c.Check(l1.Messages[0].Message, Equals, "test 1234")
	c.Assert(l1.Messages[0].Attrs, HasLen, 3)
	c.Check(l1.Messages[0].Attrs["attr1"], Equals, 1234)
	c.Check(l1.Messages[0].Attrs["attr2"], Equals, 4321)
	c.Check(l1.Messages[0].Attrs["attr3"], Equals, "val3")
	c.Check(l1.Messages[0].Level, Equals, LevelFatal)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test 1234")
	c.Assert(l2.Messages[0].Attrs, HasLen, 3)
	c.Check(l2.Messages[0].Attrs["attr1"], Equals, 1234)
	c.Check(l2.Messages[0].Attrs["attr2"], Equals, 4321)
	c.Check(l2.Messages[0].Attrs["attr3"], Equals, "val3")
	c.Check(l2.Messages[0].Level, Equals, LevelFatal)
}

func (s *GomolSuite) TestBaseDie(c *C) {
	b := NewBase()

	l1 := newDefaultMemLogger()
	l2 := newDefaultMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Die(1234, "test")

	c.Assert(l1.Messages, HasLen, 1)
	c.Check(l1.Messages[0].Message, Equals, "test")
	c.Check(l1.Messages[0].Attrs, HasLen, 0)
	c.Check(l1.Messages[0].Level, Equals, LevelFatal)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test")
	c.Check(l2.Messages[0].Attrs, HasLen, 0)
	c.Check(l2.Messages[0].Level, Equals, LevelFatal)

	c.Check(b.isInitialized, Equals, false)
	c.Check(curTestExiter.exited, Equals, true)
	c.Check(curTestExiter.code, Equals, 1234)
}

func (s *GomolSuite) TestBaseDief(c *C) {
	b := NewBase()

	l1 := newDefaultMemLogger()
	l2 := newDefaultMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Dief(1234, "test %v", 1234)

	c.Assert(l1.Messages, HasLen, 1)
	c.Check(l1.Messages[0].Message, Equals, "test 1234")
	c.Check(l1.Messages[0].Attrs, HasLen, 0)
	c.Check(l1.Messages[0].Level, Equals, LevelFatal)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test 1234")
	c.Check(l2.Messages[0].Attrs, HasLen, 0)
	c.Check(l2.Messages[0].Level, Equals, LevelFatal)

	c.Check(b.isInitialized, Equals, false)
	c.Check(curTestExiter.exited, Equals, true)
	c.Check(curTestExiter.code, Equals, 1234)
}

func (s *GomolSuite) TestBaseDiem(c *C) {
	b := NewBase()
	b.SetAttr("attr1", 1234)

	l1 := newDefaultMemLogger()
	l2 := newDefaultMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Diem(
		1234,
		map[string]interface{}{
			"attr2": 4321,
			"attr3": "val3",
		},
		"test %v",
		1234)
	b.ShutdownLoggers()

	c.Assert(l1.Messages, HasLen, 1)
	c.Check(l1.Messages[0].Message, Equals, "test 1234")
	c.Assert(l1.Messages[0].Attrs, HasLen, 3)
	c.Check(l1.Messages[0].Attrs["attr1"], Equals, 1234)
	c.Check(l1.Messages[0].Attrs["attr2"], Equals, 4321)
	c.Check(l1.Messages[0].Attrs["attr3"], Equals, "val3")
	c.Check(l1.Messages[0].Level, Equals, LevelFatal)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test 1234")
	c.Assert(l2.Messages[0].Attrs, HasLen, 3)
	c.Check(l2.Messages[0].Attrs["attr1"], Equals, 1234)
	c.Check(l2.Messages[0].Attrs["attr2"], Equals, 4321)
	c.Check(l2.Messages[0].Attrs["attr3"], Equals, "val3")
	c.Check(l2.Messages[0].Level, Equals, LevelFatal)

	c.Check(b.isInitialized, Equals, false)
	c.Check(curTestExiter.exited, Equals, true)
	c.Check(curTestExiter.code, Equals, 1234)
}

func (s *GomolSuite) TestBaseOrdering(c *C) {
	b := NewBase()
	b.SetAttr("attr1", 1234)

	l1 := newDefaultMemLogger()
	l2 := newDefaultMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Fatalm(
		map[string]interface{}{
			"attr2": 4321,
			"attr3": "val3",
		},
		"test %v",
		1234)
	b.Fatalm(
		map[string]interface{}{
			"attr4": 4321,
			"attr5": "val3",
		},
		"test %v",
		4321)
	b.ShutdownLoggers()

	c.Assert(l1.Messages, HasLen, 2)
	c.Check(l1.Messages[0].Message, Equals, "test 1234")
	c.Assert(l1.Messages[0].Attrs, HasLen, 3)
	c.Check(l1.Messages[0].Attrs["attr1"], Equals, 1234)
	c.Check(l1.Messages[0].Attrs["attr2"], Equals, 4321)
	c.Check(l1.Messages[0].Attrs["attr3"], Equals, "val3")
	c.Check(l1.Messages[0].Level, Equals, LevelFatal)
	c.Check(l1.Messages[1].Message, Equals, "test 4321")
	c.Assert(l1.Messages[1].Attrs, HasLen, 3)
	c.Check(l1.Messages[1].Attrs["attr1"], Equals, 1234)
	c.Check(l1.Messages[1].Attrs["attr4"], Equals, 4321)
	c.Check(l1.Messages[1].Attrs["attr5"], Equals, "val3")
	c.Check(l1.Messages[1].Level, Equals, LevelFatal)

	c.Assert(l2.Messages, HasLen, 2)
	c.Check(l2.Messages[0].Message, Equals, "test 1234")
	c.Assert(l2.Messages[0].Attrs, HasLen, 3)
	c.Check(l2.Messages[0].Attrs["attr1"], Equals, 1234)
	c.Check(l2.Messages[0].Attrs["attr2"], Equals, 4321)
	c.Check(l2.Messages[0].Attrs["attr3"], Equals, "val3")
	c.Check(l2.Messages[0].Level, Equals, LevelFatal)
	c.Check(l2.Messages[1].Message, Equals, "test 4321")
	c.Assert(l2.Messages[1].Attrs, HasLen, 3)
	c.Check(l2.Messages[1].Attrs["attr1"], Equals, 1234)
	c.Check(l2.Messages[1].Attrs["attr4"], Equals, 4321)
	c.Check(l2.Messages[1].Attrs["attr5"], Equals, "val3")
	c.Check(l2.Messages[1].Level, Equals, LevelFatal)
}
