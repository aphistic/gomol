package gomol

import (
	"os"
)

/*
Base holds an instance of all information needed for logging.  It is possible
to create multiple instances of Base if multiple sets of loggers or attributes
are desired.
*/
type Base struct {
	isInitialized bool
	queue         *queue
	logLevel      LogLevel
	loggers       []Logger
	BaseAttrs     map[string]interface{}
}

// NewBase creates a new instance of Base with default values set.
func NewBase() *Base {
	b := &Base{
		queue:     newQueue(),
		logLevel:  LEVEL_DEBUG,
		loggers:   make([]Logger, 0),
		BaseAttrs: make(map[string]interface{}, 0),
	}
	return b
}

type appExiter interface {
	Exit(code int)
}
type osExiter struct{}

func (exiter *osExiter) Exit(code int) {
	os.Exit(code)
}

var curExiter appExiter = &osExiter{}

func setExiter(exiter appExiter) {
	curExiter = exiter
}

/*
SetLogLevel sets the level messages will be logged at.  It will log any message
that is at the level or more severe than the level.
*/
func (b *Base) SetLogLevel(level LogLevel) {
	b.logLevel = level
}

func (b *Base) shouldLog(level LogLevel) bool {
	if level <= b.logLevel {
		return true
	}
	return false
}

// AddLogger adds a new logger instance to the Base
func (b *Base) AddLogger(logger Logger) error {
	if b.IsInitialized() && !logger.IsInitialized() {
		err := logger.InitLogger()
		if err != nil {
			return err
		}
	} else if !b.IsInitialized() && logger.IsInitialized() {
		err := logger.ShutdownLogger()
		if err != nil {
			return err
		}
	}
	b.loggers = append(b.loggers, logger)
	logger.SetBase(b)
	return nil
}

/*
ClearLoggers will shut down and remove any loggers added to the Base. If an
error occurs while shutting down one of the loggers, the list will not be
cleared but any loggers that have already been shut down before the error
occurred will remain shut down.
*/
func (b *Base) ClearLoggers() error {
	for _, logger := range b.loggers {
		err := logger.ShutdownLogger()
		if err != nil {
			return err
		}
	}
	b.loggers = make([]Logger, 0)

	return nil
}

// IsInitialized returns true if InitLoggers has been successfully run on the Base
func (b *Base) IsInitialized() bool {
	return b.isInitialized
}

/*
InitLoggers will run InitLogger on each Logger that has been added to the Base.
If an error occurs in initializing a logger, the loggers that have already been
initialized will continue to be initialized.
*/
func (b *Base) InitLoggers() error {
	for _, logger := range b.loggers {
		err := logger.InitLogger()
		if err != nil {
			return err
		}
	}

	b.queue.startQueueWorkers()
	b.isInitialized = true

	return nil
}

/*
RemoveLogger will run ShutdownLogger on the give logger and then remove the given
Logger from the list in Base
*/
func (b *Base) RemoveLogger(logger Logger) error {
	for idx, rLogger := range b.loggers {
		if rLogger == logger {
			err := rLogger.ShutdownLogger()
			if err != nil {
				return err
			}
			b.loggers[idx] = b.loggers[len(b.loggers)-1]
			b.loggers[len(b.loggers)-1] = nil
			b.loggers = b.loggers[:len(b.loggers)-1]
			return nil
		}
	}
	return nil
}

/*
ShutdownLoggers will run ShutdownLogger on each Logger in Base.  If an error occurs
while shutting down a Logger, the error will be returned and all the loggers that
were already shut down will remain shut down.
*/
func (b *Base) ShutdownLoggers() error {
	b.queue.stopQueueWorkers()

	for _, logger := range b.loggers {
		err := logger.ShutdownLogger()
		if err != nil {
			return err
		}
	}

	b.isInitialized = false

	return nil
}

/*
NewLogAdapter creates a LogAdapter using Base to log messages
*/
func (b *Base) NewLogAdapter(attrs map[string]interface{}) *LogAdapter {
	return newLogAdapter(b, attrs)
}

// ClearAttrs will remove all the attributes added to Base
func (b *Base) ClearAttrs() {
	b.BaseAttrs = make(map[string]interface{}, 0)
}

/*
SetAttr will set the value for the attribute with the name key.  If the key
already exists it will be overwritten with the new value.
*/
func (b *Base) SetAttr(key string, value interface{}) {
	b.BaseAttrs[key] = value
}

// RemoveAttr will remove the attribute with the name key.
func (b *Base) RemoveAttr(key string) {
	delete(b.BaseAttrs, key)
}

func (b *Base) log(level LogLevel, m map[string]interface{}, msg string, a ...interface{}) error {
	if !b.shouldLog(level) {
		return nil
	}
	nm := newMessage(b, level, m, msg, a...)
	return b.queue.QueueMessage(nm)
}

// Dbg is a short-hand version of Debug
func (b *Base) Dbg(msg string) error {
	return b.Debug(msg)
}

// Dbgf is a short-hand version of Debugf
func (b *Base) Dbgf(msg string, a ...interface{}) error {
	return b.Debugf(msg, a...)
}

// Dbgm is a short-hand version of Debugm
func (b *Base) Dbgm(m map[string]interface{}, msg string, a ...interface{}) error {
	return b.Debugm(m, msg, a...)
}

// Debug logs msg to all added loggers at LogLevel.LEVEL_DEBUG
func (b *Base) Debug(msg string) error {
	return b.log(LEVEL_DEBUG, nil, msg)
}

/*
Debugf uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LEVEL_DEBUG
*/
func (b *Base) Debugf(msg string, a ...interface{}) error {
	return b.log(LEVEL_DEBUG, nil, msg, a...)
}

/*
Debugm uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LEVEL_DEBUG. It will also
merge all attributes passed in m with any attributes added to Base and include them
with the message if the Logger supports it.
*/
func (b *Base) Debugm(m map[string]interface{}, msg string, a ...interface{}) error {
	return b.log(LEVEL_DEBUG, m, msg, a...)
}

// Info logs msg to all added loggers at LogLevel.LEVEL_INFO
func (b *Base) Info(msg string) error {
	return b.log(LEVEL_INFO, nil, msg)
}

/*
Infof uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LEVEL_INFO
*/
func (b *Base) Infof(msg string, a ...interface{}) error {
	return b.log(LEVEL_INFO, nil, msg, a...)
}

/*
Infom uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LEVEL_INFO. It will also
merge all attributes passed in m with any attributes added to Base and include them
with the message if the Logger supports it.
*/
func (b *Base) Infom(m map[string]interface{}, msg string, a ...interface{}) error {
	return b.log(LEVEL_INFO, m, msg, a...)
}

// Warn is a short-hand version of Warning
func (b *Base) Warn(msg string) error {
	return b.Warning(msg)
}

// Warnf is a short-hand version of Warningf
func (b *Base) Warnf(msg string, a ...interface{}) error {
	return b.Warningf(msg, a...)
}

// Warnm is a short-hand version of Warningm
func (b *Base) Warnm(m map[string]interface{}, msg string, a ...interface{}) error {
	return b.Warningm(m, msg, a...)
}

/*
Warning uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LEVEL_WARNING
*/
func (b *Base) Warning(msg string) error {
	return b.log(LEVEL_WARNING, nil, msg)
}

/*
Warningf uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LEVEL_WARNING
*/
func (b *Base) Warningf(msg string, a ...interface{}) error {
	return b.log(LEVEL_WARNING, nil, msg, a...)
}

/*
Warningm uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LEVEL_WARNING. It will also
merge all attributes passed in m with any attributes added to Base and include them
with the message if the Logger supports it.
*/
func (b *Base) Warningm(m map[string]interface{}, msg string, a ...interface{}) error {
	return b.log(LEVEL_WARNING, m, msg, a...)
}

func (b *Base) Err(msg string) error {
	return b.Error(msg)
}
func (b *Base) Errf(msg string, a ...interface{}) error {
	return b.Errorf(msg, a...)
}
func (b *Base) Errm(m map[string]interface{}, msg string, a ...interface{}) error {
	return b.Errorm(m, msg, a...)
}
func (b *Base) Error(msg string) error {
	return b.log(LEVEL_ERROR, nil, msg)
}
func (b *Base) Errorf(msg string, a ...interface{}) error {
	return b.log(LEVEL_ERROR, nil, msg, a...)
}
func (b *Base) Errorm(m map[string]interface{}, msg string, a ...interface{}) error {
	return b.log(LEVEL_ERROR, m, msg, a...)
}

func (b *Base) Fatal(msg string) error {
	return b.log(LEVEL_FATAL, nil, msg)
}
func (b *Base) Fatalf(msg string, a ...interface{}) error {
	return b.log(LEVEL_FATAL, nil, msg, a...)
}
func (b *Base) Fatalm(m map[string]interface{}, msg string, a ...interface{}) error {
	return b.log(LEVEL_FATAL, m, msg, a...)
}

func (b *Base) Die(exitCode int, msg string) {
	b.log(LEVEL_FATAL, nil, msg)
	b.ShutdownLoggers()
	curExiter.Exit(exitCode)
}
func (b *Base) Dief(exitCode int, msg string, a ...interface{}) {
	b.log(LEVEL_FATAL, nil, msg, a...)
	b.ShutdownLoggers()
	curExiter.Exit(exitCode)
}
func (b *Base) Diem(exitCode int, m map[string]interface{}, msg string, a ...interface{}) {
	b.log(LEVEL_FATAL, m, msg, a...)
	b.ShutdownLoggers()
	curExiter.Exit(exitCode)
}
