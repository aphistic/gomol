package gomol

import (
	"testing"

	. "gopkg.in/check.v1"
)

/*
This is in its own file so the line numbers don't change.
These tests are testing calling locations so putting them in their own
file will limit the number of changes to that data.
*/

func (s *GomolSuite) TestLogFilenameAttr(c *C) {
	b := NewBase()
	b.config.FilenameAttr = "filename"

	l := newDefaultMemLogger()

	b.AddLogger(l)

	b.InitLoggers()
	func() {
		// Put it deeper in the calls just to test the "breaking out" of gomol
		b.Dbg("test")
	}()
	b.ShutdownLoggers()

	c.Assert(l.Messages, HasLen, 1)
	c.Check(l.Messages[0].Message, Equals, "test")
	c.Check(l.Messages[0].Attrs, HasLen, 1)
	c.Check(l.Messages[0].Attrs["filename"], Equals, "base_runtime_test.go")
	c.Check(l.Messages[0].Attrs["line"], IsNil)
}

func (s *GomolSuite) TestLogLineNumberAttr(c *C) {
	b := NewBase()
	b.config.LineNumberAttr = "line"

	l := newDefaultMemLogger()

	b.AddLogger(l)

	b.InitLoggers()
	func() {
		// Put it deeper in the calls just to test the "breaking out" of gomol
		b.Dbg("test")
	}()
	b.ShutdownLoggers()

	c.Assert(l.Messages, HasLen, 1)
	c.Check(l.Messages[0].Message, Equals, "test")
	c.Check(l.Messages[0].Attrs, HasLen, 1)
	c.Check(l.Messages[0].Attrs["filename"], IsNil)
	c.Check(l.Messages[0].Attrs["line"], Equals, 48)

}

func (s *GomolSuite) TestLogFilenameAndLineAttr(c *C) {
	b := NewBase()
	b.config.FilenameAttr = "filename"
	b.config.LineNumberAttr = "line"

	l := newDefaultMemLogger()

	b.AddLogger(l)

	b.InitLoggers()
	func() {
		// Put it deeper in the calls just to test the "breaking out" of gomol
		b.Dbg("test")
	}()
	b.ShutdownLoggers()

	c.Assert(l.Messages, HasLen, 1)
	c.Check(l.Messages[0].Message, Equals, "test")
	c.Check(l.Messages[0].Attrs, HasLen, 2)
	c.Check(l.Messages[0].Attrs["filename"], Equals, "base_runtime_test.go")
	c.Check(l.Messages[0].Attrs["line"], Equals, 72)
}

func (s *GomolSuite) BenchmarkLogInsertionWithFilename(c *C) {
	base := NewBase()
	base.config.FilenameAttr = "filename"
	base.InitLoggers()
	for i := 0; i < c.N; i++ {
		base.log(LevelDebug, map[string]interface{}{
			"attr1": "val1",
			"attr2": "val2",
			"attr3": 3,
			"attr4": 4,
		}, "msg %v %v", "string", 1234)
	}
	base.ShutdownLoggers()
}
func (s *GomolSuite) BenchmarkLogInsertionWithLineNumber(c *C) {
	base := NewBase()
	base.config.LineNumberAttr = "line"
	base.InitLoggers()
	for i := 0; i < c.N; i++ {
		base.log(LevelDebug, map[string]interface{}{
			"attr1": "val1",
			"attr2": "val2",
			"attr3": 3,
			"attr4": 4,
		}, "msg %v %v", "string", 1234)
	}
	base.ShutdownLoggers()
}
func (s *GomolSuite) BenchmarkLogInsertionWithFilenameAndLineNumber(c *C) {
	base := NewBase()
	base.config.FilenameAttr = "filename"
	base.config.LineNumberAttr = "line"
	base.InitLoggers()
	for i := 0; i < c.N; i++ {
		base.log(LevelDebug, map[string]interface{}{
			"attr1": "val1",
			"attr2": "val2",
			"attr3": 3,
			"attr4": 4,
		}, "msg %v %v", "string", 1234)
	}
	base.ShutdownLoggers()
}

func BenchmarkIsGomolCaller(b *testing.B) {
	for i := 0; i < b.N; i++ {
		isGomolCaller("/home/gomoltest/some/sub/dir/that/is/long/filename.go")
	}
}

func BenchmarkGetCallerInfo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getCallerInfo()
	}
}
