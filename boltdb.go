package fcache

import (
	"time"

	bolt "go.etcd.io/bbolt"
)

type boltDB struct {
	bucket []byte
	*bolt.DB
}

// BoltDB returns an instance of boltDB which implements DB interface.
func BoltDB(path string) (DB, error) {
	db, err := bolt.Open(path, 0666, &bolt.Options{Timeout: time.Second})
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("cache"))
		return err
	})
	if err != nil {
		return nil, db.Close()
	}

	return &boltDB{
		DB:     db,
		bucket: []byte("cache"),
	}, err
}

func (db *boltDB) Close() error {
	return db.DB.Close()
}

func (db *boltDB) Put(key string, size int64) error {
	return db.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(db.bucket)

		// Try fetching record first. If the records presents and its not real one,
		// update its value and save it. If the records does not present, just put it
		// directly.
		data := bucket.Get([]byte(key))
		if data == nil {
			item := NewItem(key, size)
			b, err := item.Bytes()
			if err != nil {
				return err
			}
			return bucket.Put(item.Key(), b)
		}

		var item Item
		err := item.Parse(data)
		if err != nil {
			return err
		}

		if !item.Real {
			item.Real = true
			item.Size = size
			item.CreatedAt = time.Now()
			b, err := item.Bytes()
			if err != nil {
				return err
			}
			return bucket.Put(item.Key(), b)
		}

		return ErrDupKey
	})
}

func (db *boltDB) Get(key string) (item Item, err error) {
	err = db.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(db.bucket)
		data := bucket.Get([]byte(key))
		if data == nil {
			return ErrCacheMiss
		}
		return item.Parse(data)
	})
	return item, err
}

func (db *boltDB) Remove(key string) error {
	return db.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(db.bucket)
		return bucket.Delete([]byte(key))
	})
}

func (db *boltDB) IncrRef(keys ...string) error {
	return db.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(db.bucket)

		// Iterate the keys. Increment the reference count if the key exists.
		// Create dummy records if the key does not exist.
		for _, key := range keys {

			// Retrieve the data structure by given key.
			data := bucket.Get([]byte(key))

			// If the data entries does not present. Create a dummy one.
			if data == nil {
				item := NewDummyItem(key)
				item.Ref++
				item.Used++
				item.LastUsed = time.Now()
				b, err := item.Bytes()
				if err != nil {
					return err
				}
				err = bucket.Put([]byte(key), b)
				if err != nil {
					return err
				}
				continue
			}

			// If the data entry presents, increment its reference count and used count.
			var item Item
			err := item.Parse(data)
			if err != nil {
				return err
			}
			item.Ref++
			item.Used++
			item.LastUsed = time.Now()

			b, err := item.Bytes()
			if err != nil {
				return err
			}

			err = bucket.Put([]byte(key), b)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (db *boltDB) DecrRef(keys ...string) error {
	return db.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(db.bucket)

		// Iterate the keys, decrement reference count if and only if
		// the key pair presents and its reference count is larger or equal
		// than 1.
		for _, key := range keys {
			data := bucket.Get([]byte(key))
			if data == nil {
				continue
			}

			var item Item
			err := item.Parse(data)
			if err != nil {
				return err
			}

			if item.Ref > 0 {
				item.Ref--
			}

			b, err := item.Bytes()
			if err != nil {
				return err
			}

			err = bucket.Put(item.Key(), b)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (db *boltDB) Iter(iterCb func(k string, value Item) error) error {
	wrapper := func(k, v []byte) error {
		var item Item
		err := item.Parse(v)
		if err != nil {
			return err
		}
		return iterCb(item.Path, item)
	}

	return db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(db.bucket)
		return bucket.ForEach(wrapper)
	})
}

func (db *boltDB) iter(iterCb func(k, v []byte) error) error {
	return db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(db.bucket)
		return bucket.ForEach(iterCb)
	})
}
