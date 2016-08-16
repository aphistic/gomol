package gomol

var curDefault *Base

func init() {
	curDefault = NewBase()
}

func SetConfig(config *Config) {
	curDefault.SetConfig(config)
}

func SetLogLevel(level LogLevel) {
	curDefault.SetLogLevel(level)
}

func AddLogger(logger Logger) {
	curDefault.AddLogger(logger)
}

func RemoveLogger(logger Logger) error {
	return curDefault.RemoveLogger(logger)
}

func ClearLoggers() error {
	return curDefault.ClearLoggers()
}

func IsInitialized() bool {
	return curDefault.IsInitialized()
}

func InitLoggers() error {
	return curDefault.InitLoggers()
}

/*
ShutdownLoggers will call ShutdownLogger on each logger that has been added to
the current default logger.
*/
func ShutdownLoggers() error {
	return curDefault.ShutdownLoggers()
}

func ClearAttrs() {
	curDefault.ClearAttrs()
}
func SetAttr(key string, value interface{}) {
	curDefault.SetAttr(key, value)
}
func GetAttr(key string) interface{} {
	return curDefault.GetAttr(key)
}
func RemoveAttr(key string) {
	curDefault.RemoveAttr(key)
}

func NewLogAdapter(attrs *Attrs) *LogAdapter {
	return curDefault.NewLogAdapter(attrs)
}

func Dbg(msg string) error {
	return Debug(msg)
}
func Dbgf(msg string, a ...interface{}) error {
	return Debugf(msg, a...)
}
func Dbgm(m *Attrs, msg string, a ...interface{}) error {
	return Debugm(m, msg, a...)
}
func Debug(msg string) error {
	return curDefault.Debug(msg)
}
func Debugf(msg string, a ...interface{}) error {
	return curDefault.Debugf(msg, a...)
}
func Debugm(m *Attrs, msg string, a ...interface{}) error {
	return curDefault.Debugm(m, msg, a...)
}

func Info(msg string) error {
	return curDefault.Info(msg)
}
func Infof(msg string, a ...interface{}) error {
	return curDefault.Infof(msg, a...)
}
func Infom(m *Attrs, msg string, a ...interface{}) error {
	return curDefault.Infom(m, msg, a...)
}

func Warn(msg string) error {
	return curDefault.Warn(msg)
}
func Warnf(msg string, a ...interface{}) error {
	return curDefault.Warnf(msg, a...)
}
func Warnm(m *Attrs, msg string, a ...interface{}) error {
	return curDefault.Warnm(m, msg, a...)
}

func Err(msg string) error {
	return Error(msg)
}
func Errf(msg string, a ...interface{}) error {
	return Errorf(msg, a...)
}
func Errm(m *Attrs, msg string, a ...interface{}) error {
	return Errorm(m, msg, a...)
}
func Error(msg string) error {
	return curDefault.Error(msg)
}
func Errorf(msg string, a ...interface{}) error {
	return curDefault.Errorf(msg, a...)
}
func Errorm(m *Attrs, msg string, a ...interface{}) error {
	return curDefault.Errorm(m, msg, a...)
}

func Fatal(msg string) error {
	return curDefault.Fatal(msg)
}
func Fatalf(msg string, a ...interface{}) error {
	return curDefault.Fatalf(msg, a...)
}
func Fatalm(m *Attrs, msg string, a ...interface{}) error {
	return curDefault.Fatalm(m, msg, a...)
}

func Die(exitCode int, msg string) {
	curDefault.Die(exitCode, msg)
}
func Dief(exitCode int, msg string, a ...interface{}) {
	curDefault.Dief(exitCode, msg, a...)
}
func Diem(exitCode int, m *Attrs, msg string, a ...interface{}) {
	curDefault.Diem(exitCode, m, msg, a...)
}
