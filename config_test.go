package gomol

import (
	"testing"

	. "github.com/onsi/gomega"
)

func (s *GomolSuite) TestNewConfig(t *testing.T) {
	cfg := NewConfig()
	Expect(cfg.FilenameAttr).To(Equal(""))
	Expect(cfg.LineNumberAttr).To(Equal(""))
}
