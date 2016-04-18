package gomol

import (
	"os"
)

type Base struct {
	isInitialized bool
	queue         *queue
	logLevel      LogLevel
	loggers       []Logger
	BaseAttrs     map[string]interface{}
}

func newBase() *Base {
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
Sets the level to log at.  It will log any message that is at the level
or more severe than the level.
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

func (b *Base) IsInitialized() bool {
	return b.isInitialized
}

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

func (b *Base) ClearAttrs() {
	b.BaseAttrs = make(map[string]interface{}, 0)
}
func (b *Base) SetAttr(key string, value interface{}) {
	b.BaseAttrs[key] = value
}
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

func (b *Base) Dbg(msg string) error {
	return b.Debug(msg)
}
func (b *Base) Dbgf(msg string, a ...interface{}) error {
	return b.Debugf(msg, a...)
}
func (b *Base) Dbgm(m map[string]interface{}, msg string, a ...interface{}) error {
	return b.Debugm(m, msg, a...)
}
func (b *Base) Debug(msg string) error {
	return b.log(LEVEL_DEBUG, nil, msg)
}
func (b *Base) Debugf(msg string, a ...interface{}) error {
	return b.log(LEVEL_DEBUG, nil, msg, a...)
}
func (b *Base) Debugm(m map[string]interface{}, msg string, a ...interface{}) error {
	return b.log(LEVEL_DEBUG, m, msg, a...)
}

func (b *Base) Info(msg string) error {
	return b.log(LEVEL_INFO, nil, msg)
}
func (b *Base) Infof(msg string, a ...interface{}) error {
	return b.log(LEVEL_INFO, nil, msg, a...)
}
func (b *Base) Infom(m map[string]interface{}, msg string, a ...interface{}) error {
	return b.log(LEVEL_INFO, m, msg, a...)
}

func (b *Base) Warn(msg string) error {
	return b.log(LEVEL_WARNING, nil, msg)
}
func (b *Base) Warnf(msg string, a ...interface{}) error {
	return b.log(LEVEL_WARNING, nil, msg, a...)
}
func (b *Base) Warnm(m map[string]interface{}, msg string, a ...interface{}) error {
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
