package fcache

import "time"

// FIFO implements policy interface.
type fifo struct {
	validator func(item Item) bool
}

// FIFO returns a FIFO (first-in-first-out) cache replacement policy instance.
func FIFO(opts ...PolicyOption) Policy {
	opt := combine(opts...)
	return fifo{validator: opt.Validate}
}

// Emit implements FIFO cache replacement policy.
func (fifo fifo) Emit(db DB) (victim Item, err error) {
	least := time.Now()
	err = db.Iter(func(k string, v Item) error {
		if !fifo.validator(v) {
			return nil
		}
		if v.CreatedAt.Before(least) {
			victim = v
		}
		return nil
	})
	if err != nil {
		return victim, err
	}
	if victim.IsZero() {
		return victim, ErrNoEmitableCaches
	}
	return victim, nil
}
