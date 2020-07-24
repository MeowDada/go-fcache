package fcache

import (
	"github.com/pkg/errors"
)

// ErrDupKey raises when try inserting duplicated key.
var ErrDupKey = errors.New("cache key duplicates")

// ErrCacheTooLarge raises when the inserting cache item is too large.
var ErrCacheTooLarge = errors.New("cache item is too large")

// ErrNoEmitableCaches raises when all the cache item cannot be emitable.
var ErrNoEmitableCaches = errors.New("no emitable caches")

// ErrCacheMiss raises when try getting an unexist record.
var ErrCacheMiss = errors.New("cache miss")

var errRetry = errors.New("keep retrying")

var errMockErr = errors.New("mock error")
