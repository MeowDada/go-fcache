package main

import (
	"os"

	"github.com/dustin/go-humanize"
	"github.com/meowdada/go-fcache"
	"github.com/meowdada/go-fcache/backend/gomap"
	"github.com/meowdada/go-fcache/codec"
	"github.com/meowdada/go-fcache/policy"
)

// Note that this is just a example not best practice. Please
// do not use it directly.
func main() {
	// Create a cache manager to manage file caches.
	mgr := fcache.New(fcache.Options{
		Capacity:     int64(2 * humanize.GiByte),
		Codec:        codec.Gob{},
		Backend:      gomap.New(),
		CachePolicy:  policy.LRU(),
		RetryOptions: nil,
	})

	// Put a 500MiB file cache into the manager.
	err := mgr.Set("/path/to/file", int64(500*humanize.MiByte))
	if err != nil {
		panic(err)
	}

	// Get the file cache from the manager.
	item, err := mgr.Get("/path/to/file")
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
