package policy

import (
	"github.com/meowdada/go-fcache/cache"
)

// rr implements policy interface.
type rr struct {
	validator func(item cache.Item) bool
}

// RR returns a RR (random replacement) cache replacement policy instance.
func RR(opts ...Option) Policy {
	opt := combine(opts...)
	return rr{validator: opt.Validate}
}

// Emit implements MRU cache replacement policy.
func (rr rr) Evict(pool cache.Pool) (victim cache.Item, err error) {
	err = pool.Iter(func(k string, v cache.Item) error {
		if !rr.validator(v) {
			return nil
		}
		victim = v
		return nil
	})
	if victim.IsZero() && err == nil {
		return victim, ErrNoEmitableCaches
	}
	return victim, err
}
