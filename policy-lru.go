package fcache

import "time"

// LRU implements policy interface.
type lru struct {
	validator func(item Item) bool
}

// LRU returns a LRU (least recenctly used) cache replacement policy instance.
func LRU(opts ...PolicyOption) Policy {
	opt := combine(opts...)
	return lru{validator: opt.Validate}
}

// Emit implements LRU cache replacement policy.
func (lru lru) Emit(db DB) (victim Item, err error) {
	least := time.Now()
	err = db.Iter(func(k string, v Item) error {
		if !lru.validator(v) {
			return nil
		}
		if v.LastUsed.Before(least) {
			least = v.LastUsed
			victim = v
		}
		return nil
	})
	if victim.IsZero() && err == nil {
		return victim, ErrNoEmitableCaches
	}
	return victim, err
}
