package backend

import (
	"github.com/pkg/errors"
)

// ErrDupKey raises when try inserting duplicated key.
var ErrDupKey = errors.New("cache key duplicates")
