package cache

// Pool is a cache pool which stores information of caches.
type Pool interface {

	// Iter accepts a iteration callback for iterating key-value pairs, which are
	// stroed in the pool. The order could be nondeterministic, and should panics
	// when the callback function is a nil pointer. Also, an error should be raised once
	// error occurs in iteration callback.
	Iter(iterCb func(k string, v Item) error) error

	// Put puts a file cache into the pool with given key (usually same as path to the file)
	// and its size.
	Put(key string, size int64) error

	// Get gets a file cache from the pool with given key. If the key is missed, returns a
	// special error as ErrNoSuchKey.
	Get(key string) (Item, error)

	// Remove removes a key from the pool. If the key does not exist, then nothing should be done.
	Remove(key string) error

	// IncrRef increment the reference count of the file caches with given keys by one. If the
	// key does not exist, create a psudo one and increment it, too.
	IncrRef(keys ...string) error

	// DecrRef decrement the reference count fo the file caches with given keys by one. if the
	// reference count of the file caches is equal or less than 0, then do nothing. Also note
	// that if the key does not exist, do nothing, too.
	DecrRef(keys ...string) error

	// Close closes the file cache pool. The implementation depends, Usually just recycle
	// the allocated resources of the cache pool.
	Close() error
}
