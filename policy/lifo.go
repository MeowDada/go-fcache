package policy

import (
	"time"

	"github.com/meowdada/go-fcache/cache"
)

// lifo implements policy interface.
type lifo struct {
	validator func(item cache.Item) bool
}

// LIFO returns a LIFO (last-in-first-out) cache replacement policy instance.
func LIFO(opts ...Option) Policy {
	opt := combine(opts...)
	return lifo{validator: opt.Validate}
}

// Emit implements LIFO cache replacement policy.
func (lifo lifo) Emit(db cache.DB) (victim cache.Item, err error) {
	t := time.Time{}
	err = db.Iter(func(k string, v cache.Item) error {
		if !lifo.validator(v) {
			return nil
		}
		if v.CTime().After(t) {
			t = v.CTime()
			victim = v
		}
		return nil
	})
	if victim.IsZero() && err == nil {
		return victim, ErrNoEmitableCaches
	}
	return victim, err
}
