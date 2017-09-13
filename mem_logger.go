package gomol

import (
	"bytes"
	"errors"
	"sync"
	"text/template"
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

	messageLock sync.Mutex
	messages    []*memMessage

	isInitialized bool
	isShutdown    bool

	tpl *template.Template
}

type memMessage struct {
	Timestamp   time.Time
	Level       LogLevel
	Message     string
	Attrs       map[string]interface{}
	StringAttrs map[string]string
}

func newMemLogger(config *memLoggerConfig) (*memLogger, error) {
	valTpl, err := template.New("memTpl").Parse("{{ . }}")
	if err != nil {
		return nil, err
	}

	l := &memLogger{
		config: config,

		messages: make([]*memMessage, 0),

		tpl: valTpl,
	}
	return l, nil
}

func newMemMessage() *memMessage {
	msg := &memMessage{
		Level:       LevelDebug,
		Message:     "",
		Attrs:       make(map[string]interface{}),
		StringAttrs: make(map[string]string),
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
		for k, v := range l.base.BaseAttrs.Attrs() {
			nm.Attrs[k] = v

			buf := bytes.NewBufferString("")
			err := l.tpl.Execute(buf, v)
			if err != nil {
				return err
			}
			nm.StringAttrs[k] = buf.String()
		}
	}

	if m != nil {
		for k, v := range m {
			nm.Attrs[k] = v

			buf := bytes.NewBufferString("")
			err := l.tpl.Execute(buf, v)
			if err != nil {
				return err
			}
			nm.StringAttrs[k] = buf.String()
		}
	}

	l.messageLock.Lock()
	l.messages = append(l.messages, nm)
	l.messageLock.Unlock()

	return nil
}

func (l *memLogger) Messages() []*memMessage {
	l.messageLock.Lock()
	defer l.messageLock.Unlock()

	return l.messages
}

func (l *memLogger) ClearMessages() {
	l.messageLock.Lock()
	defer l.messageLock.Unlock()

	l.messages = make([]*memMessage, 0)
}
