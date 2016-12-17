package gomol

import (
	"time"

	"testing"

	. "github.com/onsi/gomega"
)

func (s *GomolSuite) TestNewLogAdapterEmpty(t *testing.T) {
	b := NewBase()

	la := b.NewLogAdapter(nil)
	Expect(la).ToNot(BeNil())
	Expect(la.base).To(Equal(b))

	Expect(la.attrs).ToNot(BeNil())
}

func (s *GomolSuite) TestNewLogAdapter(t *testing.T) {
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

func (s *GomolSuite) TestLogAdapterSetAttr(t *testing.T) {
	b := NewBase()

	la := b.NewLogAdapter(nil)
	la.SetAttr("foo", "bar")
	Expect(la.attrs.GetAttr("foo")).To(Equal("bar"))
}

func (s *GomolSuite) TestLogAdapterGetAttr(t *testing.T) {
	b := NewBase()

	la := b.NewLogAdapter(nil)
	la.SetAttr("attr1", 1)
	la.SetAttr("attr2", "val2")

	Expect(la.GetAttr("attr2")).To(Equal("val2"))
	Expect(la.GetAttr("notakey")).To(BeNil())
}

func (s *GomolSuite) TestLogAdapterRemoveAttr(t *testing.T) {
	b := NewBase()

	la := b.NewLogAdapter(NewAttrs().SetAttr("foo", "bar"))
	Expect(la.attrs.GetAttr("foo")).To(Equal("bar"))
	la.RemoveAttr("foo")
	Expect(la.attrs.GetAttr("foo")).To(BeNil())
}

func (s *GomolSuite) TestLogAdapterClearAttrs(t *testing.T) {
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

func (s *GomolSuite) TestLogAdapterLogWithTime(t *testing.T) {
	b := NewBase()
	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	b.InitLoggers()

	la := b.NewLogAdapter(NewAttrs().SetAttr("foo", "bar"))
	Expect(len(ml.Messages)).To(Equal(0))

	ts := time.Now()

	la.LogWithTime(LevelDebug, ts, NewAttrs().SetAttr("foo", "newBar"), "MessageM %d", 2)

	b.ShutdownLoggers()

	Expect(len(ml.Messages)).To(Equal(1))
	Expect(ml.Messages[0]).To(Equal(&memMessage{
		Timestamp: ts,
		Level:     LevelDebug,
		Message:   "MessageM 2",
		Attrs: map[string]interface{}{
			"foo": "newBar",
		},
	}))
}

func (s *GomolSuite) TestLogAdapterDebug(t *testing.T) {
	b := NewBase()
	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	b.InitLoggers()

	la := b.NewLogAdapter(NewAttrs().SetAttr("foo", "bar"))
	Expect(len(ml.Messages)).To(Equal(0))

	la.Dbg("Message 1")
	la.Debug("Message 2")
	la.Dbgf("MessageF %v", 1)
	la.Debugf("MessageF %v", 2)
	la.Dbgm(NewAttrs().SetAttr("attr1", "val1"), "MessageM %v", 1)
	la.Debugm(NewAttrs().SetAttr("foo", "newBar"), "MessageM %v", 2)

	b.ShutdownLoggers()

	Expect(len(ml.Messages)).To(Equal(6))
	Expect(ml.Messages[0]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelDebug,
		Message:   "Message 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	}))
	Expect(ml.Messages[1]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelDebug,
		Message:   "Message 2",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	}))
	Expect(ml.Messages[2]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelDebug,
		Message:   "MessageF 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	}))
	Expect(ml.Messages[3]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelDebug,
		Message:   "MessageF 2",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	}))
	Expect(ml.Messages[4]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelDebug,
		Message:   "MessageM 1",
		Attrs: map[string]interface{}{
			"foo":   "bar",
			"attr1": "val1",
		},
	}))
	Expect(ml.Messages[5]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelDebug,
		Message:   "MessageM 2",
		Attrs: map[string]interface{}{
			"foo": "newBar",
		},
	}))
}

func (s *GomolSuite) TestLogAdapterInfo(t *testing.T) {
	b := NewBase()
	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	b.InitLoggers()

	la := b.NewLogAdapter(NewAttrs().SetAttr("foo", "bar"))
	Expect(len(ml.Messages)).To(Equal(0))

	la.Info("Message 1")
	la.Infof("MessageF %v", 1)
	la.Infom(NewAttrs().SetAttr("attr1", "val1"), "MessageM %v", 1)

	b.ShutdownLoggers()

	Expect(len(ml.Messages)).To(Equal(3))
	Expect(ml.Messages[0]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelInfo,
		Message:   "Message 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	}))
	Expect(ml.Messages[1]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelInfo,
		Message:   "MessageF 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	}))
	Expect(ml.Messages[2]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelInfo,
		Message:   "MessageM 1",
		Attrs: map[string]interface{}{
			"foo":   "bar",
			"attr1": "val1",
		},
	}))
}

func (s *GomolSuite) TestLogAdapterWarn(t *testing.T) {
	b := NewBase()
	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	b.InitLoggers()

	la := b.NewLogAdapter(NewAttrs().SetAttr("foo", "bar"))
	Expect(len(ml.Messages)).To(Equal(0))

	la.Warn("Message 1")
	la.Warning("Message 2")
	la.Warnf("MessageF %v", 1)
	la.Warningf("MessageF %v", 2)
	la.Warnm(NewAttrs().SetAttr("attr1", "val1"), "MessageM %v", 1)
	la.Warningm(NewAttrs().SetAttr("foo", "newBar"), "MessageM %v", 2)

	b.ShutdownLoggers()

	Expect(len(ml.Messages)).To(Equal(6))
	Expect(ml.Messages[0]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelWarning,
		Message:   "Message 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	}))
	Expect(ml.Messages[1]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelWarning,
		Message:   "Message 2",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	}))
	Expect(ml.Messages[2]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelWarning,
		Message:   "MessageF 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	}))
	Expect(ml.Messages[3]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelWarning,
		Message:   "MessageF 2",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	}))
	Expect(ml.Messages[4]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelWarning,
		Message:   "MessageM 1",
		Attrs: map[string]interface{}{
			"foo":   "bar",
			"attr1": "val1",
		},
	}))
	Expect(ml.Messages[5]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelWarning,
		Message:   "MessageM 2",
		Attrs: map[string]interface{}{
			"foo": "newBar",
		},
	}))
}

func (s *GomolSuite) TestLogAdapterError(t *testing.T) {
	b := NewBase()
	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	b.InitLoggers()

	la := b.NewLogAdapter(NewAttrs().SetAttr("foo", "bar"))
	Expect(len(ml.Messages)).To(Equal(0))

	la.Err("Message 1")
	la.Error("Message 2")
	la.Errf("MessageF %v", 1)
	la.Errorf("MessageF %v", 2)
	la.Errm(NewAttrs().SetAttr("attr1", "val1"), "MessageM %v", 1)
	la.Errorm(NewAttrs().SetAttr("foo", "newBar"), "MessageM %v", 2)

	b.ShutdownLoggers()

	Expect(len(ml.Messages)).To(Equal(6))
	Expect(ml.Messages[0]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelError,
		Message:   "Message 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	}))
	Expect(ml.Messages[1]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelError,
		Message:   "Message 2",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	}))
	Expect(ml.Messages[2]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelError,
		Message:   "MessageF 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	}))
	Expect(ml.Messages[3]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelError,
		Message:   "MessageF 2",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	}))
	Expect(ml.Messages[4]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelError,
		Message:   "MessageM 1",
		Attrs: map[string]interface{}{
			"foo":   "bar",
			"attr1": "val1",
		},
	}))
	Expect(ml.Messages[5]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelError,
		Message:   "MessageM 2",
		Attrs: map[string]interface{}{
			"foo": "newBar",
		},
	}))
}

func (s *GomolSuite) TestLogAdapterFatal(t *testing.T) {
	b := NewBase()
	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	b.InitLoggers()

	la := b.NewLogAdapter(NewAttrs().SetAttr("foo", "bar"))
	Expect(len(ml.Messages)).To(Equal(0))

	la.Fatal("Message 1")
	la.Fatalf("MessageF %v", 1)
	la.Fatalm(NewAttrs().SetAttr("attr1", "val1"), "MessageM %v", 1)

	b.ShutdownLoggers()

	Expect(len(ml.Messages)).To(Equal(3))
	Expect(ml.Messages[0]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelFatal,
		Message:   "Message 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	}))
	Expect(ml.Messages[1]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelFatal,
		Message:   "MessageF 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	}))
	Expect(ml.Messages[2]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelFatal,
		Message:   "MessageM 1",
		Attrs: map[string]interface{}{
			"foo":   "bar",
			"attr1": "val1",
		},
	}))
}

func (s *GomolSuite) TestLogAdapterDie(t *testing.T) {
	b := NewBase()
	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	b.InitLoggers()

	la := b.NewLogAdapter(NewAttrs().SetAttr("foo", "bar"))
	Expect(len(ml.Messages)).To(Equal(0))

	la.Die(1234, "Message 1")

	b.ShutdownLoggers()

	Expect(len(ml.Messages)).To(Equal(1))
	Expect(ml.Messages[0]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelFatal,
		Message:   "Message 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	}))
	Expect(curTestExiter.exited).To(Equal(true))
	Expect(curTestExiter.code).To(Equal(1234))
}

func (s *GomolSuite) TestLogAdapterDief(t *testing.T) {
	b := NewBase()
	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	b.InitLoggers()

	la := b.NewLogAdapter(NewAttrs().SetAttr("foo", "bar"))
	Expect(len(ml.Messages)).To(Equal(0))

	la.Dief(1234, "MessageF %v", 1)

	b.ShutdownLoggers()

	Expect(len(ml.Messages)).To(Equal(1))
	Expect(ml.Messages[0]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelFatal,
		Message:   "MessageF 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	}))
	Expect(curTestExiter.exited).To(Equal(true))
	Expect(curTestExiter.code).To(Equal(1234))
}

func (s *GomolSuite) TestLogAdapterDiem(t *testing.T) {
	b := NewBase()
	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	b.InitLoggers()

	la := b.NewLogAdapter(NewAttrs().SetAttr("foo", "bar"))
	Expect(len(ml.Messages)).To(Equal(0))

	la.Diem(1234, NewAttrs().SetAttr("attr1", "val1"), "MessageM %v", 1)

	b.ShutdownLoggers()

	Expect(len(ml.Messages)).To(Equal(1))
	Expect(ml.Messages[0]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelFatal,
		Message:   "MessageM 1",
		Attrs: map[string]interface{}{
			"foo":   "bar",
			"attr1": "val1",
		},
	}))
	Expect(curTestExiter.exited).To(Equal(true))
	Expect(curTestExiter.code).To(Equal(1234))
}
