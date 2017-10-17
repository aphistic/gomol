package gomol_test

import (
	"fmt"

	"github.com/aphistic/gomol"
	"github.com/aphistic/gomol-console"
	"github.com/aphistic/gomol-gelf"
)

func Example() {
	// Add a console logger
	consoleCfg := gomolconsole.NewConsoleLoggerConfig()
	consoleLogger, _ := gomolconsole.NewConsoleLogger(consoleCfg)
	// Set the template to display the full message, including
	// attributes.
	consoleLogger.SetTemplate(gomolconsole.NewTemplateFull())
	gomol.AddLogger(consoleLogger)

	// Add a GELF logger
	gelfCfg := gomolgelf.NewGelfLoggerConfig()
	gelfCfg.Hostname = "localhost"
	gelfCfg.Port = 12201
	gelfLogger, _ := gomolgelf.NewGelfLogger(gelfCfg)
	gomol.AddLogger(gelfLogger)

	// Set some global attrs that will be added to all
	// messages automatically
	gomol.SetAttr("facility", "gomol.example")
	gomol.SetAttr("another_attr", 1234)

	// Configure gomol to add the filename and line number that the
	// log was generated from, the internal sequence number to help
	// with ordering events if your log system doesn't support a
	// small enough sub-second resolution, and set the size of the
	// internal queue (default is 10k messages).
	cfg := gomol.NewConfig()
	cfg.FilenameAttr = "filename"
	cfg.LineNumberAttr = "line"
	cfg.SequenceAttr = "sequence"
	cfg.MaxQueueSize = 50000
	gomol.SetConfig(cfg)

	// Initialize the loggers
	gomol.InitLoggers()
	defer gomol.ShutdownLoggers()

	// Create a channel on which to receive internal (asynchronous)
	// logger errors. This is optional, but recommended in order to
	// determine when logging may be dropping messages.
	ch := make(chan error)

	go func() {
		// This consumer is expected to be efficient as writes to
		// the channel are blocking. If this handler is slow, the
		// user should add a buffer to the channel, or manually
		// queue and batch errors for processing.

		for err := range ch {
			fmt.Printf("[Internal Error] %s\n", err.Error())
		}
	}()

	gomol.SetErrorChan(ch)

	// Log some debug messages with message-level attrs
	// that will be sent only with that message
	for idx := 1; idx <= 10; idx++ {
		gomol.Dbgm(gomol.NewAttrsFromMap(map[string]interface{}{
			"msg_attr1": 4321,
		}), "Test message %v", idx)
	}
}
