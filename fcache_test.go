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
