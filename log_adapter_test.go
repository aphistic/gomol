package gomol

import (
	"time"

	. "gopkg.in/check.v1"
)

func (s *GomolSuite) TestNewLogAdapterEmpty(c *C) {
	b := NewBase()

	la := b.NewLogAdapter(nil)
	c.Check(la, NotNil)
	c.Check(la.base, Equals, b)

	c.Check(la.attrs, NotNil)
}

func (s *GomolSuite) TestNewLogAdapter(c *C) {
	b := NewBase()

	la := b.NewLogAdapter(NewAttrs().
		SetAttr("testNum", 1234).
		SetAttr("testStr", "foo"))
	c.Check(la, NotNil)
	c.Check(la.base, Equals, b)

	c.Check(la.attrs, NotNil)
	c.Check(la.attrs.GetAttr("testNum"), Equals, 1234)
	c.Check(la.attrs.GetAttr("testStr"), Equals, "foo")
}

func (s *GomolSuite) TestLogAdapterSetAttr(c *C) {
	b := NewBase()

	la := b.NewLogAdapter(nil)
	la.SetAttr("foo", "bar")
	c.Check(la.attrs.GetAttr("foo"), Equals, "bar")
}

func (s *GomolSuite) TestLogAdapterGetAttr(c *C) {
	b := NewBase()

	la := b.NewLogAdapter(nil)
	la.SetAttr("attr1", 1)
	la.SetAttr("attr2", "val2")

	c.Check(la.GetAttr("attr2"), Equals, "val2")
	c.Check(la.GetAttr("notakey"), IsNil)
}

func (s *GomolSuite) TestLogAdapterRemoveAttr(c *C) {
	b := NewBase()

	la := b.NewLogAdapter(NewAttrs().SetAttr("foo", "bar"))
	c.Check(la.attrs.GetAttr("foo"), Equals, "bar")
	la.RemoveAttr("foo")
	c.Check(la.attrs.GetAttr("foo"), IsNil)
}

func (s *GomolSuite) TestLogAdapterClearAttrs(c *C) {
	b := NewBase()

	la := b.NewLogAdapter(NewAttrs().
		SetAttr("foo", "bar").
		SetAttr("baz", "qux"))
	c.Check(la.attrs.GetAttr("foo"), Equals, "bar")
	c.Check(la.attrs.GetAttr("baz"), Equals, "qux")
	la.ClearAttrs()
	c.Check(la.attrs.GetAttr("foo"), IsNil)
	c.Check(la.attrs.GetAttr("baz"), IsNil)
}

func (s *GomolSuite) TestLogAdapterLogWithTime(c *C) {
	b := NewBase()
	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	b.InitLoggers()

	la := b.NewLogAdapter(NewAttrs().SetAttr("foo", "bar"))
	c.Check(len(ml.Messages), Equals, 0)

	ts := time.Now()

	la.LogWithTime(LevelDebug, ts, NewAttrs().SetAttr("foo", "newBar"), "MessageM %d", 2)

	b.ShutdownLoggers()

	c.Assert(len(ml.Messages), Equals, 1)
	c.Check(ml.Messages[0], DeepEquals, &memMessage{
		Timestamp: ts,
		Level:     LevelDebug,
		Message:   "MessageM 2",
		Attrs: map[string]interface{}{
			"foo": "newBar",
		},
	})
}

func (s *GomolSuite) TestLogAdapterDebug(c *C) {
	b := NewBase()
	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	b.InitLoggers()

	la := b.NewLogAdapter(NewAttrs().SetAttr("foo", "bar"))
	c.Check(len(ml.Messages), Equals, 0)

	la.Dbg("Message 1")
	la.Debug("Message 2")
	la.Dbgf("MessageF %v", 1)
	la.Debugf("MessageF %v", 2)
	la.Dbgm(NewAttrs().SetAttr("attr1", "val1"), "MessageM %v", 1)
	la.Debugm(NewAttrs().SetAttr("foo", "newBar"), "MessageM %v", 2)

	b.ShutdownLoggers()

	c.Assert(len(ml.Messages), Equals, 6)
	c.Check(ml.Messages[0], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LevelDebug,
		Message:   "Message 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	})
	c.Check(ml.Messages[1], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LevelDebug,
		Message:   "Message 2",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	})
	c.Check(ml.Messages[2], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LevelDebug,
		Message:   "MessageF 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	})
	c.Check(ml.Messages[3], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LevelDebug,
		Message:   "MessageF 2",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	})
	c.Check(ml.Messages[4], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LevelDebug,
		Message:   "MessageM 1",
		Attrs: map[string]interface{}{
			"foo":   "bar",
			"attr1": "val1",
		},
	})
	c.Check(ml.Messages[5], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LevelDebug,
		Message:   "MessageM 2",
		Attrs: map[string]interface{}{
			"foo": "newBar",
		},
	})
}

func (s *GomolSuite) TestLogAdapterInfo(c *C) {
	b := NewBase()
	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	b.InitLoggers()

	la := b.NewLogAdapter(NewAttrs().SetAttr("foo", "bar"))
	c.Check(len(ml.Messages), Equals, 0)

	la.Info("Message 1")
	la.Infof("MessageF %v", 1)
	la.Infom(NewAttrs().SetAttr("attr1", "val1"), "MessageM %v", 1)

	b.ShutdownLoggers()

	c.Assert(len(ml.Messages), Equals, 3)
	c.Check(ml.Messages[0], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LevelInfo,
		Message:   "Message 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	})
	c.Check(ml.Messages[1], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LevelInfo,
		Message:   "MessageF 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	})
	c.Check(ml.Messages[2], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LevelInfo,
		Message:   "MessageM 1",
		Attrs: map[string]interface{}{
			"foo":   "bar",
			"attr1": "val1",
		},
	})
}

func (s *GomolSuite) TestLogAdapterWarn(c *C) {
	b := NewBase()
	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	b.InitLoggers()

	la := b.NewLogAdapter(NewAttrs().SetAttr("foo", "bar"))
	c.Check(len(ml.Messages), Equals, 0)

	la.Warn("Message 1")
	la.Warning("Message 2")
	la.Warnf("MessageF %v", 1)
	la.Warningf("MessageF %v", 2)
	la.Warnm(NewAttrs().SetAttr("attr1", "val1"), "MessageM %v", 1)
	la.Warningm(NewAttrs().SetAttr("foo", "newBar"), "MessageM %v", 2)

	b.ShutdownLoggers()

	c.Assert(len(ml.Messages), Equals, 6)
	c.Check(ml.Messages[0], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LevelWarning,
		Message:   "Message 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	})
	c.Check(ml.Messages[1], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LevelWarning,
		Message:   "Message 2",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	})
	c.Check(ml.Messages[2], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LevelWarning,
		Message:   "MessageF 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	})
	c.Check(ml.Messages[3], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LevelWarning,
		Message:   "MessageF 2",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	})
	c.Check(ml.Messages[4], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LevelWarning,
		Message:   "MessageM 1",
		Attrs: map[string]interface{}{
			"foo":   "bar",
			"attr1": "val1",
		},
	})
	c.Check(ml.Messages[5], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LevelWarning,
		Message:   "MessageM 2",
		Attrs: map[string]interface{}{
			"foo": "newBar",
		},
	})
}

func (s *GomolSuite) TestLogAdapterError(c *C) {
	b := NewBase()
	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	b.InitLoggers()

	la := b.NewLogAdapter(NewAttrs().SetAttr("foo", "bar"))
	c.Check(len(ml.Messages), Equals, 0)

	la.Err("Message 1")
	la.Error("Message 2")
	la.Errf("MessageF %v", 1)
	la.Errorf("MessageF %v", 2)
	la.Errm(NewAttrs().SetAttr("attr1", "val1"), "MessageM %v", 1)
	la.Errorm(NewAttrs().SetAttr("foo", "newBar"), "MessageM %v", 2)

	b.ShutdownLoggers()

	c.Assert(len(ml.Messages), Equals, 6)
	c.Check(ml.Messages[0], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LevelError,
		Message:   "Message 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	})
	c.Check(ml.Messages[1], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LevelError,
		Message:   "Message 2",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	})
	c.Check(ml.Messages[2], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LevelError,
		Message:   "MessageF 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	})
	c.Check(ml.Messages[3], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LevelError,
		Message:   "MessageF 2",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	})
	c.Check(ml.Messages[4], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LevelError,
		Message:   "MessageM 1",
		Attrs: map[string]interface{}{
			"foo":   "bar",
			"attr1": "val1",
		},
	})
	c.Check(ml.Messages[5], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LevelError,
		Message:   "MessageM 2",
		Attrs: map[string]interface{}{
			"foo": "newBar",
		},
	})
}

func (s *GomolSuite) TestLogAdapterFatal(c *C) {
	b := NewBase()
	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	b.InitLoggers()

	la := b.NewLogAdapter(NewAttrs().SetAttr("foo", "bar"))
	c.Check(len(ml.Messages), Equals, 0)

	la.Fatal("Message 1")
	la.Fatalf("MessageF %v", 1)
	la.Fatalm(NewAttrs().SetAttr("attr1", "val1"), "MessageM %v", 1)

	b.ShutdownLoggers()

	c.Assert(len(ml.Messages), Equals, 3)
	c.Check(ml.Messages[0], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LevelFatal,
		Message:   "Message 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	})
	c.Check(ml.Messages[1], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LevelFatal,
		Message:   "MessageF 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	})
	c.Check(ml.Messages[2], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LevelFatal,
		Message:   "MessageM 1",
		Attrs: map[string]interface{}{
			"foo":   "bar",
			"attr1": "val1",
		},
	})
}

func (s *GomolSuite) TestLogAdapterDie(c *C) {
	b := NewBase()
	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	b.InitLoggers()

	la := b.NewLogAdapter(NewAttrs().SetAttr("foo", "bar"))
	c.Check(len(ml.Messages), Equals, 0)

	la.Die(1234, "Message 1")

	b.ShutdownLoggers()

	c.Assert(len(ml.Messages), Equals, 1)
	c.Check(ml.Messages[0], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LevelFatal,
		Message:   "Message 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	})
	c.Check(curTestExiter.exited, Equals, true)
	c.Check(curTestExiter.code, Equals, 1234)
}

func (s *GomolSuite) TestLogAdapterDief(c *C) {
	b := NewBase()
	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	b.InitLoggers()

	la := b.NewLogAdapter(NewAttrs().SetAttr("foo", "bar"))
	c.Check(len(ml.Messages), Equals, 0)

	la.Dief(1234, "MessageF %v", 1)

	b.ShutdownLoggers()

	c.Assert(len(ml.Messages), Equals, 1)
	c.Check(ml.Messages[0], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LevelFatal,
		Message:   "MessageF 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	})
	c.Check(curTestExiter.exited, Equals, true)
	c.Check(curTestExiter.code, Equals, 1234)
}

func (s *GomolSuite) TestLogAdapterDiem(c *C) {
	b := NewBase()
	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	b.InitLoggers()

	la := b.NewLogAdapter(NewAttrs().SetAttr("foo", "bar"))
	c.Check(len(ml.Messages), Equals, 0)

	la.Diem(1234, NewAttrs().SetAttr("attr1", "val1"), "MessageM %v", 1)

	b.ShutdownLoggers()

	c.Assert(len(ml.Messages), Equals, 1)
	c.Check(ml.Messages[0], DeepEquals, &memMessage{
		Timestamp: clock().Now(),
		Level:     LevelFatal,
		Message:   "MessageM 1",
		Attrs: map[string]interface{}{
			"foo":   "bar",
			"attr1": "val1",
		},
	})
	c.Check(curTestExiter.exited, Equals, true)
	c.Check(curTestExiter.code, Equals, 1234)
}
