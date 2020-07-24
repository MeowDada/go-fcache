package fcache

import retry "github.com/avast/retry-go"

// Options configures file cache manager.
type Options struct {
	Capacity     int64
	Backend      DB
	CachePolicy  Policy
	RetryOptions []retry.Option
}
