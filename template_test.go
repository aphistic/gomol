package gomol

import (
	"encoding/json"
	. "gopkg.in/check.v1"
	"time"
)

func (s *GomolSuite) TestTplColors(c *C) {
	msg := newMessage(nil, LEVEL_FATAL, nil, "colors!")
	tpl, err := NewTemplate("{{color}}hascolor{{reset}} {{.Message}}")
	c.Assert(err, IsNil)

	out, err := tpl.executeInternalMsg(msg, true)
	c.Assert(err, IsNil)

	c.Check(out, Equals, "\x1b[1;31mhascolor\x1b[0m colors!")
}

func (s *GomolSuite) TestTplFuncsCase(c *C) {
	msg := newMessage(nil, LEVEL_ERROR, nil, "UPPER")
	tpl, err := NewTemplate("{{ucase .LevelName}} {{lcase .Message}} {{title .LevelName}}")
	c.Assert(err, IsNil)

	out, err := tpl.executeInternalMsg(msg, false)
	c.Assert(err, IsNil)

	c.Check(out, Equals, "ERROR upper Error")
}

func (s *GomolSuite) TestTplMsgFromInternal(c *C) {
	clock = NewTestClock(time.Now())

	b := newBase()
	b.SetAttr("baseAttr", 1234)
	b.SetAttr("overrideAttr", 1234)
	msg := newMessage(b, LEVEL_INFO, map[string]interface{}{
		"msgAttr":      4321,
		"overrideAttr": "test",
	}, "Format %v %v", 1234, "asdf")

	tplMsg, err := newTemplateMsg(msg)
	c.Assert(err, IsNil)
	c.Check(tplMsg.Timestamp, Equals, clock.Now())
	c.Check(tplMsg.Level, Equals, LEVEL_INFO)
	c.Check(tplMsg.LevelName, Equals, "info")
	c.Check(tplMsg.Message, Equals, "Format 1234 asdf")
	c.Assert(tplMsg.Attrs, HasLen, 3)
	c.Check(tplMsg.Attrs["baseAttr"], Equals, 1234)
	c.Check(tplMsg.Attrs["overrideAttr"], Equals, "test")
	c.Check(tplMsg.Attrs["msgAttr"], Equals, 4321)
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
	c.Check(tplMsg.Level, Equals, LEVEL_INFO)
	c.Check(tplMsg.LevelName, Equals, "info")
	c.Check(tplMsg.Message, Equals, "Format 1234 asdf")
	c.Assert(tplMsg.Attrs, HasLen, 3)
	c.Check(tplMsg.Attrs["baseAttr"], Equals, 1234)
	c.Check(tplMsg.Attrs["overrideAttr"], Equals, "test")
	c.Check(tplMsg.Attrs["msgAttr"], Equals, 4321)

	tpl, err := NewTemplate("{{range $key, $val := .Attrs}}{{$key}}=={{$val}}\n{{end}}")
	c.Assert(err, IsNil)

	out, err := tpl.executeInternalMsg(msg, false)
	c.Assert(err, IsNil)

	c.Check(out, Equals, "baseAttr==1234\nmsgAttr==4321\noverrideAttr==test\n")
}

func (s *GomolSuite) TestTplTimestamp(c *C) {
	clock = NewTestClock(time.Now())

	msg := newMessage(nil, LEVEL_ERROR, nil, "message")
	tpl, err := NewTemplate("{{ .Timestamp.Format \"2006-01-02T15:04:05.999999999Z07:00\" }}")
	c.Assert(err, IsNil)

	out, err := tpl.executeInternalMsg(msg, false)
	c.Assert(err, IsNil)

	c.Check(out, Equals, clock.Now().Format("2006-01-02T15:04:05.999999999Z07:00"))
}

func (s *GomolSuite) TestTplJson(c *C) {
	clock = NewTestClock(time.Unix(1000000000, 100))

	msg := newMessage(nil, LEVEL_ERROR, map[string]interface{}{
		"attr1": "val1",
		"attr2": 1234,
	}, "message")
	tpl, err := NewTemplate("{{ json . }}")
	c.Assert(err, IsNil)

	tplMsg, err := newTemplateMsg(msg)
	c.Assert(err, IsNil)

	out, err := tpl.Execute(tplMsg, false)
	c.Assert(err, IsNil)

	// Unmarshal from json and check that because on Travis the timezone is different
	// and I don't want to create a new version of time.Time to marshal the value
	// differently
	dataOut := &TemplateMsg{}
	err = json.Unmarshal([]byte(out), dataOut)
	c.Assert(err, IsNil)

	c.Check(dataOut.Timestamp.UnixNano(), Equals, msg.Timestamp.UnixNano())
	c.Check(dataOut.Level, Equals, tplMsg.Level)
	c.Check(dataOut.LevelName, Equals, tplMsg.LevelName)
	c.Check(dataOut.Message, Equals, tplMsg.Message)
	c.Check(dataOut.Attrs, HasLen, 2)
	c.Check(dataOut.Attrs["attr1"], Equals, "val1")
	c.Check(dataOut.Attrs["attr2"], Equals, float64(1234))

	tpl, err = NewTemplate("{{ json .Attrs }}")
	c.Assert(err, IsNil)

	out, err = tpl.executeInternalMsg(msg, false)
	c.Assert(err, IsNil)

	c.Check(out, Equals, "{\"attr1\":\"val1\",\"attr2\":1234}")
}
