package gomol

type Logger interface {
	SetBase(*Base)

	InitLogger() error
	ShutdownLogger() error
	IsInitialized() bool

	Dbg(string) error
	Dbgf(string, ...interface{}) error
	Dbgm(map[string]interface{}, string, ...interface{}) error
	Debug(string) error
	Debugf(string, ...interface{}) error
	Debugm(map[string]interface{}, string, ...interface{}) error

	Info(string) error
	Infof(string, ...interface{}) error
	Infom(map[string]interface{}, string, ...interface{}) error

	Warn(string) error
	Warnf(string, ...interface{}) error
	Warnm(map[string]interface{}, string, ...interface{}) error

	Err(string) error
	Errf(string, ...interface{}) error
	Errm(map[string]interface{}, string, ...interface{}) error
	Error(string) error
	Errorf(string, ...interface{}) error
	Errorm(map[string]interface{}, string, ...interface{}) error

	Fatal(string) error
	Fatalf(string, ...interface{}) error
	Fatalm(map[string]interface{}, string, ...interface{}) error
}
