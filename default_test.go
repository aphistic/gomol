package gomol

import (
	"testing"

	. "github.com/onsi/gomega"
)

func (s *GomolSuite) TestDefaultSetConfig(t *testing.T) {
	cfg := NewConfig()

	Expect(curDefault.config).To(Equal(cfg))

	cfg.FilenameAttr = "file"
	cfg.LineNumberAttr = "line"
	cfg.SequenceAttr = "seq"
	SetConfig(cfg)

	Expect(curDefault.config).To(Equal(cfg))
	Expect(curDefault.config).To(Equal(cfg))
}

func (s *GomolSuite) TestDefaultInitLogger(t *testing.T) {
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

func (s *GomolSuite) TestDefaultShutdownLogger(t *testing.T) {
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

func (s *GomolSuite) TestDefaultAddLogger(t *testing.T) {
	curDefault = NewBase()
	Expect(curDefault.loggers).To(HaveLen(0))
	AddLogger(newDefaultMemLogger())
	Expect(curDefault.loggers).To(HaveLen(1))
}

func (s *GomolSuite) TestDefaultRemoveLogger(t *testing.T) {
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

func (s *GomolSuite) TestDefaultClearLoggers(t *testing.T) {
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

func (s *GomolSuite) TestDefaultSetLogLevel(t *testing.T) {
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
	Expect(ml.Messages).To(HaveLen(3))
}

func (s *GomolSuite) TestDefaultSetAttr(t *testing.T) {
	curDefault = NewBase()
	Expect(curDefault.BaseAttrs.Attrs()).To(HaveLen(0))
	SetAttr("attr", 1234)
	Expect(curDefault.BaseAttrs.Attrs()).To(HaveLen(1))
	Expect(curDefault.BaseAttrs.GetAttr("attr")).To(Equal(1234))
}

func (s *GomolSuite) TestDefaultGetAttr(t *testing.T) {
	curDefault = NewBase()
	SetAttr("attr1", 1)
	SetAttr("attr2", "val2")

	Expect(GetAttr("attr2")).To(Equal("val2"))
	Expect(GetAttr("notakey")).To(BeNil())
}

func (s *GomolSuite) TestDefaultRemoveAttr(t *testing.T) {
	curDefault = NewBase()
	Expect(curDefault.BaseAttrs.Attrs()).To(HaveLen(0))
	SetAttr("attr", 1234)
	Expect(curDefault.BaseAttrs.Attrs()).To(HaveLen(1))
	Expect(curDefault.BaseAttrs.GetAttr("attr")).To(Equal(1234))
	RemoveAttr("attr")
	Expect(curDefault.BaseAttrs.Attrs()).To(HaveLen(0))
}

func (s *GomolSuite) TestDefaultClearAttrs(t *testing.T) {
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

func (s *GomolSuite) TestDefaultNewLogAdapter(t *testing.T) {
	la := NewLogAdapter(NewAttrs().SetAttr("foo", "bar"))
	defLogger := curDefault.loggers[0].(*memLogger)

	la.Dbgm(NewAttrs().SetAttr("attr", "val"), "test")

	ShutdownLoggers()

	Expect(len(defLogger.Messages)).To(Equal(1))
	Expect(defLogger.Messages[0]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelDebug,
		Message:   "test",
		Attrs: map[string]interface{}{
			"foo":  "bar",
			"attr": "val",
		},
	}))
}

func (s *GomolSuite) TestDefaultDbg(t *testing.T) {
	Dbg("test")
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*memLogger)
	Expect(defLogger.Messages).To(HaveLen(1))
	Expect(defLogger.Messages[0]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelDebug,
		Message:   "test",
		Attrs:     map[string]interface{}{},
	}))
}

func (s *GomolSuite) TestDefaultDbgf(t *testing.T) {
	Dbgf("test %v", 1234)
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*memLogger)
	Expect(defLogger.Messages).To(HaveLen(1))
	Expect(defLogger.Messages[0]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelDebug,
		Message:   "test 1234",
		Attrs:     map[string]interface{}{},
	}))
}

func (s *GomolSuite) TestDefaultDbgm(t *testing.T) {
	Dbgm(
		NewAttrs().SetAttr("attr1", 4321),
		"test %v", 1234)
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*memLogger)
	Expect(defLogger.Messages).To(HaveLen(1))
	Expect(defLogger.Messages[0]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelDebug,
		Message:   "test 1234",
		Attrs: map[string]interface{}{
			"attr1": 4321,
		},
	}))
}

func (s *GomolSuite) TestDefaultInfo(t *testing.T) {
	Info("test")
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*memLogger)
	Expect(defLogger.Messages).To(HaveLen(1))
	Expect(defLogger.Messages[0]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelInfo,
		Message:   "test",
		Attrs:     map[string]interface{}{},
	}))
}

func (s *GomolSuite) TestDefaultInfof(t *testing.T) {
	Infof("test %v", 1234)
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*memLogger)
	Expect(defLogger.Messages).To(HaveLen(1))
	Expect(defLogger.Messages[0]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelInfo,
		Message:   "test 1234",
		Attrs:     map[string]interface{}{},
	}))
}

func (s *GomolSuite) TestDefaultInfom(t *testing.T) {
	Infom(
		NewAttrs().SetAttr("attr1", 4321),
		"test %v", 1234)
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*memLogger)
	Expect(defLogger.Messages).To(HaveLen(1))
	Expect(defLogger.Messages[0]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelInfo,
		Message:   "test 1234",
		Attrs: map[string]interface{}{
			"attr1": 4321,
		},
	}))
}

func (s *GomolSuite) TestDefaultWarn(t *testing.T) {
	Warn("test")
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*memLogger)
	Expect(defLogger.Messages).To(HaveLen(1))
	Expect(defLogger.Messages[0]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelWarning,
		Message:   "test",
		Attrs:     map[string]interface{}{},
	}))
}

func (s *GomolSuite) TestDefaultWarnf(t *testing.T) {
	Warnf("test %v", 1234)
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*memLogger)
	Expect(defLogger.Messages).To(HaveLen(1))
	Expect(defLogger.Messages[0]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelWarning,
		Message:   "test 1234",
		Attrs:     map[string]interface{}{},
	}))
}

func (s *GomolSuite) TestDefaultWarnm(t *testing.T) {
	Warnm(
		NewAttrs().SetAttr("attr1", 4321),
		"test %v", 1234)
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*memLogger)
	Expect(defLogger.Messages).To(HaveLen(1))
	Expect(defLogger.Messages[0]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelWarning,
		Message:   "test 1234",
		Attrs: map[string]interface{}{
			"attr1": 4321,
		},
	}))
}

func (s *GomolSuite) TestDefaultErr(t *testing.T) {
	Err("test")
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*memLogger)
	Expect(defLogger.Messages).To(HaveLen(1))
	Expect(defLogger.Messages[0]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelError,
		Message:   "test",
		Attrs:     map[string]interface{}{},
	}))
}

func (s *GomolSuite) TestDefaultErrf(t *testing.T) {
	Errf("test %v", 1234)
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*memLogger)
	Expect(defLogger.Messages).To(HaveLen(1))
	Expect(defLogger.Messages[0]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelError,
		Message:   "test 1234",
		Attrs:     map[string]interface{}{},
	}))
}

func (s *GomolSuite) TestDefaultErrm(t *testing.T) {
	Errm(
		NewAttrs().SetAttr("attr1", 4321),
		"test %v", 1234)
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*memLogger)
	Expect(defLogger.Messages).To(HaveLen(1))
	Expect(defLogger.Messages[0]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelError,
		Message:   "test 1234",
		Attrs: map[string]interface{}{
			"attr1": 4321,
		},
	}))
}

func (s *GomolSuite) TestDefaultFatal(t *testing.T) {
	Fatal("test")
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*memLogger)
	Expect(defLogger.Messages).To(HaveLen(1))
	Expect(defLogger.Messages[0]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelFatal,
		Message:   "test",
		Attrs:     map[string]interface{}{},
	}))
}

func (s *GomolSuite) TestDefaultFatalf(t *testing.T) {
	Fatalf("test %v", 1234)
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*memLogger)
	Expect(defLogger.Messages).To(HaveLen(1))
	Expect(defLogger.Messages[0]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelFatal,
		Message:   "test 1234",
		Attrs:     map[string]interface{}{},
	}))
}

func (s *GomolSuite) TestDefaultFatalm(t *testing.T) {
	Fatalm(
		NewAttrs().SetAttr("attr1", 4321),
		"test %v", 1234)
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*memLogger)
	Expect(defLogger.Messages).To(HaveLen(1))
	Expect(defLogger.Messages[0]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelFatal,
		Message:   "test 1234",
		Attrs: map[string]interface{}{
			"attr1": 4321,
		},
	}))
}

func (s *GomolSuite) TestDefaultDie(t *testing.T) {
	Die(1234, "test")
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*memLogger)
	Expect(defLogger.Messages).To(HaveLen(1))
	Expect(defLogger.Messages[0]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelFatal,
		Message:   "test",
		Attrs:     map[string]interface{}{},
	}))

	Expect(curDefault.isInitialized).To(Equal(false))
	Expect(curTestExiter.exited).To(Equal(true))
	Expect(curTestExiter.code).To(Equal(1234))
}

func (s *GomolSuite) TestDefaultDief(t *testing.T) {
	Dief(1234, "test %v", 1234)
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*memLogger)
	Expect(defLogger.Messages).To(HaveLen(1))
	Expect(defLogger.Messages[0]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelFatal,
		Message:   "test 1234",
		Attrs:     map[string]interface{}{},
	}))
	Expect(defLogger.Messages[0].Level).To(Equal(LevelFatal))
	Expect(defLogger.Messages[0].Message).To(Equal("test 1234"))
	Expect(defLogger.Messages[0].Attrs).To(HaveLen(0))

	Expect(curDefault.isInitialized).To(Equal(false))
	Expect(curTestExiter.exited).To(Equal(true))
	Expect(curTestExiter.code).To(Equal(1234))
}

func (s *GomolSuite) TestDefaultDiem(t *testing.T) {
	Diem(
		1234,
		NewAttrs().SetAttr("attr1", 4321),
		"test %v", 1234)
	curDefault.queue.stopQueueWorkers()
	defLogger := curDefault.loggers[0].(*memLogger)
	Expect(defLogger.Messages).To(HaveLen(1))
	Expect(defLogger.Messages[0]).To(Equal(&memMessage{
		Timestamp: clock().Now(),
		Level:     LevelFatal,
		Message:   "test 1234",
		Attrs: map[string]interface{}{
			"attr1": 4321,
		},
	}))

	Expect(curDefault.isInitialized).To(Equal(false))
	Expect(curTestExiter.exited).To(Equal(true))
	Expect(curTestExiter.code).To(Equal(1234))
}
