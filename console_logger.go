package gomol

import (
	"errors"
	"fmt"

	"github.com/mgutz/ansi"
)

type ConsoleLoggerConfig struct {
	Colorize bool
}

type ConsoleLogger struct {
	base          *Base
	writer        consoleWriter
	tpl           *Template
	isInitialized bool
	config        *ConsoleLoggerConfig
}
type consoleWriter interface {
	Print(msg string)
}

// TTY writer for logging to the actual console
type ttyWriter struct {
}

func (w *ttyWriter) Print(msg string) {
	fmt.Print(msg)
}

func NewConsoleLoggerConfig() *ConsoleLoggerConfig {
	return &ConsoleLoggerConfig{
		Colorize: true,
	}
}

func NewConsoleLogger(config *ConsoleLoggerConfig) (*ConsoleLogger, error) {
	l := &ConsoleLogger{
		writer: &ttyWriter{},
		config: config,
	}
	tpl, err := NewTemplate("[{{color}}{{ucase .LevelName}}{{reset}}] {{.Message}}")
	if err != nil {
		return nil, err
	}
	l.tpl = tpl
	return l, nil
}

var printclean = func(msg string) string {
	return msg
}
var printdbg = ansi.ColorFunc("cyan")
var printinfo = ansi.ColorFunc("green")
var printwarn = ansi.ColorFunc("yellow")
var printerr = ansi.ColorFunc("red")
var printfatal = ansi.ColorFunc("red+b")

func (l *ConsoleLogger) setWriter(w consoleWriter) {
	l.writer = w
}

func (l *ConsoleLogger) SetBase(base *Base) {
	l.base = base
}

func (l *ConsoleLogger) SetTemplate(tpl *Template) error {
	if tpl == nil {
		return errors.New("A template must be provided")
	}
	l.tpl = tpl

	return nil
}

func (l *ConsoleLogger) InitLogger() error {
	l.isInitialized = true
	return nil
}
func (l *ConsoleLogger) IsInitialized() bool {
	return l.isInitialized
}

func (l *ConsoleLogger) ShutdownLogger() error {
	l.isInitialized = false
	return nil
}

func (l *ConsoleLogger) Logm(level LogLevel, attrs map[string]interface{}, msg string) error {
	nMsg := newMessage(l.base, level, attrs, msg)
	out, err := l.tpl.executeInternalMsg(nMsg, l.config.Colorize)
	if err != nil {
		return err
	}
	l.writer.Print(out + "\n")
	return nil
}
