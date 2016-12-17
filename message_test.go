package gomol

import (
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func (s *GomolSuite) TestLevelGetName(t *testing.T) {
	Expect(getLevelName(LevelDebug)).To(Equal("debug"))
	Expect(getLevelName(LevelInfo)).To(Equal("info"))
	Expect(getLevelName(LevelWarning)).To(Equal("warn"))
	Expect(getLevelName(LevelError)).To(Equal("error"))
	Expect(getLevelName(LevelFatal)).To(Equal("fatal"))
	Expect(getLevelName(LevelNone)).To(Equal("none"))

	Expect(getLevelName(LogLevel(-1234))).To(Equal("unknown"))
}

func (s *GomolSuite) TestLevelString(t *testing.T) {
	Expect(LevelDebug.String()).To(Equal("debug"))
	Expect(LevelInfo.String()).To(Equal("info"))
	Expect(LevelWarning.String()).To(Equal("warn"))
	Expect(LevelError.String()).To(Equal("error"))
	Expect(LevelFatal.String()).To(Equal("fatal"))
	Expect(LevelNone.String()).To(Equal("none"))
}

func (s *GomolSuite) TestToLogLevel(t *testing.T) {
	var level LogLevel
	var err error

	level, err = ToLogLevel("dBg")
	Expect(level).To(Equal(LevelDebug))
	Expect(err).To(BeNil())
	level, err = ToLogLevel("DebuG")
	Expect(level).To(Equal(LevelDebug))
	Expect(err).To(BeNil())

	level, err = ToLogLevel("InFo")
	Expect(level).To(Equal(LevelInfo))
	Expect(err).To(BeNil())

	level, err = ToLogLevel("wARn")
	Expect(level).To(Equal(LevelWarning))
	Expect(err).To(BeNil())
	level, err = ToLogLevel("WaRNiNg")
	Expect(level).To(Equal(LevelWarning))
	Expect(err).To(BeNil())

	level, err = ToLogLevel("ErR")
	Expect(level).To(Equal(LevelError))
	Expect(err).To(BeNil())
	level, err = ToLogLevel("ERRoR")
	Expect(level).To(Equal(LevelError))
	Expect(err).To(BeNil())

	level, err = ToLogLevel("FaTaL")
	Expect(level).To(Equal(LevelFatal))
	Expect(err).To(BeNil())

	level, err = ToLogLevel("NonE")
	Expect(level).To(Equal(LevelNone))
	Expect(err).To(BeNil())
}

func (s *GomolSuite) TestToLogLevelError(t *testing.T) {
	var level LogLevel
	var err error

	level, err = ToLogLevel("asdf")
	Expect(level).To(Equal(LogLevel(0)))
	Expect(err).To(Equal(ErrUnknownLevel))
}

func (s *GomolSuite) TestNewMessageAttrsNil(t *testing.T) {
	setClock(newTestClock(time.Now()))
	m := newMessage(clock().Now(), curDefault, LevelDebug, nil, "test")
	Expect(m.Base).To(Equal(curDefault))
	Expect(m.Timestamp).To(Equal(clock().Now()))
	Expect(m.Level).To(Equal(LevelDebug))
	Expect(m.Attrs).ToNot(BeNil())
	Expect(m.Attrs.Attrs()).To(HaveLen(0))
	Expect(m.Msg).To(Equal("test"))
}

func (s *GomolSuite) TestNewMessageMsgAttrsNil(t *testing.T) {
	setClock(newTestClock(time.Now()))

	ma := NewAttrs().
		SetAttr("msgAttr", "strVal").
		SetAttr("otherAttr", 4321)

	m := newMessage(clock().Now(), curDefault, LevelDebug, ma, "test")
	Expect(m.Base).To(Equal(curDefault))
	Expect(m.Timestamp).To(Equal(clock().Now()))
	Expect(m.Level).To(Equal(LevelDebug))
	Expect(m.Attrs).ToNot(BeNil())
	Expect(m.Attrs.Attrs()).To(HaveLen(2))
	Expect(m.Attrs.GetAttr("msgAttr")).To(Equal("strVal"))
	Expect(m.Attrs.GetAttr("otherAttr")).To(Equal(4321))
	Expect(m.Msg).To(Equal("test"))
}

func (s *GomolSuite) TestNewMessageFormat(t *testing.T) {
	setClock(newTestClock(time.Now()))
	m := newMessage(clock().Now(), curDefault, LevelDebug, nil, "test %v %v", "str", 1234)
	Expect(m.Base).To(Equal(curDefault))
	Expect(m.Timestamp).To(Equal(clock().Now()))
	Expect(m.Level).To(Equal(LevelDebug))
	Expect(m.Attrs).ToNot(BeNil())
	Expect(m.Attrs.Attrs()).To(HaveLen(0))
	Expect(m.Msg).To(Equal("test str 1234"))
}

func (s *GomolSuite) TestNewMessageFormatWithAttrs(t *testing.T) {
	setClock(newTestClock(time.Now()))

	ma := NewAttrs().
		SetAttr("msgAttr", "strVal").
		SetAttr("otherAttr", 4321)

	m := newMessage(clock().Now(), curDefault, LevelDebug, ma, "test %v %v", "str", 1234)
	Expect(m.Base).To(Equal(curDefault))
	Expect(m.Timestamp).To(Equal(clock().Now()))
	Expect(m.Level).To(Equal(LevelDebug))
	Expect(m.Attrs).ToNot(BeNil())
	Expect(m.Attrs.Attrs()).To(HaveLen(2))
	Expect(m.Attrs.GetAttr("msgAttr")).To(Equal("strVal"))
	Expect(m.Attrs.GetAttr("otherAttr")).To(Equal(4321))
	Expect(m.Msg).To(Equal("test str 1234"))
}
