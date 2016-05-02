package gomol

import (
	"bufio"
	"errors"
	"io"
	"sync"
)

type WriterLoggerConfig struct {
	/*
		The number of messages that will be buffered before flushing them to
		the file.
	*/
	BufferSize int
}

func NewWriterLoggerConfig() *WriterLoggerConfig {
	return &WriterLoggerConfig{
		BufferSize: 1000,
	}
}

type WriterLogger struct {
	base          *Base
	config        *WriterLoggerConfig
	writeLock     sync.Mutex
	buffer        []*message
	bufWriter     *bufio.Writer
	tpl           *Template
	isInitialized bool
}

func NewWriterLogger(w io.Writer, cfg *WriterLoggerConfig) (*WriterLogger, error) {
	if w == nil {
		return nil, errors.New("An io.Writer must be provided")
	}

	if cfg == nil {
		cfg = NewWriterLoggerConfig()
	}

	l := &WriterLogger{
		config:    cfg,
		buffer:    make([]*message, 0),
		bufWriter: bufio.NewWriter(w),
	}
	tpl, err := NewTemplate("[{{ucase .LevelName}}] {{.Message}}")
	if err != nil {
		return nil, err
	}
	l.SetTemplate(tpl)

	return l, nil
}

func (l *WriterLogger) SetBase(base *Base) {
	l.base = base
}

func (l *WriterLogger) SetTemplate(tpl *Template) error {
	if tpl == nil {
		return errors.New("A template must be provided")
	}
	l.tpl = tpl

	return nil
}

func (l *WriterLogger) InitLogger() error {
	l.isInitialized = true
	return nil
}
func (l *WriterLogger) IsInitialized() bool {
	return l.isInitialized
}
func (l *WriterLogger) ShutdownLogger() error {
	err := l.flushMessages()
	if err != nil {
		return err
	}

	l.isInitialized = false
	return nil
}

func (l *WriterLogger) flushMessages() error {
	if len(l.buffer) == 0 {
		return nil
	}

	sendMsgs := func() []*message {
		l.writeLock.Lock()
		defer l.writeLock.Unlock()

		retBuf := l.buffer
		l.buffer = make([]*message, 0)

		return retBuf
	}()

	for _, sendMsg := range sendMsgs {
		// Use colors for this because if they use colors in their
		// non-default template there's probably a reason.  This won't
		// affect any templates that don't include colors
		out, err := l.tpl.executeInternalMsg(sendMsg, true)
		if err != nil {
			// Need to make a channel or something to send logging
			// errors back to
		}
		l.bufWriter.WriteString(out + "\n")
	}
	l.bufWriter.Flush()

	return nil
}

func (l *WriterLogger) Logm(level LogLevel, m map[string]interface{}, msg string) error {
	newMsg := newMessage(l.base, level, m, msg)
	func() {
		l.writeLock.Lock()
		defer l.writeLock.Unlock()

		l.buffer = append(l.buffer, newMsg)
	}()

	if len(l.buffer) >= l.config.BufferSize {
		l.flushMessages()
	}

	return nil
}
