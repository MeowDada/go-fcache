package fcache

// rr implements policy interface.
type rr struct {
	validator func(item Item) bool
}

// RR returns a RR (random replacement) cache replacement policy instance.
func RR(opts ...PolicyOption) Policy {
	opt := combine(opts...)
	return rr{validator: opt.Validate}
}

// Emit implements MRU cache replacement policy.
func (rr rr) Emit(db DB) (victim Item, err error) {
	err = db.Iter(func(k string, v Item) error {
		if !rr.validator(v) {
			return nil
		}
		victim = v
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
