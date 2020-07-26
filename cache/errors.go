package cache

import (
	"github.com/pkg/errors"
)

// ErrNoSuchKey raises when try accessing a Store with key nonexist.
var ErrNoSuchKey = errors.New("no such key")
