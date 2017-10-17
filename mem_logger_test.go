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

type MemLoggerSuite struct{}

func (s *MemLoggerSuite) TestMemInitLogger(t sweet.T) {
	ml := newDefaultMemLogger()
	Expect(ml.IsInitialized()).To(Equal(false))
	ml.InitLogger()
	Expect(ml.IsInitialized()).To(Equal(true))
}

func (s *MemLoggerSuite) TestMemInitLoggerFail(t sweet.T) {
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

func (s *MemLoggerSuite) TestMemShutdownLogger(t sweet.T) {
	ml := newDefaultMemLogger()
	Expect(ml.isShutdown).To(Equal(false))
	ml.ShutdownLogger()
	Expect(ml.isShutdown).To(Equal(true))
}

func (s *MemLoggerSuite) TestMemShutdownLoggerFail(t sweet.T) {
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

func (s *MemLoggerSuite) TestMemClearMessages(t sweet.T) {
	ml := newDefaultMemLogger()
	Expect(ml.Messages()).To(HaveLen(0))
	ml.Logm(time.Now(), LevelDebug, nil, "test")
	Expect(ml.Messages()).To(HaveLen(1))
	ml.ClearMessages()
	Expect(ml.Messages()).To(HaveLen(0))
}

func (s *MemLoggerSuite) TestMemLogmNoAttrs(t sweet.T) {
	ml := newDefaultMemLogger()
	ml.Logm(time.Now(), LevelDebug, nil, "test")

	Expect(ml.Messages()).To(HaveLen(1))
	msg := ml.Messages()[0]
	Expect(msg.Level).To(Equal(LevelDebug))
	Expect(msg.Message).To(Equal("test"))
	Expect(msg.Attrs).To(HaveLen(0))
}

func (s *MemLoggerSuite) TestMemLogmAttrs(t sweet.T) {
	ts := time.Unix(10, 0)
	ml := newDefaultMemLogger()
	ml.Logm(
		ts,
		LevelDebug,
		map[string]interface{}{"attr1": 4321},
		"test 1234")

	Expect(ml.Messages()).To(HaveLen(1))
	msg := ml.Messages()[0]
	Expect(msg.Timestamp).To(Equal(ts))
	Expect(msg.Level).To(Equal(LevelDebug))
	Expect(msg.Message).To(Equal("test 1234"))
	Expect(msg.Attrs).To(HaveLen(1))
	Expect(msg.Attrs["attr1"]).To(Equal(4321))
}

func (s *MemLoggerSuite) TestMemBaseAttrs(t sweet.T) {
	ts := time.Unix(10, 0)

	b := NewBase()
	b.SetAttr("attr1", 7890)
	b.SetAttr("attr2", "val2")

	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	ml.Logm(
		ts,
		LevelDebug,
		map[string]interface{}{
			"attr1": 4321,
			"attr3": "val3",
		},
		"test 1234")

	Expect(ml.Messages()).To(HaveLen(1))
	msg := ml.Messages()[0]
	Expect(msg.Timestamp).To(Equal(ts))
	Expect(msg.Level).To(Equal(LevelDebug))
	Expect(msg.Message).To(Equal("test 1234"))
	Expect(msg.Attrs).To(HaveLen(3))
	Expect(msg.Attrs["attr1"]).To(Equal(4321))
	Expect(msg.Attrs["attr2"]).To(Equal("val2"))
	Expect(msg.Attrs["attr3"]).To(Equal("val3"))
}

func (s *MemLoggerSuite) TestMemStringAttrs(t sweet.T) {
	ts := time.Unix(10, 0)

	b := NewBase()
	b.SetAttr("attr1", 1234)
	b.SetAttr("attr2", "val2")

	ml := newDefaultMemLogger()
	b.AddLogger(ml)
	ml.Logm(
		ts,
		LevelDebug,
		map[string]interface{}{
			"attr1": 4321,
			"attr3": "val3",
		},
		"test 1234",
	)

	Expect(ml.Messages()).To(HaveLen(1))
	msg := ml.Messages()[0]
	Expect(msg.Timestamp).To(Equal(ts))
	Expect(msg.Level).To(Equal(LevelDebug))
	Expect(msg.Message).To(Equal("test 1234"))
	Expect(msg.Attrs).To(HaveLen(3))
	Expect(msg.Attrs["attr1"]).To(Equal(4321))
	Expect(msg.Attrs["attr2"]).To(Equal("val2"))
	Expect(msg.Attrs["attr3"]).To(Equal("val3"))
	Expect(msg.StringAttrs["attr1"]).To(Equal("4321"))
	Expect(msg.StringAttrs["attr2"]).To(Equal("val2"))
	Expect(msg.StringAttrs["attr3"]).To(Equal("val3"))
}

func (s *MemLoggerSuite) TestHealthy(t sweet.T) {
	ml := newDefaultMemLogger()

	Expect(ml.Healthy()).To(BeFalse())

	ml.SetHealthy(true)

	Expect(ml.Healthy()).To(BeTrue())

	ml.SetHealthy(false)

	Expect(ml.Healthy()).To(BeFalse())
}
