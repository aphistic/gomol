package gomol

import (
	"fmt"
	"github.com/mgutz/ansi"
)

type ConsoleLogger struct {
	base     *Base
	Colorize bool
}

func NewConsoleLogger() *ConsoleLogger {
	l := &ConsoleLogger{
		Colorize: true,
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

func (l *ConsoleLogger) logf(level int, msg string, a ...interface{}) {
	printlog := printclean
	prefix := ""

	switch {
	case level == levelDbg:
		prefix = "[DEBUG]"
		if l.Colorize {
			printlog = printdbg
		}
	case level == levelInfo:
		prefix = "[INFO]"
		if l.Colorize {
			printlog = printinfo
		}
	case level == levelWarn:
		prefix = "[WARN]"
		if l.Colorize {
			printlog = printwarn
		}
	case level == levelError:
		prefix = "[ERROR]"
		if l.Colorize {
			printlog = printerr
		}
	case level == levelFatal:
		prefix = "[FATAL]"
		if l.Colorize {
			printlog = printfatal
		}
	}

	formatted := fmt.Sprintf(prefix+" "+msg+"\n", a...)
	out := printlog(formatted)
	fmt.Print(out)
}

func (l *ConsoleLogger) SetBase(base *Base) {
	l.base = base
}

func (l *ConsoleLogger) InitLogger() error {
	return nil
}

func (l *ConsoleLogger) ShutdownLogger() error {
	return nil
}

func (l *ConsoleLogger) Dbg(msg string) error {
	l.logf(levelDbg, msg)
	return nil
}
func (l *ConsoleLogger) Dbgf(msg string, a ...interface{}) error {
	l.logf(levelDbg, msg, a...)
	return nil
}
func (l *ConsoleLogger) Dbgm(m map[string]interface{}, msg string, a ...interface{}) error {
	l.logf(levelDbg, msg, a...)
	return nil
}

func (l *ConsoleLogger) Info(msg string) error {
	l.logf(levelInfo, msg)
	return nil
}
func (l *ConsoleLogger) Infof(msg string, a ...interface{}) error {
	l.logf(levelInfo, msg, a...)
	return nil
}
func (l *ConsoleLogger) Infom(m map[string]interface{}, msg string, a ...interface{}) error {
	l.logf(levelInfo, msg, a...)
	return nil
}

func (l *ConsoleLogger) Warn(msg string) error {
	l.logf(levelWarn, msg)
	return nil
}
func (l *ConsoleLogger) Warnf(msg string, a ...interface{}) error {
	l.logf(levelWarn, msg, a...)
	return nil
}
func (l *ConsoleLogger) Warnm(m map[string]interface{}, msg string, a ...interface{}) error {
	l.logf(levelWarn, msg, a...)
	return nil
}

func (l *ConsoleLogger) Err(msg string) error {
	l.logf(levelError, msg)
	return nil
}
func (l *ConsoleLogger) Errf(msg string, a ...interface{}) error {
	l.logf(levelError, msg, a...)
	return nil
}
func (l *ConsoleLogger) Errm(m map[string]interface{}, msg string, a ...interface{}) error {
	l.logf(levelError, msg, a...)
	return nil
}

func (l *ConsoleLogger) Fatal(msg string) error {
	l.logf(levelFatal, msg)
	return nil
}
func (l *ConsoleLogger) Fatalf(msg string, a ...interface{}) error {
	l.logf(levelFatal, msg, a...)
	return nil
}
func (l *ConsoleLogger) Fatalm(m map[string]interface{}, msg string, a ...interface{}) error {
	l.logf(levelFatal, msg, a...)
	return nil
}
