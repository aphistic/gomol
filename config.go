package gomol

// Config is the runtime configuration for Gomol
type Config struct {
	FilenameAttr   string
	LineNumberAttr string
}

// NewConfig creates a new configuration with default settings
func NewConfig() *Config {
	return &Config{
		FilenameAttr:   "",
		LineNumberAttr: "",
	}
}
