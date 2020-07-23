go-fcache
=====
[![Go Report Card](https://goreportcard.com/badge/github.com/etcd-io/bbolt?style=flat-square)](https://goreportcard.com/report/github.com/MeowDada/go-fcache)

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

# Goal
* Concurrent safe
* Prfer friendly API interface than performance
* Scalabilities

# Project Status
This project is still at very early stage. DO NOT use it for production environment. Any API interface may changes during development.

# Examples
* TODO

# Performance
* TODO