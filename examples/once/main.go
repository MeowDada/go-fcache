package main

import (
	"os"

	"github.com/meowdada/go-fcache"
	"github.com/meowdada/go-fcache/backend/gomap"
	"github.com/meowdada/go-fcache/cache"
	"github.com/meowdada/go-fcache/codec"
	"github.com/meowdada/go-fcache/policy"
)

// Note that this is just a example not best practice. Please
// do not use it directly.
func main() {
	// Create a cache manager to manage file caches.
	mgr := fcache.New(fcache.Options{
		Capacity:     int64(1000),
		Codec:        codec.Gob{},
		Backend:      gomap.New(),
		CachePolicy:  policy.LRU(),
		RetryOptions: nil,
	})

	// onceHandler will be invoked when try getting a missing cache item
	// by calling Once.
	onceHandler := func(
		preconditionChecker func(item cache.Item) error,
		putCacheFn func(path string, size int64) error,
		rollback func(path string) error,
	) (item cache.Item, err error) {
		// Create a psudo cache item first.
		item = cache.New(1000, "file1.tmp", 200)

		// Check if the cache item is valid or not.
		err = preconditionChecker(item)
		if err != nil {
			return item, err
		}

		// Put the psudo cache item into the cache manager first to ensure
		// there is enough space to insert this cache.
		err = putCacheFn("file1.tmp", 200)
		if err != nil {
			return item, err
		}

		// Prepare to rollback if download file failed.
		defer func() {
			if err != nil {
				rollback("file1.tmp")
			}
		}()

		// Then, download it from the cloud.
		err = downloadFileFromS3()
		return item, err
	}

	item, err := mgr.Once("file1.tmp", onceHandler)
	if err != nil {
		panic(err)
	}

	// Get the file reader and use it as you want.
	f, err := os.Open(item.Path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Your code logic...
}

// Download a file from s3.
func downloadFileFromS3() error {
	// Your code logic...
	return nil
}
