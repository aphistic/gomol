package gomol

import (
	. "gopkg.in/check.v1"
)

func (s *GomolSuite) TestTplFuncsCase(c *C) {
	msg := newMessage(nil, LEVEL_ERROR, nil, "UPPER")
	tpl, err := NewTemplate("{{ucase .Level}} {{lcase .Message}} {{title .Level}}")
	c.Assert(err, IsNil)

	out, err := tpl.executeInternalMsg(msg)
	c.Assert(err, IsNil)

	c.Check(out, Equals, "ERROR upper Error")
}

func (s *GomolSuite) TestTplMsgFromInternal(c *C) {
	b := newBase()
	b.SetAttr("baseAttr", 1234)
	b.SetAttr("overrideAttr", 1234)
	msg := newMessage(b, LEVEL_INFO, map[string]interface{}{
		"msgAttr":      4321,
		"overrideAttr": "test",
	}, "Format %v %v", 1234, "asdf")

	tplMsg, err := newTemplateMsg(msg)
	c.Assert(err, IsNil)
	c.Check(tplMsg.Level(), Equals, "info")
	c.Check(tplMsg.Message(), Equals, "Format 1234 asdf")
	c.Assert(tplMsg.Attrs(), HasLen, 3)
	c.Check(tplMsg.Attrs()["baseAttr"], Equals, 1234)
	c.Check(tplMsg.Attrs()["overrideAttr"], Equals, "test")
	c.Check(tplMsg.Attrs()["msgAttr"], Equals, 4321)
}

func (s *GomolSuite) TestTplMsgAttrs(c *C) {
	b := newBase()
	b.SetAttr("baseAttr", 1234)
	b.SetAttr("overrideAttr", 1234)
	msg := newMessage(b, LEVEL_INFO, map[string]interface{}{
		"msgAttr":      4321,
		"overrideAttr": "test",
	}, "Format %v %v", 1234, "asdf")

	tplMsg, err := newTemplateMsg(msg)
	c.Assert(err, IsNil)
	c.Check(tplMsg.Level(), Equals, "info")
	c.Check(tplMsg.Message(), Equals, "Format 1234 asdf")
	c.Assert(tplMsg.Attrs(), HasLen, 3)
	c.Check(tplMsg.Attrs()["baseAttr"], Equals, 1234)
	c.Check(tplMsg.Attrs()["overrideAttr"], Equals, "test")
	c.Check(tplMsg.Attrs()["msgAttr"], Equals, 4321)

	tpl, err := NewTemplate("{{range $key, $val := .Attrs}}{{$key}}=={{$val}}\n{{end}}")
	c.Assert(err, IsNil)

	out, err := tpl.executeInternalMsg(msg)
	c.Assert(err, IsNil)

	c.Check(out, Equals, "baseAttr==1234\nmsgAttr==4321\noverrideAttr==test\n")
}
