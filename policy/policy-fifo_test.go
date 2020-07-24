package policy

import (
	"testing"

	"github.com/meowdada/go-fcache/backend"
	"github.com/meowdada/go-fcache/backend/gomap"
	"github.com/meowdada/go-fcache/codec"
)

func TestCacheReplacementAlgoFIFO(t *testing.T) {
	db := backend.Adapter(gomap.New(), codec.Gob{})

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
	}

	fifo := FIFO()
	item, err := fifo.Emit(db)
	if err != nil {
		t.Fatal(err)
	}

	if item.Size != pairs[0].size {
		t.Errorf("expect %v, but get %v\n", pairs[0], item)
	}
}

func TestCacheReplacementAlgoFIFOError(t *testing.T) {
	db := backend.Adapter(gomap.New(), codec.Gob{})
	db.IncrRef("123")
	fifo := FIFO()
	_, err := fifo.Emit(db)
	if err == nil {
		t.Errorf("expect err = %v, but get nil", ErrNoEmitableCaches)
	}
}
