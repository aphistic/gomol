package gomol

import (
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type GomolSuite struct{}

var _ = Suite(&GomolSuite{})

var testBase *Base

func (s *GomolSuite) SetUpTest(c *C) {
	testBase = newBase()
	testBase.AddLogger(NewMemLogger())
	testBase.InitLoggers()

	curDefault = newBase()
	curDefault.AddLogger(NewMemLogger())
	curDefault.InitLoggers()
}

func (s *GomolSuite) TearDownTest(c *C) {
	curDefault.ShutdownLoggers()

	testBase.ShutdownLoggers()
}

func (s *GomolSuite) TestAddLogger(c *C) {
	b := newBase()
	b.InitLoggers()
	c.Check(b.loggers, HasLen, 0)

	ml := NewMemLogger()
	c.Check(ml.IsInitialized(), Equals, false)
	c.Check(ml.base, IsNil)

	b.AddLogger(ml)
	c.Check(b.isInitialized, Equals, true)
	c.Assert(b.loggers, HasLen, 1)
	c.Check(b.loggers[0].IsInitialized(), Equals, true)
	c.Check(ml.base, Equals, b)
}

func (s *GomolSuite) TestAddLoggerAfterInit(c *C) {
	b := newBase()
	b.InitLoggers()

	ml := NewMemLogger()
	c.Check(ml.IsInitialized(), Equals, false)

	ret := b.AddLogger(ml)
	c.Check(ret, IsNil)
	c.Check(ml.IsInitialized(), Equals, true)
}

func (s *GomolSuite) TestAddLoggerAfterShutdown(c *C) {
	b := newBase()

	ml := NewMemLogger()
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

	ml := NewMemLoggerWithConfig(MemLoggerConfig{FailInit: true})
	c.Check(ml.IsInitialized(), Equals, false)

	ret := b.AddLogger(ml)
	c.Check(ret, NotNil)
	c.Check(ret.Error(), Equals, "Init failed")
	c.Check(ml.IsInitialized(), Equals, false)
	c.Check(b.loggers, HasLen, 0)
}

func (s *GomolSuite) TestAddLoggerAfterShutdownFail(c *C) {
	b := newBase()

	ml := NewMemLoggerWithConfig(MemLoggerConfig{FailShutdown: true})
	c.Check(ml.IsInitialized(), Equals, false)
	ml.InitLogger()
	c.Check(ml.IsInitialized(), Equals, true)

	ret := b.AddLogger(ml)
	c.Check(ret, NotNil)
	c.Check(ret.Error(), Equals, "Shutdown failed")
	c.Check(ml.IsInitialized(), Equals, true)
	c.Check(b.loggers, HasLen, 0)
}

func (s *GomolSuite) TestInitLoggers(c *C) {
	b := newBase()

	ml1 := NewMemLogger()
	ml2 := NewMemLogger()

	b.AddLogger(ml1)
	b.AddLogger(ml2)

	b.InitLoggers()

	c.Check(b.isInitialized, Equals, true)
	c.Check(ml1.IsInitialized(), Equals, true)
	c.Check(ml2.IsInitialized(), Equals, true)
}

func (s *GomolSuite) TestInitLoggersFail(c *C) {
	b := newBase()

	ml1 := NewMemLoggerWithConfig(MemLoggerConfig{FailInit: true})
	ml2 := NewMemLoggerWithConfig(MemLoggerConfig{FailInit: true})

	b.AddLogger(ml1)
	b.AddLogger(ml2)

	err := b.InitLoggers()
	c.Check(err, NotNil)
	c.Check(err.Error(), Equals, "Init failed")

	c.Check(b.isInitialized, Equals, false)
	c.Check(ml1.IsInitialized(), Equals, false)
	c.Check(ml2.IsInitialized(), Equals, false)
}

func (s *GomolSuite) TestShutdownLoggers(c *C) {
	b := newBase()

	ml1 := NewMemLogger()
	ml2 := NewMemLogger()

	b.AddLogger(ml1)
	b.AddLogger(ml2)

	b.InitLoggers()
	b.ShutdownLoggers()

	c.Check(ml1.isShutdown, Equals, true)
	c.Check(ml2.isShutdown, Equals, true)
}

func (s *GomolSuite) TestShutdownLoggersFail(c *C) {
	b := newBase()

	ml1 := NewMemLoggerWithConfig(MemLoggerConfig{FailShutdown: true})
	ml2 := NewMemLoggerWithConfig(MemLoggerConfig{FailShutdown: true})

	b.AddLogger(ml1)
	b.AddLogger(ml2)

	b.InitLoggers()
	err := b.ShutdownLoggers()
	c.Check(err, NotNil)
	c.Check(err.Error(), Equals, "Shutdown failed")

	c.Check(ml1.isShutdown, Equals, false)
	c.Check(ml2.isShutdown, Equals, false)
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

	l1 := NewMemLogger()
	l2 := NewMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Dbg("test")
	b.ShutdownLoggers()

	c.Assert(l1.Messages, HasLen, 1)
	c.Check(l1.Messages[0].Message, Equals, "test")
	c.Check(l1.Messages[0].Attrs, HasLen, 0)
	c.Check(l1.Messages[0].Level, Equals, levelDbg)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test")
	c.Check(l2.Messages[0].Attrs, HasLen, 0)
	c.Check(l2.Messages[0].Level, Equals, levelDbg)
}

func (s *GomolSuite) TestBaseDbgf(c *C) {
	b := newBase()

	l1 := NewMemLogger()
	l2 := NewMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Dbgf("test %v", 1234)
	b.ShutdownLoggers()

	c.Assert(l1.Messages, HasLen, 1)
	c.Check(l1.Messages[0].Message, Equals, "test 1234")
	c.Check(l1.Messages[0].Attrs, HasLen, 0)
	c.Check(l1.Messages[0].Level, Equals, levelDbg)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test 1234")
	c.Check(l2.Messages[0].Attrs, HasLen, 0)
	c.Check(l2.Messages[0].Level, Equals, levelDbg)
}

func (s *GomolSuite) TestBaseDbgm(c *C) {
	b := newBase()
	b.SetAttr("attr1", 1234)

	l1 := NewMemLogger()
	l2 := NewMemLogger()

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
	c.Check(l1.Messages[0].Level, Equals, levelDbg)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test 1234")
	c.Assert(l2.Messages[0].Attrs, HasLen, 3)
	c.Check(l2.Messages[0].Attrs["attr1"], Equals, 1234)
	c.Check(l2.Messages[0].Attrs["attr2"], Equals, 4321)
	c.Check(l2.Messages[0].Attrs["attr3"], Equals, "val3")
	c.Check(l2.Messages[0].Level, Equals, levelDbg)
}

func (s *GomolSuite) TestBaseInfo(c *C) {
	b := newBase()

	l1 := NewMemLogger()
	l2 := NewMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Info("test")
	b.ShutdownLoggers()

	c.Assert(l1.Messages, HasLen, 1)
	c.Check(l1.Messages[0].Message, Equals, "test")
	c.Check(l1.Messages[0].Attrs, HasLen, 0)
	c.Check(l1.Messages[0].Level, Equals, levelInfo)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test")
	c.Check(l2.Messages[0].Attrs, HasLen, 0)
	c.Check(l2.Messages[0].Level, Equals, levelInfo)
}

func (s *GomolSuite) TestBaseInfof(c *C) {
	b := newBase()

	l1 := NewMemLogger()
	l2 := NewMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Infof("test %v", 1234)
	b.ShutdownLoggers()

	c.Assert(l1.Messages, HasLen, 1)
	c.Check(l1.Messages[0].Message, Equals, "test 1234")
	c.Check(l1.Messages[0].Attrs, HasLen, 0)
	c.Check(l1.Messages[0].Level, Equals, levelInfo)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test 1234")
	c.Check(l2.Messages[0].Attrs, HasLen, 0)
	c.Check(l2.Messages[0].Level, Equals, levelInfo)
}

func (s *GomolSuite) TestBaseInfom(c *C) {
	b := newBase()
	b.SetAttr("attr1", 1234)

	l1 := NewMemLogger()
	l2 := NewMemLogger()

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
	c.Check(l1.Messages[0].Level, Equals, levelInfo)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test 1234")
	c.Assert(l2.Messages[0].Attrs, HasLen, 3)
	c.Check(l2.Messages[0].Attrs["attr1"], Equals, 1234)
	c.Check(l2.Messages[0].Attrs["attr2"], Equals, 4321)
	c.Check(l2.Messages[0].Attrs["attr3"], Equals, "val3")
	c.Check(l2.Messages[0].Level, Equals, levelInfo)
}

func (s *GomolSuite) TestBaseWarn(c *C) {
	b := newBase()

	l1 := NewMemLogger()
	l2 := NewMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Warn("test")
	b.ShutdownLoggers()

	c.Assert(l1.Messages, HasLen, 1)
	c.Check(l1.Messages[0].Message, Equals, "test")
	c.Check(l1.Messages[0].Attrs, HasLen, 0)
	c.Check(l1.Messages[0].Level, Equals, levelWarn)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test")
	c.Check(l2.Messages[0].Attrs, HasLen, 0)
	c.Check(l2.Messages[0].Level, Equals, levelWarn)
}

func (s *GomolSuite) TestBaseWarnf(c *C) {
	b := newBase()

	l1 := NewMemLogger()
	l2 := NewMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Warnf("test %v", 1234)
	b.ShutdownLoggers()

	c.Assert(l1.Messages, HasLen, 1)
	c.Check(l1.Messages[0].Message, Equals, "test 1234")
	c.Check(l1.Messages[0].Attrs, HasLen, 0)
	c.Check(l1.Messages[0].Level, Equals, levelWarn)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test 1234")
	c.Check(l2.Messages[0].Attrs, HasLen, 0)
	c.Check(l2.Messages[0].Level, Equals, levelWarn)
}

func (s *GomolSuite) TestBaseWarnm(c *C) {
	b := newBase()
	b.SetAttr("attr1", 1234)

	l1 := NewMemLogger()
	l2 := NewMemLogger()

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
	c.Check(l1.Messages[0].Level, Equals, levelWarn)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test 1234")
	c.Assert(l2.Messages[0].Attrs, HasLen, 3)
	c.Check(l2.Messages[0].Attrs["attr1"], Equals, 1234)
	c.Check(l2.Messages[0].Attrs["attr2"], Equals, 4321)
	c.Check(l2.Messages[0].Attrs["attr3"], Equals, "val3")
	c.Check(l2.Messages[0].Level, Equals, levelWarn)
}

func (s *GomolSuite) TestBaseErr(c *C) {
	b := newBase()

	l1 := NewMemLogger()
	l2 := NewMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Err("test")
	b.ShutdownLoggers()

	c.Assert(l1.Messages, HasLen, 1)
	c.Check(l1.Messages[0].Message, Equals, "test")
	c.Check(l1.Messages[0].Attrs, HasLen, 0)
	c.Check(l1.Messages[0].Level, Equals, levelError)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test")
	c.Check(l2.Messages[0].Attrs, HasLen, 0)
	c.Check(l2.Messages[0].Level, Equals, levelError)
}

func (s *GomolSuite) TestBaseErrf(c *C) {
	b := newBase()

	l1 := NewMemLogger()
	l2 := NewMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Errf("test %v", 1234)
	b.ShutdownLoggers()

	c.Assert(l1.Messages, HasLen, 1)
	c.Check(l1.Messages[0].Message, Equals, "test 1234")
	c.Check(l1.Messages[0].Attrs, HasLen, 0)
	c.Check(l1.Messages[0].Level, Equals, levelError)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test 1234")
	c.Check(l2.Messages[0].Attrs, HasLen, 0)
	c.Check(l2.Messages[0].Level, Equals, levelError)
}

func (s *GomolSuite) TestBaseErrm(c *C) {
	b := newBase()
	b.SetAttr("attr1", 1234)

	l1 := NewMemLogger()
	l2 := NewMemLogger()

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
	c.Check(l1.Messages[0].Level, Equals, levelError)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test 1234")
	c.Assert(l2.Messages[0].Attrs, HasLen, 3)
	c.Check(l2.Messages[0].Attrs["attr1"], Equals, 1234)
	c.Check(l2.Messages[0].Attrs["attr2"], Equals, 4321)
	c.Check(l2.Messages[0].Attrs["attr3"], Equals, "val3")
	c.Check(l2.Messages[0].Level, Equals, levelError)
}

func (s *GomolSuite) TestBaseFatal(c *C) {
	b := newBase()

	l1 := NewMemLogger()
	l2 := NewMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Fatal("test")
	b.ShutdownLoggers()

	c.Assert(l1.Messages, HasLen, 1)
	c.Check(l1.Messages[0].Message, Equals, "test")
	c.Check(l1.Messages[0].Attrs, HasLen, 0)
	c.Check(l1.Messages[0].Level, Equals, levelFatal)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test")
	c.Check(l2.Messages[0].Attrs, HasLen, 0)
	c.Check(l2.Messages[0].Level, Equals, levelFatal)
}

func (s *GomolSuite) TestBaseFatalf(c *C) {
	b := newBase()

	l1 := NewMemLogger()
	l2 := NewMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Fatalf("test %v", 1234)
	b.ShutdownLoggers()

	c.Assert(l1.Messages, HasLen, 1)
	c.Check(l1.Messages[0].Message, Equals, "test 1234")
	c.Check(l1.Messages[0].Attrs, HasLen, 0)
	c.Check(l1.Messages[0].Level, Equals, levelFatal)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test 1234")
	c.Check(l2.Messages[0].Attrs, HasLen, 0)
	c.Check(l2.Messages[0].Level, Equals, levelFatal)
}

func (s *GomolSuite) TestBaseFatalm(c *C) {
	b := newBase()
	b.SetAttr("attr1", 1234)

	l1 := NewMemLogger()
	l2 := NewMemLogger()

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
	c.Check(l1.Messages[0].Level, Equals, levelFatal)

	c.Assert(l2.Messages, HasLen, 1)
	c.Check(l2.Messages[0].Message, Equals, "test 1234")
	c.Assert(l2.Messages[0].Attrs, HasLen, 3)
	c.Check(l2.Messages[0].Attrs["attr1"], Equals, 1234)
	c.Check(l2.Messages[0].Attrs["attr2"], Equals, 4321)
	c.Check(l2.Messages[0].Attrs["attr3"], Equals, "val3")
	c.Check(l2.Messages[0].Level, Equals, levelFatal)
}

func (s *GomolSuite) TestBaseOrdering(c *C) {
	b := newBase()
	b.SetAttr("attr1", 1234)

	l1 := NewMemLogger()
	l2 := NewMemLogger()

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
	c.Check(l1.Messages[0].Level, Equals, levelFatal)
	c.Check(l1.Messages[1].Message, Equals, "test 4321")
	c.Assert(l1.Messages[1].Attrs, HasLen, 3)
	c.Check(l1.Messages[1].Attrs["attr1"], Equals, 1234)
	c.Check(l1.Messages[1].Attrs["attr4"], Equals, 4321)
	c.Check(l1.Messages[1].Attrs["attr5"], Equals, "val3")
	c.Check(l1.Messages[1].Level, Equals, levelFatal)

	c.Assert(l2.Messages, HasLen, 2)
	c.Check(l2.Messages[0].Message, Equals, "test 1234")
	c.Assert(l2.Messages[0].Attrs, HasLen, 3)
	c.Check(l2.Messages[0].Attrs["attr1"], Equals, 1234)
	c.Check(l2.Messages[0].Attrs["attr2"], Equals, 4321)
	c.Check(l2.Messages[0].Attrs["attr3"], Equals, "val3")
	c.Check(l2.Messages[0].Level, Equals, levelFatal)
	c.Check(l2.Messages[1].Message, Equals, "test 4321")
	c.Assert(l2.Messages[1].Attrs, HasLen, 3)
	c.Check(l2.Messages[1].Attrs["attr1"], Equals, 1234)
	c.Check(l2.Messages[1].Attrs["attr4"], Equals, 4321)
	c.Check(l2.Messages[1].Attrs["attr5"], Equals, "val3")
	c.Check(l2.Messages[1].Level, Equals, levelFatal)
}
