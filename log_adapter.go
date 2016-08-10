package gomol

/*
LogAdapter provides a way to easily override certain log attributes without
modifying the base attributes or specifying them for every log message.
*/
type LogAdapter struct {
	base  *Base
	attrs map[string]interface{}
}

func newLogAdapter(base *Base, attrs map[string]interface{}) *LogAdapter {
	newAttrs := attrs
	if attrs == nil {
		newAttrs = make(map[string]interface{})
	}

	return &LogAdapter{
		base:  base,
		attrs: newAttrs,
	}
}

// SetAttr sets the attribute key to value for this LogAdapter only
func (la *LogAdapter) SetAttr(key string, value interface{}) {
	la.attrs[key] = value
}

// GetAttr gets the attribute with the given key for this LogAdapter only. If the
// key doesn't exist on this LogAdapter it will return nil
func (la *LogAdapter) GetAttr(key string) interface{} {
	if val, ok := la.attrs[key]; ok {
		return val
	}
	return nil
}

// RemoveAttr removes the attribute key for this LogAdapter only
func (la *LogAdapter) RemoveAttr(key string) {
	delete(la.attrs, key)
}

// ClearAttrs removes all attributes for this LogAdapter only
func (la *LogAdapter) ClearAttrs() {
	la.attrs = make(map[string]interface{})
}

func (la *LogAdapter) log(level LogLevel, attrs map[string]interface{}, msg string, a ...interface{}) error {
	mergedAttrs := mergeAttrs(la.attrs, attrs)
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
func (la *LogAdapter) Dbgm(m map[string]interface{}, msg string, a ...interface{}) error {
	return la.Debugm(m, msg, a...)
}

// Debug logs msg to all added loggers at LogLevel.LEVEL_DEBUG
func (la *LogAdapter) Debug(msg string) error {
	return la.log(LEVEL_DEBUG, nil, msg)
}

/*
Debugf uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LEVEL_DEBUG
*/
func (la *LogAdapter) Debugf(msg string, a ...interface{}) error {
	return la.log(LEVEL_DEBUG, nil, msg, a...)
}

/*
Debugm uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LEVEL_DEBUG. It will also
merge all attributes passed in m with any attributes added to Base and include them
with the message if the Logger supports it.
*/
func (la *LogAdapter) Debugm(m map[string]interface{}, msg string, a ...interface{}) error {
	return la.log(LEVEL_DEBUG, m, msg, a...)
}

// Info logs msg to all added loggers at LogLevel.LEVEL_INFO
func (la *LogAdapter) Info(msg string) error {
	return la.log(LEVEL_INFO, nil, msg)
}

/*
Infof uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LEVEL_INFO
*/
func (la *LogAdapter) Infof(msg string, a ...interface{}) error {
	return la.log(LEVEL_INFO, nil, msg, a...)
}

/*
Infom uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LEVEL_INFO. It will also
merge all attributes passed in m with any attributes added to Base and include them
with the message if the Logger supports it.
*/
func (la *LogAdapter) Infom(m map[string]interface{}, msg string, a ...interface{}) error {
	return la.log(LEVEL_INFO, m, msg, a...)
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
func (la *LogAdapter) Warnm(m map[string]interface{}, msg string, a ...interface{}) error {
	return la.Warningm(m, msg, a...)
}

/*
Warning uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LEVEL_WARNING
*/
func (la *LogAdapter) Warning(msg string) error {
	return la.log(LEVEL_WARNING, nil, msg)
}

/*
Warningf uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LEVEL_WARNING
*/
func (la *LogAdapter) Warningf(msg string, a ...interface{}) error {
	return la.log(LEVEL_WARNING, nil, msg, a...)
}

/*
Warningm uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LEVEL_WARNING. It will also
merge all attributes passed in m with any attributes added to Base and include them
with the message if the Logger supports it.
*/
func (la *LogAdapter) Warningm(m map[string]interface{}, msg string, a ...interface{}) error {
	return la.log(LEVEL_WARNING, m, msg, a...)
}

func (la *LogAdapter) Err(msg string) error {
	return la.Error(msg)
}
func (la *LogAdapter) Errf(msg string, a ...interface{}) error {
	return la.Errorf(msg, a...)
}
func (la *LogAdapter) Errm(m map[string]interface{}, msg string, a ...interface{}) error {
	return la.Errorm(m, msg, a...)
}
func (la *LogAdapter) Error(msg string) error {
	return la.log(LEVEL_ERROR, nil, msg)
}
func (la *LogAdapter) Errorf(msg string, a ...interface{}) error {
	return la.log(LEVEL_ERROR, nil, msg, a...)
}
func (la *LogAdapter) Errorm(m map[string]interface{}, msg string, a ...interface{}) error {
	return la.log(LEVEL_ERROR, m, msg, a...)
}

func (la *LogAdapter) Fatal(msg string) error {
	return la.log(LEVEL_FATAL, nil, msg)
}
func (la *LogAdapter) Fatalf(msg string, a ...interface{}) error {
	return la.log(LEVEL_FATAL, nil, msg, a...)
}
func (la *LogAdapter) Fatalm(m map[string]interface{}, msg string, a ...interface{}) error {
	return la.log(LEVEL_FATAL, m, msg, a...)
}

func (la *LogAdapter) Die(exitCode int, msg string) {
	la.log(LEVEL_FATAL, nil, msg)
	la.base.ShutdownLoggers()
	curExiter.Exit(exitCode)
}
func (la *LogAdapter) Dief(exitCode int, msg string, a ...interface{}) {
	la.log(LEVEL_FATAL, nil, msg, a...)
	la.base.ShutdownLoggers()
	curExiter.Exit(exitCode)
}
func (la *LogAdapter) Diem(exitCode int, m map[string]interface{}, msg string, a ...interface{}) {
	la.log(LEVEL_FATAL, m, msg, a...)
	la.base.ShutdownLoggers()
	curExiter.Exit(exitCode)
}
