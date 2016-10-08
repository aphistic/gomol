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
