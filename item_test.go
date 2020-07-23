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

	data := item.Bytes()

	var dup Item
	dup.Parse(data)

	if !cmp.Equal(dup, item) {
		t.Errorf("dup = %v, item = %v", dup, item)
	}
}
