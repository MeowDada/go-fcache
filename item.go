package fcache

import (
	"bytes"
	"encoding/gob"
	"os"
	"reflect"
	"time"
)

// Item is a data structure representing a cache item.
type Item struct {
	Path      string
	Size      int64
	Ref       int
	Used      int
	CreatedAt time.Time
	LastUsed  time.Time
	Real      bool
}

// NewItem creates a file cache item.
func NewItem(path string, size int64) Item {
	return Item{
		Path:      path,
		Size:      size,
		Ref:       0,
		Used:      0,
		CreatedAt: time.Now(),
		Real:      true,
	}
}

// NewDummyItem creates a dummy file cache item.
func NewDummyItem(path string) Item {
	return Item{
		Path: path,
		Real: false,
	}
}

// Key returns the key of the item.
func (item *Item) Key() []byte {
	return []byte(item.Path)
}

// Bytes marshals the data structure into binaries.
func (item *Item) Bytes() ([]byte, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(item)
	return b.Bytes(), err
}

// Parse parses data fields from the input binaries.
func (item *Item) Parse(data []byte) error {
	b := bytes.NewBuffer(data)
	dec := gob.NewDecoder(b)
	return dec.Decode(item)
}

// IsZero returns this item is zero valued or not.
func (item *Item) IsZero() bool {
	el := reflect.ValueOf(item).Elem()
	for i := 0; i < el.NumField(); i++ {
		if !el.Field(i).IsZero() {
			return false
		}
	}
	return true
}

// Remove removes the file cache from cache space physically.
func (item *Item) Remove() error {
	return os.Remove(item.Path)
}
