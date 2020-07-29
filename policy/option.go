package policy

import (
	"time"

	"github.com/meowdada/go-fcache/cache"
)

type validateOption struct {
	NotEvictPsudo   bool
	EvictReferenced bool
	MinUsed         int
	MinLiveTime     time.Duration
	LastUsed        time.Duration
}

func newValidateOption() *validateOption {
	return &validateOption{}
}

// combine combines multiple options to a single, complete validateOption.
func combine(opts ...Option) *validateOption {
	ret := newValidateOption()
	for _, opt := range opts {
		opt.setValidateOption(ret)
	}
	return ret
}

// Validate validates a cache item could be picked as a
// evict cache item or not.
func (opts *validateOption) Validate(item cache.Item) bool {

	// If not allow to evict a psudo cache item and the cache
	// item is psudo, then return false.
	if opts.NotEvictPsudo && !item.IsReal() {
		return false
	}

	// If not allow to evict a referenced cache item and the
	// cache item with at least one referenced, then return false.
	if !opts.EvictReferenced && item.Reference() > 0 {
		return false
	}

	// If not allow to evict a cache item with used count smaller
	// than the setting value and the cache item does satisfies
	// this constrain, then return false.
	if opts.MinUsed > item.UsedCount() {
		return false
	}

	t := time.Now()

	// If the cache item with smaller live time than setting, then
	// return false.
	if opts.MinLiveTime > t.Sub(item.CTime()) {
		return false
	}

	// If the cache item are recently used and the interval is smaller
	// than the setting, then return false.
	if opts.LastUsed > t.Sub(item.ATime()) {
		return false
	}

	// All constrain are satisfied, the cache item is ok to be eivcted.
	return true
}

// Option configures cache replacement policy. A cache item can be evicted as a victim
// if and only if all option are satisfied.
type Option interface {
	setValidateOption(opts *validateOption)
}

type notEvictPsudo struct{}

func (notEvictPsudo) setValidateOption(opts *validateOption) {
	opts.NotEvictPsudo = true
}

// NotEvictPsudo returns a cache policy option that disallow a cacher to evict a psudo
// cache item.
func NotEvictPsudo() Option {
	return notEvictPsudo{}
}

type evictReferenced struct{}

func (evictReferenced) setValidateOption(opts *validateOption) {
	opts.EvictReferenced = true
}

// EvictReferenced returns a cache policy option that allow cacher to evict a
// referenced cache item.
func EvictReferenced() Option {
	return evictReferenced{}
}

type minimalUsed struct {
	count int
}

func (m minimalUsed) setValidateOption(opts *validateOption) {
	opts.MinUsed = m.count
}

// MinimalUsed returns a cache policy option that allow cacher to evict a
// cache item if and only if its used count has equal or greater than
// the setting value.
func MinimalUsed(count int) Option {
	return minimalUsed{count}
}

type minimalLiveTime struct {
	ttl time.Duration
}

func (m minimalLiveTime) setValidateOption(opts *validateOption) {
	opts.MinLiveTime = m.ttl
}

// MinLiveTime returns a cache policy option that allow a cacher to
// evict a cache item if and only if its lifetime has equal or greater than
// the setting value.
func MinLiveTime(duration time.Duration) Option {
	return minimalLiveTime{duration}
}

type lastUsed struct {
	du time.Duration
}

func (l lastUsed) setValidateOption(opts *validateOption) {
	opts.LastUsed = l.du
}

// LastUsed returns a cache policy option that disallow a cacher to
// evict a recently used cache item according to the setting value.
func LastUsed(duration time.Duration) Option {
	return lastUsed{duration}
}
