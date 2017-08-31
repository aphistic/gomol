package gomol

import "errors"

var (
	// ErrUnknownLevel is returned when the provided log level is not known
	ErrUnknownLevel = errors.New("unknown log level")

	// ErrMessageDropped is reported if loggers are backed up and an old log
	// message has been forgotten
	ErrMessageDropped = errors.New("queue full - dropping message")

	// ErrNotInitialized is returned when a resource has not been completely
	// initialized
	ErrNotInitialized = errors.New("not initialized")
)
