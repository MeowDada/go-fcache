package cache

import (
	"os"
	"reflect"
	"time"
)

// New creates a file with given key and size.
func New(key string, size int64) Item {
	return Item{
		Key:       key,
		Path:      key,
		Size:      size,
		Ref:       0,
		Used:      0,
		Real:      true,
		CreatedAt: time.Now(),
	}
}

// Dummy creates a file with given key.
func Dummy(key string) Item {
	return Item{
		Key:  key,
		Path: key,
		Size: 0,
		Ref:  0,
		Used: 0,
	}
}

// Item implements Item interface. It represents
// a file cache item.
type Item struct {
	Key       string
	Path      string
	Size      int64
	Ref       int
	Used      int
	Real      bool
	CreatedAt time.Time
	LastUsed  time.Time
}

// SetSize sets the field of cache size.
func (f *Item) SetSize(size int64) {
	f.Size = size
}

// IncrRef increments the reference count of the cache item.
func (f *Item) IncrRef() {
	f.Ref++
}

// DecrRef decrements the reference count of the cache item.
func (f *Item) DecrRef() {
	f.Ref--
}

// IncrUsed increments the used count of the cache item.
func (f *Item) IncrUsed() {
	f.Used++
}

// Reference returns the reference count of the cache item.
func (f *Item) Reference() int {
	return f.Ref
}

// UsedCount returns the used count of the cache item.
func (f *Item) UsedCount() int {
	return f.Used
}

// CTime returns the creation date of the cache item.
func (f *Item) CTime() time.Time {
	return f.CreatedAt
}

// ATime returns the lastest used timestamp.
func (f *Item) ATime() time.Time {
	return f.LastUsed
}

// UpdateCreatedAt updates the created timestamp.
func (f *Item) UpdateCreatedAt() {
	f.CreatedAt = time.Now()
}

// UpdateLastUsed updates the last used timestamp.
func (f *Item) UpdateLastUsed() {
	f.LastUsed = time.Now()
}

// Remove removes the cache item from disk.
func (f *Item) Remove() error {
	if f.Real {
		return os.Remove(f.Path)
	}
	return nil
}

// IsReal returns if the cache item is a dummy one or not.
func (f *Item) IsReal() bool {
	return f.Real
}

// SetReal makes the cache item become a concret one.
func (f *Item) SetReal() {
	f.Real = true
}

// IsZero returns the item is a zero value item or not.
func (f *Item) IsZero() bool {
	el := reflect.ValueOf(f).Elem()
	for i := 0; i < el.NumField(); i++ {
		if !el.Field(i).IsZero() {
			return false
		}
	}
	return true
}
