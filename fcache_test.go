package fcache

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestManager(t *testing.T) {
	db, err := BoltDB("bolt.db")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		db.Close()
		os.Remove("bolt.db")
	}()

	type pair struct {
		path string
		size int64
	}

	pairs := []pair{
		{"a1", 40},
		{"a2", 20},
		{"a3", 30},
		{"a4", 10},
		{"a5", 0},
		{"b1", 100},
		{"b2", 50},
		{"c1", 50},
		{"c2", 50},
	}
	defer func(pairs []pair) {
		for _, pair := range pairs {
			os.Remove(pair.path)
		}
	}(pairs)

	for _, pair := range pairs {
		if err := ioutil.WriteFile(pair.path, nil, 0644); err != nil {
			t.Fatal(err)
		}
	}

	mgr := New(200, db, RR())

	for _, pair := range pairs {
		if strings.Contains(pair.path, "a") {
			mgr.Register(pair.path)
		}
	}

	for _, pair := range pairs {
		err = mgr.Set(pair.path, pair.size)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestManagerCap(t *testing.T) {
	cap := int64(1024)
	mgr := New(cap, Hashmap(), LRU())
	if mgr.Cap() != cap {
		t.Errorf("expect %d, but get %d", cap, mgr.Cap())
	}
}

func TestManagerPut(t *testing.T) {
	cap := int64(1000)
	mgr := New(cap, Hashmap(), LRU())
	err := mgr.Set("123", cap+1)
	if err != ErrCacheTooLarge {
		t.Errorf("expect %v, but get %v", ErrCacheTooLarge, err)
	}

	err = mgr.Set("456", 100)
	if err != nil {
		t.Fatal(err)
	}

	_, err = mgr.Get("456")
	if err != nil {
		t.Fatal(err)
	}
}

func TestManagerRegister(t *testing.T) {
	// Mock data.
	pairs := []struct {
		path string
		size int64
	}{
		{"tmp/a", 500},
		{"tmp/e", 600},
		{"tmp/i", 700},
		{"tmp/o", 800},
		{"tmp/u", 900},
	}

	err := os.Mkdir("tmp", 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("tmp")

	for _, p := range pairs {
		err := ioutil.WriteFile(p.path, nil, 0644)
		if err != nil {
			t.Fatal(err)
		}
	}
	defer func() {
		for _, p := range pairs {
			os.Remove(p.path)
		}
	}()

	mgr := New(1000, Hashmap(), LRU())
	for _, p := range pairs {
		mgr.Register(p.path)
		mgr.Unregister(p.path)
		err := mgr.Set(p.path, p.size)
		if err != nil {
			t.Fatal(err)
		}
	}
}
