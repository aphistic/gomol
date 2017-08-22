package gomol

import (
	"time"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

const TestMaxQueueSize = 5000

func (s *GomolSuite) TestPressure(t sweet.T) {
	q := newQueue(NewBase(), TestMaxQueueSize)
	Expect(q.pressure()).To(Equal(0))
	q.queueChan <- &Message{}
	q.queueChan <- &Message{}
	Expect(q.pressure()).To(Equal(2))
}

func (s *GomolSuite) TestQueueMessageWithoutWorker(t sweet.T) {
	q := newQueue(NewBase(), TestMaxQueueSize)
	err := q.queueMessage(&Message{})
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("the logging system is not running - has InitLoggers() been executed?"))
}

func (s *GomolSuite) TestQueueStartWorker(t sweet.T) {
	q := newQueue(NewBase(), TestMaxQueueSize)
	q.startWorker()
	Expect(q.pressure()).To(Equal(0))
	q.stopWorker()
}

func (s *GomolSuite) TestQueueStartWorkerTwice(t sweet.T) {
	q := newQueue(NewBase(), TestMaxQueueSize)
	err := q.startWorker()
	Expect(err).To(BeNil())
	Expect(q.pressure()).To(Equal(0))
	err = q.startWorker()
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("workers are already running"))
	q.stopWorker()
}

func (s *GomolSuite) TestQueueStopWorker(t sweet.T) {
	q := newQueue(NewBase(), TestMaxQueueSize)
	q.startWorker()

	q.stopWorker()
	Expect(q.pressure()).To(Equal(0))
}

func (s *GomolSuite) TestQueueStopWorkerTwice(t sweet.T) {
	q := newQueue(NewBase(), TestMaxQueueSize)
	q.startWorker()

	err := q.stopWorker()
	Expect(err).To(BeNil())
	Expect(q.pressure()).To(Equal(0))
	err = q.stopWorker()
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("workers are not running"))
}

func (s *GomolSuite) TestQueueFlushMessages(t sweet.T) {
	q := newQueue(NewBase(), TestMaxQueueSize)
	q.startWorker()

	for i := 0; i < 100; i++ {
		q.queueMessage(newMessage(time.Now(), testBase, LevelDebug, nil, "test"))
	}

	q.stopWorker()
	Expect(q.pressure()).To(Equal(0))
}
