package gomol

import (
	"fmt"
	"github.com/mgutz/ansi"
)

type ConsoleLoggerConfig struct {
	Colorize bool
}

type ConsoleLogger struct {
	base   *Base
	writer consoleWriter

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

func NewConsoleLogger(config *ConsoleLoggerConfig) *ConsoleLogger {
	l := &ConsoleLogger{
		writer: &ttyWriter{},
		config: config,
	}
	return l
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

func (l *ConsoleLogger) logf(level LogLevel, msg string, a ...interface{}) {
	printlog := printclean
	prefix := ""

	switch {
	case level == LEVEL_DEBUG:
		prefix = "[DEBUG]"
		if l.config.Colorize {
			printlog = printdbg
		}
	case level == LEVEL_INFO:
		prefix = "[INFO]"
		if l.config.Colorize {
			printlog = printinfo
		}
	case level == LEVEL_WARNING:
		prefix = "[WARN]"
		if l.config.Colorize {
			printlog = printwarn
		}
	case level == LEVEL_ERROR:
		prefix = "[ERROR]"
		if l.config.Colorize {
			printlog = printerr
		}
	case level == LEVEL_FATAL:
		prefix = "[FATAL]"
		if l.config.Colorize {
			printlog = printfatal
		}
	}

	formatted := fmt.Sprintf(prefix+" "+msg+"\n", a...)
	out := printlog(formatted)
	l.writer.Print(out)
}

func (l *ConsoleLogger) SetBase(base *Base) {
	l.base = base
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
	l.logf(LEVEL_DEBUG, msg)
	return nil
}
func (l *ConsoleLogger) Dbgf(msg string, a ...interface{}) error {
	l.logf(LEVEL_DEBUG, msg, a...)
	return nil
}
func (l *ConsoleLogger) Dbgm(m map[string]interface{}, msg string, a ...interface{}) error {
	l.logf(LEVEL_DEBUG, msg, a...)
	return nil
}

func (l *ConsoleLogger) Info(msg string) error {
	l.logf(LEVEL_INFO, msg)
	return nil
}
func (l *ConsoleLogger) Infof(msg string, a ...interface{}) error {
	l.logf(LEVEL_INFO, msg, a...)
	return nil
}
func (l *ConsoleLogger) Infom(m map[string]interface{}, msg string, a ...interface{}) error {
	l.logf(LEVEL_INFO, msg, a...)
	return nil
}

func (l *ConsoleLogger) Warn(msg string) error {
	l.logf(LEVEL_WARNING, msg)
	return nil
}
func (l *ConsoleLogger) Warnf(msg string, a ...interface{}) error {
	l.logf(LEVEL_WARNING, msg, a...)
	return nil
}
func (l *ConsoleLogger) Warnm(m map[string]interface{}, msg string, a ...interface{}) error {
	l.logf(LEVEL_WARNING, msg, a...)
	return nil
}

func (l *ConsoleLogger) Err(msg string) error {
	l.logf(LEVEL_ERROR, msg)
	return nil
}
func (l *ConsoleLogger) Errf(msg string, a ...interface{}) error {
	l.logf(LEVEL_ERROR, msg, a...)
	return nil
}
func (l *ConsoleLogger) Errm(m map[string]interface{}, msg string, a ...interface{}) error {
	l.logf(LEVEL_ERROR, msg, a...)
	return nil
}

func (l *ConsoleLogger) Fatal(msg string) error {
	l.logf(LEVEL_FATAL, msg)
	return nil
}
func (l *ConsoleLogger) Fatalf(msg string, a ...interface{}) error {
	l.logf(LEVEL_FATAL, msg, a...)
	return nil
}
func (l *ConsoleLogger) Fatalm(m map[string]interface{}, msg string, a ...interface{}) error {
	l.logf(LEVEL_FATAL, msg, a...)
	return nil
}
