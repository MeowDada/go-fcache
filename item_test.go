package fcache

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestItem(t *testing.T) {
	item := Item{
		Path:      "/path/to/file",
		Size:      100,
		Ref:       0,
		Used:      5,
		CreatedAt: time.Now(),
		LastUsed:  time.Now(),
	}

	data, err := item.Bytes()
	if err != nil {
		t.Fatal(err)
	}

	var dup Item
	err = dup.Parse(data)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(dup, item) {
		t.Errorf("dup = %v, item = %v", dup, item)
	}
}
