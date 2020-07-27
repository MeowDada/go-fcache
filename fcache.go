// Package fcache provides utilities to manage file caches with limited
// local cache volume.
package fcache

import (
	"sync"

	retry "github.com/avast/retry-go"
	"github.com/meowdada/go-fcache/backend"
	"github.com/meowdada/go-fcache/cache"
	"github.com/meowdada/go-fcache/policy"
)

// OnceHandler is a handler for once method.
type OnceHandler func(
	preconditionCheck func(item cache.Item) error,
	putCacheFn func(path string, size int64) error,
	rollback func(path string) error,
) (cache.Item, error)

// Manager manages transactions of file caches.
type Manager struct {
	cap       int64
	usage     int64
	db        cache.DB
	policy    policy.Policy
	retryOpts []retry.Option
	mu        sync.RWMutex
}

// New creates an instance of file cache manager.
func New(opts Options) *Manager {
	return &Manager{
		cap:       opts.Capacity,
		db:        backend.Adapter(opts.Backend, opts.Codec),
		policy:    opts.CachePolicy,
		retryOpts: opts.RetryOptions,
	}
}

// Cap returns the capacity of the cache volume.
func (mgr *Manager) Cap() int64 {
	return mgr.cap
}

// Set sets a file as a cache record into the manager. If the cache volume is full,
// it will try emit some cache items to cleanup some cache space then insert this one.
// It is possible that no cache items could be emitted at the moment which leads to this
// operation be unavailable. To prevent waiting deadlock, by default we use timeout setting
// and retry mechanism internally to prevent this condition.
func (mgr *Manager) Set(key string, size int64) error {
	// First, we must make sure that the cache volume is able to
	// fit the item. Or it is impossible to handle this cache item.
	if size > mgr.cap {
		return ErrCacheTooLarge
	}

	// Define retry function. This retry function locks the manager and
	// keep checking that if it is possible to put the cache item into
	// cache volume. If the volume space can fit the cache item, the manager
	// will make its backend to handle it. If the cache volume does not has
	// enough space, it will try cleaning up some space for it. After that,
	// re-check if it is possible to insert the cache item.
	return mgr.retryPutCache(key, size)
}

// Get gets the cache item record from the cache volume. If it failed to
// find the cache item. It returns a zero valued Item and error as ErrCacheMiss.
func (mgr *Manager) Get(key string) (item cache.Item, err error) {
	mgr.rlockFn(func() {
		item, err = mgr.db.Get(key)
	})
	return item, err
}

// Once try get a cache item from the cache volume first. If the cache item has
// been found, it will return it immediately. If not, it will invoke the given
// lambda createFn to create the file cache, then insert it to the cache volume.
// And finally, return the inserted cache item as result.
func (mgr *Manager) Once(path string, createFn OnceHandler) (item cache.Item, err error) {
	mgr.rlockFn(func() {
		item, err = mgr.db.Get(path)
	})
	if err == nil {
		return item, err
	}
	return createFn(mgr.preconditionCheck, mgr.retryPutCache, mgr.rollback)
}

// Register register file caches with their key and increment their
// reference count. With normal cache replacement policy, a referenced
// file cache will not be pick as a victim for cache replacement. If you
// no longer need these files, Unregister it if you have once registered
// them. It is also possible to register a key which does not present so
// far.
func (mgr *Manager) Register(keys ...string) {
	mgr.lockFn(func() {
		mgr.register(keys...)
	})
}

// Unregister unregister file caches with their key and decrement their
// reference count. Any file caches with no reference count unregistered
// by this function will remains the same status.
func (mgr *Manager) Unregister(keys ...string) {
	mgr.lockFn(func() {
		mgr.unregister(keys...)
	})
}

func (mgr *Manager) register(keys ...string) {
	mgr.db.IncrRef(keys...)
}

func (mgr *Manager) unregister(keys ...string) {
	mgr.db.DecrRef(keys...)
}

func (mgr *Manager) lockFn(fn func()) {
	mgr.mu.Lock()
	fn()
	mgr.mu.Unlock()
}

func (mgr *Manager) rlockFn(fn func()) {
	mgr.mu.RLock()
	fn()
	mgr.mu.RUnlock()
}

func (mgr *Manager) preconditionCheck(item cache.Item) error {
	if item.Size > mgr.cap {
		return ErrCacheTooLarge
	}
	return nil
}

func (mgr *Manager) retryPutCache(key string, size int64) error {
	return retry.Do(func() error {
		return mgr.set(key, size)
	}, mgr.retryOpts...)
}

func (mgr *Manager) set(key string, size int64) error {
	var (
		db     = mgr.db
		policy = mgr.policy
	)

	// Cache volume is able to fit the cache item.
	if mgr.usage+size <= mgr.cap {
		err := db.Put(key, size)
		if err != nil {
			return err
		}

		// Only increment the usage if and only if the PUT action
		// finished successfully.
		mgr.usage += size
		return nil
	}

	// When cache volume is unable to fit the cache item, emit
	// a victim from the cache to cleanup some space for it.
	var item cache.Item
	item, err := policy.Emit(db)
	if err != nil {
		return err
	}

	err = item.Remove()
	if err != nil {
		return err
	}

	err = db.Remove(item.Key)
	if err != nil {
		return err
	}
	mgr.usage -= item.Size

	// If the cache volume still cannot fit the cache item. Return
	// a specific error and keep trying.
	if mgr.usage+size > mgr.cap {
		return errRetry
	}

	// If the cache volume is able to fit the cache item after emiting
	// a victim, then put it into the cache space.
	err = db.Put(key, size)
	if err != nil {
		return err
	}
	mgr.usage += size
	return nil
}

func (mgr *Manager) rollback(key string) (err error) {
	mgr.lockFn(func() {
		err = mgr.db.Remove(key)
	})
	return err
}
