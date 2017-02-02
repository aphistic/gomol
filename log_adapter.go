package gomol

import "time"

/*
LogAdapter provides a way to easily override certain log attributes without
modifying the base attributes or specifying them for every log message.
*/
type LogAdapter struct {
	base  WrappableLogger
	attrs *Attrs
}

/*
WrappableLogger is an interface for a logger which can be wrapped by a LogAdapter.
his interface is implemented by both Base and LogAdapter itself so that adapters
can stack.
*/
type WrappableLogger interface {
	// LogWithTime will log a message at the provided level to all added loggers with the
	// timestamp set to the value of ts.
	LogWithTime(level LogLevel, ts time.Time, m *Attrs, msg string, a ...interface{}) error

	// Log will log a message at the provided level to all added loggers with the timestamp
	// set to the time Log was called.
	Log(level LogLevel, m *Attrs, msg string, a ...interface{}) error

	// ShutdownLoggers will run ShutdownLogger on each Logger in Base.  If an error occurs
	// while shutting down a Logger, the error will be returned and all the loggers that
	//were already shut down will remain shut down.
	ShutdownLoggers() error
}

func newLogAdapter(base *Base, attrs *Attrs) *LogAdapter {
	newAttrs := attrs
	if attrs == nil {
		attrs = NewAttrs()
	}

	return &LogAdapter{
		base:  base,
		attrs: attrs,
	}
}

// SetAttr sets the attribute key to value for this LogAdapter only
func (la *LogAdapter) SetAttr(key string, value interface{}) {
	la.attrs.SetAttr(key, value)
}

// GetAttr gets the attribute with the given key for this LogAdapter only. If the
// key doesn't exist on this LogAdapter it will return nil
func (la *LogAdapter) GetAttr(key string) interface{} {
	return la.attrs.GetAttr(key)
}

// RemoveAttr removes the attribute key for this LogAdapter only
func (la *LogAdapter) RemoveAttr(key string) {
	la.attrs.RemoveAttr(key)
}

// ClearAttrs removes all attributes for this LogAdapter only
func (la *LogAdapter) ClearAttrs() {
	la.attrs = NewAttrs()
}

// LogWithTime will log a message at the provided level to all loggers added
// to the Base associated with this LogAdapter. It is similar to Log except
// the timestamp will be set to the value of ts.
func (la *LogAdapter) LogWithTime(level LogLevel, ts time.Time, attrs *Attrs, msg string, a ...interface{}) error {
	mergedAttrs := la.attrs.clone()
	mergedAttrs.MergeAttrs(attrs)
	return la.base.LogWithTime(level, ts, mergedAttrs, msg, a...)
}

// Log will log a message at the provided level to all loggers added
// to the Base associated with this LogAdapter
func (la *LogAdapter) Log(level LogLevel, attrs *Attrs, msg string, a ...interface{}) error {
	mergedAttrs := la.attrs.clone()
	mergedAttrs.MergeAttrs(attrs)
	return la.base.Log(level, mergedAttrs, msg, a...)
}

// Dbg is a short-hand version of Debug
func (la *LogAdapter) Dbg(msg string) error {
	return la.Debug(msg)
}

// Dbgf is a short-hand version of Debugf
func (la *LogAdapter) Dbgf(msg string, a ...interface{}) error {
	return la.Debugf(msg, a...)
}

// Dbgm is a short-hand version of Debugm
func (la *LogAdapter) Dbgm(m *Attrs, msg string, a ...interface{}) error {
	return la.Debugm(m, msg, a...)
}

// Debug logs msg to all added loggers at LogLevel.LevelDebug
func (la *LogAdapter) Debug(msg string) error {
	return la.Log(LevelDebug, nil, msg)
}

/*
Debugf uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelDebug
*/
func (la *LogAdapter) Debugf(msg string, a ...interface{}) error {
	return la.Log(LevelDebug, nil, msg, a...)
}

/*
Debugm uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelDebug. It will also
merge all attributes passed in m with any attributes added to Base and include them
with the message if the Logger supports it.
*/
func (la *LogAdapter) Debugm(m *Attrs, msg string, a ...interface{}) error {
	return la.Log(LevelDebug, m, msg, a...)
}

// Info logs msg to all added loggers at LogLevel.LevelInfo
func (la *LogAdapter) Info(msg string) error {
	return la.Log(LevelInfo, nil, msg)
}

/*
Infof uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelInfo
*/
func (la *LogAdapter) Infof(msg string, a ...interface{}) error {
	return la.Log(LevelInfo, nil, msg, a...)
}

/*
Infom uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelInfo. It will also
merge all attributes passed in m with any attributes added to Base and include them
with the message if the Logger supports it.
*/
func (la *LogAdapter) Infom(m *Attrs, msg string, a ...interface{}) error {
	return la.Log(LevelInfo, m, msg, a...)
}

// Warn is a short-hand version of Warning
func (la *LogAdapter) Warn(msg string) error {
	return la.Warning(msg)
}

// Warnf is a short-hand version of Warningf
func (la *LogAdapter) Warnf(msg string, a ...interface{}) error {
	return la.Warningf(msg, a...)
}

// Warnm is a short-hand version of Warningm
func (la *LogAdapter) Warnm(m *Attrs, msg string, a ...interface{}) error {
	return la.Warningm(m, msg, a...)
}

/*
Warning uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelWarning
*/
func (la *LogAdapter) Warning(msg string) error {
	return la.Log(LevelWarning, nil, msg)
}

/*
Warningf uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelWarning
*/
func (la *LogAdapter) Warningf(msg string, a ...interface{}) error {
	return la.Log(LevelWarning, nil, msg, a...)
}

/*
Warningm uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelWarning. It will also
merge all attributes passed in m with any attributes added to Base and include them
with the message if the Logger supports it.
*/
func (la *LogAdapter) Warningm(m *Attrs, msg string, a ...interface{}) error {
	return la.Log(LevelWarning, m, msg, a...)
}

// Err is a short-hand version of Error
func (la *LogAdapter) Err(msg string) error {
	return la.Error(msg)
}

// Errf is a short-hand version of Errorf
func (la *LogAdapter) Errf(msg string, a ...interface{}) error {
	return la.Errorf(msg, a...)
}

// Errm is a short-hand version of Errorm
func (la *LogAdapter) Errm(m *Attrs, msg string, a ...interface{}) error {
	return la.Errorm(m, msg, a...)
}

/*
Error uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelError
*/
func (la *LogAdapter) Error(msg string) error {
	return la.Log(LevelError, nil, msg)
}

/*
Errorf uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelError
*/
func (la *LogAdapter) Errorf(msg string, a ...interface{}) error {
	return la.Log(LevelError, nil, msg, a...)
}

/*
Errorm uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelError. It will also
merge all attributes passed in m with any attributes added to Base and include them
with the message if the Logger supports it.
*/
func (la *LogAdapter) Errorm(m *Attrs, msg string, a ...interface{}) error {
	return la.Log(LevelError, m, msg, a...)
}

/*
Fatal uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelFatal
*/
func (la *LogAdapter) Fatal(msg string) error {
	return la.Log(LevelFatal, nil, msg)
}

/*
Fatalf uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelFatal
*/
func (la *LogAdapter) Fatalf(msg string, a ...interface{}) error {
	return la.Log(LevelFatal, nil, msg, a...)
}

/*
Fatalm uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelFatal. It will also
merge all attributes passed in m with any attributes added to Base and include them
with the message if the Logger supports it.
*/
func (la *LogAdapter) Fatalm(m *Attrs, msg string, a ...interface{}) error {
	return la.Log(LevelFatal, m, msg, a...)
}

// Die will log a message using Fatal, call ShutdownLoggers and then exit the application with the provided exit code.
func (la *LogAdapter) Die(exitCode int, msg string) {
	la.Log(LevelFatal, nil, msg)
	la.base.ShutdownLoggers()
	curExiter.Exit(exitCode)
}

// Dief will log a message using Fatalf, call ShutdownLoggers and then exit the application with the provided exit code.
func (la *LogAdapter) Dief(exitCode int, msg string, a ...interface{}) {
	la.Log(LevelFatal, nil, msg, a...)
	la.base.ShutdownLoggers()
	curExiter.Exit(exitCode)
}

// Diem will log a message using Fatalm, call ShutdownLoggers and then exit the application with the provided exit code.
func (la *LogAdapter) Diem(exitCode int, m *Attrs, msg string, a ...interface{}) {
	la.Log(LevelFatal, m, msg, a...)
	la.base.ShutdownLoggers()
	curExiter.Exit(exitCode)
}

// ShutdownLoggers will call the wrapped logger's ShutdownLoggers method.
func (la *LogAdapter) ShutdownLoggers() error {
	return la.base.ShutdownLoggers()
}
