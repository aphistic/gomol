package gomol

import (
	"fmt"
	"math"
	"strings"
	"time"
)

type LogLevel int

const (
	LevelDebug   LogLevel = 7
	LevelInfo    LogLevel = 6
	LevelWarning LogLevel = 4
	LevelError   LogLevel = 3
	LevelFatal   LogLevel = 2
	LevelNone    LogLevel = math.MinInt64
)

// ToLogLevel will take a string and return the appropriate log level for
// the string if known.  If the string is not recognized it will return
// an ErrUnknownLevel error.
func ToLogLevel(level string) (LogLevel, error) {
	lowLevel := strings.ToLower(level)

	switch lowLevel {
	case "dbg":
		fallthrough
	case "debug":
		return LevelDebug, nil
	case "info":
		return LevelInfo, nil
	case "warn":
		fallthrough
	case "warning":
		return LevelWarning, nil
	case "err":
		fallthrough
	case "error":
		return LevelError, nil
	case "fatal":
		return LevelFatal, nil
	case "none":
		return LevelNone, nil
	}

	return 0, ErrUnknownLevel
}

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
