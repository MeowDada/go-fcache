package backend

import (
	"github.com/pkg/errors"
)

// Store is a storage backend interface which provides some
// methods to interacts with the item.
type Store interface {
	Put(k, v []byte) error
	Get(k []byte) (v []byte, e error)
	Remove(k []byte) error
	Iter(iterCb func(k, v []byte) error) error
	Close() error
}

// ErrNoSuchKey raises when try accessing a Store with key nonexist.
var ErrNoSuchKey = errors.New("no such key")

// IsNoKeyError returns true if the key reprsents ErrNoSuchKey.
func IsNoKeyError(err error) bool {
	return err == ErrNoSuchKey
}
