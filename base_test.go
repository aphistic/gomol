package gomol

import (
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type GomolSuite struct{}

var _ = Suite(&GomolSuite{})

func (s *GomolSuite) SetUpTest(c *C) {
	b := newBase()
	b.AddLogger(NewMemLogger())
	SetDefault(b)
}

func (s *GomolSuite) TestAddLogger(c *C) {
	b := newBase()
	c.Check(b.loggers, HasLen, 0)

	ml := NewMemLogger()
	c.Check(ml.base, IsNil)

	b.AddLogger(ml)
	c.Check(b.loggers, HasLen, 1)
	c.Check(ml.base, Equals, b)
}

func (s *GomolSuite) TestInitLoggers(c *C) {
	b := newBase()

	ml1 := NewMemLogger()
	ml2 := NewMemLogger()

	b.AddLogger(ml1)
	b.AddLogger(ml2)

	b.InitLoggers()

	c.Check(ml1.IsInitialized, Equals, true)
	c.Check(ml2.IsInitialized, Equals, true)
}

func (s *GomolSuite) TestShutdownLoggers(c *C) {
	b := newBase()

	ml1 := NewMemLogger()
	ml2 := NewMemLogger()

	b.AddLogger(ml1)
	b.AddLogger(ml2)

	b.ShutdownLoggers()

	c.Check(ml1.IsShutdown, Equals, true)
	c.Check(ml2.IsShutdown, Equals, true)
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

	b.Dbg("test")

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

	b.Dbgf("test %v", 1234)

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

	b.Dbgm(
		map[string]interface{}{
			"attr2": 4321,
			"attr3": "val3",
		},
		"test %v",
		1234)

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

	b.Info("test")

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

	b.Infof("test %v", 1234)

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

	b.Infom(
		map[string]interface{}{
			"attr2": 4321,
			"attr3": "val3",
		},
		"test %v",
		1234)

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

	b.Warn("test")

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

	b.Warnf("test %v", 1234)

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

	b.Warnm(
		map[string]interface{}{
			"attr2": 4321,
			"attr3": "val3",
		},
		"test %v",
		1234)

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

	b.Err("test")

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

	b.Errf("test %v", 1234)

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

	b.Errm(
		map[string]interface{}{
			"attr2": 4321,
			"attr3": "val3",
		},
		"test %v",
		1234)

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

	b.Fatal("test")

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

	b.Fatalf("test %v", 1234)

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

	b.Fatalm(
		map[string]interface{}{
			"attr2": 4321,
			"attr3": "val3",
		},
		"test %v",
		1234)

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
