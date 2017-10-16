package gomol

import (
	"fmt"
	"time"

	"github.com/aphistic/sweet"
	"github.com/efritz/glock"
	. "github.com/onsi/gomega"
)

type LogAdapterSuite struct{}

func (s *LogAdapterSuite) TestNewLogAdapterEmpty(t sweet.T) {
	b := NewBase()

	la := b.NewLogAdapter(nil)
	Expect(la).ToNot(BeNil())
	Expect(la.base).To(Equal(b))

	Expect(la.attrs).ToNot(BeNil())
}

func (s *LogAdapterSuite) TestNewLogAdapter(t sweet.T) {
	b := NewBase()

	la := b.NewLogAdapter(NewAttrs().
		SetAttr("testNum", 1234).
		SetAttr("testStr", "foo"))
	Expect(la).ToNot(BeNil())
	Expect(la.base).To(Equal(b))

	Expect(la.attrs).ToNot(BeNil())
	Expect(la.attrs.GetAttr("testNum")).To(Equal(1234))
	Expect(la.attrs.GetAttr("testStr")).To(Equal("foo"))
}

func (s *LogAdapterSuite) TestLogAdapterSetAttr(t sweet.T) {
	b := NewBase()

	la := b.NewLogAdapter(nil)
	la.SetAttr("foo", "bar")
	Expect(la.attrs.GetAttr("foo")).To(Equal("bar"))
}

func (s *LogAdapterSuite) TestLogAdapterGetAttr(t sweet.T) {
	b := NewBase()

	la := b.NewLogAdapter(nil)
	la.SetAttr("attr1", 1)
	la.SetAttr("attr2", "val2")

	Expect(la.GetAttr("attr2")).To(Equal("val2"))
	Expect(la.GetAttr("notakey")).To(BeNil())
}

func (s *LogAdapterSuite) TestLogAdapterRemoveAttr(t sweet.T) {
	b := NewBase()

	la := b.NewLogAdapter(NewAttrs().SetAttr("foo", "bar"))
	Expect(la.attrs.GetAttr("foo")).To(Equal("bar"))
	la.RemoveAttr("foo")
	Expect(la.attrs.GetAttr("foo")).To(BeNil())
}

func (s *LogAdapterSuite) TestLogAdapterClearAttrs(t sweet.T) {
	b := NewBase()

	la := b.NewLogAdapter(NewAttrs().
		SetAttr("foo", "bar").
		SetAttr("baz", "qux"))
	Expect(la.attrs.GetAttr("foo")).To(Equal("bar"))
	Expect(la.attrs.GetAttr("baz")).To(Equal("qux"))
	la.ClearAttrs()
	Expect(la.attrs.GetAttr("foo")).To(BeNil())
	Expect(la.attrs.GetAttr("baz")).To(BeNil())
}

func (s *LogAdapterSuite) TestLogAdapterLogWithTime(t sweet.T) {
	b := NewBase()
	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	b.InitLoggers()

	la := b.NewLogAdapter(NewAttrs().SetAttr("foo", "bar"))
	Expect(ml.Messages()).To(HaveLen(0))

	ts := time.Now()

	la.LogWithTime(LevelDebug, ts, NewAttrs().SetAttr("foo", "newBar"), "MessageM %d", 2)

	b.ShutdownLoggers()

	Expect(ml.Messages()).To(HaveLen(1))
	Expect(ml.Messages()[0]).To(Equal(&memMessage{
		Timestamp: ts,
		Level:     LevelDebug,
		Message:   "MessageM 2",
		Attrs: map[string]interface{}{
			"foo": "newBar",
		},
		StringAttrs: map[string]string{
			"foo": "newBar",
		},
	}))
}

func (s *LogAdapterSuite) TestLogLevelLogWithTime(t sweet.T) {
	b := NewBase()
	b.SetLogLevel(LevelInfo)
	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	b.InitLoggers()

	la := b.NewLogAdapter(NewAttrs().SetAttr("foo", "bar"))
	la.SetLogLevel(LevelError)
	Expect(ml.Messages()).To(HaveLen(0))

	ts := time.Now()

	la.LogWithTime(LevelFatal, ts, NewAttrs().SetAttr("foo", "newBar"), "MessageM %d", LevelFatal)
	la.LogWithTime(LevelError, ts, NewAttrs().SetAttr("foo", "newBar"), "MessageM %d", LevelError)
	la.LogWithTime(LevelWarning, ts, NewAttrs().SetAttr("foo", "newBar"), "MessageM %d", LevelWarning)

	b.ShutdownLoggers()

	Expect(ml.Messages()).To(HaveLen(2))
	Expect(ml.Messages()[0]).To(Equal(&memMessage{
		Timestamp: ts,
		Level:     LevelFatal,
		Message:   "MessageM 2",
		Attrs: map[string]interface{}{
			"foo": "newBar",
		},
		StringAttrs: map[string]string{
			"foo": "newBar",
		},
	}))
	Expect(ml.Messages()[1]).To(Equal(&memMessage{
		Timestamp: ts,
		Level:     LevelError,
		Message:   "MessageM 3",
		Attrs: map[string]interface{}{
			"foo": "newBar",
		},
		StringAttrs: map[string]string{
			"foo": "newBar",
		},
	}))
}

func (s *LogAdapterSuite) TestLogLevelLog(t sweet.T) {
	b := NewBase()
	b.SetLogLevel(LevelInfo)
	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	b.InitLoggers()

	la := b.NewLogAdapter(NewAttrs().SetAttr("foo", "bar"))
	la.SetLogLevel(LevelError)
	Expect(ml.Messages()).To(HaveLen(0))

	la.Log(LevelFatal, NewAttrs().SetAttr("foo", "newBar"), "MessageM %d", LevelFatal)
	la.Log(LevelError, NewAttrs().SetAttr("foo", "newBar"), "MessageM %d", LevelError)
	la.Log(LevelWarning, NewAttrs().SetAttr("foo", "newBar"), "MessageM %d", LevelWarning)

	b.ShutdownLoggers()

	Expect(ml.Messages()).To(HaveLen(2))
	Expect(ml.Messages()[0]).To(Equal(&memMessage{
		Timestamp: ml.Messages()[0].Timestamp,
		Level:     LevelFatal,
		Message:   "MessageM 2",
		Attrs: map[string]interface{}{
			"foo": "newBar",
		},
		StringAttrs: map[string]string{
			"foo": "newBar",
		},
	}))
	Expect(ml.Messages()[1]).To(Equal(&memMessage{
		Timestamp: ml.Messages()[1].Timestamp,
		Level:     LevelError,
		Message:   "MessageM 3",
		Attrs: map[string]interface{}{
			"foo": "newBar",
		},
		StringAttrs: map[string]string{
			"foo": "newBar",
		},
	}))
}

func (s *LogAdapterSuite) TestLogAdapterThing(t sweet.T) {
	clock := glock.NewMockClock()
	b := NewBase(withClock(clock))
	b.SetAttr("base_attr", "foo")
	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	b.InitLoggers()

	la := b.NewLogAdapter(NewAttrsFromMap(map[string]interface{}{
		"adapter_attr": "bar",
	}))

	la.Dbgm(NewAttrsFromMap(map[string]interface{}{
		"log_attr": "baz",
	}), "Message %d", 1)

	b.ShutdownLoggers()

	Expect(ml.Messages()).To(HaveLen(1))
	Expect(ml.Messages()[0]).To(Equal(&memMessage{
		Timestamp: clock.Now(),
		Level:     LevelDebug,
		Message:   "Message 1",
		Attrs: map[string]interface{}{
			"base_attr":    "foo",
			"adapter_attr": "bar",
			"log_attr":     "baz",
		},
		StringAttrs: map[string]string{
			"base_attr":    "foo",
			"adapter_attr": "bar",
			"log_attr":     "baz",
		},
	}))
}

func (s *LogAdapterSuite) TestLogAdapterDebug(t sweet.T) {
	clock := glock.NewMockClock()

	b := NewBase(withClock(clock))
	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	b.InitLoggers()

	la := b.NewLogAdapter(NewAttrs().SetAttr("foo", "bar"))
	Expect(ml.Messages()).To(HaveLen(0))

	la.Dbg("Message 1")
	la.Debug("Message 2")
	la.Dbgf("MessageF %v", 1)
	la.Debugf("MessageF %v", 2)
	la.Dbgm(NewAttrs().SetAttr("attr1", "val1"), "MessageM %v", 1)
	la.Debugm(NewAttrs().SetAttr("foo", "newBar"), "MessageM %v", 2)

	b.ShutdownLoggers()

	Expect(ml.Messages()).To(HaveLen(6))
	Expect(ml.Messages()[0]).To(Equal(&memMessage{
		Timestamp: clock.Now(),
		Level:     LevelDebug,
		Message:   "Message 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
		StringAttrs: map[string]string{
			"foo": "bar",
		},
	}))
	Expect(ml.Messages()[1]).To(Equal(&memMessage{
		Timestamp: clock.Now(),
		Level:     LevelDebug,
		Message:   "Message 2",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
		StringAttrs: map[string]string{
			"foo": "bar",
		},
	}))
	Expect(ml.Messages()[2]).To(Equal(&memMessage{
		Timestamp: clock.Now(),
		Level:     LevelDebug,
		Message:   "MessageF 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
		StringAttrs: map[string]string{
			"foo": "bar",
		},
	}))
	Expect(ml.Messages()[3]).To(Equal(&memMessage{
		Timestamp: clock.Now(),
		Level:     LevelDebug,
		Message:   "MessageF 2",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
		StringAttrs: map[string]string{
			"foo": "bar",
		},
	}))
	Expect(ml.Messages()[4]).To(Equal(&memMessage{
		Timestamp: clock.Now(),
		Level:     LevelDebug,
		Message:   "MessageM 1",
		Attrs: map[string]interface{}{
			"foo":   "bar",
			"attr1": "val1",
		},
		StringAttrs: map[string]string{
			"foo":   "bar",
			"attr1": "val1",
		},
	}))
	Expect(ml.Messages()[5]).To(Equal(&memMessage{
		Timestamp: clock.Now(),
		Level:     LevelDebug,
		Message:   "MessageM 2",
		Attrs: map[string]interface{}{
			"foo": "newBar",
		},
		StringAttrs: map[string]string{
			"foo": "newBar",
		},
	}))
}

func (s *LogAdapterSuite) TestLogAdapterInfo(t sweet.T) {
	clock := glock.NewMockClock()
	b := NewBase(withClock(clock))
	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	b.InitLoggers()

	la := b.NewLogAdapter(NewAttrs().SetAttr("foo", "bar"))
	Expect(ml.Messages()).To(HaveLen(0))

	la.Info("Message 1")
	la.Infof("MessageF %v", 1)
	la.Infom(NewAttrs().SetAttr("attr1", "val1"), "MessageM %v", 1)

	b.ShutdownLoggers()

	Expect(ml.Messages()).To(HaveLen(3))
	Expect(ml.Messages()[0]).To(Equal(&memMessage{
		Timestamp: clock.Now(),
		Level:     LevelInfo,
		Message:   "Message 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
		StringAttrs: map[string]string{
			"foo": "bar",
		},
	}))
	Expect(ml.Messages()[1]).To(Equal(&memMessage{
		Timestamp: clock.Now(),
		Level:     LevelInfo,
		Message:   "MessageF 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
		StringAttrs: map[string]string{
			"foo": "bar",
		},
	}))
	Expect(ml.Messages()[2]).To(Equal(&memMessage{
		Timestamp: clock.Now(),
		Level:     LevelInfo,
		Message:   "MessageM 1",
		Attrs: map[string]interface{}{
			"foo":   "bar",
			"attr1": "val1",
		},
		StringAttrs: map[string]string{
			"foo":   "bar",
			"attr1": "val1",
		},
	}))
}

func (s *LogAdapterSuite) TestLogAdapterWarn(t sweet.T) {
	clock := glock.NewMockClock()
	b := NewBase(withClock(clock))
	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	b.InitLoggers()

	la := b.NewLogAdapter(NewAttrs().SetAttr("foo", "bar"))
	Expect(ml.Messages()).To(HaveLen(0))

	la.Warn("Message 1")
	la.Warning("Message 2")
	la.Warnf("MessageF %v", 1)
	la.Warningf("MessageF %v", 2)
	la.Warnm(NewAttrs().SetAttr("attr1", "val1"), "MessageM %v", 1)
	la.Warningm(NewAttrs().SetAttr("foo", "newBar"), "MessageM %v", 2)

	b.ShutdownLoggers()

	Expect(ml.Messages()).To(HaveLen(6))
	Expect(ml.Messages()[0]).To(Equal(&memMessage{
		Timestamp: clock.Now(),
		Level:     LevelWarning,
		Message:   "Message 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
		StringAttrs: map[string]string{
			"foo": "bar",
		},
	}))
	Expect(ml.Messages()[1]).To(Equal(&memMessage{
		Timestamp: clock.Now(),
		Level:     LevelWarning,
		Message:   "Message 2",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
		StringAttrs: map[string]string{
			"foo": "bar",
		},
	}))
	Expect(ml.Messages()[2]).To(Equal(&memMessage{
		Timestamp: clock.Now(),
		Level:     LevelWarning,
		Message:   "MessageF 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
		StringAttrs: map[string]string{
			"foo": "bar",
		},
	}))
	Expect(ml.Messages()[3]).To(Equal(&memMessage{
		Timestamp: clock.Now(),
		Level:     LevelWarning,
		Message:   "MessageF 2",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
		StringAttrs: map[string]string{
			"foo": "bar",
		},
	}))
	Expect(ml.Messages()[4]).To(Equal(&memMessage{
		Timestamp: clock.Now(),
		Level:     LevelWarning,
		Message:   "MessageM 1",
		Attrs: map[string]interface{}{
			"foo":   "bar",
			"attr1": "val1",
		},
		StringAttrs: map[string]string{
			"foo":   "bar",
			"attr1": "val1",
		},
	}))
	Expect(ml.Messages()[5]).To(Equal(&memMessage{
		Timestamp: clock.Now(),
		Level:     LevelWarning,
		Message:   "MessageM 2",
		Attrs: map[string]interface{}{
			"foo": "newBar",
		},
		StringAttrs: map[string]string{
			"foo": "newBar",
		},
	}))
}

func (s *LogAdapterSuite) TestLogAdapterError(t sweet.T) {
	clock := glock.NewMockClock()
	b := NewBase(withClock(clock))
	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	b.InitLoggers()

	la := b.NewLogAdapter(NewAttrs().SetAttr("foo", "bar"))
	Expect(ml.Messages()).To(HaveLen(0))

	la.Err("Message 1")
	la.Error("Message 2")
	la.Errf("MessageF %v", 1)
	la.Errorf("MessageF %v", 2)
	la.Errm(NewAttrs().SetAttr("attr1", "val1"), "MessageM %v", 1)
	la.Errorm(NewAttrs().SetAttr("foo", "newBar"), "MessageM %v", 2)

	b.ShutdownLoggers()

	Expect(ml.Messages()).To(HaveLen(6))
	Expect(ml.Messages()[0]).To(Equal(&memMessage{
		Timestamp: clock.Now(),
		Level:     LevelError,
		Message:   "Message 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
		StringAttrs: map[string]string{
			"foo": "bar",
		},
	}))
	Expect(ml.Messages()[1]).To(Equal(&memMessage{
		Timestamp: clock.Now(),
		Level:     LevelError,
		Message:   "Message 2",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
		StringAttrs: map[string]string{
			"foo": "bar",
		},
	}))
	Expect(ml.Messages()[2]).To(Equal(&memMessage{
		Timestamp: clock.Now(),
		Level:     LevelError,
		Message:   "MessageF 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
		StringAttrs: map[string]string{
			"foo": "bar",
		},
	}))
	Expect(ml.Messages()[3]).To(Equal(&memMessage{
		Timestamp: clock.Now(),
		Level:     LevelError,
		Message:   "MessageF 2",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
		StringAttrs: map[string]string{
			"foo": "bar",
		},
	}))
	Expect(ml.Messages()[4]).To(Equal(&memMessage{
		Timestamp: clock.Now(),
		Level:     LevelError,
		Message:   "MessageM 1",
		Attrs: map[string]interface{}{
			"foo":   "bar",
			"attr1": "val1",
		},
		StringAttrs: map[string]string{
			"foo":   "bar",
			"attr1": "val1",
		},
	}))
	Expect(ml.Messages()[5]).To(Equal(&memMessage{
		Timestamp: clock.Now(),
		Level:     LevelError,
		Message:   "MessageM 2",
		Attrs: map[string]interface{}{
			"foo": "newBar",
		},
		StringAttrs: map[string]string{
			"foo": "newBar",
		},
	}))
}

func (s *LogAdapterSuite) TestLogAdapterFatal(t sweet.T) {
	clock := glock.NewMockClock()
	b := NewBase(withClock(clock))
	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	b.InitLoggers()

	la := b.NewLogAdapter(NewAttrs().SetAttr("foo", "bar"))
	Expect(ml.Messages()).To(HaveLen(0))

	la.Fatal("Message 1")
	la.Fatalf("MessageF %v", 1)
	la.Fatalm(NewAttrs().SetAttr("attr1", "val1"), "MessageM %v", 1)

	b.ShutdownLoggers()

	Expect(ml.Messages()).To(HaveLen(3))
	Expect(ml.Messages()[0]).To(Equal(&memMessage{
		Timestamp: clock.Now(),
		Level:     LevelFatal,
		Message:   "Message 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
		StringAttrs: map[string]string{
			"foo": "bar",
		},
	}))
	Expect(ml.Messages()[1]).To(Equal(&memMessage{
		Timestamp: clock.Now(),
		Level:     LevelFatal,
		Message:   "MessageF 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
		StringAttrs: map[string]string{
			"foo": "bar",
		},
	}))
	Expect(ml.Messages()[2]).To(Equal(&memMessage{
		Timestamp: clock.Now(),
		Level:     LevelFatal,
		Message:   "MessageM 1",
		Attrs: map[string]interface{}{
			"foo":   "bar",
			"attr1": "val1",
		},
		StringAttrs: map[string]string{
			"foo":   "bar",
			"attr1": "val1",
		},
	}))
}

func (s *LogAdapterSuite) TestLogAdapterDie(t sweet.T) {
	clock := glock.NewMockClock()
	b := NewBase(withClock(clock))
	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	b.InitLoggers()

	la := b.NewLogAdapter(NewAttrs().SetAttr("foo", "bar"))
	Expect(ml.Messages()).To(HaveLen(0))

	la.Die(1234, "Message 1")

	b.ShutdownLoggers()

	Expect(ml.Messages()).To(HaveLen(1))
	Expect(ml.Messages()[0]).To(Equal(&memMessage{
		Timestamp: clock.Now(),
		Level:     LevelFatal,
		Message:   "Message 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
		StringAttrs: map[string]string{
			"foo": "bar",
		},
	}))
	Expect(curTestExiter.exited).To(Equal(true))
	Expect(curTestExiter.code).To(Equal(1234))
}

func (s *LogAdapterSuite) TestLogAdapterDief(t sweet.T) {
	clock := glock.NewMockClock()
	b := NewBase(withClock(clock))
	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	b.InitLoggers()

	la := b.NewLogAdapter(NewAttrs().SetAttr("foo", "bar"))
	Expect(ml.Messages()).To(HaveLen(0))

	la.Dief(1234, "MessageF %v", 1)

	b.ShutdownLoggers()

	Expect(ml.Messages()).To(HaveLen(1))
	Expect(ml.Messages()[0]).To(Equal(&memMessage{
		Timestamp: clock.Now(),
		Level:     LevelFatal,
		Message:   "MessageF 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
		StringAttrs: map[string]string{
			"foo": "bar",
		},
	}))
	Expect(curTestExiter.exited).To(Equal(true))
	Expect(curTestExiter.code).To(Equal(1234))
}

func (s *LogAdapterSuite) TestLogAdapterDiem(t sweet.T) {
	clock := glock.NewMockClock()
	b := NewBase(withClock(clock))
	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	b.InitLoggers()

	la := b.NewLogAdapter(NewAttrs().SetAttr("foo", "bar"))
	Expect(ml.Messages()).To(HaveLen(0))

	la.Diem(1234, NewAttrs().SetAttr("attr1", "val1"), "MessageM %v", 1)

	b.ShutdownLoggers()

	Expect(ml.Messages()).To(HaveLen(1))
	Expect(ml.Messages()[0]).To(Equal(&memMessage{
		Timestamp: clock.Now(),
		Level:     LevelFatal,
		Message:   "MessageM 1",
		Attrs: map[string]interface{}{
			"foo":   "bar",
			"attr1": "val1",
		},
		StringAttrs: map[string]string{
			"foo":   "bar",
			"attr1": "val1",
		},
	}))
	Expect(curTestExiter.exited).To(Equal(true))
	Expect(curTestExiter.code).To(Equal(1234))
}

type mockBase struct{}

func (mb *mockBase) LogWithTime(level LogLevel, ts time.Time, m *Attrs, msg string, a ...interface{}) error {
	return nil
}

func (mb *mockBase) Log(level LogLevel, m *Attrs, msg string, a ...interface{}) error {
	return nil
}

func (mb *mockBase) ShutdownLoggers() error {
	return fmt.Errorf("foo")
}

func (s *LogAdapterSuite) TestLogAdapterShutdownLoggers(t sweet.T) {
	la := NewLogAdapterFor(&mockBase{}, nil)
	Expect(la.ShutdownLoggers()).To(MatchError("foo"))
}
