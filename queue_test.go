package gomol

import (
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func (s *GomolSuite) TestPressure(t *testing.T) {
	q := newQueue()
	Expect(q.pressure()).To(Equal(0))
	q.queueChan <- &Message{}
	q.queueChan <- &Message{}
	Expect(q.pressure()).To(Equal(2))
}

func (s *GomolSuite) TestQueueMessageWithoutWorker(t *testing.T) {
	q := newQueue()
	err := q.queueMessage(&Message{})
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("the logging system is not running - has InitLoggers() been executed?"))
}

func (s *GomolSuite) TestQueueStartWorker(t *testing.T) {
	q := newQueue()
	q.startWorker()
	Expect(q.pressure()).To(Equal(0))
	q.stopWorker()
}

func (s *GomolSuite) TestQueueStartWorkerTwice(t *testing.T) {
	q := newQueue()
	err := q.startWorker()
	Expect(err).To(BeNil())
	Expect(q.pressure()).To(Equal(0))
	err = q.startWorker()
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("workers are already running"))
	q.stopWorker()
}

func (s *GomolSuite) TestQueueStopWorker(t *testing.T) {
	q := newQueue()
	q.startWorker()

	q.stopWorker()
	Expect(q.pressure()).To(Equal(0))
}

func (s *GomolSuite) TestQueueStopWorkerTwice(t *testing.T) {
	q := newQueue()
	q.startWorker()

	err := q.stopWorker()
	Expect(err).To(BeNil())
	Expect(q.pressure()).To(Equal(0))
	err = q.stopWorker()
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("workers are not running"))
}

func (s *GomolSuite) TestQueueFlushMessages(t *testing.T) {
	q := newQueue()
	q.startWorker()

	for i := 0; i < 100; i++ {
		q.queueMessage(newMessage(time.Now(), testBase, LevelDebug, nil, "test"))
	}

	q.stopWorker()
	Expect(q.pressure()).To(Equal(0))
}
