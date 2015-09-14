package gomol

import (
	"fmt"
)

type message struct {
	Base  *Base
	Level int
	Attrs map[string]interface{}
	Msg   string
}

func newMessage(base *Base,
	level int,
	msgAttrs map[string]interface{},
	format string, va ...interface{}) *message {

	nm := &message{
		Base:  base,
		Level: level,
		Attrs: make(map[string]interface{}, len(msgAttrs)),
		Msg:   fmt.Sprintf(format, va...),
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
