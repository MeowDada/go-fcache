package boltdb

import (
	"github.com/meowdada/go-fcache/cache"
	bolt "go.etcd.io/bbolt"
)

// New is a factory method to creates a boltdb instance. Note that it will
// not open a boltDB immediately until used.
func New(opts Options) *BoltDB {
	return &BoltDB{
		core:   nil,
		opts:   opts,
		bucket: []byte(opts.Bucket),
	}
}

// BoltDB implements backend.Store interface.
type BoltDB struct {
	core   *bolt.DB
	opts   Options
	bucket []byte
}

// Put puts a key-value pair into the boltDB.
func (b *BoltDB) Put(k, v []byte) error {
	if err := b.init(); err != nil {
		return err
	}
	return b.core.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(b.bucket)
		return bucket.Put(k, v)
	})
}

// Get gets a value by given key from boltDB.
func (b *BoltDB) Get(k []byte) (v []byte, e error) {
	if err := b.init(); err != nil {
		return nil, err
	}
	e = b.core.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(b.bucket)
		v = bucket.Get(k)
		return nil
	})
	if v == nil {
		return nil, cache.ErrNoSuchKey
	}
	return v, e
}

// Remove removes a key-value pair from the boltDB.
func (b *BoltDB) Remove(k []byte) error {
	if err := b.init(); err != nil {
		return err
	}
	return b.core.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(b.bucket)
		return bucket.Delete(k)
	})
}

// Iter iterates all key-value pairs from the boltDB.
func (b *BoltDB) Iter(iterCb func(k, v []byte) error) error {
	if err := b.init(); err != nil {
		return err
	}
	return b.core.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(b.bucket)
		return bucket.ForEach(iterCb)
	})
}

// Close cloes the boltDB.
func (b *BoltDB) Close() error {
	if b.core == nil {
		return nil
	}
	return b.core.Close()
}

func (b *BoltDB) init() (err error) {
	if b.core == nil {
		opts := b.opts

		// Open the boltDB.
		b.core, err = bolt.Open(opts.Path, opts.Mode, opts.Options)
		if err != nil {
			return err
		}

		// Closes the database connection if any error occurs.
		defer func() {
			if err != nil && b.core != nil {
				b.core.Close()
			}
		}()

		// Create bucket if not exist.
		err = b.core.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists(b.bucket)
			return err
		})
	}
	return err
}
