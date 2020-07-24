package policy

import (
	"time"

	"github.com/meowdada/go-fcache/cache"
)

// LRU implements policy interface.
type lru struct {
	validator func(item cache.Item) bool
}

// LRU returns a LRU (least recenctly used) cache replacement policy instance.
func LRU(opts ...Option) Policy {
	opt := combine(opts...)
	return lru{validator: opt.Validate}
}

// Emit implements LRU cache replacement policy.
func (lru lru) Emit(db cache.DB) (victim cache.Item, err error) {
	least := time.Now()
	err = db.Iter(func(k string, v cache.Item) error {
		if !lru.validator(v) {
			return nil
		}
		if v.ATime().Before(least) {
			least = v.ATime()
			victim = v
		}
		return nil
	})
	if victim.IsZero() && err == nil {
		return victim, ErrNoEmitableCaches
	}
	return victim, err
}
