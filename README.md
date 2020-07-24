go-fcache
=====
[![Go Report Card](https://goreportcard.com/badge/github.com/MeowDada/go-fcache)](https://goreportcard.com/report/github.com/MeowDada/go-fcache)
[![codecov](https://codecov.io/gh/meowdada/go-fcache/branch/master/graph/badge.svg)](https://codecov.io/gh/meowdada/go-fcache)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/meowdada/go-fcache)](https://pkg.go.dev/github.com/meowdada/go-fcache)
[![Build Status](https://travis-ci.org/MeowDada/go-fcache.svg?branch=master)](https://travis-ci.org/MeowDada/go-fcache)
[![LICENSE](https://img.shields.io/github/license/meowdada/go-fcache)](https://github.com/MeowDada/go-fcache/blob/master/LICENSE)

go-fcache is a file cache library implemented by pure Go.

In most cases, we don't actually needs so called file cache. Because we usually put cache items into memory (a single map, or more complex, a database) in format of bytes or some well-wrapped data structure. Its quite easy to access them, and its very fast and efficient.

Everything goes well in most of the web application. But, when it comes to the application of file or object storgae. 
There are some problems here:
1. Lacks of memory to store these caches.
2. Many cache libraries do not allow to store large file.
3. Some of these library does not provides settings for cache capacity upper bound.

# Features
## Psudo file locking
Though we cannot lock a file physically, we still can ensures that NOT TO remove a file cache ( a data structure ) locked by this package within the scope. Which means if a caller didn't explicity remove the file outside the control of this package. The file should still remains unless we delete it by this package.

With this feature can we guarantee that any reading operations on this file should process properly (unless the disk physicall go down).

And with default setting, any referenced file cache won't be removed (emitted) at any time.
## Cache volume capacity
Setup an upper bound of cache volume capacity to make sures that the cache manager will not hold caches that exceed the limitation.

# Goal
* Concurrent safe
* Prfer friendly API interface than performance
* Scalabilities

# Project Status
This project is still at very early stage. DO NOT use it for production environment. Any API interface may changes during development.

# Examples
## Simple cache manager
Creates a file cache manager to manages file caches with limited cache volume.
```go
import (
    fcache "github.com/meowdada/go-fcache"
    humanize "github.com/dustin/go-humanize"
)

// Creates a cache manager with capacity 2GiB, using hashmap as backend, with
// LRU cache replacement policy and default retry option when failed to emit a cache.
mgr := fcache.New(fcache.Options{
    Capacity:     int64(2*humanize.GiByte)
    Backend:      fcache.Hashmap(),
    CachePolicy:  fcache.LRU(),
    RetryOptions: nil,
})

// Adds a file path/to/file with its size as a file cache.
err := mgr.Set("path/to/file", int64(humanized.GiByte))
if err != nil {
    // Handle put error here...
}

// Try getting a file cache from given path. If the cache is missing,
// it will return error as fcache.ErrCacheMiss
item, err := mgr.Get("path/to/file")
if err != nil {
    // Handle get error here...
}
```

## File cache locking
Creates a file cache manager and locks a file to make sure that it will not be
emited from the cache volume.
```go
import (
    fcache "github.com/meowdada/go-fcache"
)

mgr := fcache.New(fcache.Options{
    Capacity:     100,
    Backend:      fcache.Hashmap(),
    CachePolicy:  fcache.RR(),
    RetryOptions: nil,
})

// Adds 3 files as caches
mgr.Set("path/to/file1", 70)
mgr.Set("path/to/file2", 20)
mgr.Set("path/to/file3", 10)

// Lock the file "path/to/file1" to ensure that not to emit this file cache
// from the cache volume. (Maybe some other applications are reading "path/to/file1"'s
// content, so we need to lock it or it may be removed from the disk. )
mgr.Register("path/to/file1")

// Adds file cache with size = 30 bytes, file cache "path/to/file2" and 
// "path/to/file3" should be clean up the manager. 
mgr.Set("path/to/file4", 30)
```


# Performance
* TODO