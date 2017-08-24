package gomol

import (
	"testing"
	"time"

	. "github.com/onsi/gomega"

	"github.com/aphistic/sweet"
	junit "github.com/aphistic/sweet-junit"
)

func TestMain(m *testing.M) {
	RegisterFailHandler(sweet.GomegaFail)

	sweet.Run(m, func(s *sweet.S) {
		s.RegisterPlugin(junit.NewPlugin())

		s.AddSuite(&GomolSuite{})
		s.AddSuite(&AttrsSuite{})
		s.AddSuite(&BaseSuite{})
		s.AddSuite(&LogAdapterSuite{})
		s.AddSuite(&IssueSuite{})
		s.AddSuite(&LogLevelSuite{})
	})
}

type GomolSuite struct{}

func (s *GomolSuite) SetUpTest(t sweet.T) {
	setFakeCallerInfo("", 0)
	setClock(newTestClock(time.Now()))
	gomolFiles = map[string]fileRecord{}

	curTestExiter = &testExiter{}
	setExiter(curTestExiter)

	testBase = NewBase()
	testBase.AddLogger(newDefaultMemLogger())
	testBase.InitLoggers()

	curDefault = NewBase()
	curDefault.AddLogger(newDefaultMemLogger())
	curDefault.InitLoggers()
}

func (s *GomolSuite) TearDownTest(t sweet.T) {
	curDefault.ShutdownLoggers()

	testBase.ShutdownLoggers()
}
