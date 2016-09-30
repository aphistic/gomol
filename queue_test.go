package gomol

import (
	"time"

	. "gopkg.in/check.v1"
)

func (s *GomolSuite) TestQueueLenChan(c *C) {
	q := newQueue()

	c.Check(q.Length(), Equals, 0)
	q.queueChan <- &message{}
	c.Check(q.Length(), Equals, 1)
}

func (s *GomolSuite) TestQueueLenArray(c *C) {
	q := newQueue()

	c.Check(q.Length(), Equals, 0)
	q.queue = append(q.queue, &message{})
	c.Check(q.Length(), Equals, 1)
}

func (s *GomolSuite) TestQueueLen(c *C) {
	q := newQueue()

	c.Check(q.Length(), Equals, 0)
	q.queueChan <- &message{}
	q.queue = append(q.queue, &message{})
	c.Check(q.Length(), Equals, 2)
}

func (s *GomolSuite) TestQueueMessageWithoutWorkers(c *C) {
	q := newQueue()
	err := q.QueueMessage(&message{})
	c.Check(err, NotNil)
	c.Check(err.Error(), Equals, "The logging system is not running, has InitLoggers() been executed?")
}

func (s *GomolSuite) TestQueueStartWorkers(c *C) {
	q := newQueue()
	q.startQueueWorkers()
	c.Check(q.IsActive(), Equals, true)
	c.Check(q.queue, HasLen, 0)
	q.stopQueueWorkers()
}

func (s *GomolSuite) TestQueueStartWorkersTwice(c *C) {
	q := newQueue()
	err := q.startQueueWorkers()
	c.Check(err, IsNil)
	c.Check(q.IsActive(), Equals, true)
	c.Check(q.queue, HasLen, 0)
	err = q.startQueueWorkers()
	c.Check(err, NotNil)
	c.Check(err.Error(), Equals, "Workers are already running")
	q.stopQueueWorkers()
}

func (s *GomolSuite) TestQueueStopWorkers(c *C) {
	q := newQueue()
	q.startQueueWorkers()

	q.stopQueueWorkers()
	c.Check(q.IsActive(), Equals, false)
	c.Check(q.queue, HasLen, 0)
	c.Check(len(q.queueChan), Equals, 0)
}

func (s *GomolSuite) TestQueueStopWorkersTwice(c *C) {
	q := newQueue()
	q.startQueueWorkers()

	err := q.stopQueueWorkers()
	c.Check(err, IsNil)
	c.Check(q.IsActive(), Equals, false)
	c.Check(q.queue, HasLen, 0)
	c.Check(len(q.queueChan), Equals, 0)
	err = q.stopQueueWorkers()
	c.Check(err, NotNil)
	c.Check(err.Error(), Equals, "Workers are not running")
}

func (s *GomolSuite) TestQueueFlushMessages(c *C) {
	q := newQueue()
	q.startQueueWorkers()

	for i := 0; i < 100; i++ {
		q.QueueMessage(newMessage(time.Now(), testBase, LevelDebug, nil, "test"))
	}

	q.stopQueueWorkers()
	c.Check(q.IsActive(), Equals, false)
	c.Check(q.queue, HasLen, 0)
	c.Check(len(q.queueChan), Equals, 0)
}
