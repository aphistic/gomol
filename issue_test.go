package gomol

import (
	"time"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type IssueSuite struct{}

func (s *IssueSuite) TestIssue22(t sweet.T) {
	setClock(newTestClock(time.Unix(10, 0)))

	mlCfg := newMemLoggerConfig()
	ml, err := newMemLogger(mlCfg)
	Expect(err).To(BeNil())
	Expect(ml).ToNot(BeNil())

	base := NewBase()
	base.AddLogger(ml)
	err = base.InitLoggers()

	Expect(err).To(BeNil())
	defer base.ShutdownLoggers()

	base.SetAttr("baseKey", 1234)

	err = base.Dbgm(
		NewAttrsFromMap(map[string]interface{}{
			"mapKey": 1234,
		}),
		"Message 1",
	)
	Expect(err).To(BeNil())

	base.Flush()

	Expect(ml.Messages()).To(HaveLen(1))
	Expect(ml.Messages()[0]).To(Equal(&memMessage{
		Timestamp: time.Unix(10, 0),
		Level:     LevelDebug,
		Message:   "Message 1",
		Attrs: map[string]interface{}{
			"baseKey": 1234,
			"mapKey":  1234,
		},
		StringAttrs: map[string]string{
			"baseKey": "1234",
			"mapKey":  "1234",
		},
	}))

	err = base.Dbgm(
		NewAttrsFromMap(map[string]interface{}{
			"mapFunc": func() int {
				return 4321
			},
		}),
		"Message 2",
	)
	Expect(err).To(BeNil())

	base.Flush()

	Expect(ml.Messages()).To(HaveLen(2))
	Expect(ml.Messages()[1]).To(Equal(&memMessage{
		Timestamp: time.Unix(10, 0),
		Level:     LevelDebug,
		Message:   "Message 2",
		Attrs: map[string]interface{}{
			"baseKey": 1234,
			"mapFunc": "func() int",
		},
		StringAttrs: map[string]string{
			"baseKey": "1234",
			"mapFunc": "func() int",
		},
	}))
}
