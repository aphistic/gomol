package gomol

import (
	"fmt"
	"math"
	"time"
)

type LogLevel int

const (
	LevelUnknown LogLevel = math.MaxInt64
	LevelDebug   LogLevel = 7
	LevelInfo    LogLevel = 6
	LevelWarning LogLevel = 4
	LevelError   LogLevel = 3
	LevelFatal   LogLevel = 2
	LevelNone    LogLevel = math.MinInt64
)

func (ll LogLevel) String() string {
	return getLevelName(ll)
}

func getLevelName(level LogLevel) string {
	switch level {
	case LevelNone:
		return "none"
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarning:
		return "warn"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	default:
		return "unknown"
	}
}

type message struct {
	Base      *Base
	Level     LogLevel
	Timestamp time.Time
	Attrs     *Attrs
	Msg       string
}

func newMessage(timestamp time.Time,
	base *Base,
	level LogLevel,
	msgAttrs *Attrs,
	format string, va ...interface{}) *message {

	msgStr := format
	if len(va) > 0 {
		msgStr = fmt.Sprintf(format, va...)
	}

	var attrs *Attrs
	if msgAttrs != nil {
		attrs = msgAttrs
	} else {
		attrs = NewAttrs()
	}

	nm := &message{
		Base:      base,
		Level:     level,
		Timestamp: timestamp,
		Attrs:     attrs,
		Msg:       msgStr,
	}

	return nm
}
