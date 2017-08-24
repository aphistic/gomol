package gomol

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// LogLevel represents the level a message is logged at.
type LogLevel int

const (
	// LevelDebug designates messages that are most useful when debugging applications.
	LevelDebug LogLevel = 7
	// LevelInfo designates messages that show application progression
	LevelInfo LogLevel = 6
	// LevelWarning designates messages that could potentially cause problems
	LevelWarning LogLevel = 4
	// LevelError designates error messages that don't stop the application from running
	LevelError LogLevel = 3
	// LevelFatal designates messages for severe errors where the application cannot continue
	LevelFatal LogLevel = 2

	// LevelNone is used when configuring log levels to disable all log levels
	LevelNone LogLevel = math.MinInt32
)

func (ll *LogLevel) MarshalJSON() ([]byte, error) {
	jsonLevel := getLevelName(*ll)
	jsonData := append([]byte{'"'}, append([]byte(jsonLevel), '"')...)
	return jsonData, nil
}

func (ll *LogLevel) UnmarshalJSON(data []byte) error {
	levelStr := strings.Trim(string(data), `"`)

	jsonLevel, err := ToLogLevel(levelStr)
	if err != nil {
		return err
	}
	*ll = jsonLevel
	return nil
}

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

// Message holds the information for a log message
type Message struct {
	base      *Base
	Level     LogLevel
	Timestamp time.Time
	Attrs     *Attrs
	Msg       string
}

func newMessage(timestamp time.Time,
	base *Base,
	level LogLevel,
	msgAttrs *Attrs,
	format string, va ...interface{}) *Message {

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

	nm := &Message{
		base:      base,
		Level:     level,
		Timestamp: timestamp,
		Attrs:     attrs,
		Msg:       msgStr,
	}

	return nm
}
