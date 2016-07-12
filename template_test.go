package gomol

import (
	"encoding/json"
	"time"

	. "gopkg.in/check.v1"
)

func (s *GomolSuite) TestNewTemplate(c *C) {
	tpl, err := NewTemplate("{{bad template}")
	c.Check(tpl, IsNil)
	c.Check(err, NotNil)
}

func (s *GomolSuite) TestExecInternalMsg(c *C) {
	tpl, err := NewTemplate("test")
	c.Assert(err, IsNil)
	c.Assert(tpl, NotNil)

	data, err := tpl.executeInternalMsg(nil, false)
	c.Check(data, Equals, "")
	c.Check(err, NotNil)
}

func (s *GomolSuite) TestNewTemplateMsgError(c *C) {
	tplMsg, err := newTemplateMsgFromMessage(nil)
	c.Check(tplMsg, IsNil)
	c.Check(err, NotNil)
}

func (s *GomolSuite) TestTemplateExecuteError(c *C) {
	setClock(newTestClock(time.Unix(1000000000, 100)))

	msg := newMessage(clock().Now(), nil, LEVEL_ERROR, map[string]interface{}{
		"attr1": "val1",
		"attr2": 1234,
	}, "message")
	tpl, err := NewTemplate("{{ .ThisDoesNotExist }}")
	c.Assert(err, IsNil)

	tplMsg, err := newTemplateMsgFromMessage(msg)
	c.Assert(err, IsNil)

	out, err := tpl.Execute(tplMsg, false)
	c.Check(out, Equals, "")
	c.Check(err, NotNil)
}

func (s *GomolSuite) TestTplColorsDebug(c *C) {
	setClock(newTestClock(time.Now()))
	msg := newMessage(clock().Now(), nil, LEVEL_DEBUG, nil, "colors!")
	tpl, err := NewTemplate("{{color}}hascolor{{reset}} {{.Message}}")
	c.Assert(err, IsNil)

	out, err := tpl.executeInternalMsg(msg, true)
	c.Assert(err, IsNil)

	c.Check(out, Equals, "\x1b[36mhascolor\x1b[0m colors!")
}

func (s *GomolSuite) TestTplColorsInfo(c *C) {
	setClock(newTestClock(time.Now()))
	msg := newMessage(clock().Now(), nil, LEVEL_INFO, nil, "colors!")
	tpl, err := NewTemplate("{{color}}hascolor{{reset}} {{.Message}}")
	c.Assert(err, IsNil)

	out, err := tpl.executeInternalMsg(msg, true)
	c.Assert(err, IsNil)

	c.Check(out, Equals, "\x1b[32mhascolor\x1b[0m colors!")
}

func (s *GomolSuite) TestTplColorsWarning(c *C) {
	setClock(newTestClock(time.Now()))
	msg := newMessage(clock().Now(), nil, LEVEL_WARNING, nil, "colors!")
	tpl, err := NewTemplate("{{color}}hascolor{{reset}} {{.Message}}")
	c.Assert(err, IsNil)

	out, err := tpl.executeInternalMsg(msg, true)
	c.Assert(err, IsNil)

	c.Check(out, Equals, "\x1b[33mhascolor\x1b[0m colors!")
}

func (s *GomolSuite) TestTplColorsError(c *C) {
	setClock(newTestClock(time.Now()))
	msg := newMessage(clock().Now(), nil, LEVEL_ERROR, nil, "colors!")
	tpl, err := NewTemplate("{{color}}hascolor{{reset}} {{.Message}}")
	c.Assert(err, IsNil)

	out, err := tpl.executeInternalMsg(msg, true)
	c.Assert(err, IsNil)

	c.Check(out, Equals, "\x1b[31mhascolor\x1b[0m colors!")
}

func (s *GomolSuite) TestTplColorsFatal(c *C) {
	setClock(newTestClock(time.Now()))
	msg := newMessage(clock().Now(), nil, LEVEL_FATAL, nil, "colors!")
	tpl, err := NewTemplate("{{color}}hascolor{{reset}} {{.Message}}")
	c.Assert(err, IsNil)

	out, err := tpl.executeInternalMsg(msg, true)
	c.Assert(err, IsNil)

	c.Check(out, Equals, "\x1b[1;31mhascolor\x1b[0m colors!")
}

func (s *GomolSuite) TestTplFuncsCase(c *C) {
	setClock(newTestClock(time.Now()))
	msg := newMessage(clock().Now(), nil, LEVEL_ERROR, nil, "UPPER")
	tpl, err := NewTemplate("{{ucase .LevelName}} {{lcase .Message}} {{title .LevelName}}")
	c.Assert(err, IsNil)

	out, err := tpl.executeInternalMsg(msg, false)
	c.Assert(err, IsNil)

	c.Check(out, Equals, "ERROR upper Error")
}

func (s *GomolSuite) TestTplMsgFromInternal(c *C) {
	setClock(newTestClock(time.Now()))

	b := NewBase()
	b.SetAttr("baseAttr", 1234)
	b.SetAttr("overrideAttr", 1234)
	msg := newMessage(clock().Now(), b, LEVEL_INFO, map[string]interface{}{
		"msgAttr":      4321,
		"overrideAttr": "test",
	}, "Format %v %v", 1234, "asdf")

	tplMsg, err := newTemplateMsgFromMessage(msg)
	c.Assert(err, IsNil)
	c.Check(tplMsg.Timestamp, Equals, clock().Now())
	c.Check(tplMsg.Level, Equals, LEVEL_INFO)
	c.Check(tplMsg.LevelName, Equals, "info")
	c.Check(tplMsg.Message, Equals, "Format 1234 asdf")
	c.Assert(tplMsg.Attrs, HasLen, 3)
	c.Check(tplMsg.Attrs["baseAttr"], Equals, 1234)
	c.Check(tplMsg.Attrs["overrideAttr"], Equals, "test")
	c.Check(tplMsg.Attrs["msgAttr"], Equals, 4321)
}

func (s *GomolSuite) TestTplMsgAttrs(c *C) {
	setClock(newTestClock(time.Now()))

	b := NewBase()
	b.SetAttr("baseAttr", 1234)
	b.SetAttr("overrideAttr", 1234)
	msg := newMessage(clock().Now(), b, LEVEL_INFO, map[string]interface{}{
		"msgAttr":      4321,
		"overrideAttr": "test",
	}, "Format %v %v", 1234, "asdf")

	tplMsg, err := newTemplateMsgFromMessage(msg)
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
	setClock(newTestClock(time.Now()))

	msg := newMessage(clock().Now(), nil, LEVEL_ERROR, nil, "message")
	tpl, err := NewTemplate("{{ .Timestamp.Format \"2006-01-02T15:04:05.999999999Z07:00\" }}")
	c.Assert(err, IsNil)

	out, err := tpl.executeInternalMsg(msg, false)
	c.Assert(err, IsNil)

	c.Check(out, Equals, clock().Now().Format("2006-01-02T15:04:05.999999999Z07:00"))
}

func (s *GomolSuite) TestTplJson(c *C) {
	setClock(newTestClock(time.Unix(1000000000, 100)))

	msg := newMessage(clock().Now(), nil, LEVEL_ERROR, map[string]interface{}{
		"attr1": "val1",
		"attr2": 1234,
	}, "message")
	tpl, err := NewTemplate("{{ json . }}")
	c.Assert(err, IsNil)

	tplMsg, err := newTemplateMsgFromMessage(msg)
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

func (s *GomolSuite) TestTplAttrTemplate(c *C) {
	setClock(newTestClock(time.Now()))
	msg := newMessage(clock().Now(),
		nil,
		LEVEL_FATAL,
		map[string]interface{}{"attrName": "attrVal"},
		"test")
	tpl, err := NewTemplate("[{{.Attrs.attrName}}] {{.Message}}")
	c.Assert(err, IsNil)

	out, err := tpl.executeInternalMsg(msg, true)
	c.Assert(err, IsNil)

	c.Check(out, Equals, "[attrVal] test")
}

func (s *GomolSuite) TestTplJsonError(c *C) {
	data, err := tplJson(map[string]interface{}{
		"attr1": s.SetUpTest,
		"attr2": s.TearDownTest,
	})
	c.Check(data, Equals, "")
	c.Check(err, NotNil)
}

func (s *GomolSuite) TestNewTemplateMsgMinimal(c *C) {
	setClock(newTestClock(time.Now()))

	tmsg := NewTemplateMsg(clock().Now(), LEVEL_DEBUG, nil, "test")
	c.Check(tmsg.Timestamp, Equals, clock().Now())
	c.Check(tmsg.Level, Equals, LEVEL_DEBUG)
	c.Check(tmsg.LevelName, Equals, LEVEL_DEBUG.String())
	c.Check(tmsg.Attrs, DeepEquals, map[string]interface{}{})
	c.Check(tmsg.Message, Equals, "test")
}
