package policy

import "github.com/meowdada/go-fcache/cache"

// Policy is a cache replacement algorithm which able to Evict a cache item
// from a cache pool.
type Policy interface {
	Evict(db cache.Pool) (cache.Item, error)
}
