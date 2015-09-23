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
	isInitialized bool
	queue         *queue
	loggers       []Logger
	BaseAttrs     map[string]interface{}
}

func newBase() *Base {
	b := &Base{
		queue:     newQueue(),
		loggers:   make([]Logger, 0),
		BaseAttrs: make(map[string]interface{}, 0),
	}
	return b
}

func (b *Base) AddLogger(logger Logger) error {
	if b.isInitialized && !logger.IsInitialized() {
		err := logger.InitLogger()
		if err != nil {
			return err
		}
	} else if !b.isInitialized && logger.IsInitialized() {
		err := logger.ShutdownLogger()
		if err != nil {
			return err
		}
	}
	b.loggers = append(b.loggers, logger)
	logger.SetBase(b)
	return nil
}

func (b *Base) InitLoggers() error {
	for _, logger := range b.loggers {
		err := logger.InitLogger()
		if err != nil {
			return err
		}
	}

	b.queue.startQueueWorkers()
	b.isInitialized = true

	return nil
}
func (b *Base) ShutdownLoggers() error {
	b.queue.stopQueueWorkers()

	for _, logger := range b.loggers {
		err := logger.ShutdownLogger()
		if err != nil {
			return err
		}
	}

	b.isInitialized = false

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
	nm := newMessage(b, levelDbg, nil, msg)
	return b.queue.QueueMessage(nm)
}
func (b *Base) Dbgf(msg string, a ...interface{}) error {
	nm := newMessage(b, levelDbg, nil, msg, a...)
	return b.queue.QueueMessage(nm)
}
func (b *Base) Dbgm(m map[string]interface{}, msg string, a ...interface{}) error {
	nm := newMessage(b, levelDbg, m, msg, a...)
	return b.queue.QueueMessage(nm)
}

func (b *Base) Info(msg string) error {
	nm := newMessage(b, levelInfo, nil, msg)
	return b.queue.QueueMessage(nm)
}
func (b *Base) Infof(msg string, a ...interface{}) error {
	nm := newMessage(b, levelInfo, nil, msg, a...)
	return b.queue.QueueMessage(nm)
}
func (b *Base) Infom(m map[string]interface{}, msg string, a ...interface{}) error {
	nm := newMessage(b, levelInfo, m, msg, a...)
	return b.queue.QueueMessage(nm)
}

func (b *Base) Warn(msg string) error {
	nm := newMessage(b, levelWarn, nil, msg)
	return b.queue.QueueMessage(nm)
}
func (b *Base) Warnf(msg string, a ...interface{}) error {
	nm := newMessage(b, levelWarn, nil, msg, a...)
	return b.queue.QueueMessage(nm)
}
func (b *Base) Warnm(m map[string]interface{}, msg string, a ...interface{}) error {
	nm := newMessage(b, levelWarn, m, msg, a...)
	return b.queue.QueueMessage(nm)
}

func (b *Base) Err(msg string) error {
	nm := newMessage(b, levelError, nil, msg)
	return b.queue.QueueMessage(nm)
}
func (b *Base) Errf(msg string, a ...interface{}) error {
	nm := newMessage(b, levelError, nil, msg, a...)
	return b.queue.QueueMessage(nm)
}
func (b *Base) Errm(m map[string]interface{}, msg string, a ...interface{}) error {
	nm := newMessage(b, levelError, m, msg, a...)
	return b.queue.QueueMessage(nm)
}

func (b *Base) Fatal(msg string) error {
	nm := newMessage(b, levelFatal, nil, msg)
	return b.queue.QueueMessage(nm)
}
func (b *Base) Fatalf(msg string, a ...interface{}) error {
	nm := newMessage(b, levelFatal, nil, msg, a...)
	return b.queue.QueueMessage(nm)
}
func (b *Base) Fatalm(m map[string]interface{}, msg string, a ...interface{}) error {
	nm := newMessage(b, levelFatal, m, msg, a...)
	return b.queue.QueueMessage(nm)
}
