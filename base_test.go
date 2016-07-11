package gomol

import (
	"testing"

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
	curTestExiter = &testExiter{}
	setExiter(curTestExiter)

	testBase = newBase()
	testBase.AddLogger(newDefaultMemLogger())
	testBase.InitLoggers()

	curDefault = newBase()
	curDefault.AddLogger(newDefaultMemLogger())
	curDefault.InitLoggers()
}

func (s *GomolSuite) TearDownTest(c *C) {
	curDefault.ShutdownLoggers()

	testBase.ShutdownLoggers()
}

func (s *GomolSuite) TestShouldLog(c *C) {
	b := newBase()
	b.SetLogLevel(LEVEL_INFO)
	c.Check(b.shouldLog(LEVEL_UNKNOWN), Equals, false)
	c.Check(b.shouldLog(LEVEL_DEBUG), Equals, false)
	c.Check(b.shouldLog(LEVEL_INFO), Equals, true)
	c.Check(b.shouldLog(LEVEL_WARNING), Equals, true)
	c.Check(b.shouldLog(LEVEL_ERROR), Equals, true)
	c.Check(b.shouldLog(LEVEL_FATAL), Equals, true)

	b.SetLogLevel(LEVEL_FATAL)
	c.Check(b.shouldLog(LEVEL_UNKNOWN), Equals, false)
	c.Check(b.shouldLog(LEVEL_DEBUG), Equals, false)
	c.Check(b.shouldLog(LEVEL_INFO), Equals, false)
	c.Check(b.shouldLog(LEVEL_WARNING), Equals, false)
	c.Check(b.shouldLog(LEVEL_ERROR), Equals, false)
	c.Check(b.shouldLog(LEVEL_FATAL), Equals, true)

	b.SetLogLevel(LEVEL_NONE)
	c.Check(b.shouldLog(LEVEL_UNKNOWN), Equals, false)
	c.Check(b.shouldLog(LEVEL_DEBUG), Equals, false)
	c.Check(b.shouldLog(LEVEL_INFO), Equals, false)
	c.Check(b.shouldLog(LEVEL_WARNING), Equals, false)
	c.Check(b.shouldLog(LEVEL_ERROR), Equals, false)
	c.Check(b.shouldLog(LEVEL_FATAL), Equals, false)
}

func (s *GomolSuite) TestSetLogLevel(c *C) {
	b := newBase()
	b.InitLoggers()
	ml := newDefaultMemLogger()
	b.AddLogger(ml)

	b.SetLogLevel(LEVEL_WARNING)
	b.Dbg("test")
	b.Info("test")
	b.Warn("test")
	b.Err("test")
	b.Fatal("test")
	b.ShutdownLoggers()
	c.Check(ml.Messages, HasLen, 3)
}

func (s *GomolSuite) TestAddLogger(c *C) {
	b := newBase()
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
	b := newBase()
	b.InitLoggers()

	ml := newDefaultMemLogger()
	c.Check(ml.IsInitialized(), Equals, false)

	ret := b.AddLogger(ml)
	c.Check(ret, IsNil)
	c.Check(ml.IsInitialized(), Equals, true)
}

func (s *GomolSuite) TestAddLoggerAfterShutdown(c *C) {
	b := newBase()

	ml := newDefaultMemLogger()
	c.Check(ml.IsInitialized(), Equals, false)
	ml.InitLogger()
	c.Check(ml.IsInitialized(), Equals, true)

	ret := b.AddLogger(ml)
	c.Check(ret, IsNil)
	c.Check(ml.IsInitialized(), Equals, false)
}

func (s *GomolSuite) TestAddLoggerAfterInitFail(c *C) {
	b := newBase()
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
	b := newBase()

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
	b := newBase()

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
	b := newBase()

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
	b := newBase()

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
	b := newBase()
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
	b := newBase()
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
	b := newBase()

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
	b := newBase()

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
	b := newBase()

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
	b := newBase()

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
	b := newBase()

	b.SetAttr("attr1", 1)
	c.Check(b.BaseAttrs, HasLen, 1)
	c.Check(b.BaseAttrs["attr1"], Equals, 1)
	b.SetAttr("attr2", "val2")
	c.Check(b.BaseAttrs, HasLen, 2)
	c.Check(b.BaseAttrs["attr2"], Equals, "val2")
}

func (s *GomolSuite) TestRemoveAttr(c *C) {
	b := newBase()

	b.SetAttr("attr1", 1)
	c.Check(b.BaseAttrs, HasLen, 1)
	c.Check(b.BaseAttrs["attr1"], Equals, 1)

	b.RemoveAttr("attr1")
	c.Check(b.BaseAttrs, HasLen, 0)
}

func (s *GomolSuite) TestClearAttrs(c *C) {
	b := newBase()

	b.SetAttr("attr1", 1)
	b.SetAttr("attr2", "val2")
	c.Check(b.BaseAttrs, HasLen, 2)

	b.ClearAttrs()
	c.Check(b.BaseAttrs, HasLen, 0)
}

// Base func tests

func (s *GomolSuite) TestBaseDbg(c *C) {
	b := newBase()

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
	c.Check(l1.Messages[0].Level, Equals, LEVEL_DEBUG)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test")
	c.Check(l2.Messages[0].Attrs, HasLen, 0)
	c.Check(l2.Messages[0].Level, Equals, LEVEL_DEBUG)
}

func (s *GomolSuite) TestBaseDbgf(c *C) {
	b := newBase()

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
	c.Check(l1.Messages[0].Level, Equals, LEVEL_DEBUG)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test 1234")
	c.Check(l2.Messages[0].Attrs, HasLen, 0)
	c.Check(l2.Messages[0].Level, Equals, LEVEL_DEBUG)
}

func (s *GomolSuite) TestBaseDbgm(c *C) {
	b := newBase()
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
	c.Check(l1.Messages[0].Level, Equals, LEVEL_DEBUG)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test 1234")
	c.Assert(l2.Messages[0].Attrs, HasLen, 3)
	c.Check(l2.Messages[0].Attrs["attr1"], Equals, 1234)
	c.Check(l2.Messages[0].Attrs["attr2"], Equals, 4321)
	c.Check(l2.Messages[0].Attrs["attr3"], Equals, "val3")
	c.Check(l2.Messages[0].Level, Equals, LEVEL_DEBUG)
}

func (s *GomolSuite) TestBaseInfo(c *C) {
	b := newBase()

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
	c.Check(l1.Messages[0].Level, Equals, LEVEL_INFO)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test")
	c.Check(l2.Messages[0].Attrs, HasLen, 0)
	c.Check(l2.Messages[0].Level, Equals, LEVEL_INFO)
}

func (s *GomolSuite) TestBaseInfof(c *C) {
	b := newBase()

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
	c.Check(l1.Messages[0].Level, Equals, LEVEL_INFO)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test 1234")
	c.Check(l2.Messages[0].Attrs, HasLen, 0)
	c.Check(l2.Messages[0].Level, Equals, LEVEL_INFO)
}

func (s *GomolSuite) TestBaseInfom(c *C) {
	b := newBase()
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
	c.Check(l1.Messages[0].Level, Equals, LEVEL_INFO)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test 1234")
	c.Assert(l2.Messages[0].Attrs, HasLen, 3)
	c.Check(l2.Messages[0].Attrs["attr1"], Equals, 1234)
	c.Check(l2.Messages[0].Attrs["attr2"], Equals, 4321)
	c.Check(l2.Messages[0].Attrs["attr3"], Equals, "val3")
	c.Check(l2.Messages[0].Level, Equals, LEVEL_INFO)
}

func (s *GomolSuite) TestBaseWarn(c *C) {
	b := newBase()

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
	c.Check(l1.Messages[0].Level, Equals, LEVEL_WARNING)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test")
	c.Check(l2.Messages[0].Attrs, HasLen, 0)
	c.Check(l2.Messages[0].Level, Equals, LEVEL_WARNING)
}

func (s *GomolSuite) TestBaseWarnf(c *C) {
	b := newBase()

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
	c.Check(l1.Messages[0].Level, Equals, LEVEL_WARNING)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test 1234")
	c.Check(l2.Messages[0].Attrs, HasLen, 0)
	c.Check(l2.Messages[0].Level, Equals, LEVEL_WARNING)
}

func (s *GomolSuite) TestBaseWarnm(c *C) {
	b := newBase()
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
	c.Check(l1.Messages[0].Level, Equals, LEVEL_WARNING)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test 1234")
	c.Assert(l2.Messages[0].Attrs, HasLen, 3)
	c.Check(l2.Messages[0].Attrs["attr1"], Equals, 1234)
	c.Check(l2.Messages[0].Attrs["attr2"], Equals, 4321)
	c.Check(l2.Messages[0].Attrs["attr3"], Equals, "val3")
	c.Check(l2.Messages[0].Level, Equals, LEVEL_WARNING)
}

func (s *GomolSuite) TestBaseErr(c *C) {
	b := newBase()

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
	c.Check(l1.Messages[0].Level, Equals, LEVEL_ERROR)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test")
	c.Check(l2.Messages[0].Attrs, HasLen, 0)
	c.Check(l2.Messages[0].Level, Equals, LEVEL_ERROR)
}

func (s *GomolSuite) TestBaseErrf(c *C) {
	b := newBase()

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
	c.Check(l1.Messages[0].Level, Equals, LEVEL_ERROR)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test 1234")
	c.Check(l2.Messages[0].Attrs, HasLen, 0)
	c.Check(l2.Messages[0].Level, Equals, LEVEL_ERROR)
}

func (s *GomolSuite) TestBaseErrm(c *C) {
	b := newBase()
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
	c.Check(l1.Messages[0].Level, Equals, LEVEL_ERROR)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test 1234")
	c.Assert(l2.Messages[0].Attrs, HasLen, 3)
	c.Check(l2.Messages[0].Attrs["attr1"], Equals, 1234)
	c.Check(l2.Messages[0].Attrs["attr2"], Equals, 4321)
	c.Check(l2.Messages[0].Attrs["attr3"], Equals, "val3")
	c.Check(l2.Messages[0].Level, Equals, LEVEL_ERROR)
}

func (s *GomolSuite) TestBaseFatal(c *C) {
	b := newBase()

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
	c.Check(l1.Messages[0].Level, Equals, LEVEL_FATAL)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test")
	c.Check(l2.Messages[0].Attrs, HasLen, 0)
	c.Check(l2.Messages[0].Level, Equals, LEVEL_FATAL)
}

func (s *GomolSuite) TestBaseFatalf(c *C) {
	b := newBase()

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
	c.Check(l1.Messages[0].Level, Equals, LEVEL_FATAL)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test 1234")
	c.Check(l2.Messages[0].Attrs, HasLen, 0)
	c.Check(l2.Messages[0].Level, Equals, LEVEL_FATAL)
}

func (s *GomolSuite) TestBaseFatalm(c *C) {
	b := newBase()
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
	c.Check(l1.Messages[0].Level, Equals, LEVEL_FATAL)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test 1234")
	c.Assert(l2.Messages[0].Attrs, HasLen, 3)
	c.Check(l2.Messages[0].Attrs["attr1"], Equals, 1234)
	c.Check(l2.Messages[0].Attrs["attr2"], Equals, 4321)
	c.Check(l2.Messages[0].Attrs["attr3"], Equals, "val3")
	c.Check(l2.Messages[0].Level, Equals, LEVEL_FATAL)
}

func (s *GomolSuite) TestBaseDie(c *C) {
	b := newBase()

	l1 := newDefaultMemLogger()
	l2 := newDefaultMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Die(1234, "test")

	c.Assert(l1.Messages, HasLen, 1)
	c.Check(l1.Messages[0].Message, Equals, "test")
	c.Check(l1.Messages[0].Attrs, HasLen, 0)
	c.Check(l1.Messages[0].Level, Equals, LEVEL_FATAL)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test")
	c.Check(l2.Messages[0].Attrs, HasLen, 0)
	c.Check(l2.Messages[0].Level, Equals, LEVEL_FATAL)

	c.Check(b.isInitialized, Equals, false)
	c.Check(curTestExiter.exited, Equals, true)
	c.Check(curTestExiter.code, Equals, 1234)
}

func (s *GomolSuite) TestBaseDief(c *C) {
	b := newBase()

	l1 := newDefaultMemLogger()
	l2 := newDefaultMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Dief(1234, "test %v", 1234)

	c.Assert(l1.Messages, HasLen, 1)
	c.Check(l1.Messages[0].Message, Equals, "test 1234")
	c.Check(l1.Messages[0].Attrs, HasLen, 0)
	c.Check(l1.Messages[0].Level, Equals, LEVEL_FATAL)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test 1234")
	c.Check(l2.Messages[0].Attrs, HasLen, 0)
	c.Check(l2.Messages[0].Level, Equals, LEVEL_FATAL)

	c.Check(b.isInitialized, Equals, false)
	c.Check(curTestExiter.exited, Equals, true)
	c.Check(curTestExiter.code, Equals, 1234)
}

func (s *GomolSuite) TestBaseDiem(c *C) {
	b := newBase()
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
	c.Check(l1.Messages[0].Level, Equals, LEVEL_FATAL)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test 1234")
	c.Assert(l2.Messages[0].Attrs, HasLen, 3)
	c.Check(l2.Messages[0].Attrs["attr1"], Equals, 1234)
	c.Check(l2.Messages[0].Attrs["attr2"], Equals, 4321)
	c.Check(l2.Messages[0].Attrs["attr3"], Equals, "val3")
	c.Check(l2.Messages[0].Level, Equals, LEVEL_FATAL)

	c.Check(b.isInitialized, Equals, false)
	c.Check(curTestExiter.exited, Equals, true)
	c.Check(curTestExiter.code, Equals, 1234)
}

func (s *GomolSuite) TestBaseOrdering(c *C) {
	b := newBase()
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
	c.Check(l1.Messages[0].Level, Equals, LEVEL_FATAL)
	c.Check(l1.Messages[1].Message, Equals, "test 4321")
	c.Assert(l1.Messages[1].Attrs, HasLen, 3)
	c.Check(l1.Messages[1].Attrs["attr1"], Equals, 1234)
	c.Check(l1.Messages[1].Attrs["attr4"], Equals, 4321)
	c.Check(l1.Messages[1].Attrs["attr5"], Equals, "val3")
	c.Check(l1.Messages[1].Level, Equals, LEVEL_FATAL)

	c.Assert(l2.Messages, HasLen, 2)
	c.Check(l2.Messages[0].Message, Equals, "test 1234")
	c.Assert(l2.Messages[0].Attrs, HasLen, 3)
	c.Check(l2.Messages[0].Attrs["attr1"], Equals, 1234)
	c.Check(l2.Messages[0].Attrs["attr2"], Equals, 4321)
	c.Check(l2.Messages[0].Attrs["attr3"], Equals, "val3")
	c.Check(l2.Messages[0].Level, Equals, LEVEL_FATAL)
	c.Check(l2.Messages[1].Message, Equals, "test 4321")
	c.Assert(l2.Messages[1].Attrs, HasLen, 3)
	c.Check(l2.Messages[1].Attrs["attr1"], Equals, 1234)
	c.Check(l2.Messages[1].Attrs["attr4"], Equals, 4321)
	c.Check(l2.Messages[1].Attrs["attr5"], Equals, "val3")
	c.Check(l2.Messages[1].Level, Equals, LEVEL_FATAL)
}
