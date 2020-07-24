package fcache

import (
	retry "github.com/avast/retry-go"
	"github.com/meowdada/go-fcache/backend"
	"github.com/meowdada/go-fcache/codec"
	"github.com/meowdada/go-fcache/policy"
)

// Options configures file cache manager.
type Options struct {
	Capacity     int64
	Codec        codec.Codec
	Backend      backend.Store
	CachePolicy  policy.Policy
	RetryOptions []retry.Option
}
