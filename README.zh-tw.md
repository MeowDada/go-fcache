go-fcache
=====
[![Go Report Card](https://goreportcard.com/badge/github.com/MeowDada/go-fcache)](https://goreportcard.com/report/github.com/MeowDada/go-fcache)
[![codecov](https://codecov.io/gh/meowdada/go-fcache/branch/master/graph/badge.svg)](https://codecov.io/gh/meowdada/go-fcache)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/meowdada/go-fcache)](https://pkg.go.dev/github.com/meowdada/go-fcache)
[![Build Status](https://travis-ci.org/MeowDada/go-fcache.svg?branch=master)](https://travis-ci.org/MeowDada/go-fcache)
[![LICENSE](https://img.shields.io/github/license/meowdada/go-fcache)](https://github.com/MeowDada/go-fcache/blob/master/LICENSE)

`go-fcache` 是一個由純 Golang 實作的 package. 此套件能夠提供一些常用的檔案緩存界面以方便快速存取.

## 動機
由 Golang 實作的快取套件非常多, 諸如 [go-cache](https://github.com/patrickmn/go-cache), [bigCache](https://github.com/allegro/bigcache), [redis-cache](https://github.com/go-redis/cache) ...等. 但這些套件絕大部分都是直接將資料存放在記憶體, 且沒有提供太多更高層次的控制, 比如說快取空間上限, 快取演算法替換...等.

其中最重要的, 就是當使用者把`檔案`作為快取時, 這些套件幾乎都不能保證這些檔案在被有人讀取的情況下不被快取機制替換掉. 若要做一個檔案儲存相關的應用, 這點將會非常致命.

於是乎, `go-fcache`誕生了. 它的目標便是試圖解決以上的問題, 並專注在檔案快取這個部份. 提供呼叫者可靠的併發函式與回調. 同時提供共通界面讓開發者能夠自行定義他們想要的後端與快取演算法.

## 特徵
* 以 KV 儲存作為後端, 支援直接嵌入
* 支援邏輯上的快取空間上限
* 只會在快取空間容量不足時清出足夠空間
* 不使用記憶體作為讀取檔案內容手段
* 確保檔案緩存在有人使用的時候不會被替換掉
* 支援多種常見快取演算法
* 支援安全併發
* 界面簡潔, 具高可擴張性
* 可自行定義快取演算法與儲存後端

## 快取演算法
目前為止, 內建支援的快取演算法如下:
* FIFO (First-in-first-out)
* LIFO (Last-in-first-out)
* LRU (Least Recently Used)
* MRU (Most Recently Used)
* RR (Random Replacement)

## 儲存後端
目前為止, 內建支援的儲存後端如下:
* [gomap](https://github.com/MeowDada/go-fcache/blob/master/backend/gomap/gomap.go) (其實就是golang build-in的map, 只是加了鎖)
* [boltdb](https://github.com/MeowDada/go-fcache/blob/master/backend/boltdb/boltdb.go) (https://github.com/etcd-io/bbolt)

## 自定義 
### 如何自定義快取演算法
任何實作以下界面的資料結構, 皆可作為快取演算法
```golang
import "github.com/meowdada/go-fcache/cache"

type Policy interface {
	Emit(db cache.DB) (cache.Item, error)
}
```

而 `cache.DB` 又提供了以下APIs:
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
通常情況只須使用到 `DB.Iter` 函式就足以實作自己的快取演算法.

### 如何自定義儲存後端
任何實作以下界面的資料結構, 皆可作為儲存後端
```golang
type Store interface {
	Put(k, v []byte) error
	Get(k []byte) (v []byte, e error)
	Remove(k []byte) error
	Iter(iterCb func(k, v []byte) error) error
	Close() error
}
```
但請注意, 若是cache miss的情況下, 依照目前設計必須要回傳cache.ErrNoSuchKey 這個特定 error. 才能確保正常運作.

## 使用範例
### 最簡範例
最基本的快取檔案與取回內容
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

### 快取失敗即加入
此範例展示了如何在快取未命中時, 馬上將欲抓取的快取內容放進快取當中. 熟悉此方式之後便可省去 SET if GET miss 的麻煩流程.
直接使用 ONCE 即可.
```golang
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
```


## 注意事項
此套件還未經過大量驗證,且仍處於早期開發階段.在穩定版尚未發布之前,請勿將其導入專案使用,否則後果自負.