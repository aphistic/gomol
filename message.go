package gomol

import (
	"fmt"
	"math"
	"time"
)

type LogLevel int

const (
	LEVEL_UNKNOWN LogLevel = math.MaxInt64
	LEVEL_DEBUG   LogLevel = 7
	LEVEL_INFO    LogLevel = 6
	LEVEL_WARNING LogLevel = 4
	LEVEL_ERROR   LogLevel = 3
	LEVEL_FATAL   LogLevel = 2
	LEVEL_NONE    LogLevel = math.MinInt64
)

func (ll LogLevel) String() string {
	return getLevelName(ll)
}

func getLevelName(level LogLevel) string {
	switch level {
	case LEVEL_NONE:
		return "none"
	case LEVEL_DEBUG:
		return "debug"
	case LEVEL_INFO:
		return "info"
	case LEVEL_WARNING:
		return "warn"
	case LEVEL_ERROR:
		return "error"
	case LEVEL_FATAL:
		return "fatal"
	default:
		return "unknown"
	}
}

type message struct {
	Base      *Base
	Level     LogLevel
	Timestamp time.Time
	Attrs     map[string]interface{}
	Msg       string
}

func newMessage(timestamp time.Time,
	base *Base,
	level LogLevel,
	msgAttrs map[string]interface{},
	format string, va ...interface{}) *message {

	msgStr := format
	if len(va) > 0 {
		msgStr = fmt.Sprintf(format, va...)
	}

	nm := &message{
		Base:      base,
		Level:     level,
		Timestamp: timestamp,
		Attrs:     make(map[string]interface{}, len(msgAttrs)),
		Msg:       msgStr,
	}

	for msgKey, msgVal := range msgAttrs {
		nm.Attrs[msgKey] = msgVal
	}

	return nm
}

func mergeAttrs(baseAttrs map[string]interface{}, msgAttrs map[string]interface{}) map[string]interface{} {
	attrs := make(map[string]interface{}, len(baseAttrs)+len(msgAttrs))
	for attrKey, attrVal := range baseAttrs {
		attrs[attrKey] = attrVal
	}
	for attrKey, attrVal := range msgAttrs {
		attrs[attrKey] = attrVal
	}
	return attrs
}
