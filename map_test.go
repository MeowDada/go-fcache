package fcache

import (
	"fmt"
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

func TestHashmapIncrRef(t *testing.T) {
	m := Hashmap()
	dummies := []string{
		"d1",
		"d2",
		"d3",
		"d4",
		"d5",
	}
	err := m.IncrRef(dummies...)
	if err != nil {
		t.Fatal(err)
	}

	err = m.Iter(func(k string, v Item) error {
		if v.Ref != 1 {
			return fmt.Errorf("expect %s's item with ref = 1, but get %d", k, v.Ref)
		}
		return nil
	})
	if err != nil {
		t.Error(err)
	}
}

func TestHashmapDecrRef(t *testing.T) {
	m := Hashmap()

	// Case1: Decrement reference count of unexist item.
	// Should return no error.
	err := m.DecrRef("unexist")
	if err != nil {
		t.Error(err)
	}

	// Case2: Decrement reference count of an existing item with no reference.
	// Should return no error
	m.Put("abc", 100)
	err = m.DecrRef("abc")
	if err != nil {
		t.Error(err)
	}

	// Case3: Decrement reference count of an existing item with reference count > 1.
	m.IncrRef("san")
	m.IncrRef("san")
	m.IncrRef("san")
	err = m.DecrRef("san")
	if err != nil {
		t.Error(err)
	}

	item, err := m.Get("san")
	if err != nil {
		t.Fatal(err)
	}

	if item.Ref != 2 {
		t.Errorf("expect reference count = 2, but get %d\n", item.Ref)
	}

}
