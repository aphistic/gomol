package gomol

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type FallbackLoggerSuite struct{}

func (s *FallbackLoggerSuite) TestSetInitialized(t sweet.T) {
	b := NewBase()

	ml := newDefaultMemLogger()
	err := ml.InitLogger()
	Expect(err).To(BeNil())
	Expect(ml.IsInitialized()).To(BeTrue())

	err = b.SetFallbackLogger(ml)
	Expect(err).To(BeNil())

	Expect(ml.IsInitialized()).To(BeTrue())
}

func (s *FallbackLoggerSuite) TestSetUninitialized(t sweet.T) {
	b := NewBase()

	ml := newDefaultMemLogger()

	Expect(ml.IsInitialized()).ToNot(BeTrue())

	err := b.SetFallbackLogger(ml)
	Expect(err).To(BeNil())

	Expect(ml.IsInitialized()).To(BeTrue())

	Expect(b.fallbackLogger).To(Equal(ml))
}

func (s *FallbackLoggerSuite) TestSetUninitializedFailToInitialize(t sweet.T) {
	b := NewBase()

	ml := newDefaultMemLogger()
	ml.config.FailInit = true

	Expect(ml.IsInitialized()).ToNot(BeTrue())

	err := b.SetFallbackLogger(ml)
	Expect(err).ToNot(BeNil())

	Expect(b.fallbackLogger).To(BeNil())
}

func (s *FallbackLoggerSuite) TestShutdownLoggerWhenReplacingWithNil(t sweet.T) {
	b := NewBase()

	ml := newDefaultMemLogger()

	err := b.SetFallbackLogger(ml)
	Expect(err).To(BeNil())
	Expect(ml.IsInitialized()).To(BeTrue())

	err = b.SetFallbackLogger(nil)
	Expect(err).To(BeNil())
	Expect(ml.IsInitialized()).ToNot(BeTrue())
}

func (s *FallbackLoggerSuite) TestShutdownLoggerWhenReplacingWithNewLogger(t sweet.T) {
	b := NewBase()

	ml := newDefaultMemLogger()

	err := b.SetFallbackLogger(ml)
	Expect(err).To(BeNil())
	Expect(ml.IsInitialized()).To(BeTrue())

	ml2 := newDefaultMemLogger()
	err = b.SetFallbackLogger(ml2)
	Expect(err).To(BeNil())
	Expect(ml2.IsInitialized()).To(BeTrue())
	Expect(ml.IsInitialized()).ToNot(BeTrue())
}

func (s *FallbackLoggerSuite) TestNoReplaceOldLoggerIfNewFailsToInitialize(t sweet.T) {
	b := NewBase()

	ml := newDefaultMemLogger()

	err := b.SetFallbackLogger(ml)
	Expect(err).To(BeNil())
	Expect(ml.IsInitialized()).To(BeTrue())

	ml2 := newDefaultMemLogger()
	ml2.config.FailInit = true
	err = b.SetFallbackLogger(ml2)
	Expect(err).ToNot(BeNil())
	Expect(ml2.IsInitialized()).ToNot(BeTrue())
	Expect(ml.IsInitialized()).To(BeTrue())
	Expect(b.fallbackLogger).To(Equal(ml))
}

func (s *FallbackLoggerSuite) TestLogToFallbackWithZeroLoggers(t sweet.T) {
	b := NewBase()
	err := b.InitLoggers()
	Expect(err).To(BeNil())

	fb := newDefaultMemLogger()
	fb.SetHealthy(true)

	err = b.SetFallbackLogger(fb)
	Expect(err).To(BeNil())
	Expect(fb.IsInitialized()).To(BeTrue())

	b.Info("message 1")
	b.Flush()

	Expect(fb.Messages()).To(HaveLen(1))
}

func (s *FallbackLoggerSuite) TestLogToFallbackOnUnhealthy(t sweet.T) {
	b := NewBase()
	ml := newDefaultMemLogger()
	ml.SetHealthy(true)
	err := b.AddLogger(ml)
	Expect(err).To(BeNil())
	err = b.InitLoggers()
	Expect(err).To(BeNil())

	fb := newDefaultMemLogger()
	fb.SetHealthy(true)
	err = b.SetFallbackLogger(fb)
	Expect(err).To(BeNil())

	b.Info("message 1")
	b.Flush()

	Expect(ml.Messages()).To(HaveLen(1))
	Expect(fb.Messages()).To(HaveLen(0))

	ml.SetHealthy(false)

	b.Info("message 2")
	b.Flush()

	Expect(ml.Messages()).To(HaveLen(2))
	Expect(fb.Messages()).To(HaveLen(1))

	ml.SetHealthy(true)

	b.Info("message 3")
	b.Flush()

	Expect(ml.Messages()).To(HaveLen(3))
	Expect(fb.Messages()).To(HaveLen(1))
}

func (s *FallbackLoggerSuite) TestLogToFallbackMultipleUnhealthy(t sweet.T) {
	b := NewBase()
	ml1 := newDefaultMemLogger() // Starts unhealthy
	ml2 := newDefaultMemLogger() // Starts unhealthy
	err := b.AddLogger(ml1)
	Expect(err).To(BeNil())
	err = b.AddLogger(ml2)
	Expect(err).To(BeNil())
	err = b.InitLoggers()
	Expect(err).To(BeNil())

	fb := newDefaultMemLogger()
	fb.SetHealthy(true)
	err = b.SetFallbackLogger(fb)
	Expect(err).To(BeNil())

	b.Info("message 1")
	b.Flush()

	Expect(ml1.Messages()).To(HaveLen(1))
	Expect(ml2.Messages()).To(HaveLen(1))
	Expect(fb.Messages()).To(HaveLen(1))

	ml1.SetHealthy(true)

	b.Info("message 2")
	b.Flush()

	Expect(ml1.Messages()).To(HaveLen(2))
	Expect(ml2.Messages()).To(HaveLen(2))
	Expect(fb.Messages()).To(HaveLen(2))

	ml2.SetHealthy(true)

	b.Info("message 3")
	b.Flush()

	Expect(ml1.Messages()).To(HaveLen(3))
	Expect(ml2.Messages()).To(HaveLen(3))
	Expect(fb.Messages()).To(HaveLen(2))
}
