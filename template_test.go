package gomol

import (
	"encoding/json"
	"text/template"
	"time"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

func (s *GomolSuite) TestNewTemplate(t sweet.T) {
	tpl, err := NewTemplate("{{bad template}")
	Expect(tpl).To(BeNil())
	Expect(err).ToNot(BeNil())
}

func (s *GomolSuite) TestExecInternalMsg(t sweet.T) {
	tpl, err := NewTemplate("test")
	Expect(err).To(BeNil())
	Expect(tpl).ToNot(BeNil())

	data, err := tpl.executeInternalMsg(nil, false)
	Expect(data).To(Equal(""))
	Expect(err).ToNot(BeNil())
}

func (s *GomolSuite) TestNewTemplateMsgError(t sweet.T) {
	tplMsg, err := newTemplateMsgFromMessage(nil)
	Expect(tplMsg).To(BeNil())
	Expect(err).ToNot(BeNil())
}

func (s *GomolSuite) TestTemplateExecuteError(t sweet.T) {
	ts := time.Unix(1000000000, 100)
	msg := newMessage(ts, nil, LevelError,
		NewAttrs().
			SetAttr("attr1", "val1").
			SetAttr("attr2", 1234),
		"message")
	tpl, err := NewTemplate("{{ .ThisDoesNotExist }}")
	Expect(err).To(BeNil())

	tplMsg, err := newTemplateMsgFromMessage(msg)
	Expect(err).To(BeNil())

	out, err := tpl.Execute(tplMsg, false)
	Expect(out).To(Equal(""))
	Expect(err).ToNot(BeNil())
}

func (s *GomolSuite) TestTplColorsNone(t sweet.T) {
	msg := newMessage(time.Unix(10, 0), nil, LevelNone, nil, "colors!")
	tpl, err := NewTemplate("{{color}}hascolor{{reset}} {{.Message}}")
	Expect(err).To(BeNil())

	out, err := tpl.executeInternalMsg(msg, true)
	Expect(err).To(BeNil())

	Expect(out).To(Equal("hascolor colors!"))
}

func (s *GomolSuite) TestTplColorsDebug(t sweet.T) {
	msg := newMessage(time.Unix(10, 0), nil, LevelDebug, nil, "colors!")
	tpl, err := NewTemplate("{{color}}hascolor{{reset}} {{.Message}}")
	Expect(err).To(BeNil())

	out, err := tpl.executeInternalMsg(msg, true)
	Expect(err).To(BeNil())

	Expect(out).To(Equal("\x1b[36mhascolor\x1b[0m colors!"))
}

func (s *GomolSuite) TestTplColorsInfo(t sweet.T) {
	msg := newMessage(time.Unix(10, 0), nil, LevelInfo, nil, "colors!")
	tpl, err := NewTemplate("{{color}}hascolor{{reset}} {{.Message}}")
	Expect(err).To(BeNil())

	out, err := tpl.executeInternalMsg(msg, true)
	Expect(err).To(BeNil())

	Expect(out).To(Equal("\x1b[32mhascolor\x1b[0m colors!"))
}

func (s *GomolSuite) TestTplColorsWarning(t sweet.T) {
	msg := newMessage(time.Unix(10, 0), nil, LevelWarning, nil, "colors!")
	tpl, err := NewTemplate("{{color}}hascolor{{reset}} {{.Message}}")
	Expect(err).To(BeNil())

	out, err := tpl.executeInternalMsg(msg, true)
	Expect(err).To(BeNil())

	Expect(out).To(Equal("\x1b[33mhascolor\x1b[0m colors!"))
}

func (s *GomolSuite) TestTplColorsError(t sweet.T) {
	msg := newMessage(time.Unix(10, 0), nil, LevelError, nil, "colors!")
	tpl, err := NewTemplate("{{color}}hascolor{{reset}} {{.Message}}")
	Expect(err).To(BeNil())

	out, err := tpl.executeInternalMsg(msg, true)
	Expect(err).To(BeNil())

	Expect(out).To(Equal("\x1b[31mhascolor\x1b[0m colors!"))
}

func (s *GomolSuite) TestTplColorsFatal(t sweet.T) {
	msg := newMessage(time.Unix(10, 0), nil, LevelFatal, nil, "colors!")
	tpl, err := NewTemplate("{{color}}hascolor{{reset}} {{.Message}}")
	Expect(err).To(BeNil())

	out, err := tpl.executeInternalMsg(msg, true)
	Expect(err).To(BeNil())

	Expect(out).To(Equal("\x1b[1;31mhascolor\x1b[0m colors!"))
}

func (s *GomolSuite) TestTplColorsUnknown(t sweet.T) {
	msg := newMessage(time.Unix(10, 0), nil, LogLevel(-1000), nil, "colors!")
	tpl, err := NewTemplate("{{color}}hascolor{{reset}} {{.Message}}")
	Expect(err).To(BeNil())

	_, err = tpl.executeInternalMsg(msg, true)
	Expect(err).To(Equal(ErrUnknownLevel))
}

func (s *GomolSuite) TestTplFuncsCase(t sweet.T) {
	msg := newMessage(time.Unix(10, 0), nil, LevelError, nil, "UPPER")
	tpl, err := NewTemplate("{{ucase .LevelName}} {{lcase .Message}} {{title .LevelName}}")
	Expect(err).To(BeNil())

	out, err := tpl.executeInternalMsg(msg, false)
	Expect(err).To(BeNil())

	Expect(out).To(Equal("ERROR upper Error"))
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

func (s *GomolSuite) TestTplFuncMap(t sweet.T) {
	msg := newMessage(time.Unix(10, 0), nil, LevelDebug, nil, "forward message")
	tpl, err := NewTemplateWithFuncMap("{{reverse .Message}} {{color}}colored{{reset}}", template.FuncMap{
		"color":   func() string { return colorFatal },
		"reverse": reverseString,
	})
	Expect(err).To(BeNil())

	out, err := tpl.executeInternalMsg(msg, false)
	Expect(err).To(BeNil())

	Expect(out).To(Equal("egassem drawrof \x1b[1;31mcolored\x1b[0m"))
}

func (s *GomolSuite) TestTplMsgFromInternal(t sweet.T) {
	ts := time.Unix(10, 0)

	b := NewBase()
	b.SetAttr("baseAttr", 1234)
	b.SetAttr("overrideAttr", 1234)
	msg := newMessage(ts, b, LevelInfo,
		NewAttrs().
			SetAttr("msgAttr", 4321).
			SetAttr("overrideAttr", "test"),
		"Format %v %v", 1234, "asdf")

	tplMsg, err := newTemplateMsgFromMessage(msg)
	Expect(err).To(BeNil())
	Expect(tplMsg.Timestamp).To(Equal(ts))
	Expect(tplMsg.Level).To(Equal(LevelInfo))
	Expect(tplMsg.LevelName).To(Equal("info"))
	Expect(tplMsg.Message).To(Equal("Format 1234 asdf"))
	Expect(tplMsg.Attrs).To(HaveLen(3))
	Expect(tplMsg.Attrs["baseAttr"]).To(Equal(1234))
	Expect(tplMsg.Attrs["overrideAttr"]).To(Equal("test"))
	Expect(tplMsg.Attrs["msgAttr"]).To(Equal(4321))
}

func (s *GomolSuite) TestTplMsgAttrs(t sweet.T) {
	ts := time.Unix(10, 0)

	b := NewBase()
	b.SetAttr("baseAttr", 1234)
	b.SetAttr("overrideAttr", 1234)
	msg := newMessage(ts, b, LevelInfo,
		NewAttrs().
			SetAttr("msgAttr", 4321).
			SetAttr("overrideAttr", "test"),
		"Format %v %v", 1234, "asdf")

	tplMsg, err := newTemplateMsgFromMessage(msg)
	Expect(err).To(BeNil())
	Expect(tplMsg.Timestamp).To(Equal(ts))
	Expect(tplMsg.Level).To(Equal(LevelInfo))
	Expect(tplMsg.LevelName).To(Equal("info"))
	Expect(tplMsg.Message).To(Equal("Format 1234 asdf"))
	Expect(tplMsg.Attrs).To(HaveLen(3))
	Expect(tplMsg.Attrs["baseAttr"]).To(Equal(1234))
	Expect(tplMsg.Attrs["overrideAttr"]).To(Equal("test"))
	Expect(tplMsg.Attrs["msgAttr"]).To(Equal(4321))

	tpl, err := NewTemplate("{{range $key, $val := .Attrs}}{{$key}}=={{$val}}\n{{end}}")
	Expect(err).To(BeNil())

	out, err := tpl.executeInternalMsg(msg, false)
	Expect(err).To(BeNil())

	Expect(out).To(Equal("baseAttr==1234\nmsgAttr==4321\noverrideAttr==test\n"))
}

func (s *GomolSuite) TestTplTimestamp(t sweet.T) {
	ts := time.Unix(10, 0)

	msg := newMessage(ts, nil, LevelError, nil, "message")
	tpl, err := NewTemplate("{{ .Timestamp.Format \"2006-01-02T15:04:05.999999999Z07:00\" }}")
	Expect(err).To(BeNil())

	out, err := tpl.executeInternalMsg(msg, false)
	Expect(err).To(BeNil())

	Expect(out).To(Equal(ts.Format("2006-01-02T15:04:05.999999999Z07:00")))
}

func (s *GomolSuite) TestTplJSON(t sweet.T) {
	ts := time.Unix(1000000000, 100)

	msg := newMessage(ts, nil, LevelError,
		NewAttrs().
			SetAttr("attr1", "val1").
			SetAttr("attr2", 1234),
		"message")
	tpl, err := NewTemplate("{{ json . }}")
	Expect(err).To(BeNil())

	tplMsg, err := newTemplateMsgFromMessage(msg)
	Expect(err).To(BeNil())

	out, err := tpl.Execute(tplMsg, false)
	Expect(err).To(BeNil())

	// Unmarshal from json and check that because on Travis the timezone is different
	// and I don't want to create a new version of time.Time to marshal the value
	// differently
	dataOut := &TemplateMsg{}
	err = json.Unmarshal([]byte(out), dataOut)
	Expect(err).To(BeNil())

	Expect(dataOut.Timestamp.UnixNano()).To(Equal(msg.Timestamp.UnixNano()))
	Expect(dataOut.Level).To(Equal(tplMsg.Level))
	Expect(dataOut.LevelName).To(Equal(tplMsg.LevelName))
	Expect(dataOut.Message).To(Equal(tplMsg.Message))
	Expect(dataOut.Attrs).To(HaveLen(2))
	Expect(dataOut.Attrs["attr1"]).To(Equal("val1"))
	Expect(dataOut.Attrs["attr2"]).To(Equal(float64(1234)))

	tpl, err = NewTemplate("{{ json .Attrs }}")
	Expect(err).To(BeNil())

	out, err = tpl.executeInternalMsg(msg, false)
	Expect(err).To(BeNil())

	Expect(out).To(Equal("{\"attr1\":\"val1\",\"attr2\":1234}"))
}

func (s *GomolSuite) TestTplAttrTemplate(t sweet.T) {
	msg := newMessage(
		time.Unix(10, 0),
		nil,
		LevelFatal,
		NewAttrs().SetAttr("attrName", "attrVal"),
		"test",
	)
	tpl, err := NewTemplate("[{{.Attrs.attrName}}] {{.Message}}")
	Expect(err).To(BeNil())

	out, err := tpl.executeInternalMsg(msg, true)
	Expect(err).To(BeNil())

	Expect(out).To(Equal("[attrVal] test"))
}

func (s *GomolSuite) TestTplJSONError(t sweet.T) {
	data, err := tplJSON(map[string]interface{}{
		"attr1": s.TestTplJSONError,
	})
	Expect(data).To(Equal(""))
	Expect(err).ToNot(BeNil())
}

func (s *GomolSuite) TestNewTemplateMsgMinimal(t sweet.T) {
	ts := time.Unix(10, 0)

	tmsg := NewTemplateMsg(ts, LevelDebug, nil, "test")
	Expect(tmsg.Timestamp).To(Equal(ts))
	Expect(tmsg.Level).To(Equal(LevelDebug))
	Expect(tmsg.LevelName).To(Equal(LevelDebug.String()))
	Expect(tmsg.Attrs).To(Equal(map[string]interface{}{}))
	Expect(tmsg.Message).To(Equal("test"))
}
