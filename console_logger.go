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

func (l *ConsoleLogger) logf(level LogLevel, attrs map[string]interface{}, format string, a ...interface{}) error {
	msg := newMessage(l.base, level, attrs, format, a...)
	out, err := l.tpl.executeInternalMsg(msg, l.config.Colorize)
	if err != nil {
		return err
	}
	l.writer.Print(out + "\n")
	return nil
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

func (l *ConsoleLogger) Dbg(msg string) error {
	l.logf(LEVEL_DEBUG, nil, msg)
	return nil
}
func (l *ConsoleLogger) Dbgf(msg string, a ...interface{}) error {
	l.logf(LEVEL_DEBUG, nil, msg, a...)
	return nil
}
func (l *ConsoleLogger) Dbgm(m map[string]interface{}, msg string, a ...interface{}) error {
	l.logf(LEVEL_DEBUG, m, msg, a...)
	return nil
}

func (l *ConsoleLogger) Info(msg string) error {
	l.logf(LEVEL_INFO, nil, msg)
	return nil
}
func (l *ConsoleLogger) Infof(msg string, a ...interface{}) error {
	l.logf(LEVEL_INFO, nil, msg, a...)
	return nil
}
func (l *ConsoleLogger) Infom(m map[string]interface{}, msg string, a ...interface{}) error {
	l.logf(LEVEL_INFO, m, msg, a...)
	return nil
}

func (l *ConsoleLogger) Warn(msg string) error {
	l.logf(LEVEL_WARNING, nil, msg)
	return nil
}
func (l *ConsoleLogger) Warnf(msg string, a ...interface{}) error {
	l.logf(LEVEL_WARNING, nil, msg, a...)
	return nil
}
func (l *ConsoleLogger) Warnm(m map[string]interface{}, msg string, a ...interface{}) error {
	l.logf(LEVEL_WARNING, m, msg, a...)
	return nil
}

func (l *ConsoleLogger) Err(msg string) error {
	l.logf(LEVEL_ERROR, nil, msg)
	return nil
}
func (l *ConsoleLogger) Errf(msg string, a ...interface{}) error {
	l.logf(LEVEL_ERROR, nil, msg, a...)
	return nil
}
func (l *ConsoleLogger) Errm(m map[string]interface{}, msg string, a ...interface{}) error {
	l.logf(LEVEL_ERROR, m, msg, a...)
	return nil
}

func (l *ConsoleLogger) Fatal(msg string) error {
	l.logf(LEVEL_FATAL, nil, msg)
	return nil
}
func (l *ConsoleLogger) Fatalf(msg string, a ...interface{}) error {
	l.logf(LEVEL_FATAL, nil, msg, a...)
	return nil
}
func (l *ConsoleLogger) Fatalm(m map[string]interface{}, msg string, a ...interface{}) error {
	l.logf(LEVEL_FATAL, m, msg, a...)
	return nil
}
