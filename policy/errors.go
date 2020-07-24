package policy

import "github.com/pkg/errors"

// ErrNoEmitableCaches raises when all the cache item cannot be emitable.
var ErrNoEmitableCaches = errors.New("no emitable caches")
