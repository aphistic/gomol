package gomol

import (
	. "gopkg.in/check.v1"
)

type testConsoleWriter struct {
	Output []string
}

func newTestConsoleWriter() *testConsoleWriter {
	return &testConsoleWriter{
		Output: make([]string, 0),
	}
}

func (w *testConsoleWriter) Print(msg string) {
	w.Output = append(w.Output, msg)
}

func (s *GomolSuite) TestTestConsoleWriter(c *C) {
	w := newTestConsoleWriter()
	c.Check(w.Output, NotNil)
	c.Check(w.Output, HasLen, 0)

	w.Print("print1")
	c.Check(w.Output, HasLen, 1)

	w.Print("print2")
	c.Check(w.Output, HasLen, 2)
}

// Issue-specific tests

func (s *GomolSuite) TestIssue5StringFormatting(c *C) {
	b := newBase()
	b.InitLoggers()

	cfg := NewConsoleLoggerConfig()
	cfg.Colorize = false
	l := NewConsoleLogger(cfg)
	w := newTestConsoleWriter()
	l.setWriter(w)
	b.AddLogger(l)

	b.Dbgf("msg %v%%", 100)

	b.ShutdownLoggers()

	c.Assert(w.Output, HasLen, 1)
	c.Check(w.Output[0], Equals, "[DEBUG] msg 100%\n")
}

// General tests

func (s *GomolSuite) TestConsoleInitLogger(c *C) {
	cl := NewConsoleLogger(nil)
	c.Check(cl.IsInitialized(), Equals, false)
	cl.InitLogger()
	c.Check(cl.IsInitialized(), Equals, true)
}

func (s *GomolSuite) TestConsoleShutdownLogger(c *C) {
	cl := NewConsoleLogger(nil)
	cl.InitLogger()
	c.Check(cl.IsInitialized(), Equals, true)
	cl.ShutdownLogger()
	c.Check(cl.IsInitialized(), Equals, false)
}

func (s *GomolSuite) TestConsoleDbg(c *C) {
	cfg := NewConsoleLoggerConfig()
	cfg.Colorize = false
	cl := NewConsoleLogger(cfg)
	w := newTestConsoleWriter()
	cl.setWriter(w)
	cl.Dbg("test")
	c.Assert(w.Output, HasLen, 1)
	c.Check(w.Output[0], Equals, "[DEBUG] test\n")
}

func (s *GomolSuite) TestConsoleDbgf(c *C) {
	cfg := NewConsoleLoggerConfig()
	cfg.Colorize = false
	cl := NewConsoleLogger(cfg)
	w := newTestConsoleWriter()
	cl.setWriter(w)
	cl.Dbgf("test %v", 1234)
	c.Assert(w.Output, HasLen, 1)
	c.Check(w.Output[0], Equals, "[DEBUG] test 1234\n")
}

func (s *GomolSuite) TestConsoleDbgm(c *C) {
	cfg := NewConsoleLoggerConfig()
	cfg.Colorize = false
	cl := NewConsoleLogger(cfg)
	w := newTestConsoleWriter()
	cl.setWriter(w)
	cl.Dbgm(
		map[string]interface{}{
			"attr1": 4321,
		},
		"test %v", 1234)
	c.Assert(w.Output, HasLen, 1)
	c.Check(w.Output[0], Equals, "[DEBUG] test 1234\n")
}

func (s *GomolSuite) TestConsoleInfo(c *C) {
	cfg := NewConsoleLoggerConfig()
	cfg.Colorize = false
	cl := NewConsoleLogger(cfg)
	w := newTestConsoleWriter()
	cl.setWriter(w)
	cl.Info("test")
	c.Assert(w.Output, HasLen, 1)
	c.Check(w.Output[0], Equals, "[INFO] test\n")
}

func (s *GomolSuite) TestConsoleInfof(c *C) {
	cfg := NewConsoleLoggerConfig()
	cfg.Colorize = false
	cl := NewConsoleLogger(cfg)
	w := newTestConsoleWriter()
	cl.setWriter(w)
	cl.Infof("test %v", 1234)
	c.Assert(w.Output, HasLen, 1)
	c.Check(w.Output[0], Equals, "[INFO] test 1234\n")
}

func (s *GomolSuite) TestConsoleInfom(c *C) {
	cfg := NewConsoleLoggerConfig()
	cfg.Colorize = false
	cl := NewConsoleLogger(cfg)
	w := newTestConsoleWriter()
	cl.setWriter(w)
	cl.Infom(
		map[string]interface{}{
			"attr1": 4321,
		},
		"test %v", 1234)
	c.Assert(w.Output, HasLen, 1)
	c.Check(w.Output[0], Equals, "[INFO] test 1234\n")
}

func (s *GomolSuite) TestConsoleWarn(c *C) {
	cfg := NewConsoleLoggerConfig()
	cfg.Colorize = false
	cl := NewConsoleLogger(cfg)
	w := newTestConsoleWriter()
	cl.setWriter(w)
	cl.Warn("test")
	c.Assert(w.Output, HasLen, 1)
	c.Check(w.Output[0], Equals, "[WARN] test\n")
}

func (s *GomolSuite) TestConsoleWarnf(c *C) {
	cfg := NewConsoleLoggerConfig()
	cfg.Colorize = false
	cl := NewConsoleLogger(cfg)
	w := newTestConsoleWriter()
	cl.setWriter(w)
	cl.Warnf("test %v", 1234)
	c.Assert(w.Output, HasLen, 1)
	c.Check(w.Output[0], Equals, "[WARN] test 1234\n")
}

func (s *GomolSuite) TestConsoleWarnm(c *C) {
	cfg := NewConsoleLoggerConfig()
	cfg.Colorize = false
	cl := NewConsoleLogger(cfg)
	w := newTestConsoleWriter()
	cl.setWriter(w)
	cl.Warnm(
		map[string]interface{}{
			"attr1": 4321,
		},
		"test %v", 1234)
	c.Assert(w.Output, HasLen, 1)
	c.Check(w.Output[0], Equals, "[WARN] test 1234\n")
}

func (s *GomolSuite) TestConsoleErr(c *C) {
	cfg := NewConsoleLoggerConfig()
	cfg.Colorize = false
	cl := NewConsoleLogger(cfg)
	w := newTestConsoleWriter()
	cl.setWriter(w)
	cl.Err("test")
	c.Assert(w.Output, HasLen, 1)
	c.Check(w.Output[0], Equals, "[ERROR] test\n")
}

func (s *GomolSuite) TestConsoleErrf(c *C) {
	cfg := NewConsoleLoggerConfig()
	cfg.Colorize = false
	cl := NewConsoleLogger(cfg)
	w := newTestConsoleWriter()
	cl.setWriter(w)
	cl.Errf("test %v", 1234)
	c.Assert(w.Output, HasLen, 1)
	c.Check(w.Output[0], Equals, "[ERROR] test 1234\n")
}

func (s *GomolSuite) TestConsoleErrm(c *C) {
	cfg := NewConsoleLoggerConfig()
	cfg.Colorize = false
	cl := NewConsoleLogger(cfg)
	w := newTestConsoleWriter()
	cl.setWriter(w)
	cl.Errm(
		map[string]interface{}{
			"attr1": 4321,
		},
		"test %v", 1234)
	c.Assert(w.Output, HasLen, 1)
	c.Check(w.Output[0], Equals, "[ERROR] test 1234\n")
}

func (s *GomolSuite) TestConsoleFatal(c *C) {
	cfg := NewConsoleLoggerConfig()
	cfg.Colorize = false
	cl := NewConsoleLogger(cfg)
	w := newTestConsoleWriter()
	cl.setWriter(w)
	cl.Fatal("test")
	c.Assert(w.Output, HasLen, 1)
	c.Check(w.Output[0], Equals, "[FATAL] test\n")
}

func (s *GomolSuite) TestConsoleFatalf(c *C) {
	cfg := NewConsoleLoggerConfig()
	cfg.Colorize = false
	cl := NewConsoleLogger(cfg)
	w := newTestConsoleWriter()
	cl.setWriter(w)
	cl.Fatalf("test %v", 1234)
	c.Assert(w.Output, HasLen, 1)
	c.Check(w.Output[0], Equals, "[FATAL] test 1234\n")
}

func (s *GomolSuite) TestConsoleFatalm(c *C) {
	cfg := NewConsoleLoggerConfig()
	cfg.Colorize = false
	cl := NewConsoleLogger(cfg)
	w := newTestConsoleWriter()
	cl.setWriter(w)
	cl.Fatalm(
		map[string]interface{}{
			"attr1": 4321,
		},
		"test %v", 1234)
	c.Assert(w.Output, HasLen, 1)
	c.Check(w.Output[0], Equals, "[FATAL] test 1234\n")
}

func (s *GomolSuite) TestConsoleBaseAttrs(c *C) {
	b := newBase()
	b.SetAttr("attr1", 7890)
	b.SetAttr("attr2", "val2")

	cfg := NewConsoleLoggerConfig()
	cfg.Colorize = false
	cl := NewConsoleLogger(cfg)
	w := newTestConsoleWriter()
	cl.setWriter(w)
	b.AddLogger(cl)
	cl.Dbgm(
		map[string]interface{}{
			"attr1": 4321,
			"attr3": "val3",
		},
		"test %v", 1234)
	c.Assert(w.Output, HasLen, 1)
	c.Check(w.Output[0], Equals, "[DEBUG] test 1234\n")
}
