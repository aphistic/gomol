package gomol

import (
	"bytes"
	"fmt"
	"github.com/mgutz/ansi"
	"strings"
	"text/template"
	"time"
)

var colorDbg = ansi.ColorCode("cyan")
var colorInfo = ansi.ColorCode("green")
var colorWarn = ansi.ColorCode("yellow")
var colorErr = ansi.ColorCode("red")
var colorFatal = ansi.ColorCode("red+b")
var colorReset = ansi.ColorCode("default")

func tplColorDebug() string {
	return colorDbg
}
func tplColorInfo() string {
	return colorInfo
}
func tplColorWarn() string {
	return colorWarn
}
func tplColorError() string {
	return colorErr
}
func tplColorFatal() string {
	return colorFatal
}
func tplColorNone() string {
	return ""
}
func tplColorReset() string {
	return colorReset
}

type Template struct {
	tpls map[LogLevel]*template.Template
}

func getFuncMap(level LogLevel) template.FuncMap {
	fMap := template.FuncMap{
		"title": strings.Title,
		"lcase": strings.ToLower,
		"ucase": strings.ToUpper,
	}
	fMap["reset"] = tplColorReset
	switch level {
	case LEVEL_DEBUG:
		fMap["color"] = tplColorDebug
	case LEVEL_INFO:
		fMap["color"] = tplColorInfo
	case LEVEL_WARNING:
		fMap["color"] = tplColorWarn
	case LEVEL_ERROR:
		fMap["color"] = tplColorError
	case LEVEL_FATAL:
		fMap["color"] = tplColorFatal
	default:
		fMap["color"] = tplColorNone
		fMap["reset"] = tplColorNone
	}

	return fMap
}

func NewTemplate(tpl string) (*Template, error) {
	var levels = []LogLevel{LEVEL_NONE, LEVEL_DEBUG, LEVEL_INFO, LEVEL_WARNING, LEVEL_ERROR, LEVEL_FATAL}
	tpls := make(map[LogLevel]*template.Template, 0)
	for _, level := range levels {
		parsedTpl, err := template.New(getLevelName(level)).
			Funcs(getFuncMap(level)).
			Parse(tpl)
		if err != nil {
			return nil, err
		}
		tpls[level] = parsedTpl
	}

	newTpl := &Template{
		tpls: tpls,
	}

	return newTpl, nil
}

func (t *Template) executeInternalMsg(msg *message, colorize bool) (string, error) {
	tplMsg, err := newTemplateMsg(msg)
	if err != nil {
		return "", err
	}

	return t.Execute(tplMsg, colorize)
}

func (t *Template) Execute(msg *TemplateMsg, colorize bool) (string, error) {
	tplLevel := msg.level
	if !colorize {
		tplLevel = LEVEL_NONE
	}
	var buf bytes.Buffer
	err := t.tpls[tplLevel].Execute(&buf, msg)
	if err != nil {
		return "", nil
	}

	return buf.String(), nil
}

type TemplateMsg struct {
	timestamp time.Time
	level     LogLevel
	levelStr  string
	message   string
	attrs     map[string]interface{}
}

func newTemplateMsg(msg *message) (*TemplateMsg, error) {
	tplMsg := &TemplateMsg{
		attrs: make(map[string]interface{}, 0),
	}
	tplMsg.timestamp = msg.Timestamp
	tplMsg.message = fmt.Sprintf(msg.MsgFormat, msg.MsgParams...)
	tplMsg.level = msg.Level
	tplMsg.levelStr = getLevelName(msg.Level)
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

func (tm *TemplateMsg) Timestamp() time.Time {
	return tm.timestamp
}
func (tm *TemplateMsg) Level() string {
	return tm.levelStr
}
func (tm *TemplateMsg) Message() string {
	return tm.message
}
func (tm *TemplateMsg) Attrs() map[string]interface{} {
	return tm.attrs
}
