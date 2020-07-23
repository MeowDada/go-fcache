package fcache

import (
	"sync"
	"time"
)

// Map implements DB interface.
type hashmap struct {
	m  map[string]Item
	mu sync.RWMutex
}

// Hashmap returns a DB instance implemented by a hashmap.
func Hashmap() DB {
	return &hashmap{
		m: make(map[string]Item),
	}
}

func (hm *hashmap) Iter(iterCb func(k string, v Item) error) (err error) {
	hm.rlockFn(func() {
		err = hm.iter(iterCb)
	})
	return err
}

func (hm *hashmap) Put(key string, size int64) (err error) {
	hm.lockFn(func() {
		err = hm.put(key, size)
	})
	return err
}

func (hm *hashmap) Get(key string) (item Item, err error) {
	hm.rlockFn(func() {
		item, err = hm.get(key)
	})
	return item, err
}

func (hm *hashmap) Remove(key string) error {
	hm.lockFn(func() {
		delete(hm.m, key)
	})
	return nil
}

func (hm *hashmap) IncrRef(keys ...string) (err error) {
	hm.lockFn(func() {
		err = hm.incrRef(keys...)
	})
	return err
}

func (hm *hashmap) DecrRef(keys ...string) (err error) {
	hm.lockFn(func() {
		err = hm.decrRef(keys...)
	})
	return err
}

func (hm *hashmap) Close() error {
	return nil
}

func (hm *hashmap) incrRef(keys ...string) error {
	for _, key := range keys {
		item, ok := hm.m[key]
		if !ok {
			item = NewDummyItem(key)
		}
		item.Ref++
		item.Used++
		item.LastUsed = time.Now()
		hm.m[key] = item
	}
	return nil
}

func (hm *hashmap) decrRef(keys ...string) error {
	for _, key := range keys {
		item, ok := hm.m[key]
		if !ok {
			continue
		}
		if item.Ref > 0 {
			item.Ref--
			hm.m[key] = item
		}
	}
	return nil
}

func (hm *hashmap) put(key string, size int64) error {
	item, ok := hm.m[key]
	if !ok {
		item = NewItem(key, size)
		hm.m[key] = item
		return nil
	}

	if !item.Real {
		item.Real = true
		item.Size = size
		item.CreatedAt = time.Now()
		hm.m[key] = item
		return nil
	}
	return ErrDupKey
}

func (hm *hashmap) get(key string) (Item, error) {
	item, ok := hm.m[key]
	if ok {
		return item, nil
	}
	return Item{}, ErrCacheMiss
}

func (hm *hashmap) iter(iterCb func(k string, v Item) error) error {
	for k, v := range hm.m {
		err := iterCb(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (hm *hashmap) lockFn(fn func()) {
	hm.mu.Lock()
	fn()
	hm.mu.Unlock()
}

func (hm *hashmap) rlockFn(fn func()) {
	hm.mu.RLock()
	fn()
	hm.mu.RUnlock()
}
