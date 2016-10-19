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
func tplJSON(data interface{}) (string, error) {
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
		"json":  tplJSON,
	}
	fMap["reset"] = tplColorReset
	switch level {
	case LevelDebug:
		fMap["color"] = tplColorDebug
	case LevelInfo:
		fMap["color"] = tplColorInfo
	case LevelWarning:
		fMap["color"] = tplColorWarn
	case LevelError:
		fMap["color"] = tplColorError
	case LevelFatal:
		fMap["color"] = tplColorFatal
	default:
		fMap["color"] = tplColorNone
		fMap["reset"] = tplColorNone
	}

	return fMap
}

func NewTemplate(tpl string) (*Template, error) {
	var levels = []LogLevel{LevelNone, LevelDebug, LevelInfo, LevelWarning, LevelError, LevelFatal}
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
	tplMsg, err := newTemplateMsgFromMessage(msg)
	if err != nil {
		return "", err
	}

	return t.Execute(tplMsg, colorize)
}

func (t *Template) Execute(msg *TemplateMsg, colorize bool) (string, error) {
	tplLevel := msg.Level
	if !colorize {
		tplLevel = LevelNone
	}
	var buf bytes.Buffer
	execTpl := t.tpls[tplLevel]
	if execTpl == nil {
		return "", ErrUnknownLevel
	}
	err := execTpl.Execute(&buf, msg)
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

func NewTemplateMsg(timestamp time.Time, level LogLevel, m map[string]interface{}, msg string) *TemplateMsg {
	msgAttrs := m
	if msgAttrs == nil {
		msgAttrs = make(map[string]interface{})
	}
	tplMsg := &TemplateMsg{
		Timestamp: timestamp,
		Message:   msg,
		Level:     level,
		LevelName: level.String(),
		Attrs:     msgAttrs,
	}
	return tplMsg
}

func newTemplateMsgFromMessage(msg *message) (*TemplateMsg, error) {
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

	tplAttrs := NewAttrs()
	if msg.Base != nil {
		tplAttrs.MergeAttrs(msg.Base.BaseAttrs)
	}
	tplAttrs.MergeAttrs(msg.Attrs)
	tplMsg.Attrs = tplAttrs.Attrs()

	return tplMsg, nil
}
