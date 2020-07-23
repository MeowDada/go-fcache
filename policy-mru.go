package fcache

import "time"

// mru implements policy interface.
type mru struct {
	validator func(item Item) bool
}

// MRU returns a MRU (most-recently-used) cache replacement policy instance.
func MRU(opts ...PolicyOption) Policy {
	opt := combine(opts...)
	return mru{validator: opt.Validate}
}

// Emit implements MRU cache replacement policy.
func (mru mru) Emit(db DB) (victim Item, err error) {
	least := time.Time{}
	err = db.Iter(func(k string, v Item) error {
		if !mru.validator(v) {
			return nil
		}
		if v.LastUsed.After(least) {
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
