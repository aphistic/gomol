package gomol

/*
Logger is an interface libraries can implement to create their own loggers to be
used with gomol.
*/
type Logger interface {
	SetBase(*Base)

	InitLogger() error
	ShutdownLogger() error
	IsInitialized() bool

	Logm(LogLevel, map[string]interface{}, string) error
}
