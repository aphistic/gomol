package gomol

import (
	"fmt"
	"github.com/segmentio/go-loggly"
)

var logglyClients map[string]*loggly.Client

func init() {
	logglyClients = make(map[string]*loggly.Client, 0)
}

type LogglyLogger struct {
	base          *Base
	isInitialized bool
	Token         string
}

func NewLogglyLogger(token string) *LogglyLogger {
	l := &LogglyLogger{
		Token: token,
	}
	return l
}

func (l *LogglyLogger) getClient() *loggly.Client {
	c := logglyClients[l.Token]
	if c == nil {
		c = loggly.New(l.Token)
		c.Level = loggly.DEBUG
		logglyClients[l.Token] = c
	}
	return c
}

func (l *LogglyLogger) getFacility(m map[string]interface{}) string {
	if m != nil {
		if facility, ok := m["facility"]; ok {
			return fmt.Sprint(facility)
		}
	}
	if facility, ok := l.base.BaseAttrs["facility"]; ok {
		return fmt.Sprint(facility)
	}
	return ""
}

func (l *LogglyLogger) getMsg(m map[string]interface{}, msg string, a ...interface{}) loggly.Message {
	lm := loggly.Message{"message": fmt.Sprintf(msg, a...)}

	for key, val := range l.base.BaseAttrs {
		lm[key] = val
	}
	if m != nil {
		for key, val := range m {
			lm[key] = fmt.Sprintf("%v", val)
		}
	}

	return lm
}

func (l *LogglyLogger) SetBase(base *Base) {
	l.base = base
}

func (l *LogglyLogger) IsInitialized() bool {
	return l.isInitialized
}

func (l *LogglyLogger) InitLogger() error {
	l.isInitialized = true
	return nil
}

func (l *LogglyLogger) ShutdownLogger() error {
	c := logglyClients[l.Token]
	if c != nil {
		err := c.Flush()
		if err != nil {
			return err
		}

		delete(logglyClients, l.Token)
	}

	l.isInitialized = false

	return nil
}

func (l *LogglyLogger) Dbg(msg string) error {
	lm := l.getMsg(nil, msg)
	l.getClient().Debug(l.getFacility(nil), lm)
	return nil
}
func (l *LogglyLogger) Dbgf(msg string, a ...interface{}) error {
	lm := l.getMsg(nil, msg, a...)
	l.getClient().Debug(l.getFacility(nil), lm)
	return nil
}
func (l *LogglyLogger) Dbgm(m map[string]interface{}, msg string, a ...interface{}) error {
	lm := l.getMsg(m, msg, a...)
	l.getClient().Debug(l.getFacility(m), lm)
	return nil
}

func (l *LogglyLogger) Info(msg string) error {
	lm := l.getMsg(nil, msg)
	l.getClient().Info(l.getFacility(nil), lm)
	return nil
}
func (l *LogglyLogger) Infof(msg string, a ...interface{}) error {
	lm := l.getMsg(nil, msg, a...)
	l.getClient().Info(l.getFacility(nil), lm)
	return nil
}
func (l *LogglyLogger) Infom(m map[string]interface{}, msg string, a ...interface{}) error {
	lm := l.getMsg(m, msg, a...)
	l.getClient().Info(l.getFacility(m), lm)
	return nil
}

func (l *LogglyLogger) Warn(msg string) error {
	lm := l.getMsg(nil, msg)
	l.getClient().Warn(l.getFacility(nil), lm)
	return nil
}
func (l *LogglyLogger) Warnf(msg string, a ...interface{}) error {
	lm := l.getMsg(nil, msg, a...)
	l.getClient().Warn(l.getFacility(nil), lm)
	return nil
}
func (l *LogglyLogger) Warnm(m map[string]interface{}, msg string, a ...interface{}) error {
	lm := l.getMsg(m, msg, a...)
	l.getClient().Warn(l.getFacility(m), lm)
	return nil
}

func (l *LogglyLogger) Err(msg string) error {
	lm := l.getMsg(nil, msg)
	l.getClient().Error(l.getFacility(nil), lm)
	return nil
}
func (l *LogglyLogger) Errf(msg string, a ...interface{}) error {
	lm := l.getMsg(nil, msg, a...)
	l.getClient().Error(l.getFacility(nil), lm)
	return nil
}
func (l *LogglyLogger) Errm(m map[string]interface{}, msg string, a ...interface{}) error {
	lm := l.getMsg(m, msg, a...)
	l.getClient().Error(l.getFacility(m), lm)
	return nil
}

func (l *LogglyLogger) Fatal(msg string) error {
	lm := l.getMsg(nil, msg)
	l.getClient().Critical(l.getFacility(nil), lm)
	return nil
}
func (l *LogglyLogger) Fatalf(msg string, a ...interface{}) error {
	lm := l.getMsg(nil, msg, a...)
	l.getClient().Critical(l.getFacility(nil), lm)
	return nil
}
func (l *LogglyLogger) Fatalm(m map[string]interface{}, msg string, a ...interface{}) error {
	lm := l.getMsg(m, msg, a...)
	l.getClient().Critical(l.getFacility(m), lm)
	return nil
}
