package gomol

import (
	"time"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

func newDefaultMemLogger() *memLogger {
	cfg := newMemLoggerConfig()
	l, _ := newMemLogger(cfg)
	return l
}

func (s *GomolSuite) TestMemInitLogger(t sweet.T) {
	ml := newDefaultMemLogger()
	Expect(ml.IsInitialized()).To(Equal(false))
	ml.InitLogger()
	Expect(ml.IsInitialized()).To(Equal(true))
}

func (s *GomolSuite) TestMemInitLoggerFail(t sweet.T) {
	mlCfg := newMemLoggerConfig()
	mlCfg.FailInit = true
	ml, err := newMemLogger(mlCfg)
	Expect(err).To(BeNil())
	Expect(ml.IsInitialized()).To(Equal(false))
	err = ml.InitLogger()
	Expect(ml.IsInitialized()).To(Equal(false))
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Init failed"))
}

func (s *GomolSuite) TestMemShutdownLogger(t sweet.T) {
	ml := newDefaultMemLogger()
	Expect(ml.isShutdown).To(Equal(false))
	ml.ShutdownLogger()
	Expect(ml.isShutdown).To(Equal(true))
}

func (s *GomolSuite) TestMemShutdownLoggerFail(t sweet.T) {
	mlCfg := newMemLoggerConfig()
	mlCfg.FailShutdown = true
	ml, err := newMemLogger(mlCfg)
	Expect(err).To(BeNil())
	Expect(ml.isShutdown).To(Equal(false))
	err = ml.ShutdownLogger()
	Expect(ml.isShutdown).To(Equal(false))
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Shutdown failed"))
}

func (s *GomolSuite) TestMemClearMessages(t sweet.T) {
	ml := newDefaultMemLogger()
	Expect(ml.Messages).To(HaveLen(0))
	ml.Logm(time.Now(), LevelDebug, nil, "test")
	Expect(ml.Messages).To(HaveLen(1))
	ml.ClearMessages()
	Expect(ml.Messages).To(HaveLen(0))
}

func (s *GomolSuite) TestMemLogmNoAttrs(t sweet.T) {
	ml := newDefaultMemLogger()
	ml.Logm(time.Now(), LevelDebug, nil, "test")
	Expect(ml.Messages).To(HaveLen(1))
	Expect(ml.Messages[0].Level).To(Equal(LevelDebug))
	Expect(ml.Messages[0].Message).To(Equal("test"))
	Expect(ml.Messages[0].Attrs).To(HaveLen(0))
}

func (s *GomolSuite) TestMemLogmAttrs(t sweet.T) {
	setClock(newTestClock(time.Now()))
	ml := newDefaultMemLogger()
	ml.Logm(
		clock().Now(),
		LevelDebug,
		map[string]interface{}{"attr1": 4321},
		"test 1234")
	Expect(ml.Messages).To(HaveLen(1))
	Expect(ml.Messages[0].Timestamp).To(Equal(clock().Now()))
	Expect(ml.Messages[0].Level).To(Equal(LevelDebug))
	Expect(ml.Messages[0].Message).To(Equal("test 1234"))
	Expect(ml.Messages[0].Attrs).To(HaveLen(1))
	Expect(ml.Messages[0].Attrs["attr1"]).To(Equal(4321))
}

func (s *GomolSuite) TestMemBaseAttrs(t sweet.T) {
	setClock(newTestClock(time.Now()))

	b := NewBase()
	b.SetAttr("attr1", 7890)
	b.SetAttr("attr2", "val2")

	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	ml.Logm(
		clock().Now(),
		LevelDebug,
		map[string]interface{}{
			"attr1": 4321,
			"attr3": "val3",
		},
		"test 1234")
	Expect(ml.Messages).To(HaveLen(1))
	Expect(ml.Messages[0].Timestamp).To(Equal(clock().Now()))
	Expect(ml.Messages[0].Level).To(Equal(LevelDebug))
	Expect(ml.Messages[0].Message).To(Equal("test 1234"))
	Expect(ml.Messages[0].Attrs).To(HaveLen(3))
	Expect(ml.Messages[0].Attrs["attr1"]).To(Equal(4321))
	Expect(ml.Messages[0].Attrs["attr2"]).To(Equal("val2"))
	Expect(ml.Messages[0].Attrs["attr3"]).To(Equal("val3"))
}

func (s *GomolSuite) TestMemStringAttrs(t sweet.T) {
	setClock(newTestClock(time.Now()))

	b := NewBase()
	b.SetAttr("attr1", 1234)
	b.SetAttr("attr2", "val2")

	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	ml.Logm(
		clock().Now(),
		LevelDebug,
		map[string]interface{}{
			"attr1": 4321,
			"attr3": "val3",
		},
		"test 1234",
	)
	Expect(ml.Messages).To(HaveLen(1))
	Expect(ml.Messages[0].Timestamp).To(Equal(clock().Now()))
	Expect(ml.Messages[0].Level).To(Equal(LevelDebug))
	Expect(ml.Messages[0].Message).To(Equal("test 1234"))
	Expect(ml.Messages[0].Attrs).To(HaveLen(3))
	Expect(ml.Messages[0].Attrs["attr1"]).To(Equal(4321))
	Expect(ml.Messages[0].Attrs["attr2"]).To(Equal("val2"))
	Expect(ml.Messages[0].Attrs["attr3"]).To(Equal("val3"))
	Expect(ml.Messages[0].StringAttrs["attr1"]).To(Equal("4321"))
	Expect(ml.Messages[0].StringAttrs["attr2"]).To(Equal("val2"))
	Expect(ml.Messages[0].StringAttrs["attr3"]).To(Equal("val3"))
}
