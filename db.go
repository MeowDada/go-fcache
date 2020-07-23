package fcache

import (
	"github.com/pkg/errors"
)

// ErrCacheMiss raises when try getting an unexist record.
var ErrCacheMiss = errors.New("cache miss")

// DB is a database which ables to record the information of file caches.
type DB interface {
	Iter(iterCb func(k string, v Item) error) error
	Put(key string, size int64) error
	Get(key string) (Item, error)
	Remove(key string) error
	IncrRef(keys ...string) error
	DecrRef(keys ...string) error
	Close() error
}
