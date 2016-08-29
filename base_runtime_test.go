package gomol

import . "gopkg.in/check.v1"

/*
This is in its own file so the line numbers don't change.
These tests are testing calling locations so putting them in their own
file will limit the number of changes to that data.
*/

func (s *GomolSuite) TestIsGomolCaller(c *C) {
	res, file := isGomolCaller("/home/gomoltest/some/sub/dir/that/is/long/filename.go")
	c.Check(res, Equals, false)
	c.Check(file, Equals, "filename.go")
}

func (s *GomolSuite) TestIsGomolCallerCached(c *C) {
	c.Check(len(gomolFiles), Equals, 0)

	res, file := isGomolCaller("/home/gomoltest/some/sub/dir/that/is/long/filename.go")
	c.Check(len(gomolFiles), Equals, 1)
	c.Check(res, Equals, false)
	c.Check(file, Equals, "filename.go")

	res, file = isGomolCaller("/home/gomoltest/some/sub/dir/that/is/long/filename.go")
	c.Check(gomolFiles, HasLen, 1)
	c.Check(res, Equals, false)
	c.Check(file, Equals, "filename.go")
}

func (s *GomolSuite) TestIsGomolCallerDirTooShort(c *C) {
	res, file := isGomolCaller("1234/thiscanbesuperlong.go")
	c.Check(len(gomolFiles), Equals, 1)
	c.Check(res, Equals, false)
	c.Check(file, Equals, "thiscanbesuperlong.go")
}

func (s *GomolSuite) TestIsGomolCallerFileShort(c *C) {
	res, file := isGomolCaller("gomol/s.go")
	c.Check(len(gomolFiles), Equals, 1)
	c.Check(res, Equals, true)
	c.Check(file, Equals, "s.go")
}

func (s *GomolSuite) TestIsGomolCallerFileTest(c *C) {
	res, file := isGomolCaller("gomol/s_test.go")
	c.Check(len(gomolFiles), Equals, 1)
	c.Check(res, Equals, false)
	c.Check(file, Equals, "s_test.go")
}

func (s *GomolSuite) TestLogWithRuntimeInfo(c *C) {
	setFakeCallerInfo("fakefile.go", 1234)

	b := NewBase()
	b.config.FilenameAttr = "filename"
	b.config.LineNumberAttr = "line"

	l := newDefaultMemLogger()

	b.AddLogger(l)

	b.InitLoggers()
	b.Info("test")
	b.ShutdownLoggers()

	c.Assert(l.Messages, HasLen, 1)
	c.Check(l.Messages[0].Message, Equals, "test")
	c.Check(l.Messages[0].Attrs, HasLen, 2)
	c.Check(l.Messages[0].Attrs["filename"], Equals, "fakefile.go")
	c.Check(l.Messages[0].Attrs["line"], Equals, 1234)
	c.Check(l.Messages[0].Level, Equals, LevelInfo)
}
