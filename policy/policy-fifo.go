package policy

import (
	"time"

	"github.com/meowdada/go-fcache/cache"
)

// FIFO implements policy interface.
type fifo struct {
	validator func(item cache.Item) bool
}

// FIFO returns a FIFO (first-in-first-out) cache replacement policy instance.
func FIFO(opts ...Option) Policy {
	opt := combine(opts...)
	return fifo{validator: opt.Validate}
}

// Emit implements FIFO cache replacement policy.
func (fifo fifo) Emit(db cache.DB) (victim cache.Item, err error) {
	least := time.Now()
	err = db.Iter(func(k string, v cache.Item) error {
		if !fifo.validator(v) {
			return nil
		}
		if v.CTime().Before(least) {
			least = v.CTime()
			victim = v
		}
		return nil
	})
	if victim.IsZero() && err == nil {
		return victim, ErrNoEmitableCaches
	}
	return victim, err
}
