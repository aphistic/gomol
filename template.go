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

type TemplateMsg struct {
	Level   string
	Message string
	Attrs   map[string]interface{}
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

func newTemplateMsg(msg *message) (*TemplateMsg, error) {
	tplMsg := &TemplateMsg{
		Attrs: make(map[string]interface{}, 0),
	}
	tplMsg.Message = fmt.Sprintf(msg.MsgFormat, msg.MsgParams...)
	switch msg.Level {
	case LEVEL_NONE:
		tplMsg.Level = "none"
	case LEVEL_DEBUG:
		tplMsg.Level = "debug"
	case LEVEL_INFO:
		tplMsg.Level = "info"
	case LEVEL_WARNING:
		tplMsg.Level = "warn"
	case LEVEL_ERROR:
		tplMsg.Level = "error"
	case LEVEL_FATAL:
		tplMsg.Level = "fatal"
	default:
		tplMsg.Level = "unknown"
	}
	if msg.Base != nil {
		for key, val := range msg.Base.BaseAttrs {
			tplMsg.Attrs[key] = val
		}
	}
	for key, val := range msg.Attrs {
		tplMsg.Attrs[key] = val
	}
	return tplMsg, nil
}
