package fcache

import (
	"testing"
)

func TestCacheReplacementAlgoMRU(t *testing.T) {
	db := Hashmap()

	pairs := []struct {
		path string
		size int64
	}{
		{"a", 100},
		{"b", 200},
		{"c", 300},
		{"d", 400},
		{"e", 500},
		{"f", 600},
		{"g", 700},
		{"h", 800},
		{"i", 900},
		{"j", 1000},
	}

	for _, pair := range pairs {
		err := db.Put(pair.path, pair.size)
		if err != nil {
			t.Fatal(err)
		}
		err = db.IncrRef(pair.path)
		if err != nil {
			t.Fatal(err)
		}
		err = db.DecrRef(pair.path)
		if err != nil {
			t.Fatal(err)
		}
	}

	mru := MRU()
	item, err := mru.Emit(db)
	if err != nil {
		t.Fatal(err)
	}

	if item.Size != pairs[len(pairs)-1].size {
		t.Errorf("expect %v, but get %v\n", pairs[len(pairs)-1], item)
	}
}