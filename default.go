package gomol

var curDefault *Base

func init() {
	curDefault = NewBase()
}

// Default will return the current default gomol Base logger
func Default() *Base {
	return curDefault
}

// SetConfig executes the same function on the default Base instance
func SetConfig(config *Config) {
	curDefault.SetConfig(config)
}

// SetErrorChan executes the same function on the default Base instance
func SetErrorChan(ch chan<- error) {
	curDefault.SetErrorChan(ch)
}

// SetLogLevel executes the same function on the default Base instance
func SetLogLevel(level LogLevel) {
	curDefault.SetLogLevel(level)
}

// AddLogger executes the same function on the default Base instance
func AddLogger(logger Logger) {
	curDefault.AddLogger(logger)
}

// RemoveLogger executes the same function on the default Base instance
func RemoveLogger(logger Logger) error {
	return curDefault.RemoveLogger(logger)
}

// ClearLoggers executes the same function on the default Base instance
func ClearLoggers() error {
	return curDefault.ClearLoggers()
}

// IsInitialized executes the same function on the default Base instance
func IsInitialized() bool {
	return curDefault.IsInitialized()
}

// InitLoggers executes the same function on the default Base instance
func InitLoggers() error {
	return curDefault.InitLoggers()
}

// Flush will wait until all messages currently queued are distributed to
// all initialized loggers
func Flush() {
	curDefault.Flush()
}

// ShutdownLoggers executes the same function on the default Base instance
func ShutdownLoggers() error {
	return curDefault.ShutdownLoggers()
}

// ClearAttrs executes the same function on the default Base instance
func ClearAttrs() {
	curDefault.ClearAttrs()
}

// SetAttr executes the same function on the default Base instance
func SetAttr(key string, value interface{}) {
	curDefault.SetAttr(key, value)
}

// GetAttr executes the same function on the default Base instance
func GetAttr(key string) interface{} {
	return curDefault.GetAttr(key)
}

// RemoveAttr executes the same function on the default Base instance
func RemoveAttr(key string) {
	curDefault.RemoveAttr(key)
}

// NewLogAdapter executes the same function on the default Base instance
func NewLogAdapter(attrs *Attrs) *LogAdapter {
	return curDefault.NewLogAdapter(attrs)
}

// Dbg executes the same function on the default Base instance
func Dbg(msg string) error {
	return Debug(msg)
}

// Dbgf executes the same function on the default Base instance
func Dbgf(msg string, a ...interface{}) error {
	return Debugf(msg, a...)
}

// Dbgm executes the same function on the default Base instance
func Dbgm(m *Attrs, msg string, a ...interface{}) error {
	return Debugm(m, msg, a...)
}

// Debug executes the same function on the default Base instance
func Debug(msg string) error {
	return curDefault.Debug(msg)
}

// Debugf executes the same function on the default Base instance
func Debugf(msg string, a ...interface{}) error {
	return curDefault.Debugf(msg, a...)
}

// Debugm executes the same function on the default Base instance
func Debugm(m *Attrs, msg string, a ...interface{}) error {
	return curDefault.Debugm(m, msg, a...)
}

// Info executes the same function on the default Base instance
func Info(msg string) error {
	return curDefault.Info(msg)
}

// Infof executes the same function on the default Base instance
func Infof(msg string, a ...interface{}) error {
	return curDefault.Infof(msg, a...)
}

// Infom executes the same function on the default Base instance
func Infom(m *Attrs, msg string, a ...interface{}) error {
	return curDefault.Infom(m, msg, a...)
}

// Warn executes the same function on the default Base instance
func Warn(msg string) error {
	return Warning(msg)
}

// Warnf executes the same function on the default Base instance
func Warnf(msg string, a ...interface{}) error {
	return Warningf(msg, a...)
}

// Warnm executes the same function on the default Base instance
func Warnm(m *Attrs, msg string, a ...interface{}) error {
	return Warningm(m, msg, a...)
}

// Warning executes the same function on the default Base instance
func Warning(msg string) error {
	return curDefault.Warning(msg)
}

// Warningf executes the same function on the default Base instance
func Warningf(msg string, a ...interface{}) error {
	return curDefault.Warningf(msg, a...)
}

// Warningm executes the same function on the default Base instance
func Warningm(m *Attrs, msg string, a ...interface{}) error {
	return curDefault.Warningm(m, msg, a...)
}

// Err executes the same function on the default Base instance
func Err(msg string) error {
	return Error(msg)
}

// Errf executes the same function on the default Base instance
func Errf(msg string, a ...interface{}) error {
	return Errorf(msg, a...)
}

// Errm executes the same function on the default Base instance
func Errm(m *Attrs, msg string, a ...interface{}) error {
	return Errorm(m, msg, a...)
}

// Error executes the same function on the default Base instance
func Error(msg string) error {
	return curDefault.Error(msg)
}

// Errorf executes the same function on the default Base instance
func Errorf(msg string, a ...interface{}) error {
	return curDefault.Errorf(msg, a...)
}

// Errorm executes the same function on the default Base instance
func Errorm(m *Attrs, msg string, a ...interface{}) error {
	return curDefault.Errorm(m, msg, a...)
}

// Fatal executes the same function on the default Base instance
func Fatal(msg string) error {
	return curDefault.Fatal(msg)
}

// Fatalf executes the same function on the default Base instance
func Fatalf(msg string, a ...interface{}) error {
	return curDefault.Fatalf(msg, a...)
}

// Fatalm executes the same function on the default Base instance
func Fatalm(m *Attrs, msg string, a ...interface{}) error {
	return curDefault.Fatalm(m, msg, a...)
}

// Die executes the same function on the default Base instance
func Die(exitCode int, msg string) {
	curDefault.Die(exitCode, msg)
}

// Dief executes the same function on the default Base instance
func Dief(exitCode int, msg string, a ...interface{}) {
	curDefault.Dief(exitCode, msg, a...)
}

// Diem executes the same function on the default Base instance
func Diem(exitCode int, m *Attrs, msg string, a ...interface{}) {
	curDefault.Diem(exitCode, m, msg, a...)
}
