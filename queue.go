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

func (queue *queue) startQueueWorkers() {
	if queue.running {
		return
	}
	queue.running = true
	queue.workersStart.Add(2)
	go queue.queueWorker()
	go queue.senderWorker()
	queue.workersStart.Wait()
}

func (queue *queue) stopQueueWorkers() {
	if queue.running {
		queue.queueCtl <- 1

		queue.workersDone.Wait()
		queue.running = false
	}
}

func (queue *queue) queueWorker() {
	queue.workersDone.Add(1)
	queue.workersStart.Done()
	exiting := false
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

func (queue *queue) senderWorker() {
	queue.workersDone.Add(1)
	queue.workersStart.Done()
	exiting := false
	for {
		if exiting {
			queue.queueMut.Lock()
			if len(queue.queue) == 0 {
				queue.queueMut.Unlock()
				break
			}
			queue.queueMut.Unlock()
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
				case levelDbg:
					l.Dbgm(msg.Attrs, msg.Msg)
				case levelInfo:
					l.Infom(msg.Attrs, msg.Msg)
				case levelWarn:
					l.Warnm(msg.Attrs, msg.Msg)
				case levelError:
					l.Errm(msg.Attrs, msg.Msg)
				case levelFatal:
					l.Fatalm(msg.Attrs, msg.Msg)
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
