package gomol

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"text/template"
	"time"

	"github.com/mgutz/ansi"
)

var colorDbg = ansi.ColorCode("cyan")
var colorInfo = ansi.ColorCode("green")
var colorWarn = ansi.ColorCode("yellow")
var colorErr = ansi.ColorCode("red")
var colorFatal = ansi.ColorCode("red+b")
var colorReset = ansi.ColorCode("reset")

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
func tplJson(data interface{}) (string, error) {
	json, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(json), nil
}

type Template struct {
	tpls map[LogLevel]*template.Template
}

func getFuncMap(level LogLevel) template.FuncMap {
	fMap := template.FuncMap{
		"title": strings.Title,
		"lcase": strings.ToLower,
		"ucase": strings.ToUpper,
		"json":  tplJson,
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
	tplLevel := msg.Level
	if !colorize {
		tplLevel = LEVEL_NONE
	}
	var buf bytes.Buffer
	err := t.tpls[tplLevel].Execute(&buf, msg)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

type TemplateMsg struct {
	Timestamp time.Time              `json:"timestamp"`
	Level     LogLevel               `json:"level"`
	LevelName string                 `json:"level_name"`
	Message   string                 `json:"message"`
	Attrs     map[string]interface{} `json:"attrs"`
}

func newTemplateMsg(msg *message) (*TemplateMsg, error) {
	if msg == nil {
		return nil, errors.New("msg cannot be nil")
	}

	tplMsg := &TemplateMsg{
		Attrs: make(map[string]interface{}, 0),
	}
	tplMsg.Timestamp = msg.Timestamp
	tplMsg.Message = msg.Msg
	tplMsg.Level = msg.Level
	tplMsg.LevelName = getLevelName(msg.Level)
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
