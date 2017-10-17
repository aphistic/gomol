package gomol_test

import (
	"github.com/aphistic/gomol"
	"github.com/aphistic/gomol-console"
	"github.com/aphistic/gomol-json"
)

// ExampleFallbackLogger demonstrates how to use a fallback logger.
func Example_fallbackLogger() {
	// Create a logger that logs over TCP using JSON
	jsonCfg := gomoljson.NewJSONLoggerConfig("tcp://192.0.2.125:4321")
	// Continue startup even if we can't connect initially
	jsonCfg.AllowDisconnectedInit = true
	jsonLogger, _ := gomoljson.NewJSONLogger(jsonCfg)
	gomol.AddLogger(jsonLogger)

	// Create a logger that logs to the console
	consoleCfg := gomolconsole.NewConsoleLoggerConfig()
	consoleLogger, _ := gomolconsole.NewConsoleLogger(consoleCfg)

	// Set the fallback logger to the console so if the
	// TCP JSON logger is unhealthy we still get logs
	// to stdout.
	_ = gomol.SetFallbackLogger(consoleLogger)

	gomol.InitLoggers()
	defer gomol.ShutdownLoggers()

	gomol.Debug("This is my message!")
}
