package gomol

import (
	"github.com/aphistic/sweet"
	"github.com/efritz/glock"
	. "github.com/onsi/gomega"
)

// These tests cannot be run in parallel because they depend on global state (the default logger)
type DefaultSuite struct {
	currentClock glock.Clock
}

func (s *DefaultSuite) SetUpTest(t sweet.T) {
	s.currentClock = glock.NewMockClock()
	curDefault = NewBase(
		withClock(s.currentClock),
	)
	cfg := newMemLoggerConfig()
	memL, _ := newMemLogger(cfg)
	curDefault.AddLogger(memL)
	curDefault.InitLoggers()
}

func (s *DefaultSuite) TearDownTest(t sweet.T) {
	curDefault.ShutdownLoggers()
}

func (s *DefaultSuite) TestDefaultSetConfig(t sweet.T) {
	cfg := NewConfig()

	Expect(curDefault.config).To(Equal(cfg))

	cfg.FilenameAttr = "file"
	cfg.LineNumberAttr = "line"
	cfg.SequenceAttr = "seq"
	SetConfig(cfg)

	Expect(curDefault.config).To(Equal(cfg))
	Expect(curDefault.config).To(Equal(cfg))
}

func (s *DefaultSuite) TestDefaultSetErrorChan(t sweet.T) {
	ch := make(chan error)
	Expect(curDefault.errorChan).To(BeNil())
	SetErrorChan(ch)
	Expect(curDefault.errorChan).To(Equal((chan<- error)(ch)))
}

func (s *DefaultSuite) TestDefaultInitLogger(t sweet.T) {
	curDefault = NewBase()
	Expect(IsInitialized()).To(Equal(false))
	AddLogger(newDefaultMemLogger())
	defLogger := curDefault.loggers[0].(*memLogger)
	Expect(defLogger.IsInitialized()).To(Equal(false))
	InitLoggers()
	Expect(IsInitialized()).To(Equal(true))
	Expect(defLogger.IsInitialized()).To(Equal(true))
	ShutdownLoggers()
}

func (s *DefaultSuite) TestDefaultShutdownLogger(t sweet.T) {
	curDefault = NewBase()
	Expect(IsInitialized()).To(Equal(false))
	AddLogger(newDefaultMemLogger())
	InitLoggers()
	Expect(IsInitialized()).To(Equal(true))
	defLogger := curDefault.loggers[0].(*memLogger)
	Expect(defLogger.isShutdown).To(Equal(false))
	ShutdownLoggers()
	Expect(defLogger.isShutdown).To(Equal(true))
	Expect(IsInitialized()).To(Equal(false))
}

func (s *DefaultSuite) TestSetFallbackLogger(t sweet.T) {
	curDefault = NewBase()
	Expect(curDefault.fallbackLogger).To(BeNil())

	ml := newDefaultMemLogger()
	SetFallbackLogger(ml)
	Expect(curDefault.fallbackLogger).To(Equal(ml))
}

func (s *DefaultSuite) TestDefaultAddLogger(t sweet.T) {
	curDefault = NewBase()
	Expect(curDefault.loggers).To(HaveLen(0))
	AddLogger(newDefaultMemLogger())
	Expect(curDefault.loggers).To(HaveLen(1))
}

func (s *DefaultSuite) TestDefaultRemoveLogger(t sweet.T) {
	curDefault = NewBase()

	ml1 := newDefaultMemLogger()
	ml2 := newDefaultMemLogger()
	ml3 := newDefaultMemLogger()
	AddLogger(ml1)
	AddLogger(ml2)
	AddLogger(ml3)

	InitLoggers()

	Expect(ml1.IsInitialized()).To(Equal(true))
	Expect(ml2.IsInitialized()).To(Equal(true))
	Expect(ml3.IsInitialized()).To(Equal(true))
	Expect(curDefault.loggers).To(HaveLen(3))

	err := RemoveLogger(ml2)
	Expect(err).To(BeNil())
	Expect(ml1.IsInitialized()).To(Equal(true))
	Expect(ml2.IsInitialized()).To(Equal(false))
	Expect(ml3.IsInitialized()).To(Equal(true))
	Expect(curDefault.loggers).To(HaveLen(2))
}

func (s *DefaultSuite) TestDefaultClearLoggers(t sweet.T) {
	curDefault = NewBase()

	ml1 := newDefaultMemLogger()
	ml2 := newDefaultMemLogger()
	ml3 := newDefaultMemLogger()
	AddLogger(ml1)
	AddLogger(ml2)
	AddLogger(ml3)

	InitLoggers()

	Expect(ml1.IsInitialized()).To(Equal(true))
	Expect(ml2.IsInitialized()).To(Equal(true))
	Expect(ml3.IsInitialized()).To(Equal(true))
	Expect(curDefault.loggers).To(HaveLen(3))

	err := ClearLoggers()
	Expect(err).To(BeNil())
	Expect(ml1.IsInitialized()).To(Equal(false))
	Expect(ml2.IsInitialized()).To(Equal(false))
	Expect(ml3.IsInitialized()).To(Equal(false))
	Expect(curDefault.loggers).To(HaveLen(0))
}

func (s *DefaultSuite) TestDefaultSetLogLevel(t sweet.T) {
	curDefault = NewBase()
	InitLoggers()
	ml := newDefaultMemLogger()
	AddLogger(ml)

	SetLogLevel(LevelWarning)
	Dbg("test")
	Info("test")
	Warn("test")
	Err("test")
	Fatal("test")
	ShutdownLoggers()
	Expect(ml.Messages()).To(HaveLen(3))
}

func (s *DefaultSuite) TestDefaultSetAttr(t sweet.T) {
	curDefault = NewBase()
	Expect(curDefault.BaseAttrs.Attrs()).To(HaveLen(0))
	SetAttr("attr", 1234)
	Expect(curDefault.BaseAttrs.Attrs()).To(HaveLen(1))
	Expect(curDefault.BaseAttrs.GetAttr("attr")).To(Equal(1234))
}

func (s *DefaultSuite) TestDefaultGetAttr(t sweet.T) {
	curDefault = NewBase()
	SetAttr("attr1", 1)
	SetAttr("attr2", "val2")

	Expect(GetAttr("attr2")).To(Equal("val2"))
	Expect(GetAttr("notakey")).To(BeNil())
}

func (s *DefaultSuite) TestDefaultRemoveAttr(t sweet.T) {
	curDefault = NewBase()
	Expect(curDefault.BaseAttrs.Attrs()).To(HaveLen(0))
	SetAttr("attr", 1234)
	Expect(curDefault.BaseAttrs.Attrs()).To(HaveLen(1))
	Expect(curDefault.BaseAttrs.GetAttr("attr")).To(Equal(1234))
	RemoveAttr("attr")
	Expect(curDefault.BaseAttrs.Attrs()).To(HaveLen(0))
}

func (s *DefaultSuite) TestDefaultClearAttrs(t sweet.T) {
	curDefault = NewBase()
	Expect(curDefault.BaseAttrs.Attrs()).To(HaveLen(0))
	SetAttr("attr", 1234)
	Expect(curDefault.BaseAttrs.Attrs()).To(HaveLen(1))
	Expect(curDefault.BaseAttrs.GetAttr("attr")).To(Equal(1234))
	SetAttr("attr2", 1234)
	Expect(curDefault.BaseAttrs.Attrs()).To(HaveLen(2))
	Expect(curDefault.BaseAttrs.GetAttr("attr2")).To(Equal(1234))
	ClearAttrs()
	Expect(curDefault.BaseAttrs.Attrs()).To(HaveLen(0))
}

func (s *DefaultSuite) TestDefaultNewLogAdapter(t sweet.T) {
	la := NewLogAdapter(NewAttrs().SetAttr("foo", "bar"))
	defLogger := curDefault.loggers[0].(*memLogger)

	la.Dbgm(NewAttrs().SetAttr("attr", "val"), "test")

	ShutdownLoggers()

	Expect(defLogger.Messages()).To(HaveLen(1))
	Expect(defLogger.Messages()[0]).To(Equal(&memMessage{
		Timestamp: s.currentClock.Now(),
		Level:     LevelDebug,
		Message:   "test",
		Attrs: map[string]interface{}{
			"foo":  "bar",
			"attr": "val",
		},
		StringAttrs: map[string]string{
			"foo":  "bar",
			"attr": "val",
		},
	}))
}

func (s *DefaultSuite) TestDefaultDbg(t sweet.T) {
	Dbg("test")
	curDefault.Flush()
	defLogger := curDefault.loggers[0].(*memLogger)
	Expect(defLogger.Messages()).To(HaveLen(1))
	Expect(defLogger.Messages()[0]).To(Equal(&memMessage{
		Timestamp:   s.currentClock.Now(),
		Level:       LevelDebug,
		Message:     "test",
		Attrs:       map[string]interface{}{},
		StringAttrs: map[string]string{},
	}))
}

func (s *DefaultSuite) TestDefaultDbgf(t sweet.T) {
	Dbgf("test %v", 1234)
	curDefault.queue.stopWorker()
	defLogger := curDefault.loggers[0].(*memLogger)
	Expect(defLogger.Messages()).To(HaveLen(1))
	Expect(defLogger.Messages()[0]).To(Equal(&memMessage{
		Timestamp:   s.currentClock.Now(),
		Level:       LevelDebug,
		Message:     "test 1234",
		Attrs:       map[string]interface{}{},
		StringAttrs: map[string]string{},
	}))
}

func (s *DefaultSuite) TestDefaultDbgm(t sweet.T) {
	Dbgm(
		NewAttrs().SetAttr("attr1", 4321),
		"test %v", 1234)
	curDefault.queue.stopWorker()
	defLogger := curDefault.loggers[0].(*memLogger)
	Expect(defLogger.Messages()).To(HaveLen(1))
	Expect(defLogger.Messages()[0]).To(Equal(&memMessage{
		Timestamp: s.currentClock.Now(),
		Level:     LevelDebug,
		Message:   "test 1234",
		Attrs: map[string]interface{}{
			"attr1": 4321,
		},
		StringAttrs: map[string]string{
			"attr1": "4321",
		},
	}))
}

func (s *DefaultSuite) TestDefaultInfo(t sweet.T) {
	Info("test")
	curDefault.Flush()
	defLogger := curDefault.loggers[0].(*memLogger)
	Expect(defLogger.Messages()).To(HaveLen(1))
	Expect(defLogger.Messages()[0]).To(Equal(&memMessage{
		Timestamp:   s.currentClock.Now(),
		Level:       LevelInfo,
		Message:     "test",
		Attrs:       map[string]interface{}{},
		StringAttrs: map[string]string{},
	}))
}

func (s *DefaultSuite) TestDefaultInfof(t sweet.T) {
	Infof("test %v", 1234)
	curDefault.Flush()
	defLogger := curDefault.loggers[0].(*memLogger)
	Expect(defLogger.Messages()).To(HaveLen(1))
	Expect(defLogger.Messages()[0]).To(Equal(&memMessage{
		Timestamp:   s.currentClock.Now(),
		Level:       LevelInfo,
		Message:     "test 1234",
		Attrs:       map[string]interface{}{},
		StringAttrs: map[string]string{},
	}))
}

func (s *DefaultSuite) TestDefaultInfom(t sweet.T) {
	Infom(
		NewAttrs().SetAttr("attr1", 4321),
		"test %v",
		1234,
	)
	curDefault.Flush()
	defLogger := curDefault.loggers[0].(*memLogger)
	Expect(defLogger.Messages()).To(HaveLen(1))
	Expect(defLogger.Messages()[0]).To(Equal(&memMessage{
		Timestamp: s.currentClock.Now(),
		Level:     LevelInfo,
		Message:   "test 1234",
		Attrs: map[string]interface{}{
			"attr1": 4321,
		},
		StringAttrs: map[string]string{
			"attr1": "4321",
		},
	}))
}

func (s *DefaultSuite) TestDefaultWarn(t sweet.T) {
	Warn("test")
	curDefault.Flush()
	defLogger := curDefault.loggers[0].(*memLogger)
	Expect(defLogger.Messages()).To(HaveLen(1))
	Expect(defLogger.Messages()[0]).To(Equal(&memMessage{
		Timestamp:   s.currentClock.Now(),
		Level:       LevelWarning,
		Message:     "test",
		Attrs:       map[string]interface{}{},
		StringAttrs: map[string]string{},
	}))
}

func (s *DefaultSuite) TestDefaultWarnf(t sweet.T) {
	Warnf("test %v", 1234)
	curDefault.Flush()
	defLogger := curDefault.loggers[0].(*memLogger)
	Expect(defLogger.Messages()).To(HaveLen(1))
	Expect(defLogger.Messages()[0]).To(Equal(&memMessage{
		Timestamp:   s.currentClock.Now(),
		Level:       LevelWarning,
		Message:     "test 1234",
		Attrs:       map[string]interface{}{},
		StringAttrs: map[string]string{},
	}))
}

func (s *DefaultSuite) TestDefaultWarnm(t sweet.T) {
	Warnm(
		NewAttrs().SetAttr("attr1", 4321),
		"test %v",
		1234,
	)
	curDefault.Flush()
	defLogger := curDefault.loggers[0].(*memLogger)
	Expect(defLogger.Messages()).To(HaveLen(1))
	Expect(defLogger.Messages()[0]).To(Equal(&memMessage{
		Timestamp: s.currentClock.Now(),
		Level:     LevelWarning,
		Message:   "test 1234",
		Attrs: map[string]interface{}{
			"attr1": 4321,
		},
		StringAttrs: map[string]string{
			"attr1": "4321",
		},
	}))
}

func (s *DefaultSuite) TestDefaultErr(t sweet.T) {
	Err("test")
	curDefault.Flush()
	defLogger := curDefault.loggers[0].(*memLogger)
	Expect(defLogger.Messages()).To(HaveLen(1))
	Expect(defLogger.Messages()[0]).To(Equal(&memMessage{
		Timestamp:   s.currentClock.Now(),
		Level:       LevelError,
		Message:     "test",
		Attrs:       map[string]interface{}{},
		StringAttrs: map[string]string{},
	}))
}

func (s *DefaultSuite) TestDefaultErrf(t sweet.T) {
	Errf("test %v", 1234)
	curDefault.queue.stopWorker()
	defLogger := curDefault.loggers[0].(*memLogger)
	Expect(defLogger.Messages()).To(HaveLen(1))
	Expect(defLogger.Messages()[0]).To(Equal(&memMessage{
		Timestamp:   s.currentClock.Now(),
		Level:       LevelError,
		Message:     "test 1234",
		Attrs:       map[string]interface{}{},
		StringAttrs: map[string]string{},
	}))
}

func (s *DefaultSuite) TestDefaultErrm(t sweet.T) {
	Errm(
		NewAttrs().SetAttr("attr1", 4321),
		"test %v", 1234)
	curDefault.queue.stopWorker()
	defLogger := curDefault.loggers[0].(*memLogger)
	Expect(defLogger.Messages()).To(HaveLen(1))
	Expect(defLogger.Messages()[0]).To(Equal(&memMessage{
		Timestamp: s.currentClock.Now(),
		Level:     LevelError,
		Message:   "test 1234",
		Attrs: map[string]interface{}{
			"attr1": 4321,
		},
		StringAttrs: map[string]string{
			"attr1": "4321",
		},
	}))
}

func (s *DefaultSuite) TestDefaultFatal(t sweet.T) {
	Fatal("test")
	curDefault.queue.stopWorker()
	defLogger := curDefault.loggers[0].(*memLogger)
	Expect(defLogger.Messages()).To(HaveLen(1))
	Expect(defLogger.Messages()[0]).To(Equal(&memMessage{
		Timestamp:   s.currentClock.Now(),
		Level:       LevelFatal,
		Message:     "test",
		Attrs:       map[string]interface{}{},
		StringAttrs: map[string]string{},
	}))
}

func (s *DefaultSuite) TestDefaultFatalf(t sweet.T) {
	Fatalf("test %v", 1234)
	curDefault.queue.stopWorker()
	defLogger := curDefault.loggers[0].(*memLogger)
	Expect(defLogger.Messages()).To(HaveLen(1))
	Expect(defLogger.Messages()[0]).To(Equal(&memMessage{
		Timestamp:   s.currentClock.Now(),
		Level:       LevelFatal,
		Message:     "test 1234",
		Attrs:       map[string]interface{}{},
		StringAttrs: map[string]string{},
	}))
}

func (s *DefaultSuite) TestDefaultFatalm(t sweet.T) {
	Fatalm(
		NewAttrs().SetAttr("attr1", 4321),
		"test %v", 1234)
	curDefault.queue.stopWorker()
	defLogger := curDefault.loggers[0].(*memLogger)
	Expect(defLogger.Messages()).To(HaveLen(1))
	Expect(defLogger.Messages()[0]).To(Equal(&memMessage{
		Timestamp: s.currentClock.Now(),
		Level:     LevelFatal,
		Message:   "test 1234",
		Attrs: map[string]interface{}{
			"attr1": 4321,
		},
		StringAttrs: map[string]string{
			"attr1": "4321",
		},
	}))
}

func (s *DefaultSuite) TestDefaultDie(t sweet.T) {
	Die(1234, "test")
	curDefault.queue.stopWorker()
	defLogger := curDefault.loggers[0].(*memLogger)
	Expect(defLogger.Messages()).To(HaveLen(1))
	Expect(defLogger.Messages()[0]).To(Equal(&memMessage{
		Timestamp:   s.currentClock.Now(),
		Level:       LevelFatal,
		Message:     "test",
		Attrs:       map[string]interface{}{},
		StringAttrs: map[string]string{},
	}))

	Expect(curDefault.isInitialized).To(Equal(false))
	Expect(curTestExiter.exited).To(Equal(true))
	Expect(curTestExiter.code).To(Equal(1234))
}

func (s *DefaultSuite) TestDefaultDief(t sweet.T) {
	Dief(1234, "test %v", 1234)
	curDefault.queue.stopWorker()
	defLogger := curDefault.loggers[0].(*memLogger)
	Expect(defLogger.Messages()).To(HaveLen(1))
	Expect(defLogger.Messages()[0]).To(Equal(&memMessage{
		Timestamp:   s.currentClock.Now(),
		Level:       LevelFatal,
		Message:     "test 1234",
		Attrs:       map[string]interface{}{},
		StringAttrs: map[string]string{},
	}))
	Expect(defLogger.Messages()[0].Level).To(Equal(LevelFatal))
	Expect(defLogger.Messages()[0].Message).To(Equal("test 1234"))
	Expect(defLogger.Messages()[0].Attrs).To(HaveLen(0))

	Expect(curDefault.isInitialized).To(Equal(false))
	Expect(curTestExiter.exited).To(Equal(true))
	Expect(curTestExiter.code).To(Equal(1234))
}

func (s *DefaultSuite) TestDefaultDiem(t sweet.T) {
	Diem(
		1234,
		NewAttrs().SetAttr("attr1", 4321),
		"test %v", 1234)
	curDefault.queue.stopWorker()
	defLogger := curDefault.loggers[0].(*memLogger)
	Expect(defLogger.Messages()).To(HaveLen(1))
	Expect(defLogger.Messages()[0]).To(Equal(&memMessage{
		Timestamp: s.currentClock.Now(),
		Level:     LevelFatal,
		Message:   "test 1234",
		Attrs: map[string]interface{}{
			"attr1": 4321,
		},
		StringAttrs: map[string]string{
			"attr1": "4321",
		},
	}))

	Expect(curDefault.isInitialized).To(Equal(false))
	Expect(curTestExiter.exited).To(Equal(true))
	Expect(curTestExiter.code).To(Equal(1234))
}
