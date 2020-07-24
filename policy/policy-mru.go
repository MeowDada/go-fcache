package policy

import (
	"time"

	"github.com/meowdada/go-fcache/cache"
)

// mru implements policy interface.
type mru struct {
	validator func(item cache.Item) bool
}

// MRU returns a MRU (most-recently-used) cache replacement policy instance.
func MRU(opts ...Option) Policy {
	opt := combine(opts...)
	return mru{validator: opt.Validate}
}

// Emit implements MRU cache replacement policy.
func (mru mru) Emit(db cache.DB) (victim cache.Item, err error) {
	least := time.Time{}
	err = db.Iter(func(k string, v cache.Item) error {
		if !mru.validator(v) {
			return nil
		}
		if v.ATime().After(least) {
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
