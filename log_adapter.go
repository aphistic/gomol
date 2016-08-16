package gomol

/*
LogAdapter provides a way to easily override certain log attributes without
modifying the base attributes or specifying them for every log message.
*/
type LogAdapter struct {
	base  *Base
	attrs *Attrs
}

func newLogAdapter(base *Base, attrs *Attrs) *LogAdapter {
	newAttrs := attrs
	if attrs == nil {
		newAttrs = NewAttrs()
	}

	return &LogAdapter{
		base:  base,
		attrs: newAttrs,
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

func (la *LogAdapter) log(level LogLevel, attrs *Attrs, msg string, a ...interface{}) error {
	mergedAttrs := la.attrs.clone()
	mergedAttrs.mergeAttrs(attrs)
	return la.base.log(level, mergedAttrs, msg, a...)
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
	return la.log(LevelDebug, nil, msg)
}

/*
Debugf uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelDebug
*/
func (la *LogAdapter) Debugf(msg string, a ...interface{}) error {
	return la.log(LevelDebug, nil, msg, a...)
}

/*
Debugm uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelDebug. It will also
merge all attributes passed in m with any attributes added to Base and include them
with the message if the Logger supports it.
*/
func (la *LogAdapter) Debugm(m *Attrs, msg string, a ...interface{}) error {
	return la.log(LevelDebug, m, msg, a...)
}

// Info logs msg to all added loggers at LogLevel.LevelInfo
func (la *LogAdapter) Info(msg string) error {
	return la.log(LevelInfo, nil, msg)
}

/*
Infof uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelInfo
*/
func (la *LogAdapter) Infof(msg string, a ...interface{}) error {
	return la.log(LevelInfo, nil, msg, a...)
}

/*
Infom uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelInfo. It will also
merge all attributes passed in m with any attributes added to Base and include them
with the message if the Logger supports it.
*/
func (la *LogAdapter) Infom(m *Attrs, msg string, a ...interface{}) error {
	return la.log(LevelInfo, m, msg, a...)
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
	return la.log(LevelWarning, nil, msg)
}

/*
Warningf uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelWarning
*/
func (la *LogAdapter) Warningf(msg string, a ...interface{}) error {
	return la.log(LevelWarning, nil, msg, a...)
}

/*
Warningm uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelWarning. It will also
merge all attributes passed in m with any attributes added to Base and include them
with the message if the Logger supports it.
*/
func (la *LogAdapter) Warningm(m *Attrs, msg string, a ...interface{}) error {
	return la.log(LevelWarning, m, msg, a...)
}

func (la *LogAdapter) Err(msg string) error {
	return la.Error(msg)
}
func (la *LogAdapter) Errf(msg string, a ...interface{}) error {
	return la.Errorf(msg, a...)
}
func (la *LogAdapter) Errm(m *Attrs, msg string, a ...interface{}) error {
	return la.Errorm(m, msg, a...)
}
func (la *LogAdapter) Error(msg string) error {
	return la.log(LevelError, nil, msg)
}
func (la *LogAdapter) Errorf(msg string, a ...interface{}) error {
	return la.log(LevelError, nil, msg, a...)
}
func (la *LogAdapter) Errorm(m *Attrs, msg string, a ...interface{}) error {
	return la.log(LevelError, m, msg, a...)
}

func (la *LogAdapter) Fatal(msg string) error {
	return la.log(LevelFatal, nil, msg)
}
func (la *LogAdapter) Fatalf(msg string, a ...interface{}) error {
	return la.log(LevelFatal, nil, msg, a...)
}
func (la *LogAdapter) Fatalm(m *Attrs, msg string, a ...interface{}) error {
	return la.log(LevelFatal, m, msg, a...)
}

func (la *LogAdapter) Die(exitCode int, msg string) {
	la.log(LevelFatal, nil, msg)
	la.base.ShutdownLoggers()
	curExiter.Exit(exitCode)
}
func (la *LogAdapter) Dief(exitCode int, msg string, a ...interface{}) {
	la.log(LevelFatal, nil, msg, a...)
	la.base.ShutdownLoggers()
	curExiter.Exit(exitCode)
}
func (la *LogAdapter) Diem(exitCode int, m *Attrs, msg string, a ...interface{}) {
	la.log(LevelFatal, m, msg, a...)
	la.base.ShutdownLoggers()
	curExiter.Exit(exitCode)
}
