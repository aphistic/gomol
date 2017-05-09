package gomol

import (
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func (s *GomolSuite) TestQueueLenChan(t *testing.T) {
	q := newQueue()

	Expect(q.Length()).To(Equal(0))
	q.queueChan <- &Message{}
	Expect(q.Length()).To(Equal(1))
}

func (s *GomolSuite) TestQueueLenArray(t *testing.T) {
	q := newQueue()

	Expect(q.Length()).To(Equal(0))
	q.queue = append(q.queue, &Message{})
	Expect(q.Length()).To(Equal(1))
}

func (s *GomolSuite) TestQueueLen(t *testing.T) {
	q := newQueue()

	Expect(q.Length()).To(Equal(0))
	q.queueChan <- &Message{}
	q.queue = append(q.queue, &Message{})
	Expect(q.Length()).To(Equal(2))
}

func (s *GomolSuite) TestQueueMessageWithoutWorkers(t *testing.T) {
	q := newQueue()
	err := q.QueueMessage(&Message{})
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("The logging system is not running, has InitLoggers() been executed?"))
}

func (s *GomolSuite) TestQueueStartWorkers(t *testing.T) {
	q := newQueue()
	q.startQueueWorkers()
	Expect(q.IsActive()).To(Equal(true))
	Expect(q.queue).To(HaveLen(0))
	q.stopQueueWorkers()
}

func (s *GomolSuite) TestQueueStartWorkersTwice(t *testing.T) {
	q := newQueue()
	err := q.startQueueWorkers()
	Expect(err).To(BeNil())
	Expect(q.IsActive()).To(Equal(true))
	Expect(q.queue).To(HaveLen(0))
	err = q.startQueueWorkers()
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Workers are already running"))
	q.stopQueueWorkers()
}

func (s *GomolSuite) TestQueueStopWorkers(t *testing.T) {
	q := newQueue()
	q.startQueueWorkers()

	q.stopQueueWorkers()
	Expect(q.IsActive()).To(Equal(false))
	Expect(q.queue).To(HaveLen(0))
	Expect(len(q.queueChan)).To(Equal(0))
}

func (s *GomolSuite) TestQueueStopWorkersTwice(t *testing.T) {
	q := newQueue()
	q.startQueueWorkers()

	err := q.stopQueueWorkers()
	Expect(err).To(BeNil())
	Expect(q.IsActive()).To(Equal(false))
	Expect(q.queue).To(HaveLen(0))
	Expect(len(q.queueChan)).To(Equal(0))
	err = q.stopQueueWorkers()
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Workers are not running"))
}

func (s *GomolSuite) TestQueueFlushMessages(t *testing.T) {
	q := newQueue()
	q.startQueueWorkers()

	for i := 0; i < 100; i++ {
		q.QueueMessage(newMessage(time.Now(), testBase, LevelDebug, nil, "test"))
	}

	q.stopQueueWorkers()
	Expect(q.IsActive()).To(Equal(false))
	Expect(q.queue).To(HaveLen(0))
	Expect(len(q.queueChan)).To(Equal(0))
}
