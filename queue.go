package gomol

import (
	"errors"
	"sync"
)

type queue struct {
	running      bool
	queueCtl     chan int
	senderCtl    chan int
	workersStart sync.WaitGroup
	workersDone  sync.WaitGroup

	queueChan chan *message

	queue        []*message
	msgAddedChan chan int
	queueMut     sync.Mutex
}

func newQueue() *queue {
	return &queue{
		running:      false,
		queueChan:    make(chan *message, 500),
		queueCtl:     make(chan int, 1),
		senderCtl:    make(chan int, 1),
		queue:        make([]*message, 0),
		msgAddedChan: make(chan int, 1),
	}
}

func (queue *queue) startQueueWorkers() error {
	if queue.running {
		return errors.New("Workers are already running")
	}
	queue.running = true
	queue.workersStart.Add(2)
	go queue.queueWorker(false)
	go queue.senderWorker(false)
	queue.workersStart.Wait()

	return nil
}

func (queue *queue) stopQueueWorkers() error {
	if queue.running {
		queue.queueCtl <- 1

		queue.workersDone.Wait()
		queue.running = false

		return nil
	}

	return errors.New("Workers are not running")
}

func (queue *queue) queueWorker(exiting bool) {
	queue.workersDone.Add(1)
	queue.workersStart.Done()
	for {
		if exiting {
			queue.queueMut.Lock()
			if len(queue.queueChan) == 0 {
				queue.queueMut.Unlock()
				break
			}
			queue.queueMut.Unlock()
		}

		select {
		case msg := <-queue.queueChan:
			queue.queueMut.Lock()
			queue.queue = append(queue.queue, msg)
			select {
			case queue.msgAddedChan <- 1:
			default:
			}
			queue.queueMut.Unlock()
		case <-queue.queueCtl:
			exiting = true
		}
	}
	queue.workersDone.Done()
	queue.senderCtl <- 1
}

func (queue *queue) senderWorker(exiting bool) {
	queue.workersDone.Add(1)
	queue.workersStart.Done()
	for {
		if exiting {
			queue.queueMut.Lock()
			done := false
			if len(queue.queue) == 0 {
				done = true
			}
			queue.queueMut.Unlock()
			if done {
				break
			}
		}

		select {
		case <-queue.senderCtl:
			exiting = true
		case <-queue.msgAddedChan:
		}

		for {
			msg := queue.NextMessage()

			if msg == nil {
				break
			}

			for _, l := range msg.Base.loggers {
				switch msg.Level {
				case LEVEL_DEBUG:
					l.Dbgm(msg.Attrs, msg.MsgFormat, msg.MsgParams...)
				case LEVEL_INFO:
					l.Infom(msg.Attrs, msg.MsgFormat, msg.MsgParams...)
				case LEVEL_WARNING:
					l.Warnm(msg.Attrs, msg.MsgFormat, msg.MsgParams...)
				case LEVEL_ERROR:
					l.Errm(msg.Attrs, msg.MsgFormat, msg.MsgParams...)
				case LEVEL_FATAL:
					l.Fatalm(msg.Attrs, msg.MsgFormat, msg.MsgParams...)
				}
			}
		}
	}
	queue.workersDone.Done()
}

func (queue *queue) QueueMessage(msg *message) error {
	if !queue.running {
		return errors.New("The logging system is not running, has InitLoggers() been executed?")
	}
	queue.queueChan <- msg
	return nil
}

func (queue *queue) NextMessage() *message {
	var msg *message
	queue.queueMut.Lock()
	if len(queue.queue) > 0 {
		msg, queue.queue = queue.queue[0], queue.queue[1:]
	} else {
		msg = nil
	}
	queue.queueMut.Unlock()

	return msg
}

func (queue *queue) Length() int {
	queue.queueMut.Lock()
	defer queue.queueMut.Unlock()
	return len(queue.queue) + len(queue.queueChan)
}
