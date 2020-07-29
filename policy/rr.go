package policy

import (
	"github.com/meowdada/go-fcache/cache"
)

// rr implements policy interface.
type rr struct {
	validator func(item cache.Item) bool
	onEvict   func(item cache.Item) error
}

// RR returns a RR (random replacement) cache replacement policy instance.
func RR(opts ...Option) Policy {
	opt := combine(opts...)
	return rr{
		validator: opt.Validate,
	}
}

// Evict implements RR cache replacement policy. It will iterates the
// cache pool and return the first cache item which meet all evict
// policy.
func (rr rr) Evict(pool cache.Pool) (victim cache.Item, err error) {
	err = pool.Iter(func(k string, v cache.Item) error {
		if !rr.validator(v) {
			return nil
		}
		victim = v
		return nil
	})

	// If find no victim and without any error, then returns a special
	// error as ErrNoEmitableCaches.
	if victim.IsZero() && err == nil {
		return victim, ErrNoEmitableCaches
	}

	// If any error raises when evicting a cache item. Then return it directly.
	return victim, err
}
