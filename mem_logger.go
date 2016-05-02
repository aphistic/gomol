package gomol

import "errors"

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

func (l *MemLogger) Logm(level LogLevel, m map[string]interface{}, msg string) error {
	nm := NewMemMessage()
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

func (l *MemLogger) ClearMessages() {
	l.Messages = make([]*MemMessage, 0)
}
