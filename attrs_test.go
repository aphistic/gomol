package gomol

import . "gopkg.in/check.v1"

func (s *GomolSuite) TestAttrsMergeNilAttrs(c *C) {
	attrs := NewAttrs()
	attrs.mergeAttrs(nil)
}

func (s *GomolSuite) TestAttrsChaining(c *C) {
	attrs := NewAttrs().
		SetAttr("attr1", "val1").
		SetAttr("attr2", "val2").
		SetAttr("attr3", 3).
		SetAttr("attr4", 4)

	c.Check(attrs.attrs[getAttrHash("attr1")], Equals, "val1")
	c.Check(attrs.attrs[getAttrHash("attr2")], Equals, "val2")
	c.Check(attrs.attrs[getAttrHash("attr3")], Equals, 3)
	c.Check(attrs.attrs[getAttrHash("attr4")], Equals, 4)
}