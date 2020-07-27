go-fcache
=====
[![Go Report Card](https://goreportcard.com/badge/github.com/MeowDada/go-fcache)](https://goreportcard.com/report/github.com/MeowDada/go-fcache)
[![codecov](https://codecov.io/gh/meowdada/go-fcache/branch/master/graph/badge.svg)](https://codecov.io/gh/meowdada/go-fcache)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/meowdada/go-fcache)](https://pkg.go.dev/github.com/meowdada/go-fcache)
[![Build Status](https://travis-ci.org/MeowDada/go-fcache.svg?branch=master)](https://travis-ci.org/MeowDada/go-fcache)
[![LICENSE](https://img.shields.io/github/license/meowdada/go-fcache)](https://github.com/MeowDada/go-fcache/blob/master/LICENSE)

go-fcache is a file cache library implemented by pure Go. This package provides simple interface for file caching which allows you to peform some operations to these file caches.

## Motivation
There are many existing and awesome caching libraries implemented in Go, such as [go-cache](https://github.com/patrickmn/go-cache), [bigCache](https://github.com/allegro/bigcache), [redis-cache](https://github.com/go-redis/cache) ...etc.

But most of them stored data in memory instead of a file. Some of them do not handle the upper bound of what they can store. Furthermore, cache replacement mechanism is also hide from user, so we cannot substitute it easily. And the most important thing is that they can merely guarantee that not to evit the caches which are used by others. This might leads severe errors when developing a storage related application.

Due to the reasons above, `go-fcache` comes up with
focusing file cache only, and guaranteens that not to evict the referenced cache item from cache volume.

## Features
* Using Key-Value store as backend
* Support upper bound of a cache volume
* Only evict file caches when it needs space to store new caches.
* Ensures that not to evict file caches which are being referenced.
* Build-in common cache replacement algorithms
* Support concurrent usage
* Simple interface with high scalabilities
* Cache replacement alogrithm and storing backend are customizable

## Built-in cache replacement algorithms
* FIFO (First-in-first-out)
* LRU (Least Recently Used)
* MRU (Most Recently Used)
* RR (Random Replacement)

## Built-in backend
* [gomap](https://github.com/MeowDada/go-fcache/blob/master/backend/gomap/gomap.go) (its actually a golang build-in map with locking)
* [boltdb](https://github.com/MeowDada/go-fcache/blob/master/backend/boltdb/boltdb.go) (https://github.com/etcd-io/bbolt)

## Customization
### How to customize a cache replacement algorithm
Every object which implements cache.Policy interface could be used as a cache replacement algorithm.
```golang
import "github.com/meowdada/go-fcache/cache"

type Policy interface {
	Emit(db cache.DB) (cache.Item, error)
}
```

And `cache.DB` provides following APIs:
```golang
type DB interface {
	Iter(iterCb func(k string, v Item) error) error
	Put(key string, size int64) error
	Get(key string) (Item, error)
	Remove(key string) error
	IncrRef(keys ...string) error
	DecrRef(keys ...string) error
	Close() error
}
```
In most cases, only `DB.Iter` needs to be invoked to implement a cache replacement algorithm.

### How to customize a storing backend.
Every object which implements `backend.Store` interface could be refered as a storing backend.
```golang
type Store interface {
	Put(k, v []byte) error
	Get(k []byte) (v []byte, e error)
	Remove(k []byte) error
	Iter(iterCb func(k, v []byte) error) error
	Close() error
}
```
Note that you must return cache.ErrNoSuchKey when cache is missing, or the functionalities might break up.

## Examples
### Simple cache manager
Creates a file cache manager to manages file caches with limited cache volume.
```golang

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
```