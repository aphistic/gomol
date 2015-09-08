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
