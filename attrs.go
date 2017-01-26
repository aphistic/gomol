package gomol

import (
	"github.com/spaolacci/murmur3"

	"fmt"
	"sync"
)

// Attrs represents a collection of key/value attributes
type Attrs struct {
	attrs     map[uint32]interface{}
	attrsLock sync.RWMutex
}

// NewAttrs will create a new Attrs struct with an empty set of attributes.
func NewAttrs() *Attrs {
	return &Attrs{
		attrs: make(map[uint32]interface{}),
	}
}

// NewAttrsFromMap will create a new Attrs struct with the given attributes pre-populated
func NewAttrsFromMap(attrs map[string]interface{}) *Attrs {
	newAttrs := NewAttrs()
	for attrKey, attrVal := range attrs {
		newAttrs.SetAttr(attrKey, attrVal)
	}
	return newAttrs
}

// NewAttrsFromAttrs is a convenience function that will accept zero or more existing Attrs, create
// a new Attrs and then merge all the supplied Attrs values into the new Attrs instance.
func NewAttrsFromAttrs(attrs ...*Attrs) *Attrs {
	newAttrs := NewAttrs()
	for _, attr := range attrs {
		newAttrs.MergeAttrs(attr)
	}
	return newAttrs
}

// MergeAttrs accepts another existing Attrs and merges the attributes into its own.
func (a *Attrs) MergeAttrs(attrs *Attrs) {
	if attrs == nil {
		return
	}
	a.attrsLock.Lock()
	defer a.attrsLock.Unlock()
	for hash, val := range attrs.attrs {
		a.attrs[hash] = val
	}
}

func (a *Attrs) clone() *Attrs {
	attrs := NewAttrs()
	for hash, val := range a.attrs {
		attrs.attrs[hash] = val
	}
	return attrs
}

// SetAttr will set key to the provided value.  If the attribute already exists the value will
// be replaced with the new value.
func (a *Attrs) SetAttr(key string, value interface{}) *Attrs {
	a.attrsLock.Lock()
	defer a.attrsLock.Unlock()

	hash := getAttrHash(key)
	a.attrs[hash] = value
	return a
}

// GetAttr gets the value of the attribute with the provided name.  If the attribute does not
// exist, nil will be returned
func (a *Attrs) GetAttr(key string) interface{} {
	a.attrsLock.RLock()
	defer a.attrsLock.RUnlock()

	return a.attrs[getAttrHash(key)]
}

// RemoveAttr will remove the attribute with the provided name.
func (a *Attrs) RemoveAttr(key string) {
	a.attrsLock.Lock()
	defer a.attrsLock.Unlock()

	delete(a.attrs, getAttrHash(key))
}

// Attrs will return a map of the attributes added to the struct.
func (a *Attrs) Attrs() map[string]interface{} {
	a.attrsLock.RLock()
	defer a.attrsLock.RUnlock()

	attrs := make(map[string]interface{})
	for hash, val := range a.attrs {
		key, _ := getHashAttr(hash)
		attrs[key] = val
	}
	return attrs
}

type logAttr struct {
	Name  string
	Value interface{}
}

var hashMutex sync.RWMutex
var attrHashes = make(map[string]uint32)
var hashAttrs = make(map[uint32]string)

func getAttrHash(attr string) uint32 {
	// First try to acquire a read lock to see if we even need to hash
	// the string at all
	hashMutex.RLock()
	if hash, ok := attrHashes[attr]; ok {
		hashMutex.RUnlock()
		return hash
	}

	// We do need to hash it so release the read lock and acquire a write lock
	hashMutex.RUnlock()

	hashMutex.Lock()
	defer hashMutex.Unlock()
	if hash, ok := attrHashes[attr]; ok {
		return hash
	}

	hasher := murmur3.New32()
	hasher.Write([]byte(attr))

	hash := hasher.Sum32()
	hashAttrs[hash] = attr
	attrHashes[attr] = hash

	return hash
}

func getHashAttr(hash uint32) (string, error) {
	hashMutex.RLock()
	defer hashMutex.RUnlock()

	if attr, ok := hashAttrs[hash]; ok {
		return attr, nil
	}

	return "", fmt.Errorf("Could not find attr for hash %v", hash)
}
