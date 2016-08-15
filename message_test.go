package gomol

import (
	"time"

	. "gopkg.in/check.v1"
)

func (s *GomolSuite) TestLevelGetName(c *C) {
	c.Check(getLevelName(LevelDebug), Equals, "debug")
	c.Check(getLevelName(LevelInfo), Equals, "info")
	c.Check(getLevelName(LevelWarning), Equals, "warn")
	c.Check(getLevelName(LevelError), Equals, "error")
	c.Check(getLevelName(LevelFatal), Equals, "fatal")
	c.Check(getLevelName(LevelNone), Equals, "none")
	c.Check(getLevelName(LevelUnknown), Equals, "unknown")
}

func (s *GomolSuite) TestLevelString(c *C) {
	c.Check(LevelDebug.String(), Equals, "debug")
	c.Check(LevelInfo.String(), Equals, "info")
	c.Check(LevelWarning.String(), Equals, "warn")
	c.Check(LevelError.String(), Equals, "error")
	c.Check(LevelFatal.String(), Equals, "fatal")
	c.Check(LevelNone.String(), Equals, "none")
	c.Check(LevelUnknown.String(), Equals, "unknown")
}

func (s *GomolSuite) TestNewMessageAttrsNil(c *C) {
	setClock(newTestClock(time.Now()))
	m := newMessage(clock().Now(), curDefault, LevelDebug, nil, "test")
	c.Check(m.Base, DeepEquals, curDefault)
	c.Check(m.Timestamp, Equals, clock().Now())
	c.Check(m.Level, Equals, LevelDebug)
	c.Check(m.Attrs, NotNil)
	c.Check(m.Attrs, HasLen, 0)
	c.Check(m.Msg, Equals, "test")
}

func (s *GomolSuite) TestNewMessageMsgAttrsNil(c *C) {
	setClock(newTestClock(time.Now()))

	ma := map[string]interface{}{
		"msgAttr":   "strVal",
		"otherAttr": 4321,
	}

	m := newMessage(clock().Now(), curDefault, LevelDebug, ma, "test")
	c.Check(m.Base, DeepEquals, curDefault)
	c.Check(m.Timestamp, Equals, clock().Now())
	c.Check(m.Level, Equals, LevelDebug)
	c.Check(m.Attrs, NotNil)
	c.Check(m.Attrs, HasLen, 2)
	c.Check(m.Attrs["msgAttr"], Equals, "strVal")
	c.Check(m.Attrs["otherAttr"], Equals, 4321)
	c.Check(m.Msg, Equals, "test")
}

func (s *GomolSuite) TestNewMessageFormat(c *C) {
	setClock(newTestClock(time.Now()))
	m := newMessage(clock().Now(), curDefault, LevelDebug, nil, "test %v %v", "str", 1234)
	c.Check(m.Base, DeepEquals, curDefault)
	c.Check(m.Timestamp, Equals, clock().Now())
	c.Check(m.Level, Equals, LevelDebug)
	c.Check(m.Attrs, NotNil)
	c.Check(m.Attrs, HasLen, 0)
	c.Check(m.Msg, Equals, "test str 1234")
}

func (s *GomolSuite) TestNewMessageFormatWithAttrs(c *C) {
	setClock(newTestClock(time.Now()))

	ma := map[string]interface{}{
		"msgAttr":   "strVal",
		"otherAttr": 4321,
	}

	m := newMessage(clock().Now(), curDefault, LevelDebug, ma, "test %v %v", "str", 1234)
	c.Check(m.Base, DeepEquals, curDefault)
	c.Check(m.Timestamp, Equals, clock().Now())
	c.Check(m.Level, Equals, LevelDebug)
	c.Check(m.Attrs, NotNil)
	c.Check(m.Attrs, HasLen, 2)
	c.Check(m.Attrs["msgAttr"], Equals, "strVal")
	c.Check(m.Attrs["otherAttr"], Equals, 4321)
	c.Check(m.Msg, Equals, "test str 1234")
}

func (s *GomolSuite) TestMergeAttrsNil(c *C) {
	attrs := mergeAttrs(nil, nil)
	c.Check(attrs, NotNil)
	c.Check(attrs, HasLen, 0)
}

func (s *GomolSuite) TestMergeAttrsBaseAttrs(c *C) {
	ba := map[string]interface{}{
		"baseAttr":  "strVal",
		"otherAttr": 1234,
	}
	attrs := mergeAttrs(ba, nil)
	c.Check(attrs, NotNil)
	c.Check(attrs, HasLen, 2)
	c.Check(attrs["baseAttr"], Equals, "strVal")
	c.Check(attrs["otherAttr"], Equals, 1234)
}

func (s *GomolSuite) TestMergeAttrsMsgAttrs(c *C) {
	ma := map[string]interface{}{
		"msgAttr":   "strVal",
		"otherAttr": 4321,
	}
	attrs := mergeAttrs(nil, ma)
	c.Check(attrs, NotNil)
	c.Check(attrs, HasLen, 2)
	c.Check(attrs["msgAttr"], Equals, "strVal")
	c.Check(attrs["otherAttr"], Equals, 4321)
}

func (s *GomolSuite) TestMergeAttrs(c *C) {
	ba := map[string]interface{}{
		"baseAttr":  "strVal",
		"otherAttr": 1234,
	}
	ma := map[string]interface{}{
		"msgAttr":   "strVal",
		"otherAttr": 4321,
	}
	attrs := mergeAttrs(ba, ma)
	c.Check(attrs, NotNil)
	c.Check(attrs, HasLen, 3)
	c.Check(attrs["msgAttr"], Equals, "strVal")
	c.Check(attrs["baseAttr"], Equals, "strVal")
	c.Check(attrs["otherAttr"], Equals, 4321)
}
