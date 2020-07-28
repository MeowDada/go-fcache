package policy

import "github.com/meowdada/go-fcache/cache"

// Mock implements policy interface. It's used for testing.
type Mock struct {
	EvictFn func(pool cache.Pool) (cache.Item, error)
}

// Evict implements policy interface.
func (m Mock) Evict(pool cache.Pool) (cache.Item, error) {
	return m.EvictFn(pool)
}
