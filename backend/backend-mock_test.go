package backend

import (
	"testing"
)

func TestMockPut(t *testing.T) {
	m := mock{put: func(k, v []byte) error { return nil }}
	err := m.Put(nil, nil)
	if err != nil {
		t.Error(err)
	}
}

func TestMockGet(t *testing.T) {
	m := mock{get: func(k []byte) ([]byte, error) { return nil, nil }}
	_, err := m.Get(nil)
	if err != nil {
		t.Error(err)
	}
}

func TestMockIter(t *testing.T) {
	m := mock{iter: func(func(k, v []byte) error) error { return nil }}
	err := m.Iter(nil)
	if err != nil {
		t.Error(err)
	}
}

func TestMockRemove(t *testing.T) {
	m := mock{rm: func(k []byte) error { return nil }}
	err := m.Remove(nil)
	if err != nil {
		t.Error(err)
	}
}

func TestMockClose(t *testing.T) {
	m := mock{close: func() error { return nil }}
	err := m.Close()
	if err != nil {
		t.Error(err)
	}
}
