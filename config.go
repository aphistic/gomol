package gomol

// Config is the runtime configuration for Gomol
type Config struct {
	// FilenameAttr is the name of the attribute to put the log location's
	// filename in.  This comes at a slight performance penalty.
	FilenameAttr string

	// LineNumberAttr is the name of the attribute to put the log location's
	// line number in.  This comes at a slight performance penalty.
	LineNumberAttr string

	// SequenceAttr is the name of the attribute to put the log message's sequence
	// number in.  The sequence number is an incrementing number for each log message
	// processed by a Base.
	SequenceAttr string

	// MaxQueueSize is the number of log messages which will be queued before old
	// messages are discarded.  This value takes effect once InitLoggers is called.
	// Further changes to this value will not increase or decrease the queue size.
	MaxQueueSize uint
}

// NewConfig creates a new configuration with default settings
func NewConfig() *Config {
	return &Config{
		FilenameAttr:   "",
		LineNumberAttr: "",
		SequenceAttr:   "",
		MaxQueueSize:   10000,
	}
}
