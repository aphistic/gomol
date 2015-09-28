package gomol

import (
	"math"
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

type message struct {
	Base      *Base
	Level     LogLevel
	Attrs     map[string]interface{}
	MsgFormat string
	MsgParams []interface{}
}

func newMessage(base *Base,
	level LogLevel,
	msgAttrs map[string]interface{},
	format string, va ...interface{}) *message {

	nm := &message{
		Base:      base,
		Level:     level,
		Attrs:     make(map[string]interface{}, len(msgAttrs)),
		MsgFormat: format,
		MsgParams: va,
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
