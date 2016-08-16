package gomol

import . "gopkg.in/check.v1"

func (s *GomolSuite) BenchmarkAttrChaining(c *C) {
	for idx := 0; idx < c.N; idx++ {
		NewAttrs().
			SetAttr("attr1", "val1").
			SetAttr("attr2", "val2").
			SetAttr("attr3", 3).
			SetAttr("attr4", 4)
	}
}
