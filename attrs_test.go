package gomol

import . "gopkg.in/check.v1"

func (s *GomolSuite) TestNewAttrsFromMap(c *C) {
	attrs := NewAttrsFromMap(map[string]interface{}{
		"attr1": "val1",
		"attr2": 1234,
	})
	c.Check(attrs.attrs, HasLen, 2)
	c.Check(attrs.attrs[getAttrHash("attr1")], Equals, "val1")
	c.Check(attrs.attrs[getAttrHash("attr2")], Equals, 1234)
}

func (s *GomolSuite) TestAttrsMergeNilAttrs(c *C) {
	attrs := NewAttrs()
	attrs.MergeAttrs(nil)
}

func (s *GomolSuite) TestAttrsGetMissing(c *C) {
	attrs := NewAttrs()
	c.Check(attrs.GetAttr("not an attr"), IsNil)
}

func (s *GomolSuite) TestAttrsRemoveMissing(c *C) {
	attrs := NewAttrs()
	// Just run it to make sure it doesn't panic
	attrs.RemoveAttr("not an attr")
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

func (s *GomolSuite) TestGetHashAttrMissing(c *C) {
	res, err := getHashAttr(1234)

	c.Check(res, Equals, "")
	c.Check(err.Error(), Equals, "Could not find attr for hash 1234")
}
