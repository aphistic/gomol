package gomol

import (
	"testing"

	. "github.com/onsi/gomega"
)

func (s *GomolSuite) TestNewAttrsFromMap(t *testing.T) {
	attrs := NewAttrsFromMap(map[string]interface{}{
		"attr1": "val1",
		"attr2": 1234,
	})
	Expect(attrs.attrs).To(HaveLen(2))
	Expect(attrs.attrs[getAttrHash("attr1")]).To(Equal("val1"))
	Expect(attrs.attrs[getAttrHash("attr2")]).To(Equal(1234))
}

func (s *GomolSuite) TestAttrsMergeNilAttrs(t *testing.T) {
	attrs := NewAttrs()
	attrs.MergeAttrs(nil)
}

func (s *GomolSuite) TestNewAttrsFromAttrsNone(t *testing.T) {
	testAttrs := NewAttrsFromAttrs()
	Expect(testAttrs).ToNot(BeNil())
	Expect(testAttrs.Attrs()).To(HaveLen(0))
}

func (s *GomolSuite) TestNewAttrsFromAttrsSingle(t *testing.T) {
	attrs1 := NewAttrsFromMap(map[string]interface{}{
		"attr1": "val1",
		"attr3": 1234,
	})

	testAttrs := NewAttrsFromAttrs(attrs1)
	Expect(testAttrs).ToNot(BeNil())
	Expect(testAttrs.GetAttr("attr1")).To(Equal("val1"))
	Expect(testAttrs.GetAttr("attr3")).To(Equal(1234))
}

func (s *GomolSuite) TestNewAttrsFromAttrsMultiple(t *testing.T) {
	attrs1 := NewAttrsFromMap(map[string]interface{}{
		"attr1": "val1",
		"attr3": 1234,
	})
	attrs2 := NewAttrsFromMap(map[string]interface{}{
		"attr2": true,
		"attr4": 1234.4321,
	})

	testAttrs := NewAttrsFromAttrs(attrs1, attrs2)
	Expect(testAttrs).ToNot(BeNil())
	Expect(testAttrs.GetAttr("attr1")).To(Equal("val1"))
	Expect(testAttrs.GetAttr("attr2")).To(Equal(true))
	Expect(testAttrs.GetAttr("attr3")).To(Equal(1234))
	Expect(testAttrs.GetAttr("attr4")).To(Equal(1234.4321))
}

func (s *GomolSuite) TestNewAttrsFromAttrsNil(t *testing.T) {
	testAttrs := NewAttrsFromAttrs(nil, nil, nil)
	Expect(testAttrs).ToNot(BeNil())
	Expect(testAttrs.Attrs()).To(HaveLen(0))
}

func (s *GomolSuite) TestAttrsGetMissing(t *testing.T) {
	attrs := NewAttrs()
	Expect(attrs.GetAttr("not an attr")).To(BeNil())
}

func (s *GomolSuite) TestAttrsRemoveMissing(t *testing.T) {
	attrs := NewAttrs()
	// Just run it to make sure it doesn't panic
	attrs.RemoveAttr("not an attr")
}

func (s *GomolSuite) TestAttrsChaining(t *testing.T) {
	attrs := NewAttrs().
		SetAttr("attr1", "val1").
		SetAttr("attr2", "val2").
		SetAttr("attr3", 3).
		SetAttr("attr4", 4)

	Expect(attrs.attrs[getAttrHash("attr1")]).To(Equal("val1"))
	Expect(attrs.attrs[getAttrHash("attr2")]).To(Equal("val2"))
	Expect(attrs.attrs[getAttrHash("attr3")]).To(Equal(3))
	Expect(attrs.attrs[getAttrHash("attr4")]).To(Equal(4))
}

func (s *GomolSuite) TestGetHashAttrMissing(t *testing.T) {
	res, err := getHashAttr(1234)

	Expect(res).To(Equal(""))
	Expect(err.Error()).To(Equal("Could not find attr for hash 1234"))
}
