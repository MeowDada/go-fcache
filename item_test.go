package fcache

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestItemConvertion(t *testing.T) {
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

func TestItemRemove(t *testing.T) {
	item := NewDummyItem("nothing")
	err := item.Remove()
	if err != nil {
		t.Error(err)
	}

	err = ioutil.WriteFile("test", nil, 0644)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove("test")

	realItem := NewItem("test", 0)
	err = realItem.Remove()
	if err != nil {
		t.Error(err)
	}
}

func TestItemIsZero(t *testing.T) {
	item := Item{}
	if !item.IsZero() {
		t.Errorf("expect IsZero() return true, but get false: %v\n", item)
	}

	real := NewDummyItem("haha")
	if real.IsZero() {
		t.Errorf("expect IsZero() return false, but get true: %v\n", real)
	}
}
