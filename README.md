gomol
=====

[![GoDoc](https://godoc.org/github.com/aphistic/gomol?status.svg)](https://godoc.org/github.com/aphistic/gomol)
[![Build Status](https://img.shields.io/travis/aphistic/gomol.svg)](https://travis-ci.org/aphistic/gomol)
[![Code Coverage](https://img.shields.io/codecov/c/github/aphistic/gomol.svg)](http://codecov.io/github/aphistic/gomol?branch=master)

Gomol (Go Multi-Output Logger) is an MIT-licensed Go logging library.  The documentation at this point is thin but will be improving over time.

Features
========

* Attach meta-data to each log message with attributes
* Multiple outputs at the same time

Installation
============

The recommended way to install is via http://gopkg.in

    go get gopkg.in/aphistic/gomol.v0
    ...
    import "gopkg.in/aphistic/gomol.v0"

Gomol can also be installed the standard way as well

    go get github.com/aphistic/gomol
    ...
    import "github.com/aphistic/gomol"

Loggers
=======

Right now there are a number of loggers built directly into the gomol library.  This
adds a number of dependencies that I'd like to remove so I'm in the process of splitting
each logger out to its own package.  The supported log outputs are listed below.  If
you have a logger you've written to support gomol and you'd like to add it to this list
please either submit a pull request with the updated document or let me know and I
can add it!

* **Console** - https://github.com/aphistic/gomol-console
* **Graylog Extended Log Format (GELF)** - built into gomol at the moment
* **Loggly** - https://github.com/aphistic/gomol-loggly
* **io.Writer** - https://github.com/aphistic/gomol-writer

Examples
========

For brevity a lot of error checking has been omitted, be sure you do your checks!

This is a super basic example of adding a number of loggers and then logging a few messages:

```go
package main

import (
	"github.com/aphistic/gomol"
	gc "github.com/aphistic/gomol-console"
)

func main() {
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
```
