package fcache

import "time"

type validateOption struct {
	AllowPsudo      bool
	AllowReferenced bool
	MinimalUsed     int
	MinimalLiveTime time.Duration
}

func newValidateOption() *validateOption {
	return &validateOption{}
}

func combine(opts ...PolicyOption) *validateOption {
	ret := newValidateOption()
	for _, opt := range opts {
		opt.setValidateOption(ret)
	}
	return ret
}

func (opts *validateOption) Validate(item Item) bool {
	if !opts.AllowPsudo && !item.Real {
		return false
	}
	if !opts.AllowReferenced && item.Ref > 0 {
		return false
	}
	if opts.MinimalUsed > item.Used {
		return false
	}
	if opts.MinimalLiveTime > time.Now().Sub(item.CreatedAt) {
		return false
	}
	return true
}

// PolicyOption configures cache replacement policy.
type PolicyOption interface {
	setValidateOption(opts *validateOption)
}

type allowPsudo struct{}

func (allowPsudo) setValidateOption(opts *validateOption) {
	opts.AllowPsudo = true
}

// AllowPsudo returns a cache policy option that allows cacher to emit a psudo
// cache item.
func AllowPsudo() PolicyOption {
	return allowPsudo{}
}

type allowReferenced struct{}

func (allowReferenced) setValidateOption(opts *validateOption) {
	opts.AllowReferenced = true
}

// AllowReferenced returns a cache policy option that allow cacher to emit a
// referenced cache item.
func AllowReferenced() PolicyOption {
	return allowReferenced{}
}

type minimalUsed struct {
	count int
}

func (m minimalUsed) setValidateOption(opts *validateOption) {
	opts.MinimalUsed = m.count
}

// MinimalUsed returns a cache policy option that allow cacher to emit a
// cache item if and only if its used count has equal or larger than
// the specific value.
func MinimalUsed(count int) PolicyOption {
	return minimalUsed{count}
}

type minimalLiveTime struct {
	ttl time.Duration
}

func (m minimalLiveTime) setValidateOption(opts *validateOption) {
	opts.MinimalLiveTime = m.ttl
}

// MinimalLiveTime returns a cache policy option that allow a cacher to
// emit a cache item if and only if its lifetime has equal or greater than
// the specific value.
func MinimalLiveTime(duration time.Duration) PolicyOption {
	return minimalLiveTime{duration}
}
