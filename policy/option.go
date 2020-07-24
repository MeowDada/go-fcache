package policy

import (
	"time"

	"github.com/meowdada/go-fcache/cache"
)

type validateOption struct {
	AllowPsudo      bool
	AllowReferenced bool
	MinimalUsed     int
	MinimalLiveTime time.Duration
}

func newValidateOption() *validateOption {
	return &validateOption{}
}

func combine(opts ...Option) *validateOption {
	ret := newValidateOption()
	for _, opt := range opts {
		opt.setValidateOption(ret)
	}
	return ret
}

func (opts *validateOption) Validate(item cache.Item) bool {
	if !opts.AllowPsudo && !item.IsReal() {
		return false
	}
	if !opts.AllowReferenced && item.Reference() > 0 {
		return false
	}
	if opts.MinimalUsed > item.UsedCount() {
		return false
	}
	if opts.MinimalLiveTime > time.Now().Sub(item.CTime()) {
		return false
	}
	return true
}

// Option configures cache replacement policy.
type Option interface {
	setValidateOption(opts *validateOption)
}

type allowPsudo struct{}

func (allowPsudo) setValidateOption(opts *validateOption) {
	opts.AllowPsudo = true
}

// AllowPsudo returns a cache policy option that allows cacher to emit a psudo
// cache item.
func AllowPsudo() Option {
	return allowPsudo{}
}

type allowReferenced struct{}

func (allowReferenced) setValidateOption(opts *validateOption) {
	opts.AllowReferenced = true
}

// AllowReferenced returns a cache policy option that allow cacher to emit a
// referenced cache item.
func AllowReferenced() Option {
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
func MinimalUsed(count int) Option {
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
func MinimalLiveTime(duration time.Duration) Option {
	return minimalLiveTime{duration}
}
