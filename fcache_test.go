package fcache

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/avast/retry-go"
	"github.com/meowdada/go-fcache/backend/gomap"
	"github.com/meowdada/go-fcache/cache"
	"github.com/meowdada/go-fcache/codec"
	"github.com/meowdada/go-fcache/policy"
)

func TestManagerCap(t *testing.T) {
	cap := int64(1250)
	m := New(Options{Capacity: cap})
	if m.Cap() != cap {
		t.Errorf("expect %v, but get %v", cap, m.Cap())
	}
}

func TestManagerSet(t *testing.T) {
	testcases := []struct {
		description string
		scenario    func() error
		expectErr   error
	}{
		{
			"cache item too large",
			func() error {
				cap := int64(1000)
				m := New(Options{
					Capacity:     cap,
					Codec:        codec.Gob{},
					Backend:      gomap.New(),
					CachePolicy:  policy.LRU(),
					RetryOptions: nil,
				})
				return m.Set("123", cap+1)
			},
			ErrCacheTooLarge,
		},
		{
			"able to fit the cache item directly",
			func() error {
				cap := int64(1000)
				m := New(Options{
					Capacity:     cap,
					Codec:        codec.Gob{},
					Backend:      gomap.New(),
					CachePolicy:  policy.LRU(),
					RetryOptions: nil,
				})
				return m.Set("123", cap/2)
			},
			nil,
		},
		{
			"valid set with emitting",
			func() error {
				cap := int64(1000)
				m := New(Options{
					Capacity:    cap,
					Codec:       codec.Gob{},
					Backend:     gomap.New(),
					CachePolicy: policy.LRU(),
					RetryOptions: []retry.Option{
						retry.MaxDelay(time.Millisecond),
						retry.Attempts(10),
						retry.LastErrorOnly(true),
					},
				})

				err := ioutil.WriteFile("123", nil, 0644)
				if err != nil {
					return err
				}
				defer os.Remove("123")

				err = m.Set("123", 500)
				if err != nil {
					return err
				}
				return m.Set("456", 800)
			},
			nil,
		},
	}

	for idx, tc := range testcases {
		err := tc.scenario()
		if err != tc.expectErr {
			t.Errorf("[#Case%d] %s: expect %v, but get %v", idx, tc.description, tc.expectErr, err)
		}
	}
}

func TestManagerGet(t *testing.T) {
	testcases := []struct {
		description string
		scenario    func() error
		expectErr   error
	}{
		{
			"valid get",
			func() error {
				cap := int64(1000)
				m := New(Options{
					Capacity:     cap,
					Codec:        codec.Gob{},
					Backend:      gomap.New(),
					CachePolicy:  policy.LRU(),
					RetryOptions: nil,
				})
				err := m.Set("123", cap/2)
				if err != nil {
					return err
				}

				_, err = m.Get("123")
				return err
			},
			nil,
		},
		{
			"valid get",
			func() error {
				cap := int64(1000)
				m := New(Options{
					Capacity:     cap,
					Codec:        codec.Gob{},
					Backend:      gomap.New(),
					CachePolicy:  policy.LRU(),
					RetryOptions: nil,
				})
				_, err := m.Get("123")
				return err
			},
			cache.ErrNoSuchKey,
		},
	}

	for idx, tc := range testcases {
		err := tc.scenario()
		if err != tc.expectErr {
			t.Errorf("[#Case%d] %s: expect %v, but get %v", idx, tc.description, tc.expectErr, err)
		}
	}
}

func TestManagerOnce(t *testing.T) {
	testcases := []struct {
		description string
		scenario    func() error
		expectErr   error
	}{
		{
			"valid once",
			func() error {
				cap := int64(1000)
				m := New(Options{
					Capacity:     cap,
					Codec:        codec.Gob{},
					Backend:      gomap.New(),
					CachePolicy:  policy.LRU(),
					RetryOptions: nil,
				})
				err := m.Set("123", cap/2)
				if err != nil {
					return err
				}

				_, err = m.Once("123", nil)
				return err
			},
			nil,
		},
		{
			"valid once put get",
			func() error {
				cap := int64(1000)
				m := New(Options{
					Capacity:     cap,
					Codec:        codec.Gob{},
					Backend:      gomap.New(),
					CachePolicy:  policy.LRU(),
					RetryOptions: nil,
				})
				handler := func(
					preconditionCheck func(cache.Item) error,
					putCacheFn func(path string, size int64) error,
				) (cache.Item, error) {
					item := cache.New(123, "123", 456)
					return item, nil
				}
				_, err := m.Once("123", handler)
				return err
			},
			nil,
		},
	}

	for idx, tc := range testcases {
		err := tc.scenario()
		if err != tc.expectErr {
			t.Errorf("[#Case%d] %s: expect %v, but get %v", idx, tc.description, tc.expectErr, err)
		}
	}
}

func TestManagerRegister(t *testing.T) {
	cap := int64(1000)
	m := New(Options{
		Capacity:     cap,
		Codec:        codec.Gob{},
		Backend:      gomap.New(),
		CachePolicy:  policy.LRU(),
		RetryOptions: nil,
	})
	m.Register("123", "456")
}

func TestManagerUnregister(t *testing.T) {
	cap := int64(1000)
	m := New(Options{
		Capacity:     cap,
		Codec:        codec.Gob{},
		Backend:      gomap.New(),
		CachePolicy:  policy.LRU(),
		RetryOptions: nil,
	})
	m.Unregister("123", "456")
}
