package gomol

import . "gopkg.in/check.v1"

func (s *GomolSuite) BenchmarkNewAttrFromMap(c *C) {
	for idx := 0; idx < c.N; idx++ {
		NewAttrsFromMap(map[string]interface{}{
			"attr1": "val1",
			"attr2": "val2",
			"attr3": map[string]interface{}{
				"attr31": "val1",
				"attr32": 1234,
			},
			"attr4": 4321,
		})
	}
}

func (s *GomolSuite) BenchmarkAttrChaining(c *C) {
	for idx := 0; idx < c.N; idx++ {
		NewAttrs().
			SetAttr("attr1", "val1").
			SetAttr("attr2", "val2").
			SetAttr("attr3", 3).
			SetAttr("attr4", 4)
	}
}
