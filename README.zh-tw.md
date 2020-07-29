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

其中最重要的, 就是當使用者把`檔案`作為快取時, 這些套件幾乎都不能保證這些檔案在被有人讀取的情況下不被快取機制替換掉. 若要做一個檔案儲存相關的應用, 這點將會非常致命. 此外, 若是當 cache miss 時想要將檔案重新加入 cache, 在某些情況下這樣做的成本將會非常高. 比如說這個遺失的快取必須由 S3 雲拉下來才能放入快取, 整個速度與金錢成本將大幅提昇.

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
	Evict(pool cache.Pool) (cache.Item, error)
}
```

而 `cache.Pool` 又提供了以下APIs:
```golang
type Pool interface {
	Iter(iterCb func(k string, v Item) error) error
	Put(key string, size int64) error
	Get(key string) (Item, error)
	Remove(key string) error
	IncrRef(keys ...string) error
	DecrRef(keys ...string) error
	Close() error
}
```
通常情況只須使用到 `Pool.Iter` 函式就足以實作自己的快取演算法.

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

	mgr.Register("file1.tmp")
	
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
	mgr.Unregister("file1.tmp")
}

// Download a file from s3.
func downloadFileFromS3() error {
	// Your code logic...
	return nil
}
```
## Q & A
### 為何需要 go-fcache ?
毫不意外, 最多人想問的問題大概就是這個. 大家的疑問是正確的, 其實絕大部分的情況下, 我們並不需要 go-fcache. 因為常見的快取函式庫只要稍加包裝就可以滿足大家的需求.

但若是你對於快取的機制要求更多的控制權,以及主要希望以檔案的形式來作快取,那便有足夠的理由使用此函式庫.因為go-fcache提供了快取後端與演算法的通用界面.只要能滿足該界面,便可直接套用至go-fcache.除此之外,最大的用處還是在於保證檔案快取不會在有人使用時被刪除導致嚴重錯誤.

### 為何需要 Reigster 與 Unregister ?
為了實作`不替換掉有人正在使用的快取`. 最直接方式的方式便是預先宣告說`我想使用哪些快取檔案,請不要刪除它`. 有了此前提,我們在每次透過快取替換演算法想要踢除快取時,便會先檢查是否有人事前宣告過此快取,若是有人正在使用,則跳過剔除該快取,改找下一個.反之,若是沒人正在使用,那此快取便可被直接剔除.

### 是否有可能因為快取宣告的緣故而導致死結(deadlock) ?
某種程度來說, 有的. 當大部分甚至全部的快取都被宣告時, 十分有可能無法騰出足夠空間放入新的快取.但若是在創建快取管理者時有設定好 RetryOptions ,可參考如下:
```golang
fcache.New(fcache.Options{
	...
	RetryOptions: []retry.Option{
		retry.Attempts(100),
		retry.MaxDelay(time.Millisecond),
	}
})
```
在配置之下, 快取在嘗試插入新快取時, 若是遇到失敗,將會不斷重試直到次數達到100次,若100次之內無法完成插入,那將會回傳錯誤.如此一來便可避免死結發生.

### Register 與 Unregister 的使用情境 ?
假設現在有一個應用是會將S3雲的部份檔案快取在本機端,以提高使用者檔案存取效率.當使用者試圖存取一份檔案,而此檔案在本機端快取未命中,此時必須從S3雲將該檔案拉下來並存為快取.但檔案快取的容量有限,若是在將該檔案放入快取並讓使用者存取的同時,有其他人也試圖存取其他快取未命中的檔案.此時就有可能發生快取演算法將這個正在被使用者讀取的檔案剔除的情況.因而導致使用者讀取到不完整的內容.

想要避免這個情況,就要在下載檔案前先Register該檔案,同時確保快取空間有足夠容量能夠放下此cache,當這些前提被滿足時,便可開始下載檔案至快取空間並將該檔案加入快取中.這便是範例`快取失敗即加入`所展示的內容

### 為何用 Pool.Iter 來找出可被踢出的快取 ?
很明顯的,遍歷是非常慢的.但由於`go-fcache`與其他常見的快取機制不同, 具備自定義挑選剔除快取的選項. 在許多情況下, 直接找出可被剔除的對象, 都很有可能無法被直接踢出(因為前述自定義選項的關係). 最糟的情況, 便是即使遍歷完整個資料結構都還無法找出合適的剔除對象. 因此就算是針對 LRU 提供給一個查詢 least frequent used cache item 的 API. 可用性仍大大侷限. 最終還是得遍歷整個資料結構.

## 注意事項
此套件還未經過大量驗證,且仍處於早期開發階段.在穩定版尚未發布之前,請勿將其導入專案使用,否則後果自負.