package gomol

import "time"

/*
Logger is an interface libraries can implement to create their own loggers to be
used with gomol.
*/
type Logger interface {
	SetBase(*Base)

	InitLogger() error
	ShutdownLogger() error
	IsInitialized() bool

	Logm(time.Time, LogLevel, map[string]interface{}, string) error
}

// HookPreQueue is an interface a Logger can implement to be able to inspect
// and modify a Message before it is added to the queue
type HookPreQueue interface {
	Logger

	PreQueue(msg *Message) error
}
