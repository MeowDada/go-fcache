package gomap

import (
	"sync"

	"github.com/meowdada/go-fcache/backend"
	"github.com/meowdada/go-fcache/pkg/ioutil"
)

// New is a factory method to create a instance of Map.
func New() *Map {
	return &Map{ma: make(map[string][]byte)}
}

// Map implements backend.Store interface.
type Map struct {
	ma map[string][]byte
	mu sync.RWMutex
}

// Put is a concurrent safe method which puts a byte array k as key
// and byte array v as value into this map. If the key duplicates,
// the new one will replace the old one and returns with no error.
func (m *Map) Put(k, v []byte) error {
	m.mu.Lock()
	m.ma[ioutil.Bytes2Str(k)] = v
	m.mu.Unlock()
	return nil
}

// Get is a concurrent safe method which gets a value of byte array v
// from this map. If the key does not present, it will return with a
// nil value v and an error e as backend.ErrNoSuchKey.
func (m *Map) Get(k []byte) (v []byte, e error) {
	m.mu.RLock()
	v, ok := m.ma[ioutil.Bytes2Str(k)]
	m.mu.RUnlock()
	if ok {
		return v, nil
	}
	return nil, backend.ErrNoSuchKey
}

// Remove is a concurrent safe method which removes an entry
// from the map. It will return no error even if the key does
// not present in the map.
func (m *Map) Remove(k []byte) error {
	m.mu.Lock()
	delete(m.ma, ioutil.Bytes2Str(k))
	m.mu.Unlock()
	return nil
}

// Iter is a concurrent safe method which iterates key-value pairs
// stored in the map. If any error occurs when iterating, the for
// loop will halt and returns the underlying error. Note that do
// not modify the content stored in the map during iterating or
// the behavior will be undefined.
func (m *Map) Iter(iterCb func(k, v []byte) error) (err error) {
	m.rlockFn(func() {
		for k, v := range m.ma {
			err = iterCb(ioutil.Str2Bytes(k), v)
			if err != nil {
				return
			}
		}
	})
	return err
}

// Close closes the map, actually it does nothing.
func (m *Map) Close() error {
	return nil
}

func (m *Map) rlockFn(fn func()) {
	m.mu.RLock()
	fn()
	m.mu.RUnlock()
}
