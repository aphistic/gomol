package gomol

import (
	"bytes"
	. "gopkg.in/check.v1"
	"io/ioutil"
	"os"
	"strings"
)

func (s *GomolSuite) TestWriterSetTemplate(c *C) {
	var b bytes.Buffer
	wl, err := NewWriterLogger(&b, nil)
	c.Check(wl.tpl, NotNil)

	err = wl.SetTemplate(nil)
	c.Check(err, NotNil)

	tpl, err := NewTemplate("")
	c.Assert(err, IsNil)
	err = wl.SetTemplate(tpl)
	c.Check(err, IsNil)
}

func (s *GomolSuite) TestWriterInitLoggerNoWriter(c *C) {
	wl, err := NewWriterLogger(nil, nil)
	c.Assert(err, NotNil)
	c.Check(err.Error(), Equals, "An io.Writer must be provided")
	c.Check(wl, IsNil)
}

func (s *GomolSuite) TestWriterInitLoggerNoConfig(c *C) {
	var b bytes.Buffer
	wl, err := NewWriterLogger(&b, nil)
	c.Check(err, IsNil)
	c.Check(wl, NotNil)
	c.Assert(wl.config, NotNil)
	c.Check(wl.config.BufferSize, Equals, 1000)
}

func (s *GomolSuite) TestWriterInitLogger(c *C) {
	var b bytes.Buffer
	wl, err := NewWriterLogger(&b, nil)
	c.Assert(err, IsNil)
	c.Check(wl.IsInitialized(), Equals, false)
	wl.InitLogger()
	c.Check(wl.IsInitialized(), Equals, true)
}

func (s *GomolSuite) TestWriterShutdownLogger(c *C) {
	var b bytes.Buffer
	wl, err := NewWriterLogger(&b, nil)
	c.Assert(err, IsNil)
	c.Check(wl.IsInitialized(), Equals, false)
	wl.InitLogger()
	c.Check(wl.IsInitialized(), Equals, true)
	wl.ShutdownLogger()
	c.Check(wl.IsInitialized(), Equals, false)
}

func (s *GomolSuite) TestWriterWithConfig(c *C) {
	var b bytes.Buffer
	cfg := NewWriterLoggerConfig()
	cfg.BufferSize = 1
	wl, err := NewWriterLogger(&b, cfg)
	c.Assert(err, IsNil)
	c.Check(wl.config.BufferSize, Equals, 1)
}

func (s *GomolSuite) TestWriterMultipleMessages(c *C) {
	var b bytes.Buffer
	cfg := NewWriterLoggerConfig()
	wl, _ := NewWriterLogger(&b, cfg)
	wl.Dbgf("dbg 1234")
	wl.Warnf("warn 4321")

	wl.flushMessages()

	c.Check(b.String(), Equals, "[DEBUG] dbg 1234\n[WARN] warn 4321\n")
}

func (s *GomolSuite) TestWriterFlushOnBufferSize(c *C) {
	var b bytes.Buffer
	cfg := NewWriterLoggerConfig()
	cfg.BufferSize = 2
	wl, _ := NewWriterLogger(&b, cfg)

	c.Check(wl.buffer, HasLen, 0)

	wl.Dbg("Message 1")
	c.Check(wl.buffer, HasLen, 1)

	wl.Dbg("Message 2")
	c.Check(wl.buffer, HasLen, 0)

	c.Check(strings.Count(b.String(), "\n"), Equals, 2)
	c.Check(b.String(), Equals, "[DEBUG] Message 1\n[DEBUG] Message 2\n")
}

func (s *GomolSuite) TestWriterToFile(c *C) {
	f, err := ioutil.TempFile("", "gomol_test_")
	if err != nil {
		c.Fatal("Unable to create temp file to test writer logger")
	}
	defer f.Close()
	defer os.Remove(f.Name())

	cfg := NewWriterLoggerConfig()
	wl, _ := NewWriterLogger(f, cfg)
	wl.InitLogger()
	wl.Dbg("Message 1")
	wl.Fatal("Message 2")
	wl.ShutdownLogger()

	fData, err := ioutil.ReadFile(f.Name())
	if err != nil {
		c.Fatal("Could not read from writer logger test file")
	}
	fStr := string(fData)
	c.Check(fStr, Equals, "[DEBUG] Message 1\n[FATAL] Message 2\n")
}

func (s *GomolSuite) TestWriterDbg(c *C) {
	var b bytes.Buffer
	cfg := NewWriterLoggerConfig()
	wl, _ := NewWriterLogger(&b, cfg)
	wl.Dbg("test")

	wl.flushMessages()

	c.Check(b.String(), Equals, "[DEBUG] test\n")
}

func (s *GomolSuite) TestWriterDbgf(c *C) {
	var b bytes.Buffer
	cfg := NewWriterLoggerConfig()
	wl, _ := NewWriterLogger(&b, cfg)
	wl.Dbgf("test %v", 1234)

	wl.flushMessages()

	c.Check(b.String(), Equals, "[DEBUG] test 1234\n")
}

func (s *GomolSuite) TestWriterDbgm(c *C) {
	var b bytes.Buffer
	cfg := NewWriterLoggerConfig()
	wl, _ := NewWriterLogger(&b, cfg)
	wl.Dbgm(
		map[string]interface{}{
			"attr1": 4321,
		},
		"test %v", 1234)

	wl.flushMessages()

	c.Check(b.String(), Equals, "[DEBUG] test 1234\n")
}

func (s *GomolSuite) TestWriterInfo(c *C) {
	var b bytes.Buffer
	cfg := NewWriterLoggerConfig()
	wl, _ := NewWriterLogger(&b, cfg)
	wl.Info("test")

	wl.flushMessages()

	c.Check(b.String(), Equals, "[INFO] test\n")
}

func (s *GomolSuite) TestWriterInfof(c *C) {
	var b bytes.Buffer
	cfg := NewWriterLoggerConfig()
	wl, _ := NewWriterLogger(&b, cfg)
	wl.Infof("test %v", 1234)

	wl.flushMessages()

	c.Check(b.String(), Equals, "[INFO] test 1234\n")
}

func (s *GomolSuite) TestWriterInfom(c *C) {
	var b bytes.Buffer
	cfg := NewWriterLoggerConfig()
	wl, _ := NewWriterLogger(&b, cfg)
	wl.Infom(
		map[string]interface{}{
			"attr1": 4321,
		},
		"test %v", 1234)

	wl.flushMessages()

	c.Check(b.String(), Equals, "[INFO] test 1234\n")
}

func (s *GomolSuite) TestWriterWarn(c *C) {
	var b bytes.Buffer
	cfg := NewWriterLoggerConfig()
	wl, _ := NewWriterLogger(&b, cfg)
	wl.Warn("test")

	wl.flushMessages()

	c.Check(b.String(), Equals, "[WARN] test\n")
}

func (s *GomolSuite) TestWriterWarnf(c *C) {
	var b bytes.Buffer
	cfg := NewWriterLoggerConfig()
	wl, _ := NewWriterLogger(&b, cfg)
	wl.Warnf("test %v", 1234)

	wl.flushMessages()

	c.Check(b.String(), Equals, "[WARN] test 1234\n")
}

func (s *GomolSuite) TestWriterWarnm(c *C) {
	var b bytes.Buffer
	cfg := NewWriterLoggerConfig()
	wl, _ := NewWriterLogger(&b, cfg)
	wl.Warnm(
		map[string]interface{}{
			"attr1": 4321,
		},
		"test %v", 1234)

	wl.flushMessages()

	c.Check(b.String(), Equals, "[WARN] test 1234\n")
}

func (s *GomolSuite) TestWriterErr(c *C) {
	var b bytes.Buffer
	cfg := NewWriterLoggerConfig()
	wl, _ := NewWriterLogger(&b, cfg)
	wl.Err("test")

	wl.flushMessages()

	c.Check(b.String(), Equals, "[ERROR] test\n")
}

func (s *GomolSuite) TestWriterErrf(c *C) {
	var b bytes.Buffer
	cfg := NewWriterLoggerConfig()
	wl, _ := NewWriterLogger(&b, cfg)
	wl.Errf("test %v", 1234)

	wl.flushMessages()

	c.Check(b.String(), Equals, "[ERROR] test 1234\n")
}

func (s *GomolSuite) TestWriterErrm(c *C) {
	var b bytes.Buffer
	cfg := NewWriterLoggerConfig()
	wl, _ := NewWriterLogger(&b, cfg)
	wl.Errm(
		map[string]interface{}{
			"attr1": 4321,
		},
		"test %v", 1234)

	wl.flushMessages()

	c.Check(b.String(), Equals, "[ERROR] test 1234\n")
}

func (s *GomolSuite) TestWriterFatal(c *C) {
	var b bytes.Buffer
	cfg := NewWriterLoggerConfig()
	wl, _ := NewWriterLogger(&b, cfg)
	wl.Fatal("test")

	wl.flushMessages()

	c.Check(b.String(), Equals, "[FATAL] test\n")
}

func (s *GomolSuite) TestWriterFatalf(c *C) {
	var b bytes.Buffer
	cfg := NewWriterLoggerConfig()
	wl, _ := NewWriterLogger(&b, cfg)
	wl.Fatalf("test %v", 1234)

	wl.flushMessages()

	c.Check(b.String(), Equals, "[FATAL] test 1234\n")
}

func (s *GomolSuite) TestWriterFatalm(c *C) {
	var b bytes.Buffer
	cfg := NewWriterLoggerConfig()
	wl, _ := NewWriterLogger(&b, cfg)
	wl.Fatalm(
		map[string]interface{}{
			"attr1": 4321,
		},
		"test %v", 1234)

	wl.flushMessages()

	c.Check(b.String(), Equals, "[FATAL] test 1234\n")
}

func (s *GomolSuite) TestWriterBaseAttrs(c *C) {
	var buf bytes.Buffer
	b := newBase()
	b.SetAttr("attr1", 7890)
	b.SetAttr("attr2", "val2")

	cfg := NewWriterLoggerConfig()
	wl, _ := NewWriterLogger(&buf, cfg)
	b.AddLogger(wl)
	wl.Dbgm(
		map[string]interface{}{
			"attr1": 4321,
			"attr3": "val3",
		},
		"test %v", 1234)

	wl.flushMessages()

	c.Check(buf.String(), Equals, "[DEBUG] test 1234\n")
}
