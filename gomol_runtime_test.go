package gomol

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

/*
This is in its own file so the line numbers don't change.
These tests are testing calling locations so putting them in their own
file will limit the number of changes to that data.
*/

func (s *GomolSuite) TestIsGomolCaller(t sweet.T) {
	res, file := isGomolCaller("/home/gomoltest/some/sub/dir/that/is/long/filename.go")
	Expect(res).To(Equal(false))
	Expect(file).To(Equal("filename.go"))
}

func (s *GomolSuite) TestIsGomolCallerCached(t sweet.T) {
	Expect(len(gomolFiles)).To(Equal(0))

	res, file := isGomolCaller("/home/gomoltest/some/sub/dir/that/is/long/filename.go")
	Expect(len(gomolFiles)).To(Equal(1))
	Expect(res).To(Equal(false))
	Expect(file).To(Equal("filename.go"))

	res, file = isGomolCaller("/home/gomoltest/some/sub/dir/that/is/long/filename.go")
	Expect(gomolFiles).To(HaveLen(1))
	Expect(res).To(Equal(false))
	Expect(file).To(Equal("filename.go"))
}

func (s *GomolSuite) TestIsGomolCallerDirTooShort(t sweet.T) {
	res, file := isGomolCaller("1234/thiscanbesuperlong.go")
	Expect(len(gomolFiles)).To(Equal(1))
	Expect(res).To(Equal(false))
	Expect(file).To(Equal("thiscanbesuperlong.go"))
}

func (s *GomolSuite) TestIsGomolCallerFileShort(t sweet.T) {
	res, file := isGomolCaller("gomol/s.go")
	Expect(len(gomolFiles)).To(Equal(1))
	Expect(res).To(Equal(true))
	Expect(file).To(Equal("s.go"))
}

func (s *GomolSuite) TestIsGomolCallerFileTest(t sweet.T) {
	res, file := isGomolCaller("gomol/s_test.go")
	Expect(len(gomolFiles)).To(Equal(1))
	Expect(res).To(Equal(false))
	Expect(file).To(Equal("s_test.go"))
}

func (s *GomolSuite) TestLogWithRuntimeInfo(t sweet.T) {
	setFakeCallerInfo("fakefile.go", 1234)

	b := NewBase()
	b.config.FilenameAttr = "filename"
	b.config.LineNumberAttr = "line"

	l := newDefaultMemLogger()

	b.AddLogger(l)

	b.InitLoggers()
	b.Info("test")
	b.ShutdownLoggers()

	Expect(l.Messages).To(HaveLen(1))
	Expect(l.Messages[0].Message).To(Equal("test"))
	Expect(l.Messages[0].Attrs).To(HaveLen(2))
	Expect(l.Messages[0].Attrs["filename"]).To(Equal("fakefile.go"))
	Expect(l.Messages[0].Attrs["line"]).To(Equal(1234))
	Expect(l.Messages[0].Level).To(Equal(LevelInfo))
}
