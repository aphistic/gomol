package gomol

import (
	"errors"
	"time"
)

type memLoggerConfig struct {
	FailInit     bool
	FailShutdown bool
}

func newMemLoggerConfig() *memLoggerConfig {
	return &memLoggerConfig{}
}

type memLogger struct {
	base   *Base
	config *memLoggerConfig

	Messages []*memMessage

	isInitialized bool
	isShutdown    bool
}

type memMessage struct {
	Timestamp time.Time
	Level     LogLevel
	Message   string
	Attrs     map[string]interface{}
}

func newMemLogger(config *memLoggerConfig) (*memLogger, error) {
	l := &memLogger{
		config:   config,
		Messages: make([]*memMessage, 0),
	}
	return l, nil
}

func newMemMessage() *memMessage {
	msg := &memMessage{
		Level:   LevelUnknown,
		Message: "",
		Attrs:   make(map[string]interface{}, 0),
	}
	return msg
}

func (l *memLogger) SetBase(base *Base) {
	l.base = base
}

func (l *memLogger) InitLogger() error {
	if l.config.FailInit {
		return errors.New("Init failed")
	}
	l.isInitialized = true
	l.isShutdown = false
	return nil
}
func (l *memLogger) IsInitialized() bool {
	return l.isInitialized
}
func (l *memLogger) ShutdownLogger() error {
	if l.config.FailShutdown {
		return errors.New("Shutdown failed")
	}
	l.isInitialized = false
	l.isShutdown = true
	return nil
}

func (l *memLogger) Logm(timestamp time.Time, level LogLevel, m map[string]interface{}, msg string) error {
	nm := newMemMessage()
	nm.Timestamp = timestamp
	nm.Level = level
	nm.Message = msg

	if l.base != nil {
		for k, v := range l.base.BaseAttrs {
			nm.Attrs[k] = v
		}
	}

	if m != nil {
		for k, v := range m {
			nm.Attrs[k] = v
		}
	}

	l.Messages = append(l.Messages, nm)

	return nil
}

func (l *memLogger) ClearMessages() {
	l.Messages = make([]*memMessage, 0)
}
