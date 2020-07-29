package policy

import "github.com/meowdada/go-fcache/cache"

// Policy is a cache replacement algorithm which able to Evict a cache item
// from a cache pool. If it is unable to evict any cache item so far, return
// a special errors as ErrNoEmitableCaches.
type Policy interface {

	// Evict evicts a cache item from the cache pool. If not cache item available,
	// it should return a special error as ErrNoEmitableCaches.
	Evict(pool cache.Pool) (cache.Item, error)
}
