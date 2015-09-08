package gomol

const (
	levelDbg     = 7
	levelInfo    = 6
	levelWarn    = 4
	levelError   = 3
	levelFatal   = 2
	levelUnknown = -1
)

type Base struct {
	loggers   []Logger
	BaseAttrs map[string]interface{}
}

func newBase() *Base {
	b := &Base{
		loggers:   make([]Logger, 0),
		BaseAttrs: make(map[string]interface{}, 0),
	}
	return b
}

func (b *Base) AddLogger(logger Logger) {
	b.loggers = append(b.loggers, logger)
	logger.SetBase(b)
}

func (b *Base) InitLoggers() error {
	for _, logger := range b.loggers {
		err := logger.InitLogger()
		if err != nil {
			return err
		}
	}
	return nil
}
func (b *Base) ShutdownLoggers() error {
	for _, logger := range b.loggers {
		err := logger.ShutdownLogger()
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Base) ClearAttrs() {
	b.BaseAttrs = make(map[string]interface{}, 0)
}
func (b *Base) SetAttr(key string, value interface{}) {
	b.BaseAttrs[key] = value
}
func (b *Base) RemoveAttr(key string) {
	delete(b.BaseAttrs, key)
}

func (b *Base) Dbg(msg string) error {
	for _, l := range b.loggers {
		err := l.Dbg(msg)
		if err != nil {
			return err
		}
	}
	return nil
}
func (b *Base) Dbgf(msg string, a ...interface{}) error {
	for _, l := range b.loggers {
		err := l.Dbgf(msg, a...)
		if err != nil {
			return err
		}
	}
	return nil
}
func (b *Base) Dbgm(m map[string]interface{}, msg string, a ...interface{}) error {
	for _, l := range b.loggers {
		err := l.Dbgm(m, msg, a...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Base) Info(msg string) error {
	for _, l := range b.loggers {
		err := l.Info(msg)
		if err != nil {
			return err
		}
	}
	return nil
}
func (b *Base) Infof(msg string, a ...interface{}) error {
	for _, l := range b.loggers {
		err := l.Infof(msg, a...)
		if err != nil {
			return err
		}
	}
	return nil
}
func (b *Base) Infom(m map[string]interface{}, msg string, a ...interface{}) error {
	for _, l := range b.loggers {
		err := l.Infom(m, msg, a...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Base) Warn(msg string) error {
	for _, l := range b.loggers {
		err := l.Warn(msg)
		if err != nil {
			return err
		}
	}
	return nil
}
func (b *Base) Warnf(msg string, a ...interface{}) error {
	for _, l := range b.loggers {
		err := l.Warnf(msg, a...)
		if err != nil {
			return err
		}
	}
	return nil
}
func (b *Base) Warnm(m map[string]interface{}, msg string, a ...interface{}) error {
	for _, l := range b.loggers {
		err := l.Warnm(m, msg, a...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Base) Err(msg string) error {
	for _, l := range b.loggers {
		err := l.Err(msg)
		if err != nil {
			return err
		}
	}
	return nil
}
func (b *Base) Errf(msg string, a ...interface{}) error {
	for _, l := range b.loggers {
		err := l.Errf(msg, a...)
		if err != nil {
			return err
		}
	}
	return nil
}
func (b *Base) Errm(m map[string]interface{}, msg string, a ...interface{}) error {
	for _, l := range b.loggers {
		err := l.Errm(m, msg, a...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Base) Fatal(msg string) error {
	for _, l := range b.loggers {
		err := l.Fatal(msg)
		if err != nil {
			return err
		}
	}
	return nil
}
func (b *Base) Fatalf(msg string, a ...interface{}) error {
	for _, l := range b.loggers {
		err := l.Fatalf(msg, a...)
		if err != nil {
			return err
		}
	}
	return nil
}
func (b *Base) Fatalm(m map[string]interface{}, msg string, a ...interface{}) error {
	for _, l := range b.loggers {
		err := l.Fatalm(m, msg, a...)
		if err != nil {
			return err
		}
	}
	return nil
}
