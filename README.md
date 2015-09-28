gomol
=====

[![GoDoc](https://godoc.org/github.com/aphistic/gomol?status.svg)](https://godoc.org/github.com/aphistic/gomol)
[![Build Status](https://img.shields.io/travis/aphistic/gomol.svg)](https://travis-ci.org/aphistic/gomol)
[![Code Coverage](https://img.shields.io/codecov/c/github/aphistic/gomol.svg)](http://codecov.io/github/aphistic/gomol?branch=master)

Gomol (Go Multi-Output Logger) is an MIT-licensed Go logging library.  The documentation and test coverage at this point is thin but will be improving over time. 

Features
========

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

Examples
========

For brevity a lot of error checking has been omitted, be sure you do your checks!

This is a super basic example of adding a number of loggers and then logging a few messages:

```go
package main

import (
	"github.com/aphistic/gomol"
)

func main() {
	consoleCfg := gomol.NewConsoleLoggerConfig()
	gomol.AddLogger(gomol.NewConsoleLogger(consoleCfg))
	gomol.AddLogger(gomol.NewLogglyLogger("1234"))
	gomol.AddLogger(gomol.NewGelfLogger("localhost", 12201))

	gomol.SetAttr("facility", "gomol.example")
	gomol.SetAttr("another_attr", 1234)

	gomol.InitLoggers()

	for idx := 1; idx <= 10; idx++ {
		gomol.Dbgm(map[string]interface{}{
			"msg_attr1": 4321,
		}, "Test message %v", idx)
	}

	gomol.ShutdownLoggers()
}
```
