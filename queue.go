package gomol

import (
	"sync"
)

type queue struct {
	running     bool
	queueCtl    chan int
	senderCtl   chan int
	workersDone sync.WaitGroup

	queueChan chan *message

	queue        []*message
	msgAddedChan chan int
	queueMut     sync.RWMutex
}

var curQueue *queue

func init() {
	startQueueWorkers()
}

func newQueue() *queue {
	return &queue{
		running:      false,
		queueChan:    make(chan *message, 1000),
		queueCtl:     make(chan int),
		senderCtl:    make(chan int),
		queue:        make([]*message, 0),
		msgAddedChan: make(chan int, 1),
	}
}

func startQueueWorkers() {
	if curQueue != nil && curQueue.running {
		return
	}
	curQueue = newQueue()
	curQueue.running = true
	go curQueue.queueWorker()
	go curQueue.senderWorker()
}

func (queue *queue) queueWorker() {
	queue.workersDone.Add(1)
	exiting := false
	for {
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

		if exiting && len(queue.queueChan) == 0 {
			break
		}
	}
	queue.workersDone.Done()
}

func (queue *queue) senderWorker() {
	queue.workersDone.Add(1)
	exiting := false
	for {
		if exiting && len(queue.queue) == 0 {
			break
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

func (queue *queue) QueueMessage(msg *message) {
	queue.queueChan <- msg
}

func (queue *queue) NextMessage() *message {
	var msg *message
	queue.queueMut.Lock()
	if len(queue.queue) > 0 {
		q := queue.queue
		msg, q = q[0], q[1:]
		queue.queue = q
	} else {
		msg = nil
	}
	queue.queueMut.Unlock()

	return msg
}

func (queue *queue) Length() int {
	queue.queueMut.RLock()
	defer queue.queueMut.RUnlock()
	return len(queue.queue) + len(queue.queueChan)
}

func queueLen() int {
	return curQueue.Length()
}

func queueMessage(msg *message) {
	curQueue.QueueMessage(msg)
}

/*
Blocks until all messages in the queue have been processed, then returns.
*/
func FlushMessages() {
	if curQueue.running {
		curQueue.queueCtl <- 1
		curQueue.senderCtl <- 1

		curQueue.workersDone.Wait()
		curQueue.running = false
	}
}
