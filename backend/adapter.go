package backend

import (
	"github.com/meowdada/go-fcache/cache"
	"github.com/meowdada/go-fcache/codec"
	"github.com/meowdada/go-fcache/pkg/ioutil"
)

// Adapter creates an adapter for cache.DB.
func Adapter(store Store, codec codec.Codec) cache.Pool {
	return &adapter{
		backend: store,
		codec:   codec,
		idgen:   newSnowflake(0),
	}
}

type adapter struct {
	backend Store
	codec   codec.Codec
	idgen   IDGenerator
}

func (ada *adapter) Iter(iterCb func(k string, v cache.Item) error) error {
	var (
		b = ada.backend
	)
	return b.Iter(func(k, v []byte) error {
		item := ada.parse(v)
		return iterCb(ioutil.Bytes2Str(k), item)
	})
}

func (ada *adapter) Put(key string, size int64) error {
	var (
		b = ada.backend
		k = ioutil.Str2Bytes(key)
	)

	// If the key does not present. Create a new item and
	// insert it into the backend.
	v, err := b.Get(k)
	if IsNoKeyError(err) {
		return b.Put(k, ada.newMarshaledItem(key, size))
	}

	// Deserialize the stored data.
	item := ada.parse(v)

	// If it is a psudo item, then convert it into
	// a real one.
	if !item.IsReal() {
		item.SetReal()
		item.SetSize(size)
		item.UpdateCreatedAt()
		v = ada.marshalItem(item)
		return b.Put(k, v)
	}

	return ErrDupKey
}

func (ada *adapter) Get(key string) (cache.Item, error) {
	var (
		b = ada.backend
		k = ioutil.Str2Bytes(key)
	)

	v, err := b.Get(k)
	if err != nil {
		return cache.Item{}, err
	}

	item := ada.parse(v)
	return item, nil
}

func (ada *adapter) Remove(key string) error {
	var (
		b = ada.backend
		k = ioutil.Str2Bytes(key)
	)
	return b.Remove(k)
}

func (ada *adapter) IncrRef(keys ...string) error {
	var (
		b = ada.backend
	)

	for _, key := range keys {
		k := ioutil.Str2Bytes(key)
		v, err := b.Get(k)

		// If the key presents, then update this cache item.
		if err == nil {
			item := ada.parse(v)
			item.IncrRef()
			item.IncrUsed()
			item.UpdateLastUsed()

			if err := b.Put(k, ada.marshalItem(item)); err != nil {
				return err
			}
			continue
		}

		// If the key does not present, then create a dummy one.
		id := ada.idgen.Get()
		item := cache.Dummy(id, key)
		item.IncrRef()
		item.IncrUsed()
		item.UpdateLastUsed()

		err = b.Put(k, ada.marshalItem(item))
		if err != nil {
			return err
		}
	}
	return nil
}

func (ada *adapter) DecrRef(keys ...string) error {
	var (
		b = ada.backend
	)

	for _, key := range keys {
		k := ioutil.Str2Bytes(key)
		v, err := b.Get(k)

		// If the key does not present, then ignore it.
		if err == cache.ErrNoSuchKey {
			continue
		}
		if err != nil && err != cache.ErrNoSuchKey {
			return err
		}

		item := ada.parse(v)
		if item.Reference() > 0 {
			item.DecrRef()
		}

		err = b.Put(k, ada.marshalItem(item))
		if err != nil {
			return err
		}
	}
	return nil
}

func (ada *adapter) Close() error {
	return ada.backend.Close()
}

func (ada *adapter) parse(data []byte) cache.Item {
	var item cache.Item
	ada.codec.Unmarshal(data, &item)
	return item
}

func (ada *adapter) marshalItem(item cache.Item) []byte {
	data, _ := ada.codec.Marshal(item)
	return data
}

func (ada *adapter) newMarshaledItem(key string, size int64) []byte {
	id := ada.idgen.Get()
	item := cache.New(id, key, size)
	data, _ := ada.codec.Marshal(item)
	return data
}
