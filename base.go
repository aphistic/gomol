package gomol

import (
	"os"
	"sync/atomic"
	"time"

	"github.com/efritz/glock"
)

type baseConfigFunc func(b *Base)

func withClock(clock glock.Clock) baseConfigFunc {
	return func(b *Base) {
		b.clock = clock
	}
}

/*
Base holds an instance of all information needed for logging.  It is possible
to create multiple instances of Base if multiple sets of loggers or attributes
are desired.
*/
type Base struct {
	clock glock.Clock

	isInitialized bool
	config        *Config
	errorChan     chan<- error
	queue         *queue
	logLevel      LogLevel
	sequence      uint64
	BaseAttrs     *Attrs

	loggers        []Logger
	fallbackLogger Logger
	hookPreQueue   []HookPreQueue
}

// NewBase creates a new instance of Base with default values set.
func NewBase(configs ...baseConfigFunc) *Base {
	b := &Base{
		clock: glock.NewRealClock(),

		config:    NewConfig(),
		logLevel:  LevelDebug,
		sequence:  0,
		BaseAttrs: NewAttrs(),

		loggers:      make([]Logger, 0),
		hookPreQueue: make([]HookPreQueue, 0),
	}

	for _, f := range configs {
		f(b)
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

// SetConfig will set the configuration for the Base to the given Config
func (b *Base) SetConfig(config *Config) {
	b.config = config
}

// SetErrorChan will register a channel as the consumer of internal error
// events.  This channel will be closed once ShutdownLoggers has finished.
// The consumer of this channel is expected to be efficient as writing to
// this channel will block.
func (b *Base) SetErrorChan(ch chan<- error) {
	b.errorChan = ch
}

func (b *Base) report(err error) {
	if b.errorChan != nil {
		b.errorChan <- err
	}
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

// SetFallbackLogger sets a Logger to be used if there aren't any loggers added or any of
// the added loggers are in a degraded or unhealthy state.  A Logger passed to SetFallbackLogger
// will be initialized if it hasn't been already.  In addition, if the Logger fails to initialize
// completely the fallback logger will fail to be set.
func (b *Base) SetFallbackLogger(logger Logger) error {
	if logger == nil {
		if b.fallbackLogger != nil && b.fallbackLogger.IsInitialized() {
			b.fallbackLogger.ShutdownLogger()
		}
		b.fallbackLogger = nil
		return nil
	}

	if !logger.IsInitialized() {
		err := logger.InitLogger()
		if err != nil {
			return err
		}
	}

	// Shut down any old logger we might already have a reference to
	if b.fallbackLogger != nil && b.fallbackLogger.IsInitialized() {
		b.fallbackLogger.ShutdownLogger()
	}

	b.fallbackLogger = logger

	return nil
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

	if hook, ok := logger.(HookPreQueue); ok {
		b.hookPreQueue = append(b.hookPreQueue, hook)
	}

	logger.SetBase(b)
	return nil
}

/*
RemoveLogger will run ShutdownLogger on the given logger and then remove the given
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

	// Remove any hook instances the logger has
	for idx, hookLogger := range b.hookPreQueue {
		if hookLogger == logger {
			b.hookPreQueue[idx] = b.hookPreQueue[len(b.hookPreQueue)-1]
			b.hookPreQueue[len(b.hookPreQueue)-1] = nil
			b.hookPreQueue = b.hookPreQueue[:len(b.hookPreQueue)-1]
		}
	}

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
	b.hookPreQueue = make([]HookPreQueue, 0)

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
	if b.queue == nil {
		b.queue = newQueue(b, b.config.MaxQueueSize)
	}

	for _, logger := range b.loggers {
		err := logger.InitLogger()
		if err != nil {
			return err
		}
	}

	b.queue.startWorker()
	b.isInitialized = true

	return nil
}

// Flush will wait until all messages currently queued are distributed to
// all initialized loggers
func (b *Base) Flush() {
	if b.queue != nil {
		b.queue.flush()
	}
}

/*
ShutdownLoggers will run ShutdownLogger on each Logger in Base.  If an error occurs
while shutting down a Logger, the error will be returned and all the loggers that
were already shut down will remain shut down.
*/
func (b *Base) ShutdownLoggers() error {
	// Before shutting down we should flush all the messsages
	b.Flush()

	for _, logger := range b.loggers {
		err := logger.ShutdownLogger()
		if err != nil {
			return err
		}
	}

	if b.queue != nil {
		b.queue.stopWorker()
	}

	if b.errorChan != nil {
		close(b.errorChan)
		b.errorChan = nil
	}

	b.isInitialized = false
	return nil
}

/*
NewLogAdapter creates a LogAdapter using Base to log messages
*/
func (b *Base) NewLogAdapter(attrs *Attrs) *LogAdapter {
	return NewLogAdapterFor(b, attrs)
}

// ClearAttrs will remove all the attributes added to Base
func (b *Base) ClearAttrs() {
	b.BaseAttrs = NewAttrs()
}

/*
SetAttr will set the value for the attribute with the name key.  If the key
already exists it will be overwritten with the new value.
*/
func (b *Base) SetAttr(key string, value interface{}) {
	b.BaseAttrs.SetAttr(key, value)
}

/*
GetAttr will return the current value for the given attribute key.  If the key
isn't set this will return nil
*/
func (b *Base) GetAttr(key string) interface{} {
	return b.BaseAttrs.GetAttr(key)
}

// RemoveAttr will remove the attribute with the name key.
func (b *Base) RemoveAttr(key string) {
	b.BaseAttrs.RemoveAttr(key)
}

// LogWithTime will log a message at the provided level to all added loggers with the timestamp set to the
// value of ts.
func (b *Base) LogWithTime(level LogLevel, ts time.Time, m *Attrs, msg string, a ...interface{}) error {
	if !b.shouldLog(level) {
		return nil
	}

	if !b.isInitialized {
		return ErrNotInitialized
	}

	if len(b.config.FilenameAttr) > 0 || len(b.config.LineNumberAttr) > 0 {
		file, line := getCallerInfo()
		if m == nil {
			m = NewAttrs()
		}
		if len(b.config.FilenameAttr) > 0 {
			m.SetAttr(b.config.FilenameAttr, file)
		}
		if len(b.config.LineNumberAttr) > 0 {
			m.SetAttr(b.config.LineNumberAttr, line)
		}
	}

	if len(b.config.SequenceAttr) > 0 {
		if m == nil {
			m = NewAttrs()
		}
		seq := atomic.AddUint64(&b.sequence, 1)
		m.SetAttr(b.config.SequenceAttr, seq)
	}

	nm := newMessage(ts, b, level, m, msg, a...)

	for _, hook := range b.hookPreQueue {
		err := hook.PreQueue(nm)
		if err != nil {
			return err
		}
	}

	return b.queue.queueMessage(nm)
}

// Log will log a message at the provided level to all added loggers with the timestamp set to the time
// Log was called.
func (b *Base) Log(level LogLevel, m *Attrs, msg string, a ...interface{}) error {
	return b.LogWithTime(level, b.clock.Now(), m, msg, a...)
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
func (b *Base) Dbgm(m *Attrs, msg string, a ...interface{}) error {
	return b.Debugm(m, msg, a...)
}

// Debug logs msg to all added loggers at LogLevel.LevelDebug
func (b *Base) Debug(msg string) error {
	return b.Log(LevelDebug, nil, msg)
}

/*
Debugf uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelDebug
*/
func (b *Base) Debugf(msg string, a ...interface{}) error {
	return b.Log(LevelDebug, nil, msg, a...)
}

/*
Debugm uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelDebug. It will also
merge all attributes passed in m with any attributes added to Base and include them
with the message if the Logger supports it.
*/
func (b *Base) Debugm(m *Attrs, msg string, a ...interface{}) error {
	return b.Log(LevelDebug, m, msg, a...)
}

// Info logs msg to all added loggers at LogLevel.LevelInfo
func (b *Base) Info(msg string) error {
	return b.Log(LevelInfo, nil, msg)
}

/*
Infof uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelInfo
*/
func (b *Base) Infof(msg string, a ...interface{}) error {
	return b.Log(LevelInfo, nil, msg, a...)
}

/*
Infom uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelInfo. It will also
merge all attributes passed in m with any attributes added to Base and include them
with the message if the Logger supports it.
*/
func (b *Base) Infom(m *Attrs, msg string, a ...interface{}) error {
	return b.Log(LevelInfo, m, msg, a...)
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
func (b *Base) Warnm(m *Attrs, msg string, a ...interface{}) error {
	return b.Warningm(m, msg, a...)
}

/*
Warning uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelWarning
*/
func (b *Base) Warning(msg string) error {
	return b.Log(LevelWarning, nil, msg)
}

/*
Warningf uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelWarning
*/
func (b *Base) Warningf(msg string, a ...interface{}) error {
	return b.Log(LevelWarning, nil, msg, a...)
}

/*
Warningm uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelWarning. It will also
merge all attributes passed in m with any attributes added to Base and include them
with the message if the Logger supports it.
*/
func (b *Base) Warningm(m *Attrs, msg string, a ...interface{}) error {
	return b.Log(LevelWarning, m, msg, a...)
}

// Err is a short-hand version of Error
func (b *Base) Err(msg string) error {
	return b.Error(msg)
}

// Errf is a short-hand version of Errorf
func (b *Base) Errf(msg string, a ...interface{}) error {
	return b.Errorf(msg, a...)
}

// Errm is a short-hand version of Errorm
func (b *Base) Errm(m *Attrs, msg string, a ...interface{}) error {
	return b.Errorm(m, msg, a...)
}

/*
Error uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelError
*/
func (b *Base) Error(msg string) error {
	return b.Log(LevelError, nil, msg)
}

/*
Errorf uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelError
*/
func (b *Base) Errorf(msg string, a ...interface{}) error {
	return b.Log(LevelError, nil, msg, a...)
}

/*
Errorm uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelError. It will also
merge all attributes passed in m with any attributes added to Base and include them
with the message if the Logger supports it.
*/
func (b *Base) Errorm(m *Attrs, msg string, a ...interface{}) error {
	return b.Log(LevelError, m, msg, a...)
}

/*
Fatal uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelFatal
*/
func (b *Base) Fatal(msg string) error {
	return b.Log(LevelFatal, nil, msg)
}

/*
Fatalf uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelFatal
*/
func (b *Base) Fatalf(msg string, a ...interface{}) error {
	return b.Log(LevelFatal, nil, msg, a...)
}

/*
Fatalm uses msg as a format string with subsequent parameters as values and logs
the resulting message to all added loggers at LogLevel.LevelFatal. It will also
merge all attributes passed in m with any attributes added to Base and include them
with the message if the Logger supports it.
*/
func (b *Base) Fatalm(m *Attrs, msg string, a ...interface{}) error {
	return b.Log(LevelFatal, m, msg, a...)
}

// Die will log a message using Fatal, call ShutdownLoggers and then exit the application with the provided exit code.
func (b *Base) Die(exitCode int, msg string) {
	b.Log(LevelFatal, nil, msg)
	b.ShutdownLoggers()
	curExiter.Exit(exitCode)
}

// Dief will log a message using Fatalf, call ShutdownLoggers and then exit the application with the provided exit code.
func (b *Base) Dief(exitCode int, msg string, a ...interface{}) {
	b.Log(LevelFatal, nil, msg, a...)
	b.ShutdownLoggers()
	curExiter.Exit(exitCode)
}

// Diem will log a message using Fatalm, call ShutdownLoggers and then exit the application with the provided exit code.
func (b *Base) Diem(exitCode int, m *Attrs, msg string, a ...interface{}) {
	b.Log(LevelFatal, m, msg, a...)
	b.ShutdownLoggers()
	curExiter.Exit(exitCode)
}
