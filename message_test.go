package gomol

import (
	. "gopkg.in/check.v1"
)

func (s *GomolSuite) TestNewMessageAttrsNil(c *C) {
	m := newMessage(curDefault, LEVEL_DEBUG, nil, "test")
	c.Check(m.Base, DeepEquals, curDefault)
	c.Check(m.Level, Equals, LEVEL_DEBUG)
	c.Check(m.Attrs, NotNil)
	c.Check(m.Attrs, HasLen, 0)
	c.Check(m.Msg, Equals, "test")
}

func (s *GomolSuite) TestNewMessageMsgAttrsNil(c *C) {
	ma := map[string]interface{}{
		"msgAttr":   "strVal",
		"otherAttr": 4321,
	}

	m := newMessage(curDefault, LEVEL_DEBUG, ma, "test")
	c.Check(m.Base, DeepEquals, curDefault)
	c.Check(m.Level, Equals, LEVEL_DEBUG)
	c.Check(m.Attrs, NotNil)
	c.Check(m.Attrs, HasLen, 2)
	c.Check(m.Attrs["msgAttr"], Equals, "strVal")
	c.Check(m.Attrs["otherAttr"], Equals, 4321)
	c.Check(m.Msg, Equals, "test")
}

func (s *GomolSuite) TestNewMessageFormat(c *C) {
	m := newMessage(curDefault, LEVEL_DEBUG, nil, "test %v %v", "str", 1234)
	c.Check(m.Base, DeepEquals, curDefault)
	c.Check(m.Level, Equals, LEVEL_DEBUG)
	c.Check(m.Attrs, NotNil)
	c.Check(m.Attrs, HasLen, 0)
	c.Check(m.Msg, Equals, "test str 1234")
}

func (s *GomolSuite) TestNewMessageFormatWithAttrs(c *C) {
	ma := map[string]interface{}{
		"msgAttr":   "strVal",
		"otherAttr": 4321,
	}

	m := newMessage(curDefault, LEVEL_DEBUG, ma, "test %v %v", "str", 1234)
	c.Check(m.Base, DeepEquals, curDefault)
	c.Check(m.Level, Equals, LEVEL_DEBUG)
	c.Check(m.Attrs, NotNil)
	c.Check(m.Attrs, HasLen, 2)
	c.Check(m.Attrs["msgAttr"], Equals, "strVal")
	c.Check(m.Attrs["otherAttr"], Equals, 4321)
	c.Check(m.Msg, Equals, "test str 1234")
}

func (s *GomolSuite) TestMergeAttrsNil(c *C) {
	attrs := mergeAttrs(nil, nil)
	c.Check(attrs, NotNil)
	c.Check(attrs, HasLen, 0)
}

func (s *GomolSuite) TestMergeAttrsBaseAttrs(c *C) {
	ba := map[string]interface{}{
		"baseAttr":  "strVal",
		"otherAttr": 1234,
	}
	attrs := mergeAttrs(ba, nil)
	c.Check(attrs, NotNil)
	c.Check(attrs, HasLen, 2)
	c.Check(attrs["baseAttr"], Equals, "strVal")
	c.Check(attrs["otherAttr"], Equals, 1234)
}

func (s *GomolSuite) TestMergeAttrsMsgAttrs(c *C) {
	ma := map[string]interface{}{
		"msgAttr":   "strVal",
		"otherAttr": 4321,
	}
	attrs := mergeAttrs(nil, ma)
	c.Check(attrs, NotNil)
	c.Check(attrs, HasLen, 2)
	c.Check(attrs["msgAttr"], Equals, "strVal")
	c.Check(attrs["otherAttr"], Equals, 4321)
}

func (s *GomolSuite) TestMergeAttrs(c *C) {
	ba := map[string]interface{}{
		"baseAttr":  "strVal",
		"otherAttr": 1234,
	}
	ma := map[string]interface{}{
		"msgAttr":   "strVal",
		"otherAttr": 4321,
	}
	attrs := mergeAttrs(ba, ma)
	c.Check(attrs, NotNil)
	c.Check(attrs, HasLen, 3)
	c.Check(attrs["msgAttr"], Equals, "strVal")
	c.Check(attrs["baseAttr"], Equals, "strVal")
	c.Check(attrs["otherAttr"], Equals, 4321)
}
