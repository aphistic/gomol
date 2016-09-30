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

	c.Check(getLevelName(LogLevel(-1234)), Equals, "unknown")
}

func (s *GomolSuite) TestLevelString(c *C) {
	c.Check(LevelDebug.String(), Equals, "debug")
	c.Check(LevelInfo.String(), Equals, "info")
	c.Check(LevelWarning.String(), Equals, "warn")
	c.Check(LevelError.String(), Equals, "error")
	c.Check(LevelFatal.String(), Equals, "fatal")
	c.Check(LevelNone.String(), Equals, "none")
}

func (s *GomolSuite) TestToLogLevel(c *C) {
	var level LogLevel
	var err error

	level, err = ToLogLevel("dBg")
	c.Check(level, Equals, LevelDebug)
	c.Check(err, IsNil)
	level, err = ToLogLevel("DebuG")
	c.Check(level, Equals, LevelDebug)
	c.Check(err, IsNil)

	level, err = ToLogLevel("InFo")
	c.Check(level, Equals, LevelInfo)
	c.Check(err, IsNil)

	level, err = ToLogLevel("wARn")
	c.Check(level, Equals, LevelWarning)
	c.Check(err, IsNil)
	level, err = ToLogLevel("WaRNiNg")
	c.Check(level, Equals, LevelWarning)
	c.Check(err, IsNil)

	level, err = ToLogLevel("ErR")
	c.Check(level, Equals, LevelError)
	c.Check(err, IsNil)
	level, err = ToLogLevel("ERRoR")
	c.Check(level, Equals, LevelError)
	c.Check(err, IsNil)

	level, err = ToLogLevel("FaTaL")
	c.Check(level, Equals, LevelFatal)
	c.Check(err, IsNil)

	level, err = ToLogLevel("NonE")
	c.Check(level, Equals, LevelNone)
	c.Check(err, IsNil)
}

func (s *GomolSuite) TestToLogLevelError(c *C) {
	var level LogLevel
	var err error

	level, err = ToLogLevel("asdf")
	c.Check(level, Equals, LogLevel(0))
	c.Check(err, Equals, ErrUnknownLevel)
}

func (s *GomolSuite) TestNewMessageAttrsNil(c *C) {
	setClock(newTestClock(time.Now()))
	m := newMessage(clock().Now(), curDefault, LevelDebug, nil, "test")
	c.Check(m.Base, DeepEquals, curDefault)
	c.Check(m.Timestamp, Equals, clock().Now())
	c.Check(m.Level, Equals, LevelDebug)
	c.Check(m.Attrs, NotNil)
	c.Check(m.Attrs.Attrs(), HasLen, 0)
	c.Check(m.Msg, Equals, "test")
}

func (s *GomolSuite) TestNewMessageMsgAttrsNil(c *C) {
	setClock(newTestClock(time.Now()))

	ma := NewAttrs().
		SetAttr("msgAttr", "strVal").
		SetAttr("otherAttr", 4321)

	m := newMessage(clock().Now(), curDefault, LevelDebug, ma, "test")
	c.Check(m.Base, DeepEquals, curDefault)
	c.Check(m.Timestamp, Equals, clock().Now())
	c.Check(m.Level, Equals, LevelDebug)
	c.Check(m.Attrs, NotNil)
	c.Check(m.Attrs.Attrs(), HasLen, 2)
	c.Check(m.Attrs.GetAttr("msgAttr"), Equals, "strVal")
	c.Check(m.Attrs.GetAttr("otherAttr"), Equals, 4321)
	c.Check(m.Msg, Equals, "test")
}

func (s *GomolSuite) TestNewMessageFormat(c *C) {
	setClock(newTestClock(time.Now()))
	m := newMessage(clock().Now(), curDefault, LevelDebug, nil, "test %v %v", "str", 1234)
	c.Check(m.Base, DeepEquals, curDefault)
	c.Check(m.Timestamp, Equals, clock().Now())
	c.Check(m.Level, Equals, LevelDebug)
	c.Check(m.Attrs, NotNil)
	c.Check(m.Attrs.Attrs(), HasLen, 0)
	c.Check(m.Msg, Equals, "test str 1234")
}

func (s *GomolSuite) TestNewMessageFormatWithAttrs(c *C) {
	setClock(newTestClock(time.Now()))

	ma := NewAttrs().
		SetAttr("msgAttr", "strVal").
		SetAttr("otherAttr", 4321)

	m := newMessage(clock().Now(), curDefault, LevelDebug, ma, "test %v %v", "str", 1234)
	c.Check(m.Base, DeepEquals, curDefault)
	c.Check(m.Timestamp, Equals, clock().Now())
	c.Check(m.Level, Equals, LevelDebug)
	c.Check(m.Attrs, NotNil)
	c.Check(m.Attrs.Attrs(), HasLen, 2)
	c.Check(m.Attrs.GetAttr("msgAttr"), Equals, "strVal")
	c.Check(m.Attrs.GetAttr("otherAttr"), Equals, 4321)
	c.Check(m.Msg, Equals, "test str 1234")
}
