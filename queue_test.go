package gomol

import (
	. "gopkg.in/check.v1"
)

func (s *GomolSuite) TestQueueLenChan(c *C) {
	// Flush all messages so the workers leave things
	// in the channel/queue
	FlushMessages()

	q := newQueue()

	c.Check(q.Length(), Equals, 0)
	q.queueChan <- &message{}
	c.Check(q.Length(), Equals, 1)
}

func (s *GomolSuite) TestQueueLenArray(c *C) {
	// Flush all messages so the workers leave things
	// in the channel/queue
	FlushMessages()

	q := newQueue()

	c.Check(q.Length(), Equals, 0)
	q.queue = append(q.queue, &message{})
	c.Check(q.Length(), Equals, 1)
}

func (s *GomolSuite) TestQueueLen(c *C) {
	// Flush all messages so the workers leave things
	// in the channel/queue
	FlushMessages()

	q := newQueue()

	c.Check(q.Length(), Equals, 0)
	q.queueChan <- &message{}
	q.queue = append(q.queue, &message{})
	c.Check(q.Length(), Equals, 2)
}
