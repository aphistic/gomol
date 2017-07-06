package gomol

import "errors"

const MaxQueueSize = 30000

type queue struct {
	running   bool
	finished  chan struct{}
	queueChan chan *Message
}

func newQueue() *queue {
	return &queue{
		running:   false,
		finished:  make(chan struct{}),
		queueChan: make(chan *Message, MaxQueueSize),
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

	queue.queueChan <- msg
	return nil
}

func (queue *queue) pressure() int {
	return len(queue.queueChan)
}
