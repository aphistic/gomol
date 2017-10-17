gomol
=====

[![GoDoc](https://godoc.org/github.com/aphistic/gomol?status.svg)](https://godoc.org/github.com/aphistic/gomol)
[![Build Status](https://img.shields.io/travis/aphistic/gomol.svg)](https://travis-ci.org/aphistic/gomol)
[![Code Coverage](https://img.shields.io/codecov/c/github/aphistic/gomol.svg)](http://codecov.io/github/aphistic/gomol?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/aphistic/gomol)](https://goreportcard.com/report/github.com/aphistic/gomol)

Gomol (Go Multi-Output Logger) is an MIT-licensed structured logging library for Go.  Gomol grew
from a desire to have a structured logging library that could write to any number of outputs
while also keeping a small in-band footprint.

Features
========

* Attach meta-data to each log message with attributes
* Multiple outputs at the same time
* Pluggable Logger interface
* Asynchronous logging so slow loggers won't slow down your application

Installation
============

Gomol can also be installed the standard way for Go:

    go get github.com/aphistic/gomol
    ...
    import "github.com/aphistic/gomol"

Vendoring is recommended!

Loggers
=======

Gomol has a growing list of supported logging formats.  The known loggers are listed
below.  If you have a logger you've written to support gomol and you'd like to add it
to this list please either submit a pull request with the updated document or let me
know and I can add it!

* **Console** - https://github.com/aphistic/gomol-console
* **Graylog Extended Log Format (GELF)** - https://github.com/aphistic/gomol-gelf
* **io.Writer** - https://github.com/aphistic/gomol-writer
* **JSON** - https://github.com/aphistic/gomol-json
* **Loggly** - https://github.com/aphistic/gomol-loggly

Other Usages
============

In addition to the loggers listed above, gomol can be used with other projects as well.

* **negroni-gomol** (https://github.com/aphistic/negroni-gomol) - Negroni logging middleware
	using gomol.

Examples
========

For brevity a lot of error checking has been omitted, be sure you do your checks!

This is a super basic example of adding a number of loggers and then logging a few messages:

```go
package main

import (
	"github.com/aphistic/gomol"
	gc "github.com/aphistic/gomol-console"
	gg "github.com/aphistic/gomol-gelf"
)

func main() {
	// Add a console logger
	consoleCfg := gc.NewConsoleLoggerConfig()
	consoleLogger, _ := gc.NewConsoleLogger(consoleCfg)
	// Set the template to display the full message, including
	// attributes.
	consoleLogger.SetTemplate(gc.NewTemplateFull())
	gomol.AddLogger(consoleLogger)

	// Add a GELF logger
	gelfCfg := gg.NewGelfLoggerConfig()
	gelfCfg.Hostname = "localhost"
	gelfCfg.Port = 12201
	gelfLogger, _ := gg.NewGelfLogger(gelfCfg)
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
```

Fallback Logger
===============

One feature gomol supports is the concept of a fallback logger.  In some cases a logger
may go unhealthy (a logger to a remote server, for example) and you want to log to a
different logger to ensure log messages are not lost.  The SetFallbackLogger method is
available for these such instances.

A fallback logger will be triggered if any of the primary loggers goes unhealthy even if
all others are fine. It's recommended the fallback logger not be added to the primary
loggers or you may see duplicate messages.  This does mean if multiple loggers are added
as primary loggers and just one is unhealthy the fallback logger logger will be triggered.

To add a fallback logger there are two options.  One is to use the default gomol instance
(`gomol.SetFallbackLogger()`) and the other is to use the method on a Base instance:

```go
import (
	"github.com/aphistic/gomol"
	"github.com/aphistic/gomol-console"
	"github.com/aphistic/gomol-json"
)

func main() {
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
```