package gomol

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/aphistic/sweet"
	junit "github.com/aphistic/sweet-junit"
)

func TestMain(m *testing.M) {
	RegisterFailHandler(sweet.GomegaFail)

	sweet.Run(m, func(s *sweet.S) {
		s.RegisterPlugin(junit.NewPlugin())

		s.AddSuite(&AttrsSuite{})
		s.AddSuite(&BaseSuite{})
		s.AddSuite(&DefaultSuite{})
		s.AddSuite(&FallbackLoggerSuite{})
		s.AddSuite(&GomolSuite{})
		s.AddSuite(&IssueSuite{})
		s.AddSuite(&LogAdapterSuite{})
		s.AddSuite(&LogLevelSuite{})
		s.AddSuite(&MemLoggerSuite{})
	})
}

type GomolSuite struct{}

func (s *GomolSuite) SetUpTest(t sweet.T) {
	setFakeCallerInfo("", 0)
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
