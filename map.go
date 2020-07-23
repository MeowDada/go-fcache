package fcache

import "time"

// Map implements DB interface.
type hashmap struct {
	m map[string]Item
}

// Hashmap returns a DB instance implemented by a hashmap.
func Hashmap() DB {
	return &hashmap{
		m: make(map[string]Item),
	}
}

func (hm *hashmap) Iter(iterCb func(k string, v Item) error) error {
	for k, v := range hm.m {
		err := iterCb(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (hm *hashmap) Put(key string, size int64) error {
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

func (hm *hashmap) Get(key string) (Item, error) {
	item, ok := hm.m[key]
	if ok {
		return item, nil
	}
	return Item{}, ErrCacheMiss
}

func (hm *hashmap) Remove(key string) error {
	delete(hm.m, key)
	return nil
}

func (hm *hashmap) IncrRef(keys ...string) error {
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

func (hm *hashmap) DecrRef(keys ...string) error {
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

func (hm *hashmap) Close() error {
	return nil
}
