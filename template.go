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

/*
Template represents a Go template (See text/template) that can be used when
logging messages.  Template includes a few useful template functions:

	title
		Title cases a string
	lcase
		Lower cases a string
	ucase
		Upper cases a string
	json
		JSON marshals an object
	color
		Changes the color of any text after it to the log level's color
	reset
		Resets the current color to the default color
*/
type Template struct {
	tpls map[LogLevel]*template.Template
}

func getFuncMap(level LogLevel, forceReset bool) template.FuncMap {
	fMap := template.FuncMap{
		"title": strings.Title,
		"lcase": strings.ToLower,
		"ucase": strings.ToUpper,
		"json":  tplJSON,
		"reset": tplColorReset,
	}

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

		if !forceReset {
			fMap["reset"] = tplColorNone
		}
	}

	return fMap
}

// NewTemplate creates a new Template from the given string. An error is returned if the template fails to compile.
func NewTemplate(tpl string) (*Template, error) {
	return NewTemplateWithFuncMap(tpl, nil)
}

// NewTemplateWithFuncMap creates a new Template from the given string and a template FuncMap. The FuncMap available
// to the template during evaluation will also include the default values, if not overridden. An error is returned
// if the template fails to compile.
func NewTemplateWithFuncMap(tpl string, funcMap template.FuncMap) (*Template, error) {
	var levels = []LogLevel{LevelNone, LevelDebug, LevelInfo, LevelWarning, LevelError, LevelFatal}
	tpls := make(map[LogLevel]*template.Template, 0)
	for _, level := range levels {
		// If color is overridden, we need to ensure that {{reset}} resets for all levels.
		_, forceReset := funcMap["color"]
		fMap := getFuncMap(level, forceReset)
		for name, f := range funcMap {
			fMap[name] = f
		}

		parsedTpl, err := template.New(getLevelName(level)).
			Funcs(fMap).
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

func (t *Template) executeInternalMsg(msg *Message, colorize bool) (string, error) {
	tplMsg, err := newTemplateMsgFromMessage(msg)
	if err != nil {
		return "", err
	}

	return t.Execute(tplMsg, colorize)
}

// Execute takes a TemplateMsg and applies it to the Go template.  If colorize is true the template
// will insert ANSI color codes within the resulting string.
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

// TemplateMsg represents the parts of a message required to render a template
type TemplateMsg struct {
	Timestamp time.Time              `json:"timestamp"`
	Level     LogLevel               `json:"level"`
	LevelName string                 `json:"level_name"`
	Message   string                 `json:"message"`
	Attrs     map[string]interface{} `json:"attrs"`
}

// NewTemplateMsg will create a new TemplateMsg with values from the given parameters
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

func newTemplateMsgFromMessage(msg *Message) (*TemplateMsg, error) {
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
	if msg.base != nil {
		tplAttrs.MergeAttrs(msg.base.BaseAttrs)
	}
	tplAttrs.MergeAttrs(msg.Attrs)
	tplMsg.Attrs = tplAttrs.Attrs()

	return tplMsg, nil
}
