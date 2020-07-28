package policy

import (
	"testing"

	"github.com/meowdada/go-fcache/cache"
)

func TestMock(t *testing.T) {
	m := Mock{EvictFn: func(cache.Pool) (cache.Item, error) { return cache.Item{}, nil }}
	_, err := m.Evict(nil)
	if err != nil {
		t.Fatal(err)
	}
}
