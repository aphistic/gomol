package gomol

import (
	"errors"
	"fmt"
)

type MemLoggerConfig struct {
	FailInit     bool
	FailShutdown bool
}

func NewMemLoggerConfig() *MemLoggerConfig {
	return &MemLoggerConfig{}
}

type MemLogger struct {
	base   *Base
	config *MemLoggerConfig

	Messages []*MemMessage

	isInitialized bool
	isShutdown    bool
}

type MemMessage struct {
	Level   LogLevel
	Message string
	Attrs   map[string]interface{}
}

func NewMemLogger(config *MemLoggerConfig) (*MemLogger, error) {
	l := &MemLogger{
		config:   config,
		Messages: make([]*MemMessage, 0),
	}
	return l, nil
}

func NewMemMessage() *MemMessage {
	msg := &MemMessage{
		Level:   LEVEL_UNKNOWN,
		Message: "",
		Attrs:   make(map[string]interface{}, 0),
	}
	return msg
}

func (l *MemLogger) SetBase(base *Base) {
	l.base = base
}

func (l *MemLogger) InitLogger() error {
	if l.config.FailInit {
		return errors.New("Init failed")
	}
	l.isInitialized = true
	l.isShutdown = false
	return nil
}
func (l *MemLogger) IsInitialized() bool {
	return l.isInitialized
}
func (l *MemLogger) ShutdownLogger() error {
	if l.config.FailShutdown {
		return errors.New("Shutdown failed")
	}
	l.isInitialized = false
	l.isShutdown = true
	return nil
}

func (l *MemLogger) logm(level LogLevel, m map[string]interface{}, msg string, args ...interface{}) error {
	nm := NewMemMessage()
	nm.Level = level
	nm.Message = fmt.Sprintf(msg, args...)

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

func (l *MemLogger) ClearMessages() {
	l.Messages = make([]*MemMessage, 0)
}

func (l *MemLogger) Dbg(msg string) error {
	return l.logm(LEVEL_DEBUG, nil, msg)
}
func (l *MemLogger) Dbgf(msg string, args ...interface{}) error {
	return l.logm(LEVEL_DEBUG, nil, msg, args...)
}
func (l *MemLogger) Dbgm(m map[string]interface{}, msg string, args ...interface{}) error {
	return l.logm(LEVEL_DEBUG, m, msg, args...)
}

func (l *MemLogger) Info(msg string) error {
	return l.logm(LEVEL_INFO, nil, msg)
}
func (l *MemLogger) Infof(msg string, args ...interface{}) error {
	return l.logm(LEVEL_INFO, nil, msg, args...)
}
func (l *MemLogger) Infom(m map[string]interface{}, msg string, args ...interface{}) error {
	return l.logm(LEVEL_INFO, m, msg, args...)
}

func (l *MemLogger) Warn(msg string) error {
	return l.logm(LEVEL_WARNING, nil, msg)
}
func (l *MemLogger) Warnf(msg string, args ...interface{}) error {
	return l.logm(LEVEL_WARNING, nil, msg, args...)
}
func (l *MemLogger) Warnm(m map[string]interface{}, msg string, args ...interface{}) error {
	return l.logm(LEVEL_WARNING, m, msg, args...)
}

func (l *MemLogger) Err(msg string) error {
	return l.logm(LEVEL_ERROR, nil, msg)
}
func (l *MemLogger) Errf(msg string, args ...interface{}) error {
	return l.logm(LEVEL_ERROR, nil, msg, args...)
}
func (l *MemLogger) Errm(m map[string]interface{}, msg string, args ...interface{}) error {
	return l.logm(LEVEL_ERROR, m, msg, args...)
}

func (l *MemLogger) Fatal(msg string) error {
	return l.logm(LEVEL_FATAL, nil, msg)
}
func (l *MemLogger) Fatalf(msg string, args ...interface{}) error {
	return l.logm(LEVEL_FATAL, nil, msg, args...)
}
func (l *MemLogger) Fatalm(m map[string]interface{}, msg string, args ...interface{}) error {
	return l.logm(LEVEL_FATAL, m, msg, args...)
}
