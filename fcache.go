package fcache

import (
	"sync"

	"github.com/avast/retry-go"
)

// Manager manages transactions of file caches.
type Manager struct {
	cap    int64
	usage  int64
	mu     sync.RWMutex
	db     DB
	policy Policy
}

// New creates an instance of file cache manager.
func New(cap int64, db DB, policy Policy) *Manager {
	return &Manager{
		cap:    cap,
		db:     db,
		policy: policy,
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
func (mgr *Manager) Set(path string, size int64) error {
	var (
		db     = mgr.db
		policy = mgr.policy
	)

	// First, we must make sure that the cache volume is able to
	// fit the item. Or it is impossible to handle this cache item.
	if size > mgr.cap {
		return ErrCacheTooLarge
	}

	return retry.Do(func() (err error) {
		mgr.lockFn(func() {
			if mgr.usage+size > mgr.cap {
				var item Item
				item, err = policy.Emit(db)
				if err != nil {
					return
				}

				err = item.Remove()
				if err != nil {
					return
				}

				err = db.Remove(item.Path)
				if err != nil {
					return
				}

				mgr.usage -= item.Size
			} else {
				err = db.Put(path, size)
				if err != nil {
					return
				}
				mgr.usage += size
				return
			}
			if mgr.usage+size > mgr.cap {
				err = errRetry
				return
			}
		})
		return err
	})
}

// Get gets the cache item record from the cache volume.
func (mgr *Manager) Get(path string) (item Item, err error) {
	var (
		db = mgr.db
	)

	mgr.rlockFn(func() {
		item, err = db.Get(path)
	})
	return item, err
}

// Register register file caches with their key and increment their
// reference count.
func (mgr *Manager) Register(keys ...string) {
	mgr.lockFn(func() {
		mgr.register(keys...)
	})
}

// Unregister unregister file caches with their key and decrement their
// reference count.
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
