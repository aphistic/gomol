package gomol

import (
	"fmt"
	"github.com/aphistic/golf"
)

type GelfLogger struct {
	base          *Base
	client        *golf.Client
	logger        *golf.Logger
	isInitialized bool

	Hostname string
	Port     int
}

func NewGelfLogger(hostname string, port int) *GelfLogger {
	l := &GelfLogger{
		Hostname: hostname,
		Port:     port,
	}
	return l
}

func (l *GelfLogger) getLogger() *golf.Logger {
	return l.logger
}

func (l *GelfLogger) getAttrs(m map[string]interface{}) map[string]interface{} {
	attrs := make(map[string]interface{}, 0)
	for key, val := range l.base.BaseAttrs {
		attrs[key] = val
	}
	if m != nil {
		for key, val := range m {
			attrs[key] = fmt.Sprintf("%v", val)
		}
	}

	return attrs
}

func (l *GelfLogger) SetBase(base *Base) {
	l.base = base
}

func (l *GelfLogger) IsInitialized() bool {
	return l.isInitialized
}

func (l *GelfLogger) InitLogger() error {
	c, err := golf.NewClient()
	if err != nil {
		return err
	}
	err = c.Dial(fmt.Sprintf("udp://%v:%v", l.Hostname, l.Port))
	if err != nil {
		return err
	}
	nl, err := c.NewLogger()
	if err != nil {
		return err
	}

	l.client = c
	l.logger = nl
	l.isInitialized = true

	return nil
}

func (l *GelfLogger) ShutdownLogger() error {
	err := l.client.Close()
	if err != nil {
		return err
	}

	l.isInitialized = false
	return nil
}

func (l *GelfLogger) Dbg(msg string) error {
	attrs := l.getAttrs(nil)
	l.getLogger().Dbgm(attrs, msg)
	return nil
}
func (l *GelfLogger) Dbgf(msg string, a ...interface{}) error {
	attrs := l.getAttrs(nil)
	l.getLogger().Dbgm(attrs, msg, a...)
	return nil
}
func (l *GelfLogger) Dbgm(m map[string]interface{}, msg string, a ...interface{}) error {
	attrs := l.getAttrs(m)
	l.getLogger().Dbgm(attrs, msg, a...)
	return nil
}

func (l *GelfLogger) Info(msg string) error {
	attrs := l.getAttrs(nil)
	l.getLogger().Infom(attrs, msg)
	return nil
}
func (l *GelfLogger) Infof(msg string, a ...interface{}) error {
	attrs := l.getAttrs(nil)
	l.getLogger().Infom(attrs, msg, a...)
	return nil
}
func (l *GelfLogger) Infom(m map[string]interface{}, msg string, a ...interface{}) error {
	attrs := l.getAttrs(m)
	l.getLogger().Infom(attrs, msg, a...)
	return nil
}

func (l *GelfLogger) Warn(msg string) error {
	attrs := l.getAttrs(nil)
	l.getLogger().Warnm(attrs, msg)
	return nil
}
func (l *GelfLogger) Warnf(msg string, a ...interface{}) error {
	attrs := l.getAttrs(nil)
	l.getLogger().Warnm(attrs, msg, a...)
	return nil
}
func (l *GelfLogger) Warnm(m map[string]interface{}, msg string, a ...interface{}) error {
	attrs := l.getAttrs(m)
	l.getLogger().Warnm(attrs, msg, a...)
	return nil
}

func (l *GelfLogger) Err(msg string) error {
	attrs := l.getAttrs(nil)
	l.getLogger().Errm(attrs, msg)
	return nil
}
func (l *GelfLogger) Errf(msg string, a ...interface{}) error {
	attrs := l.getAttrs(nil)
	l.getLogger().Errm(attrs, msg, a...)
	return nil
}
func (l *GelfLogger) Errm(m map[string]interface{}, msg string, a ...interface{}) error {
	attrs := l.getAttrs(m)
	l.getLogger().Errm(attrs, msg, a...)
	return nil
}

func (l *GelfLogger) Fatal(msg string) error {
	attrs := l.getAttrs(nil)
	l.getLogger().Emergm(attrs, msg)
	return nil
}
func (l *GelfLogger) Fatalf(msg string, a ...interface{}) error {
	attrs := l.getAttrs(nil)
	l.getLogger().Emergm(attrs, msg, a...)
	return nil
}
func (l *GelfLogger) Fatalm(m map[string]interface{}, msg string, a ...interface{}) error {
	attrs := l.getAttrs(m)
	l.getLogger().Emergm(attrs, msg, a...)
	return nil
}
