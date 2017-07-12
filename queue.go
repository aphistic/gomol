package gomol

import "errors"

type queue struct {
	base      *Base
	running   bool
	finished  chan struct{}
	queueChan chan *Message
}

func newQueue(base *Base, maxQueueSize uint) *queue {
	return &queue{
		base:      base,
		running:   false,
		finished:  make(chan struct{}),
		queueChan: make(chan *Message, maxQueueSize),
	}
}

func (queue *queue) startWorker() error {
	if queue.running {
		return errors.New("workers are already running")
	}

	queue.running = true
	go queue.work()

	return nil
}

func (queue *queue) stopWorker() error {
	if !queue.running {
		return errors.New("workers are not running")
	}

	queue.running = false
	close(queue.queueChan)
	queue.flush()

	return nil
}

func (queue *queue) work() {
	defer close(queue.finished)

	for msg := range queue.queueChan {
		if msg == nil {
			continue
		}

		for _, l := range msg.base.loggers {
			l.Logm(msg.Timestamp, msg.Level, msg.Attrs.Attrs(), msg.Msg)
		}
	}
}

func (queue *queue) flush() {
	<-queue.finished
}

func (queue *queue) queueMessage(msg *Message) error {
	if !queue.running {
		return errors.New("the logging system is not running - has InitLoggers() been executed?")
	}

loop:
	for {
		// Attempt to queue the message immediately to
		// the channel.

		select {
		case queue.queueChan <- msg:
			break loop
		default:
		}

		// The queue was full. Try to read one message
		// from it (which will be the oldest) to make room
		// for another attempt to append. We do this in a
		// loop in case there's some contention - we'll keep
		// eating from the front until we finally make it in.

		select {
		case <-queue.queueChan:
			queue.base.report(ErrMessageDropped)
		default:
		}
	}

	return nil
}

func (queue *queue) pressure() int {
	return len(queue.queueChan)
}
