package gomol

import (
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

func (s *GomolSuite) TestQueueStartWorkers(c *C) {
	q := newQueue()
	q.startQueueWorkers()
	c.Check(q.running, Equals, true)
	c.Check(q.queue, HasLen, 0)
	q.stopQueueWorkers()
}

func (s *GomolSuite) TestQueueStopWorkers(c *C) {
	q := newQueue()
	q.startQueueWorkers()

	q.stopQueueWorkers()
	c.Check(q.running, Equals, false)
	c.Check(q.queue, HasLen, 0)
	c.Check(len(q.queueChan), Equals, 0)
}

func (s *GomolSuite) TestQueueFlushMessages(c *C) {
	q := newQueue()
	q.startQueueWorkers()

	for i := 0; i < 100; i++ {
		q.QueueMessage(newMessage(testBase, levelDbg, nil, "test"))
	}

	q.stopQueueWorkers()
	c.Check(q.running, Equals, false)
	c.Check(q.queue, HasLen, 0)
	c.Check(len(q.queueChan), Equals, 0)
}
