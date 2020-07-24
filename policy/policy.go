package policy

import "github.com/meowdada/go-fcache/cache"

// Policy is a cache replacement algorithm which able to emit a cache item.
type Policy interface {
	Emit(db cache.DB) (cache.Item, error)
}
