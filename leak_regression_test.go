package gomol

import (
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func (s *GomolSuite) TestLeakRegressionTest(t *testing.T) {
	var (
		blocker = make(chan struct{})
		l1      = newDefaultMemLogger()
		l2      = &BlockingLogger{ch: blocker}
		ch      = make(chan error)
		errors  = make(chan int)
	)

	go func() {
		defer close(errors)

		count := 0
		for range ch {
			count++
		}

		fmt.Printf("[LRT] error channel closed\n")
		errors <- count
	}()

	fmt.Printf("[LRT] setting up base wiith queue size of %d\n", TestMaxQueueSize)

	testBase = NewBase()
	testBase.SetConfig(&Config{MaxQueueSize: TestMaxQueueSize})
	testBase.SetErrorChan(ch)
	testBase.AddLogger(l1)
	testBase.AddLogger(l2)
	testBase.InitLoggers()

	fmt.Printf("[LRT] sending first chunk\n")

	for i := 0; i < TestMaxQueueSize; i++ {
		fmt.Printf("[LRT] i=%d\n", i)
		testBase.Infof("test %d", i)
	}

	// l1 gets a message, but then l2 blocks immediately
	// after. We should not have any additional messages
	// sent to the first logger.
	Expect(l1.Messages).To(HaveLen(1))
	Expect(l1.Messages[0].Message).To(Equal("test 0"))

	// Send additional messages while the loggers are blocked
	// and the queue is full. This should NOT block the main
	// app routine.

	fmt.Printf("[LRT] sending next two chunk\n")

	for i := TestMaxQueueSize; i < TestMaxQueueSize*3; i++ {
		fmt.Printf("[LRT] i=%d\n", i)
		testBase.Infof("test %d", i)
	}

	// Now, unblock the logger, publish another chunk, and
	// wait for the messages to drain so we can inspect what
	// messages made it through during the time the logger
	// was naughty.

	fmt.Printf("[LRT] unblocking\n")

	close(blocker)
	<-time.After(time.Millisecond * 100)

	fmt.Printf("[LRT] sending fourth chunk\n")

	for i := TestMaxQueueSize * 3; i < TestMaxQueueSize*4; i++ {
		fmt.Printf("[LRT] i=%d\n", i)
		testBase.Infof("test %d", i)
	}

	fmt.Printf("[LRT] shutting down loggers\n")

	testBase.ShutdownLoggers()

	fmt.Printf("[LRT] epilogue\n")

	// We ran into a situation where in order to keep the bound
	// of the queue some messages had to be dropped - these must
	// be the _oldest_ messages and should. In this case, it is
	// the first two chunks (except the first message). The next
	// message we should see that wasn't dropped should be the
	// first message in the third chunk.

	Expect(len(l1.Messages)).To(Equal(2*TestMaxQueueSize + 1))

	// Skip checking "test 0" message

	for i := 0; i < len(l1.Messages)-1; i++ {
		Expect(l1.Messages[i+1].Message).To(Equal(fmt.Sprintf("test %d", i+2*TestMaxQueueSize)))
	}

	Eventually(errors).Should(Receive(Equal(2*TestMaxQueueSize - 1)))
}

//
// Logger that blocks all messages

type BlockingLogger struct {
	ch chan struct{}
}

func (l *BlockingLogger) SetBase(base *Base)    {}
func (l *BlockingLogger) InitLogger() error     { return nil }
func (l *BlockingLogger) IsInitialized() bool   { return true }
func (l *BlockingLogger) ShutdownLogger() error { return nil }

func (l *BlockingLogger) Logm(timestamp time.Time, level LogLevel, attrs map[string]interface{}, msg string) error {
	fmt.Printf("[LRT] blocking\n")
	<-l.ch
	fmt.Printf("[LRT] unblocked\n")
	return nil
}
