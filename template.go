package gomol

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

type Template struct {
	tpl *template.Template
}

func getFuncMap() template.FuncMap {
	return template.FuncMap{
		"title": strings.Title,
		"lcase": strings.ToLower,
		"ucase": strings.ToUpper,
	}
}

func NewTemplate(tpl string) (*Template, error) {
	parsedTpl, err := template.New("test").Funcs(getFuncMap()).Parse(tpl)
	if err != nil {
		return nil, err
	}

	newTpl := &Template{
		tpl: parsedTpl,
	}

	return newTpl, nil
}

func (t *Template) executeInternalMsg(msg *message) (string, error) {
	tplMsg, err := newTemplateMsg(msg)
	if err != nil {
		return "", err
	}

	return t.Execute(tplMsg)
}

func (t *Template) Execute(msg *TemplateMsg) (string, error) {
	var buf bytes.Buffer
	err := t.tpl.Execute(&buf, msg)
	if err != nil {
		return "", nil
	}

	return buf.String(), nil
}

type TemplateMsg struct {
	level   string
	message string
	attrs   map[string]interface{}
}

func newTemplateMsg(msg *message) (*TemplateMsg, error) {
	tplMsg := &TemplateMsg{
		attrs: make(map[string]interface{}, 0),
	}
	tplMsg.message = fmt.Sprintf(msg.MsgFormat, msg.MsgParams...)
	switch msg.Level {
	case LEVEL_NONE:
		tplMsg.level = "none"
	case LEVEL_DEBUG:
		tplMsg.level = "debug"
	case LEVEL_INFO:
		tplMsg.level = "info"
	case LEVEL_WARNING:
		tplMsg.level = "warn"
	case LEVEL_ERROR:
		tplMsg.level = "error"
	case LEVEL_FATAL:
		tplMsg.level = "fatal"
	default:
		tplMsg.level = "unknown"
	}
	if msg.Base != nil {
		for key, val := range msg.Base.BaseAttrs {
			tplMsg.attrs[key] = val
		}
	}
	for key, val := range msg.Attrs {
		tplMsg.attrs[key] = val
	}
	return tplMsg, nil
}

func (tm *TemplateMsg) Level() string {
	return tm.level
}
func (tm *TemplateMsg) Message() string {
	return tm.message
}
func (tm *TemplateMsg) Attrs() map[string]interface{} {
	return tm.attrs
}
