package gomol

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

func (s *GomolSuite) TestNewConfig(t sweet.T) {
	cfg := NewConfig()
	Expect(cfg.FilenameAttr).To(Equal(""))
	Expect(cfg.LineNumberAttr).To(Equal(""))
}
