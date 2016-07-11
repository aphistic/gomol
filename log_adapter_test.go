package gomol

import (
	. "gopkg.in/check.v1"
)

func (s *GomolSuite) TestNewLogAdapterEmpty(c *C) {
	b := newBase()

	la := b.NewLogAdapter(nil)
	c.Check(la, NotNil)
	c.Check(la.base, Equals, b)

	c.Check(la.attrs, NotNil)
	c.Check(len(la.attrs), Equals, 0)
}

func (s *GomolSuite) TestNewLogAdapter(c *C) {
	b := newBase()

	la := b.NewLogAdapter(map[string]interface{}{
		"testNum": 1234,
		"testStr": "foo",
	})
	c.Check(la, NotNil)
	c.Check(la.base, Equals, b)

	c.Check(la.attrs, NotNil)
	c.Assert(len(la.attrs), Equals, 2)
	c.Check(la.attrs["testNum"], Equals, 1234)
	c.Check(la.attrs["testStr"], Equals, "foo")
}

func (s *GomolSuite) TestLogAdapterSetAttr(c *C) {
	b := newBase()

	la := b.NewLogAdapter(nil)
	c.Check(len(la.attrs), Equals, 0)
	la.SetAttr("foo", "bar")
	c.Assert(len(la.attrs), Equals, 1)
	c.Check(la.attrs["foo"], Equals, "bar")
}

func (s *GomolSuite) TestLogAdapterRemoveAttr(c *C) {
	b := newBase()

	la := b.NewLogAdapter(map[string]interface{}{
		"foo": "bar",
	})
	c.Assert(len(la.attrs), Equals, 1)
	c.Check(la.attrs["foo"], Equals, "bar")
	la.RemoveAttr("foo")
	c.Check(len(la.attrs), Equals, 0)
}

func (s *GomolSuite) TestLogAdapterClearAttrs(c *C) {
	b := newBase()

	la := b.NewLogAdapter(map[string]interface{}{
		"foo": "bar",
		"baz": "qux",
	})
	c.Assert(len(la.attrs), Equals, 2)
	c.Check(la.attrs["foo"], Equals, "bar")
	c.Check(la.attrs["baz"], Equals, "qux")
	la.ClearAttrs()
	c.Check(len(la.attrs), Equals, 0)
}

func (s *GomolSuite) TestLogAdapterDebug(c *C) {
	b := newBase()
	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	b.InitLoggers()

	la := b.NewLogAdapter(map[string]interface{}{"foo": "bar"})
	c.Check(len(ml.Messages), Equals, 0)

	la.Dbg("Message 1")
	la.Debug("Message 2")
	la.Dbgf("MessageF %v", 1)
	la.Debugf("MessageF %v", 2)
	la.Dbgm(map[string]interface{}{
		"attr1": "val1",
	}, "MessageM %v", 1)
	la.Debugm(map[string]interface{}{
		"foo": "newBar",
	}, "MessageM %v", 2)

	b.ShutdownLoggers()

	c.Assert(len(ml.Messages), Equals, 6)
	c.Check(ml.Messages[0], DeepEquals, &memMessage{
		Level:   LEVEL_DEBUG,
		Message: "Message 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	})
	c.Check(ml.Messages[1], DeepEquals, &memMessage{
		Level:   LEVEL_DEBUG,
		Message: "Message 2",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	})
	c.Check(ml.Messages[2], DeepEquals, &memMessage{
		Level:   LEVEL_DEBUG,
		Message: "MessageF 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	})
	c.Check(ml.Messages[3], DeepEquals, &memMessage{
		Level:   LEVEL_DEBUG,
		Message: "MessageF 2",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	})
	c.Check(ml.Messages[4], DeepEquals, &memMessage{
		Level:   LEVEL_DEBUG,
		Message: "MessageM 1",
		Attrs: map[string]interface{}{
			"foo":   "bar",
			"attr1": "val1",
		},
	})
	c.Check(ml.Messages[5], DeepEquals, &memMessage{
		Level:   LEVEL_DEBUG,
		Message: "MessageM 2",
		Attrs: map[string]interface{}{
			"foo": "newBar",
		},
	})
}

func (s *GomolSuite) TestLogAdapterInfo(c *C) {
	b := newBase()
	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	b.InitLoggers()

	la := b.NewLogAdapter(map[string]interface{}{"foo": "bar"})
	c.Check(len(ml.Messages), Equals, 0)

	la.Info("Message 1")
	la.Infof("MessageF %v", 1)
	la.Infom(map[string]interface{}{
		"attr1": "val1",
	}, "MessageM %v", 1)

	b.ShutdownLoggers()

	c.Assert(len(ml.Messages), Equals, 3)
	c.Check(ml.Messages[0], DeepEquals, &memMessage{
		Level:   LEVEL_INFO,
		Message: "Message 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	})
	c.Check(ml.Messages[1], DeepEquals, &memMessage{
		Level:   LEVEL_INFO,
		Message: "MessageF 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	})
	c.Check(ml.Messages[2], DeepEquals, &memMessage{
		Level:   LEVEL_INFO,
		Message: "MessageM 1",
		Attrs: map[string]interface{}{
			"foo":   "bar",
			"attr1": "val1",
		},
	})
}

func (s *GomolSuite) TestLogAdapterWarn(c *C) {
	b := newBase()
	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	b.InitLoggers()

	la := b.NewLogAdapter(map[string]interface{}{"foo": "bar"})
	c.Check(len(ml.Messages), Equals, 0)

	la.Warn("Message 1")
	la.Warning("Message 2")
	la.Warnf("MessageF %v", 1)
	la.Warningf("MessageF %v", 2)
	la.Warnm(map[string]interface{}{
		"attr1": "val1",
	}, "MessageM %v", 1)
	la.Warningm(map[string]interface{}{
		"foo": "newBar",
	}, "MessageM %v", 2)

	b.ShutdownLoggers()

	c.Assert(len(ml.Messages), Equals, 6)
	c.Check(ml.Messages[0], DeepEquals, &memMessage{
		Level:   LEVEL_WARNING,
		Message: "Message 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	})
	c.Check(ml.Messages[1], DeepEquals, &memMessage{
		Level:   LEVEL_WARNING,
		Message: "Message 2",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	})
	c.Check(ml.Messages[2], DeepEquals, &memMessage{
		Level:   LEVEL_WARNING,
		Message: "MessageF 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	})
	c.Check(ml.Messages[3], DeepEquals, &memMessage{
		Level:   LEVEL_WARNING,
		Message: "MessageF 2",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	})
	c.Check(ml.Messages[4], DeepEquals, &memMessage{
		Level:   LEVEL_WARNING,
		Message: "MessageM 1",
		Attrs: map[string]interface{}{
			"foo":   "bar",
			"attr1": "val1",
		},
	})
	c.Check(ml.Messages[5], DeepEquals, &memMessage{
		Level:   LEVEL_WARNING,
		Message: "MessageM 2",
		Attrs: map[string]interface{}{
			"foo": "newBar",
		},
	})
}

func (s *GomolSuite) TestLogAdapterError(c *C) {
	b := newBase()
	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	b.InitLoggers()

	la := b.NewLogAdapter(map[string]interface{}{"foo": "bar"})
	c.Check(len(ml.Messages), Equals, 0)

	la.Err("Message 1")
	la.Error("Message 2")
	la.Errf("MessageF %v", 1)
	la.Errorf("MessageF %v", 2)
	la.Errm(map[string]interface{}{
		"attr1": "val1",
	}, "MessageM %v", 1)
	la.Errorm(map[string]interface{}{
		"foo": "newBar",
	}, "MessageM %v", 2)

	b.ShutdownLoggers()

	c.Assert(len(ml.Messages), Equals, 6)
	c.Check(ml.Messages[0], DeepEquals, &memMessage{
		Level:   LEVEL_ERROR,
		Message: "Message 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	})
	c.Check(ml.Messages[1], DeepEquals, &memMessage{
		Level:   LEVEL_ERROR,
		Message: "Message 2",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	})
	c.Check(ml.Messages[2], DeepEquals, &memMessage{
		Level:   LEVEL_ERROR,
		Message: "MessageF 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	})
	c.Check(ml.Messages[3], DeepEquals, &memMessage{
		Level:   LEVEL_ERROR,
		Message: "MessageF 2",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	})
	c.Check(ml.Messages[4], DeepEquals, &memMessage{
		Level:   LEVEL_ERROR,
		Message: "MessageM 1",
		Attrs: map[string]interface{}{
			"foo":   "bar",
			"attr1": "val1",
		},
	})
	c.Check(ml.Messages[5], DeepEquals, &memMessage{
		Level:   LEVEL_ERROR,
		Message: "MessageM 2",
		Attrs: map[string]interface{}{
			"foo": "newBar",
		},
	})
}

func (s *GomolSuite) TestLogAdapterFatal(c *C) {
	b := newBase()
	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	b.InitLoggers()

	la := b.NewLogAdapter(map[string]interface{}{"foo": "bar"})
	c.Check(len(ml.Messages), Equals, 0)

	la.Fatal("Message 1")
	la.Fatalf("MessageF %v", 1)
	la.Fatalm(map[string]interface{}{
		"attr1": "val1",
	}, "MessageM %v", 1)

	b.ShutdownLoggers()

	c.Assert(len(ml.Messages), Equals, 3)
	c.Check(ml.Messages[0], DeepEquals, &memMessage{
		Level:   LEVEL_FATAL,
		Message: "Message 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	})
	c.Check(ml.Messages[1], DeepEquals, &memMessage{
		Level:   LEVEL_FATAL,
		Message: "MessageF 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	})
	c.Check(ml.Messages[2], DeepEquals, &memMessage{
		Level:   LEVEL_FATAL,
		Message: "MessageM 1",
		Attrs: map[string]interface{}{
			"foo":   "bar",
			"attr1": "val1",
		},
	})
}

func (s *GomolSuite) TestLogAdapterDie(c *C) {
	b := newBase()
	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	b.InitLoggers()

	la := b.NewLogAdapter(map[string]interface{}{"foo": "bar"})
	c.Check(len(ml.Messages), Equals, 0)

	la.Die(1234, "Message 1")

	b.ShutdownLoggers()

	c.Assert(len(ml.Messages), Equals, 1)
	c.Check(ml.Messages[0], DeepEquals, &memMessage{
		Level:   LEVEL_FATAL,
		Message: "Message 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	})
	c.Check(curTestExiter.exited, Equals, true)
	c.Check(curTestExiter.code, Equals, 1234)
}

func (s *GomolSuite) TestLogAdapterDief(c *C) {
	b := newBase()
	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	b.InitLoggers()

	la := b.NewLogAdapter(map[string]interface{}{"foo": "bar"})
	c.Check(len(ml.Messages), Equals, 0)

	la.Dief(1234, "MessageF %v", 1)

	b.ShutdownLoggers()

	c.Assert(len(ml.Messages), Equals, 1)
	c.Check(ml.Messages[0], DeepEquals, &memMessage{
		Level:   LEVEL_FATAL,
		Message: "MessageF 1",
		Attrs: map[string]interface{}{
			"foo": "bar",
		},
	})
	c.Check(curTestExiter.exited, Equals, true)
	c.Check(curTestExiter.code, Equals, 1234)
}

func (s *GomolSuite) TestLogAdapterDiem(c *C) {
	b := newBase()
	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	b.InitLoggers()

	la := b.NewLogAdapter(map[string]interface{}{"foo": "bar"})
	c.Check(len(ml.Messages), Equals, 0)

	la.Diem(1234, map[string]interface{}{
		"attr1": "val1",
	}, "MessageM %v", 1)

	b.ShutdownLoggers()

	c.Assert(len(ml.Messages), Equals, 1)
	c.Check(ml.Messages[0], DeepEquals, &memMessage{
		Level:   LEVEL_FATAL,
		Message: "MessageM 1",
		Attrs: map[string]interface{}{
			"foo":   "bar",
			"attr1": "val1",
		},
	})
	c.Check(curTestExiter.exited, Equals, true)
	c.Check(curTestExiter.code, Equals, 1234)
}