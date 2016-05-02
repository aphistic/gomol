package gomol

type Logger interface {
	SetBase(*Base)

	InitLogger() error
	ShutdownLogger() error
	IsInitialized() bool

	Logm(LogLevel, map[string]interface{}, string) error
}
