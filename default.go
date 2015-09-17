package gomol

var curDefault *Base

func init() {
	curDefault = newBase()
}

func AddLogger(logger Logger) {
	curDefault.AddLogger(logger)
}
func InitLoggers() error {
	return curDefault.InitLoggers()
}

/*
Calls ShutdownLogger() on all loggers that are part of the current default
logger base, then calls FlushMessages() to wait for all messages to be logged
before returning.
*/
func ShutdownLoggers() error {
	err := curDefault.ShutdownLoggers()
	if err != nil {
		return err
	}
	return nil
}

func ClearAttrs() {
	curDefault.ClearAttrs()
}
func SetAttr(key string, value interface{}) {
	curDefault.SetAttr(key, value)
}
func RemoveAttr(key string) {
	curDefault.RemoveAttr(key)
}

func Dbg(msg string) error {
	return curDefault.Dbg(msg)
}
func Dbgf(msg string, a ...interface{}) error {
	return curDefault.Dbgf(msg, a...)
}
func Dbgm(m map[string]interface{}, msg string, a ...interface{}) error {
	return curDefault.Dbgm(m, msg, a...)
}

func Info(msg string) error {
	return curDefault.Info(msg)
}
func Infof(msg string, a ...interface{}) error {
	return curDefault.Infof(msg, a...)
}
func Infom(m map[string]interface{}, msg string, a ...interface{}) error {
	return curDefault.Infom(m, msg, a...)
}

func Warn(msg string) error {
	return curDefault.Warn(msg)
}
func Warnf(msg string, a ...interface{}) error {
	return curDefault.Warnf(msg, a...)
}
func Warnm(m map[string]interface{}, msg string, a ...interface{}) error {
	return curDefault.Warnm(m, msg, a...)
}

func Err(msg string) error {
	return curDefault.Err(msg)
}
func Errf(msg string, a ...interface{}) error {
	return curDefault.Errf(msg, a...)
}
func Errm(m map[string]interface{}, msg string, a ...interface{}) error {
	return curDefault.Errm(m, msg, a...)
}

func Fatal(msg string) error {
	return curDefault.Fatal(msg)
}
func Fatalf(msg string, a ...interface{}) error {
	return curDefault.Fatalf(msg, a...)
}
func Fatalm(m map[string]interface{}, msg string, a ...interface{}) error {
	return curDefault.Fatalm(m, msg, a...)
}
