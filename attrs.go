package gomol

import (
	"github.com/spaolacci/murmur3"

	"fmt"
	"sync"
)

type Attrs struct {
	attrs map[uint32]interface{}
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

func (a *Attrs) mergeAttrs(attrs *Attrs) {
	if attrs == nil {
		return
	}
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

func (a *Attrs) SetAttr(key string, value interface{}) *Attrs {
	hash := getAttrHash(key)
	a.attrs[hash] = value
	return a
}

func (a *Attrs) GetAttr(key string) interface{} {
	return a.attrs[getAttrHash(key)]
}

func (a *Attrs) RemoveAttr(key string) {
	delete(a.attrs, getAttrHash(key))
}

func (a *Attrs) Attrs() map[string]interface{} {
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
	if attr, ok := hashAttrs[hash]; ok {
		return attr, nil
	}

	return "", fmt.Errorf("Could not find attr for hash %v", hash)
}
