package gomol

import "errors"

var (
	// ErrUnknownLevel is returned when the provided log level is not known
	ErrUnknownLevel = errors.New("Unknown log level")
)
