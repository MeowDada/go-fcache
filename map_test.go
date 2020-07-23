package fcache

import (
	"testing"
)

func TestHashmapPut(t *testing.T) {
	pairs := []struct {
		path string
		size int64
	}{
		{"/path/to/file", 100},
		{"hello-world", 125},
		{"fcache", 125},
		{"lala3", 200},
		{"test-put", 300},
		{"wonder4", 10},
		{"fortune", 55},
		{"solider", 76},
		{"yakusoku", 999},
	}

	m := Hashmap()
	for _, p := range pairs {
		err := m.Put(p.path, p.size)
		if err != nil {
			t.Errorf("exepct no error, but get %v when putting %v", err, p)
		}
	}
}

func TestHashmapGet(t *testing.T) {
	pairs := []struct {
		path string
		size int64
	}{
		{"/path/to/file", 100},
		{"hello-world", 125},
		{"fcache", 125},
		{"lala3", 200},
		{"test-put", 300},
		{"wonder4", 10},
		{"fortune", 55},
		{"solider", 76},
		{"yakusoku", 999},
	}

	m := Hashmap()
	for _, p := range pairs {
		err := m.Put(p.path, p.size)
		if err != nil {
			t.Errorf("exepct no error, but get %v when putting %v", err, p)
		}
	}

	// Get existing items.
	for _, p := range pairs {
		_, err := m.Get(p.path)
		if err != nil {
			t.Errorf("expect no error, but get %v when get %s", err, p.path)
		}
	}

	// Get unexist item.
	_, err := m.Get("abc")
	if err != ErrCacheMiss {
		t.Errorf("expect %v, but get no error", ErrCacheMiss)
	}
}

func TestHashmapClose(t *testing.T) {
	m := Hashmap()
	err := m.Close()
	if err != nil {
		t.Error(err)
	}
}
