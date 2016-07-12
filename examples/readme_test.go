package gomolexamples

import (
	"github.com/aphistic/gomol"
	gc "github.com/aphistic/gomol-console"
)

// Code for the README example to make sure it still builds!
func Example() {
	// Add a console logger
	consoleCfg := gc.NewConsoleLoggerConfig()
	consoleLogger, _ := gc.NewConsoleLogger(consoleCfg)
	gomol.AddLogger(consoleLogger)

	// Add a GELF logger
	gelfCfg := gomol.NewGelfLoggerConfig()
	gelfCfg.Hostname = "localhost"
	gelfCfg.Port = 12201
	gomol.AddLogger(gomol.NewGelfLogger(gelfCfg))

	// Set some global attrs that will be added to all
	// messages automatically
	gomol.SetAttr("facility", "gomol.example")
	gomol.SetAttr("another_attr", 1234)

	// Initialize the loggers
	gomol.InitLoggers()
	defer gomol.ShutdownLoggers()

	// Log some debug messages with message-level attrs
	// that will be sent only with that message
	for idx := 1; idx <= 10; idx++ {
		gomol.Dbgm(map[string]interface{}{
			"msg_attr1": 4321,
		}, "Test message %v", idx)
	}
}
