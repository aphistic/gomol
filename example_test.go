package gomol

import (
	gomol "."
)

// Code for the README example to make sure it still builds!
func ExampleCode() {
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
