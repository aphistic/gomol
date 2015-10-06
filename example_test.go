package gomol

import (
	gomol "."
)

// Code for the README example to make sure it still builds!
func ExampleCode() {
	// Add a console logger
	consoleCfg := gomol.NewConsoleLoggerConfig()
	consoleLogger, _ := gomol.NewConsoleLogger(consoleCfg)
	gomol.AddLogger(consoleLogger)

	// Add a Loggly logger
	logglyCfg := gomol.NewLogglyLoggerConfig()
	logglyCfg.Token = "1234"
	logglyLogger, _ := gomol.NewLogglyLogger(logglyCfg)
	gomol.AddLogger(logglyLogger)

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
