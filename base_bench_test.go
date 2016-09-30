package gomol

import . "gopkg.in/check.v1"

func (s *GomolSuite) BenchmarkBasicLogInsertion(c *C) {
	b := NewBase()
	b.SetAttr("attr1", 1234)

	l := newDefaultMemLogger()
	b.AddLogger(l)
	b.InitLoggers()

	for idx := 0; idx < c.N; idx++ {
		b.log(
			LevelDebug,
			NewAttrs().
				SetAttr("attr2", 4321).
				SetAttr("attr3", "val3").
				SetAttr("attr4", 1234),
			"test %v",
			1234)
	}
	b.ShutdownLoggers()
}

func (s *GomolSuite) BenchmarkBaseDbgm(c *C) {
	b := NewBase()
	b.SetAttr("attr1", 1234)

	l := newDefaultMemLogger()
	b.AddLogger(l)
	b.InitLoggers()

	for idx := 0; idx < c.N; idx++ {
		b.Dbgm(
			NewAttrs().
				SetAttr("attr2", 4321).
				SetAttr("attr3", "val3").
				SetAttr("attr4", 1234),
			"test %v",
			1234)
	}
	b.ShutdownLoggers()
}

func (s *GomolSuite) BenchmarkBaseStaticStringNoLogger(c *C) {
	b := NewBase()

	b.InitLoggers()

	c.ResetTimer()
	for idx := 0; idx < c.N; idx++ {
		b.Info("Hi I'm a log!")
	}
	b.ShutdownLoggers()
}

func (s *GomolSuite) BenchmarkBaseStaticStringMemLogger(c *C) {
	b := NewBase()

	l := newDefaultMemLogger()
	b.AddLogger(l)
	b.InitLoggers()

	c.ResetTimer()
	for idx := 0; idx < c.N; idx++ {
		b.Info("Hi I'm a log!")
	}
	b.ShutdownLoggers()
}

func (s *GomolSuite) BenchmarkLogInsertionWithFilename(c *C) {
	base := NewBase()
	base.config.FilenameAttr = "filename"
	base.InitLoggers()
	for i := 0; i < c.N; i++ {
		base.log(LevelDebug,
			NewAttrs().
				SetAttr("attr1", "val1").
				SetAttr("attr2", "val2").
				SetAttr("attr3", 3).
				SetAttr("attr4", 4),
			"msg %v %v", "string", 1234)
	}
	base.ShutdownLoggers()
}
func (s *GomolSuite) BenchmarkLogInsertionWithLineNumber(c *C) {
	base := NewBase()
	base.config.LineNumberAttr = "line"
	base.InitLoggers()
	for i := 0; i < c.N; i++ {
		base.log(LevelDebug,
			NewAttrs().
				SetAttr("attr1", "val1").
				SetAttr("attr2", "val2").
				SetAttr("attr3", 3).
				SetAttr("attr4", 4),
			"msg %v %v", "string", 1234)
	}
	base.ShutdownLoggers()
}
func (s *GomolSuite) BenchmarkLogInsertionWithFilenameAndLineNumber(c *C) {
	base := NewBase()
	base.config.FilenameAttr = "filename"
	base.config.LineNumberAttr = "line"
	base.InitLoggers()
	for i := 0; i < c.N; i++ {
		base.log(LevelDebug,
			NewAttrs().
				SetAttr("attr1", "val1").
				SetAttr("attr2", "val2").
				SetAttr("attr3", 3).
				SetAttr("attr4", 4),
			"msg %v %v", "string", 1234)
	}
	base.ShutdownLoggers()
}

func (s *GomolSuite) BenchmarkIsGomolCaller(c *C) {
	for i := 0; i < c.N; i++ {
		isGomolCaller("/home/gomoltest/some/sub/dir/that/is/long/filename.go")
	}
}

func (s *GomolSuite) BenchmarkGetCallerInfo(c *C) {
	for i := 0; i < c.N; i++ {
		getCallerInfo()
	}
}
