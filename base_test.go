package gomol

import (
	"testing"
	"time"

	. "github.com/onsi/gomega"

	"github.com/aphistic/sweet"
	"github.com/aphistic/sweet-junit"
)

type GomolSuite struct{}

func Test(t *testing.T) {
	RegisterFailHandler(sweet.GomegaFail)

	sweet.T(func(s *sweet.S) {
		s.RegisterPlugin(junit.NewPlugin())

		s.RunSuite(t, &GomolSuite{})
	})
}

var testBase *Base

type testExiter struct {
	exited bool
	code   int
}

func (exiter *testExiter) Exit(code int) {
	exiter.code = code
	exiter.exited = true
}

var curTestExiter *testExiter

func (s *GomolSuite) SetUpTest(t *testing.T) {
	setFakeCallerInfo("", 0)
	setClock(newTestClock(time.Now()))
	gomolFiles = map[string]fileRecord{}

	curTestExiter = &testExiter{}
	setExiter(curTestExiter)

	testBase = NewBase()
	testBase.AddLogger(newDefaultMemLogger())
	testBase.InitLoggers()

	curDefault = NewBase()
	curDefault.AddLogger(newDefaultMemLogger())
	curDefault.InitLoggers()
}

func (s *GomolSuite) TearDownTest(t *testing.T) {
	curDefault.ShutdownLoggers()

	testBase.ShutdownLoggers()
}

func (s *GomolSuite) TestShouldLog(t *testing.T) {
	b := NewBase()
	b.SetLogLevel(LevelInfo)
	Expect(b.shouldLog(LevelDebug)).To(Equal(false))
	Expect(b.shouldLog(LevelInfo)).To(Equal(true))
	Expect(b.shouldLog(LevelWarning)).To(Equal(true))
	Expect(b.shouldLog(LevelError)).To(Equal(true))
	Expect(b.shouldLog(LevelFatal)).To(Equal(true))

	b.SetLogLevel(LevelFatal)
	Expect(b.shouldLog(LevelDebug)).To(Equal(false))
	Expect(b.shouldLog(LevelInfo)).To(Equal(false))
	Expect(b.shouldLog(LevelWarning)).To(Equal(false))
	Expect(b.shouldLog(LevelError)).To(Equal(false))
	Expect(b.shouldLog(LevelFatal)).To(Equal(true))

	b.SetLogLevel(LevelNone)
	Expect(b.shouldLog(LevelDebug)).To(Equal(false))
	Expect(b.shouldLog(LevelInfo)).To(Equal(false))
	Expect(b.shouldLog(LevelWarning)).To(Equal(false))
	Expect(b.shouldLog(LevelError)).To(Equal(false))
	Expect(b.shouldLog(LevelFatal)).To(Equal(false))
}

func (s *GomolSuite) TestNewBase(t *testing.T) {
	b := NewBase()
	Expect(b.isInitialized).To(Equal(false))
	Expect(b.config).ToNot(BeNil())
	Expect(b.config.FilenameAttr).To(Equal(""))
	Expect(b.config.LineNumberAttr).To(Equal(""))
	Expect(b.queue).ToNot(BeNil())
	Expect(b.logLevel).To(Equal(LevelDebug))
	Expect(b.loggers).To(HaveLen(0))
	Expect(b.BaseAttrs.Attrs()).To(HaveLen(0))
}

func (s *GomolSuite) TestSetConfig(t *testing.T) {
	b := NewBase()

	Expect(b.config).ToNot(BeNil())
	Expect(b.config.FilenameAttr).To(Equal(""))
	Expect(b.config.LineNumberAttr).To(Equal(""))

	cfg := NewConfig()
	cfg.FilenameAttr = "filename"
	cfg.LineNumberAttr = "line_number"

	b.SetConfig(cfg)
	Expect(b.config).ToNot(BeNil())
	Expect(b.config.FilenameAttr).To(Equal("filename"))
	Expect(b.config.LineNumberAttr).To(Equal("line_number"))
}

func (s *GomolSuite) TestSetLogLevel(t *testing.T) {
	b := NewBase()
	b.InitLoggers()
	ml := newDefaultMemLogger()
	b.AddLogger(ml)

	b.SetLogLevel(LevelWarning)
	b.Dbg("test")
	b.Info("test")
	b.Warn("test")
	b.Err("test")
	b.Fatal("test")
	b.ShutdownLoggers()
	Expect(ml.Messages).To(HaveLen(3))
}

func (s *GomolSuite) TestAddLogger(t *testing.T) {
	b := NewBase()
	b.InitLoggers()
	Expect(b.loggers).To(HaveLen(0))

	ml := newDefaultMemLogger()
	Expect(ml.IsInitialized()).To(Equal(false))
	Expect(ml.base).To(BeNil())

	b.AddLogger(ml)
	Expect(b.IsInitialized()).To(Equal(true))
	Expect(b.loggers).To(HaveLen(1))
	Expect(b.loggers[0].IsInitialized()).To(Equal(true))
	Expect(ml.base).To(Equal(b))
}

func (s *GomolSuite) TestAddLoggerAfterInit(t *testing.T) {
	b := NewBase()
	b.InitLoggers()

	ml := newDefaultMemLogger()
	Expect(ml.IsInitialized()).To(Equal(false))

	ret := b.AddLogger(ml)
	Expect(ret).To(BeNil())
	Expect(ml.IsInitialized()).To(Equal(true))
}

func (s *GomolSuite) TestAddLoggerAfterShutdown(t *testing.T) {
	b := NewBase()

	ml := newDefaultMemLogger()
	Expect(ml.IsInitialized()).To(Equal(false))
	ml.InitLogger()
	Expect(ml.IsInitialized()).To(Equal(true))

	ret := b.AddLogger(ml)
	Expect(ret).To(BeNil())
	Expect(ml.IsInitialized()).To(Equal(false))
}

func (s *GomolSuite) TestAddLoggerAfterInitFail(t *testing.T) {
	b := NewBase()
	b.InitLoggers()

	mlCfg := newMemLoggerConfig()
	mlCfg.FailInit = true
	ml, err := newMemLogger(mlCfg)
	Expect(err).To(BeNil())
	Expect(ml.IsInitialized()).To(Equal(false))

	ret := b.AddLogger(ml)
	Expect(ret).ToNot(BeNil())
	Expect(ret.Error()).To(Equal("Init failed"))
	Expect(ml.IsInitialized()).To(Equal(false))
	Expect(b.loggers).To(HaveLen(0))
}

func (s *GomolSuite) TestAddLoggerAfterShutdownFail(t *testing.T) {
	b := NewBase()

	mlCfg := newMemLoggerConfig()
	mlCfg.FailShutdown = true
	ml, err := newMemLogger(mlCfg)
	Expect(err).To(BeNil())
	Expect(ml.IsInitialized()).To(Equal(false))
	ml.InitLogger()
	Expect(ml.IsInitialized()).To(Equal(true))

	ret := b.AddLogger(ml)
	Expect(ret).ToNot(BeNil())
	Expect(ret.Error()).To(Equal("Shutdown failed"))
	Expect(ml.IsInitialized()).To(Equal(true))
	Expect(b.loggers).To(HaveLen(0))
}

func (s *GomolSuite) TestBaseRemoveLogger(t *testing.T) {
	b := NewBase()

	ml1 := newDefaultMemLogger()
	ml2 := newDefaultMemLogger()
	ml3 := newDefaultMemLogger()
	b.AddLogger(ml1)
	b.AddLogger(ml2)
	b.AddLogger(ml3)

	b.InitLoggers()

	Expect(ml1.IsInitialized()).To(Equal(true))
	Expect(ml2.IsInitialized()).To(Equal(true))
	Expect(ml3.IsInitialized()).To(Equal(true))
	Expect(b.loggers).To(HaveLen(3))

	err := b.RemoveLogger(ml2)
	Expect(err).To(BeNil())
	Expect(ml1.IsInitialized()).To(Equal(true))
	Expect(ml2.IsInitialized()).To(Equal(false))
	Expect(ml3.IsInitialized()).To(Equal(true))
	Expect(b.loggers).To(HaveLen(2))
}

func (s *GomolSuite) TestBaseRemoveLoggerNonExistent(t *testing.T) {
	b := NewBase()

	ml1 := newDefaultMemLogger()
	ml2 := newDefaultMemLogger()
	b.AddLogger(ml1)

	b.InitLoggers()

	Expect(ml1.IsInitialized()).To(Equal(true))
	Expect(b.loggers).To(HaveLen(1))

	err := b.RemoveLogger(ml2)
	Expect(err).To(BeNil())
}

func (s *GomolSuite) TestBaseClearLoggers(t *testing.T) {
	b := NewBase()

	ml1 := newDefaultMemLogger()
	ml2 := newDefaultMemLogger()
	ml3 := newDefaultMemLogger()
	b.AddLogger(ml1)
	b.AddLogger(ml2)
	b.AddLogger(ml3)

	b.InitLoggers()

	Expect(ml1.IsInitialized()).To(Equal(true))
	Expect(ml2.IsInitialized()).To(Equal(true))
	Expect(ml3.IsInitialized()).To(Equal(true))
	Expect(b.loggers).To(HaveLen(3))

	err := b.ClearLoggers()
	Expect(err).To(BeNil())
	Expect(ml1.IsInitialized()).To(Equal(false))
	Expect(ml2.IsInitialized()).To(Equal(false))
	Expect(ml3.IsInitialized()).To(Equal(false))
	Expect(b.loggers).To(HaveLen(0))
}

func (s *GomolSuite) TestInitLoggers(t *testing.T) {
	b := NewBase()
	Expect(b.IsInitialized()).To(Equal(false))

	ml1 := newDefaultMemLogger()
	ml2 := newDefaultMemLogger()

	b.AddLogger(ml1)
	b.AddLogger(ml2)

	b.InitLoggers()

	Expect(b.IsInitialized()).To(Equal(true))
	Expect(ml1.IsInitialized()).To(Equal(true))
	Expect(ml2.IsInitialized()).To(Equal(true))
}

func (s *GomolSuite) TestInitLoggersTwice(t *testing.T) {
	b := NewBase()
	Expect(b.IsInitialized()).To(Equal(false))

	ml1 := newDefaultMemLogger()
	ml2 := newDefaultMemLogger()

	b.AddLogger(ml1)
	b.AddLogger(ml2)

	b.InitLoggers()
	b.InitLoggers()

	Expect(b.IsInitialized()).To(Equal(true))
	Expect(ml1.IsInitialized()).To(Equal(true))
	Expect(ml2.IsInitialized()).To(Equal(true))
}

func (s *GomolSuite) TestInitLoggersFail(t *testing.T) {
	b := NewBase()

	mlCfg := newMemLoggerConfig()
	mlCfg.FailInit = true
	ml1, err := newMemLogger(mlCfg)
	Expect(err).To(BeNil())
	ml2, err := newMemLogger(mlCfg)
	Expect(err).To(BeNil())

	b.AddLogger(ml1)
	b.AddLogger(ml2)

	err = b.InitLoggers()
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Init failed"))

	Expect(b.IsInitialized()).To(Equal(false))
	Expect(ml1.IsInitialized()).To(Equal(false))
	Expect(ml2.IsInitialized()).To(Equal(false))
}

func (s *GomolSuite) TestShutdownLoggers(t *testing.T) {
	b := NewBase()

	ml1 := newDefaultMemLogger()
	ml2 := newDefaultMemLogger()

	b.AddLogger(ml1)
	b.AddLogger(ml2)

	b.InitLoggers()
	b.ShutdownLoggers()

	Expect(ml1.isShutdown).To(Equal(true))
	Expect(ml2.isShutdown).To(Equal(true))
}

func (s *GomolSuite) TestShutdownLoggersFail(t *testing.T) {
	b := NewBase()

	mlCfg := newMemLoggerConfig()
	mlCfg.FailShutdown = true
	ml1, err := newMemLogger(mlCfg)
	Expect(err).To(BeNil())
	ml2, err := newMemLogger(mlCfg)
	Expect(err).To(BeNil())

	b.AddLogger(ml1)
	b.AddLogger(ml2)

	b.InitLoggers()
	err = b.ShutdownLoggers()
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Shutdown failed"))

	Expect(ml1.isShutdown).To(Equal(false))
	Expect(ml2.isShutdown).To(Equal(false))
}

func (s *GomolSuite) TestShutdownLoggersTwice(t *testing.T) {
	b := NewBase()

	ml1 := newDefaultMemLogger()
	ml2 := newDefaultMemLogger()

	b.AddLogger(ml1)
	b.AddLogger(ml2)

	b.InitLoggers()
	b.ShutdownLoggers()
	b.ShutdownLoggers()

	Expect(ml1.isShutdown).To(Equal(true))
	Expect(ml2.isShutdown).To(Equal(true))
}

func (s *GomolSuite) TestSetAttr(t *testing.T) {
	b := NewBase()

	b.SetAttr("attr1", 1)
	Expect(b.BaseAttrs.Attrs()).To(HaveLen(1))
	Expect(b.BaseAttrs.GetAttr("attr1")).To(Equal(1))
	b.SetAttr("attr2", "val2")
	Expect(b.BaseAttrs.Attrs()).To(HaveLen(2))
	Expect(b.BaseAttrs.GetAttr("attr2")).To(Equal("val2"))
}

func (s *GomolSuite) TestGetAttr(t *testing.T) {
	b := NewBase()

	b.SetAttr("attr1", 1)
	b.SetAttr("attr2", "val2")

	Expect(b.GetAttr("attr2")).To(Equal("val2"))
	Expect(b.GetAttr("notakey")).To(BeNil())
}

func (s *GomolSuite) TestRemoveAttr(t *testing.T) {
	b := NewBase()

	b.SetAttr("attr1", 1)
	Expect(b.BaseAttrs.Attrs()).To(HaveLen(1))
	Expect(b.BaseAttrs.GetAttr("attr1")).To(Equal(1))

	b.RemoveAttr("attr1")
	Expect(b.BaseAttrs.Attrs()).To(HaveLen(0))
}

func (s *GomolSuite) TestClearAttrs(t *testing.T) {
	b := NewBase()

	b.SetAttr("attr1", 1)
	b.SetAttr("attr2", "val2")
	Expect(b.BaseAttrs.Attrs()).To(HaveLen(2))

	b.ClearAttrs()
	Expect(b.BaseAttrs.Attrs()).To(HaveLen(0))
}

func (s *GomolSuite) TestSequenceDisabled(t *testing.T) {
	b := NewBase()

	b.InitLoggers()

	Expect(b.sequence).To(Equal(uint64(0)))
	b.Dbg("test")
	Expect(b.sequence).To(Equal(uint64(0)))
	b.Dbg("test")
	Expect(b.sequence).To(Equal(uint64(0)))

	b.ShutdownLoggers()
}

func (s *GomolSuite) TestSequence(t *testing.T) {
	b := NewBase()
	b.config.SequenceAttr = "seq"

	l := newDefaultMemLogger()
	b.AddLogger(l)

	b.InitLoggers()

	Expect(b.sequence).To(Equal(uint64(0)))
	b.Dbg("test")
	Expect(b.sequence).To(Equal(uint64(1)))
	b.Dbg("test")
	Expect(b.sequence).To(Equal(uint64(2)))

	b.ShutdownLoggers()

	Expect(l.Messages).To(HaveLen(2))
	Expect(l.Messages[0].Message).To(Equal("test"))
	Expect(l.Messages[0].Attrs).To(HaveLen(1))
	Expect(l.Messages[0].Attrs["seq"]).To(Equal(uint64(1)))
	Expect(l.Messages[0].Level).To(Equal(LevelDebug))
	Expect(l.Messages[1].Message).To(Equal("test"))
	Expect(l.Messages[1].Attrs).To(HaveLen(1))
	Expect(l.Messages[1].Attrs["seq"]).To(Equal(uint64(2)))
	Expect(l.Messages[1].Level).To(Equal(LevelDebug))
}

// Base func tests
func (s *GomolSuite) TestBaseDbgfWithFormattingParams(t *testing.T) {
	b := NewBase()

	l1 := newDefaultMemLogger()

	b.AddLogger(l1)

	b.InitLoggers()
	b.Dbgf("LOG %s", "%2b")
	b.ShutdownLoggers()

	Expect(l1.Messages).To(HaveLen(1))
	Expect(l1.Messages[0].Message).To(Equal("LOG %2b"))
	Expect(l1.Messages[0].Attrs).To(HaveLen(0))
	Expect(l1.Messages[0].Level).To(Equal(LevelDebug))
}

func (s *GomolSuite) TestBaseDbg(t *testing.T) {
	b := NewBase()

	l1 := newDefaultMemLogger()
	l2 := newDefaultMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Dbg("test")
	b.ShutdownLoggers()

	Expect(l1.Messages).To(HaveLen(1))
	Expect(l1.Messages[0].Message).To(Equal("test"))
	Expect(l1.Messages[0].Attrs).To(HaveLen(0))
	Expect(l1.Messages[0].Level).To(Equal(LevelDebug))

	Expect(l2.Messages).To(HaveLen(1))
	Expect(l2.Messages[0].Message).To(Equal("test"))
	Expect(l2.Messages[0].Attrs).To(HaveLen(0))
	Expect(l2.Messages[0].Level).To(Equal(LevelDebug))
}

func (s *GomolSuite) TestBaseDbgf(t *testing.T) {
	b := NewBase()

	l1 := newDefaultMemLogger()
	l2 := newDefaultMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Dbgf("test %v", 1234)
	b.ShutdownLoggers()

	Expect(l1.Messages).To(HaveLen(1))
	Expect(l1.Messages[0].Message).To(Equal("test 1234"))
	Expect(l1.Messages[0].Attrs).To(HaveLen(0))
	Expect(l1.Messages[0].Level).To(Equal(LevelDebug))

	Expect(l2.Messages).To(HaveLen(1))
	Expect(l2.Messages[0].Message).To(Equal("test 1234"))
	Expect(l2.Messages[0].Attrs).To(HaveLen(0))
	Expect(l2.Messages[0].Level).To(Equal(LevelDebug))
}

func (s *GomolSuite) TestBaseDbgm(t *testing.T) {
	b := NewBase()
	b.SetAttr("attr1", 1234)

	l1 := newDefaultMemLogger()
	l2 := newDefaultMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Dbgm(
		NewAttrs().
			SetAttr("attr2", 4321).
			SetAttr("attr3", "val3"),
		"test %v",
		1234)
	b.ShutdownLoggers()

	Expect(l1.Messages).To(HaveLen(1))
	Expect(l1.Messages[0].Message).To(Equal("test 1234"))
	Expect(l1.Messages[0].Attrs).To(HaveLen(3))
	Expect(l1.Messages[0].Attrs["attr1"]).To(Equal(1234))
	Expect(l1.Messages[0].Attrs["attr2"]).To(Equal(4321))
	Expect(l1.Messages[0].Attrs["attr3"]).To(Equal("val3"))
	Expect(l1.Messages[0].Level).To(Equal(LevelDebug))

	Expect(l2.Messages).To(HaveLen(1))
	Expect(l2.Messages[0].Message).To(Equal("test 1234"))
	Expect(l2.Messages[0].Attrs).To(HaveLen(3))
	Expect(l2.Messages[0].Attrs["attr1"]).To(Equal(1234))
	Expect(l2.Messages[0].Attrs["attr2"]).To(Equal(4321))
	Expect(l2.Messages[0].Attrs["attr3"]).To(Equal("val3"))
	Expect(l2.Messages[0].Level).To(Equal(LevelDebug))
}

func (s *GomolSuite) TestBaseInfo(t *testing.T) {
	b := NewBase()

	l1 := newDefaultMemLogger()
	l2 := newDefaultMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Info("test")
	b.ShutdownLoggers()

	Expect(l1.Messages).To(HaveLen(1))
	Expect(l1.Messages[0].Message).To(Equal("test"))
	Expect(l1.Messages[0].Attrs).To(HaveLen(0))
	Expect(l1.Messages[0].Level).To(Equal(LevelInfo))

	Expect(l2.Messages).To(HaveLen(1))
	Expect(l2.Messages[0].Message).To(Equal("test"))
	Expect(l2.Messages[0].Attrs).To(HaveLen(0))
	Expect(l2.Messages[0].Level).To(Equal(LevelInfo))
}

func (s *GomolSuite) TestBaseInfof(t *testing.T) {
	b := NewBase()

	l1 := newDefaultMemLogger()
	l2 := newDefaultMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Infof("test %v", 1234)
	b.ShutdownLoggers()

	Expect(l1.Messages).To(HaveLen(1))
	Expect(l1.Messages[0].Message).To(Equal("test 1234"))
	Expect(l1.Messages[0].Attrs).To(HaveLen(0))
	Expect(l1.Messages[0].Level).To(Equal(LevelInfo))

	Expect(l2.Messages).To(HaveLen(1))
	Expect(l2.Messages[0].Message).To(Equal("test 1234"))
	Expect(l2.Messages[0].Attrs).To(HaveLen(0))
	Expect(l2.Messages[0].Level).To(Equal(LevelInfo))
}

func (s *GomolSuite) TestBaseInfom(t *testing.T) {
	b := NewBase()
	b.SetAttr("attr1", 1234)

	l1 := newDefaultMemLogger()
	l2 := newDefaultMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Infom(
		NewAttrs().
			SetAttr("attr2", 4321).
			SetAttr("attr3", "val3"),
		"test %v",
		1234)
	b.ShutdownLoggers()

	Expect(l1.Messages).To(HaveLen(1))
	Expect(l1.Messages[0].Message).To(Equal("test 1234"))
	Expect(l1.Messages[0].Attrs).To(HaveLen(3))
	Expect(l1.Messages[0].Attrs["attr1"]).To(Equal(1234))
	Expect(l1.Messages[0].Attrs["attr2"]).To(Equal(4321))
	Expect(l1.Messages[0].Attrs["attr3"]).To(Equal("val3"))
	Expect(l1.Messages[0].Level).To(Equal(LevelInfo))

	Expect(l2.Messages).To(HaveLen(1))
	Expect(l2.Messages[0].Message).To(Equal("test 1234"))
	Expect(l2.Messages[0].Attrs).To(HaveLen(3))
	Expect(l2.Messages[0].Attrs["attr1"]).To(Equal(1234))
	Expect(l2.Messages[0].Attrs["attr2"]).To(Equal(4321))
	Expect(l2.Messages[0].Attrs["attr3"]).To(Equal("val3"))
	Expect(l2.Messages[0].Level).To(Equal(LevelInfo))
}

func (s *GomolSuite) TestBaseWarn(t *testing.T) {
	b := NewBase()

	l1 := newDefaultMemLogger()
	l2 := newDefaultMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Warn("test")
	b.ShutdownLoggers()

	Expect(l1.Messages).To(HaveLen(1))
	Expect(l1.Messages[0].Message).To(Equal("test"))
	Expect(l1.Messages[0].Attrs).To(HaveLen(0))
	Expect(l1.Messages[0].Level).To(Equal(LevelWarning))

	Expect(l2.Messages).To(HaveLen(1))
	Expect(l2.Messages[0].Message).To(Equal("test"))
	Expect(l2.Messages[0].Attrs).To(HaveLen(0))
	Expect(l2.Messages[0].Level).To(Equal(LevelWarning))
}

func (s *GomolSuite) TestBaseWarnf(t *testing.T) {
	b := NewBase()

	l1 := newDefaultMemLogger()
	l2 := newDefaultMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Warnf("test %v", 1234)
	b.ShutdownLoggers()

	Expect(l1.Messages).To(HaveLen(1))
	Expect(l1.Messages[0].Message).To(Equal("test 1234"))
	Expect(l1.Messages[0].Attrs).To(HaveLen(0))
	Expect(l1.Messages[0].Level).To(Equal(LevelWarning))

	Expect(l2.Messages).To(HaveLen(1))
	Expect(l2.Messages[0].Message).To(Equal("test 1234"))
	Expect(l2.Messages[0].Attrs).To(HaveLen(0))
	Expect(l2.Messages[0].Level).To(Equal(LevelWarning))
}

func (s *GomolSuite) TestBaseWarnm(t *testing.T) {
	b := NewBase()
	b.SetAttr("attr1", 1234)

	l1 := newDefaultMemLogger()
	l2 := newDefaultMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Warnm(
		NewAttrs().
			SetAttr("attr2", 4321).
			SetAttr("attr3", "val3"),
		"test %v",
		1234)
	b.ShutdownLoggers()

	Expect(l1.Messages).To(HaveLen(1))
	Expect(l1.Messages[0].Message).To(Equal("test 1234"))
	Expect(l1.Messages[0].Attrs).To(HaveLen(3))
	Expect(l1.Messages[0].Attrs["attr1"]).To(Equal(1234))
	Expect(l1.Messages[0].Attrs["attr2"]).To(Equal(4321))
	Expect(l1.Messages[0].Attrs["attr3"]).To(Equal("val3"))
	Expect(l1.Messages[0].Level).To(Equal(LevelWarning))

	Expect(l2.Messages).To(HaveLen(1))
	Expect(l2.Messages[0].Message).To(Equal("test 1234"))
	Expect(l2.Messages[0].Attrs).To(HaveLen(3))
	Expect(l2.Messages[0].Attrs["attr1"]).To(Equal(1234))
	Expect(l2.Messages[0].Attrs["attr2"]).To(Equal(4321))
	Expect(l2.Messages[0].Attrs["attr3"]).To(Equal("val3"))
	Expect(l2.Messages[0].Level).To(Equal(LevelWarning))
}

func (s *GomolSuite) TestBaseErr(t *testing.T) {
	b := NewBase()

	l1 := newDefaultMemLogger()
	l2 := newDefaultMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Err("test")
	b.ShutdownLoggers()

	Expect(l1.Messages).To(HaveLen(1))
	Expect(l1.Messages[0].Message).To(Equal("test"))
	Expect(l1.Messages[0].Attrs).To(HaveLen(0))
	Expect(l1.Messages[0].Level).To(Equal(LevelError))

	Expect(l2.Messages).To(HaveLen(1))
	Expect(l2.Messages[0].Message).To(Equal("test"))
	Expect(l2.Messages[0].Attrs).To(HaveLen(0))
	Expect(l2.Messages[0].Level).To(Equal(LevelError))
}

func (s *GomolSuite) TestBaseErrf(t *testing.T) {
	b := NewBase()

	l1 := newDefaultMemLogger()
	l2 := newDefaultMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Errf("test %v", 1234)
	b.ShutdownLoggers()

	Expect(l1.Messages).To(HaveLen(1))
	Expect(l1.Messages[0].Message).To(Equal("test 1234"))
	Expect(l1.Messages[0].Attrs).To(HaveLen(0))
	Expect(l1.Messages[0].Level).To(Equal(LevelError))

	Expect(l2.Messages).To(HaveLen(1))
	Expect(l2.Messages[0].Message).To(Equal("test 1234"))
	Expect(l2.Messages[0].Attrs).To(HaveLen(0))
	Expect(l2.Messages[0].Level).To(Equal(LevelError))
}

func (s *GomolSuite) TestBaseErrm(t *testing.T) {
	b := NewBase()
	b.SetAttr("attr1", 1234)

	l1 := newDefaultMemLogger()
	l2 := newDefaultMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Errm(
		NewAttrs().
			SetAttr("attr2", 4321).
			SetAttr("attr3", "val3"),
		"test %v",
		1234)
	b.ShutdownLoggers()

	Expect(l1.Messages).To(HaveLen(1))
	Expect(l1.Messages[0].Message).To(Equal("test 1234"))
	Expect(l1.Messages[0].Attrs).To(HaveLen(3))
	Expect(l1.Messages[0].Attrs["attr1"]).To(Equal(1234))
	Expect(l1.Messages[0].Attrs["attr2"]).To(Equal(4321))
	Expect(l1.Messages[0].Attrs["attr3"]).To(Equal("val3"))
	Expect(l1.Messages[0].Level).To(Equal(LevelError))

	Expect(l2.Messages).To(HaveLen(1))
	Expect(l2.Messages[0].Message).To(Equal("test 1234"))
	Expect(l2.Messages[0].Attrs).To(HaveLen(3))
	Expect(l2.Messages[0].Attrs["attr1"]).To(Equal(1234))
	Expect(l2.Messages[0].Attrs["attr2"]).To(Equal(4321))
	Expect(l2.Messages[0].Attrs["attr3"]).To(Equal("val3"))
	Expect(l2.Messages[0].Level).To(Equal(LevelError))
}

func (s *GomolSuite) TestBaseFatal(t *testing.T) {
	b := NewBase()

	l1 := newDefaultMemLogger()
	l2 := newDefaultMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Fatal("test")
	b.ShutdownLoggers()

	Expect(l1.Messages).To(HaveLen(1))
	Expect(l1.Messages[0].Message).To(Equal("test"))
	Expect(l1.Messages[0].Attrs).To(HaveLen(0))
	Expect(l1.Messages[0].Level).To(Equal(LevelFatal))

	Expect(l2.Messages).To(HaveLen(1))
	Expect(l2.Messages[0].Message).To(Equal("test"))
	Expect(l2.Messages[0].Attrs).To(HaveLen(0))
	Expect(l2.Messages[0].Level).To(Equal(LevelFatal))
}

func (s *GomolSuite) TestBaseFatalf(t *testing.T) {
	b := NewBase()

	l1 := newDefaultMemLogger()
	l2 := newDefaultMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Fatalf("test %v", 1234)
	b.ShutdownLoggers()

	Expect(l1.Messages).To(HaveLen(1))
	Expect(l1.Messages[0].Message).To(Equal("test 1234"))
	Expect(l1.Messages[0].Attrs).To(HaveLen(0))
	Expect(l1.Messages[0].Level).To(Equal(LevelFatal))

	Expect(l2.Messages).To(HaveLen(1))
	Expect(l2.Messages[0].Message).To(Equal("test 1234"))
	Expect(l2.Messages[0].Attrs).To(HaveLen(0))
	Expect(l2.Messages[0].Level).To(Equal(LevelFatal))
}

func (s *GomolSuite) TestBaseFatalm(t *testing.T) {
	b := NewBase()
	b.SetAttr("attr1", 1234)

	l1 := newDefaultMemLogger()
	l2 := newDefaultMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Fatalm(
		NewAttrs().
			SetAttr("attr2", 4321).
			SetAttr("attr3", "val3"),
		"test %v",
		1234)
	b.ShutdownLoggers()

	Expect(l1.Messages).To(HaveLen(1))
	Expect(l1.Messages[0].Message).To(Equal("test 1234"))
	Expect(l1.Messages[0].Attrs).To(HaveLen(3))
	Expect(l1.Messages[0].Attrs["attr1"]).To(Equal(1234))
	Expect(l1.Messages[0].Attrs["attr2"]).To(Equal(4321))
	Expect(l1.Messages[0].Attrs["attr3"]).To(Equal("val3"))
	Expect(l1.Messages[0].Level).To(Equal(LevelFatal))

	Expect(l2.Messages).To(HaveLen(1))
	Expect(l2.Messages[0].Message).To(Equal("test 1234"))
	Expect(l2.Messages[0].Attrs).To(HaveLen(3))
	Expect(l2.Messages[0].Attrs["attr1"]).To(Equal(1234))
	Expect(l2.Messages[0].Attrs["attr2"]).To(Equal(4321))
	Expect(l2.Messages[0].Attrs["attr3"]).To(Equal("val3"))
	Expect(l2.Messages[0].Level).To(Equal(LevelFatal))
}

func (s *GomolSuite) TestBaseDie(t *testing.T) {
	b := NewBase()

	l1 := newDefaultMemLogger()
	l2 := newDefaultMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Die(1234, "test")

	Expect(l1.Messages).To(HaveLen(1))
	Expect(l1.Messages[0].Message).To(Equal("test"))
	Expect(l1.Messages[0].Attrs).To(HaveLen(0))
	Expect(l1.Messages[0].Level).To(Equal(LevelFatal))

	Expect(l2.Messages).To(HaveLen(1))
	Expect(l2.Messages[0].Message).To(Equal("test"))
	Expect(l2.Messages[0].Attrs).To(HaveLen(0))
	Expect(l2.Messages[0].Level).To(Equal(LevelFatal))

	Expect(b.isInitialized).To(Equal(false))
	Expect(curTestExiter.exited).To(Equal(true))
	Expect(curTestExiter.code).To(Equal(1234))
}

func (s *GomolSuite) TestBaseDief(t *testing.T) {
	b := NewBase()

	l1 := newDefaultMemLogger()
	l2 := newDefaultMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Dief(1234, "test %v", 1234)

	Expect(l1.Messages).To(HaveLen(1))
	Expect(l1.Messages[0].Message).To(Equal("test 1234"))
	Expect(l1.Messages[0].Attrs).To(HaveLen(0))
	Expect(l1.Messages[0].Level).To(Equal(LevelFatal))

	Expect(l2.Messages).To(HaveLen(1))
	Expect(l2.Messages[0].Message).To(Equal("test 1234"))
	Expect(l2.Messages[0].Attrs).To(HaveLen(0))
	Expect(l2.Messages[0].Level).To(Equal(LevelFatal))

	Expect(b.isInitialized).To(Equal(false))
	Expect(curTestExiter.exited).To(Equal(true))
	Expect(curTestExiter.code).To(Equal(1234))
}

func (s *GomolSuite) TestBaseDiem(t *testing.T) {
	b := NewBase()
	b.SetAttr("attr1", 1234)

	l1 := newDefaultMemLogger()
	l2 := newDefaultMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Diem(
		1234,
		NewAttrs().
			SetAttr("attr2", 4321).
			SetAttr("attr3", "val3"),
		"test %v",
		1234)
	b.ShutdownLoggers()

	Expect(l1.Messages).To(HaveLen(1))
	Expect(l1.Messages[0].Message).To(Equal("test 1234"))
	Expect(l1.Messages[0].Attrs).To(HaveLen(3))
	Expect(l1.Messages[0].Attrs["attr1"]).To(Equal(1234))
	Expect(l1.Messages[0].Attrs["attr2"]).To(Equal(4321))
	Expect(l1.Messages[0].Attrs["attr3"]).To(Equal("val3"))
	Expect(l1.Messages[0].Level).To(Equal(LevelFatal))

	Expect(l2.Messages).To(HaveLen(1))
	Expect(l2.Messages[0].Message).To(Equal("test 1234"))
	Expect(l2.Messages[0].Attrs).To(HaveLen(3))
	Expect(l2.Messages[0].Attrs["attr1"]).To(Equal(1234))
	Expect(l2.Messages[0].Attrs["attr2"]).To(Equal(4321))
	Expect(l2.Messages[0].Attrs["attr3"]).To(Equal("val3"))
	Expect(l2.Messages[0].Level).To(Equal(LevelFatal))

	Expect(b.isInitialized).To(Equal(false))
	Expect(curTestExiter.exited).To(Equal(true))
	Expect(curTestExiter.code).To(Equal(1234))
}

func (s *GomolSuite) TestBaseOrdering(t *testing.T) {
	b := NewBase()
	b.SetAttr("attr1", 1234)

	l1 := newDefaultMemLogger()
	l2 := newDefaultMemLogger()

	b.AddLogger(l1)
	b.AddLogger(l2)

	b.InitLoggers()
	b.Fatalm(
		NewAttrs().
			SetAttr("attr2", 4321).
			SetAttr("attr3", "val3"),
		"test %v",
		1234)
	b.Fatalm(
		NewAttrs().
			SetAttr("attr4", 4321).
			SetAttr("attr5", "val3"),
		"test %v",
		4321)
	b.ShutdownLoggers()

	Expect(l1.Messages).To(HaveLen(2))
	Expect(l1.Messages[0].Message).To(Equal("test 1234"))
	Expect(l1.Messages[0].Attrs).To(HaveLen(3))
	Expect(l1.Messages[0].Attrs["attr1"]).To(Equal(1234))
	Expect(l1.Messages[0].Attrs["attr2"]).To(Equal(4321))
	Expect(l1.Messages[0].Attrs["attr3"]).To(Equal("val3"))
	Expect(l1.Messages[0].Level).To(Equal(LevelFatal))
	Expect(l1.Messages[1].Message).To(Equal("test 4321"))
	Expect(l1.Messages[1].Attrs).To(HaveLen(3))
	Expect(l1.Messages[1].Attrs["attr1"]).To(Equal(1234))
	Expect(l1.Messages[1].Attrs["attr4"]).To(Equal(4321))
	Expect(l1.Messages[1].Attrs["attr5"]).To(Equal("val3"))
	Expect(l1.Messages[1].Level).To(Equal(LevelFatal))

	Expect(l2.Messages).To(HaveLen(2))
	Expect(l2.Messages[0].Message).To(Equal("test 1234"))
	Expect(l2.Messages[0].Attrs).To(HaveLen(3))
	Expect(l2.Messages[0].Attrs["attr1"]).To(Equal(1234))
	Expect(l2.Messages[0].Attrs["attr2"]).To(Equal(4321))
	Expect(l2.Messages[0].Attrs["attr3"]).To(Equal("val3"))
	Expect(l2.Messages[0].Level).To(Equal(LevelFatal))
	Expect(l2.Messages[1].Message).To(Equal("test 4321"))
	Expect(l2.Messages[1].Attrs).To(HaveLen(3))
	Expect(l2.Messages[1].Attrs["attr1"]).To(Equal(1234))
	Expect(l2.Messages[1].Attrs["attr4"]).To(Equal(4321))
	Expect(l2.Messages[1].Attrs["attr5"]).To(Equal("val3"))
	Expect(l2.Messages[1].Level).To(Equal(LevelFatal))
}
