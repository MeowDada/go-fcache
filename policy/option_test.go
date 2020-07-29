package policy

import (
	"testing"
	"time"

	"github.com/meowdada/go-fcache/backend"
	"github.com/meowdada/go-fcache/backend/gomap"
	"github.com/meowdada/go-fcache/codec"
)

func TestPolicyOptionNotAllowPsudo(t *testing.T) {
	testcases := []struct {
		description string
		scenario    func() error
		expectErr   bool
	}{
		{
			"not allow to evict psudo cache",
			func() error {
				pool := backend.Adapter(gomap.New(), codec.Gob{})
				pool.IncrRef("123")
				pool.DecrRef("123")
				rr := RR(NotEvictPsudo())
				_, err := rr.Evict(pool)
				return err
			},
			true,
		},
		{
			"allow to evict psudo cache",
			func() error {
				pool := backend.Adapter(gomap.New(), codec.Gob{})
				pool.IncrRef("123")
				pool.DecrRef("123")
				rr := RR()
				_, err := rr.Evict(pool)
				return err
			},
			false,
		},
	}

	for idx, tc := range testcases {
		err := tc.scenario()
		if err != nil && !tc.expectErr {
			t.Errorf("[#Case%d]: %s, expect no error, but get %v", idx, tc.description, err)
		}
		if err == nil && tc.expectErr {
			t.Errorf("[#Case%d]: %s, expect an error, but get no errors", idx, tc.description)
		}
	}
}

func TestPolicyOptionEvictReferenced(t *testing.T) {
	testcases := []struct {
		description string
		scenario    func() error
		expectErr   bool
	}{
		{
			"not allow to evict a referenced cache",
			func() error {
				pool := backend.Adapter(gomap.New(), codec.Gob{})
				pool.IncrRef("123")
				rr := RR()
				_, err := rr.Evict(pool)
				return err
			},
			true,
		},
		{
			"allow to evict a referenced cache",
			func() error {
				pool := backend.Adapter(gomap.New(), codec.Gob{})
				pool.IncrRef("123")
				rr := RR(EvictReferenced())
				_, err := rr.Evict(pool)
				return err
			},
			false,
		},
	}

	for idx, tc := range testcases {
		err := tc.scenario()
		if err != nil && !tc.expectErr {
			t.Errorf("[#Case%d]: %s, expect no error, but get %v", idx, tc.description, err)
		}
		if err == nil && tc.expectErr {
			t.Errorf("[#Case%d]: %s, expect an error, but get no errors", idx, tc.description)
		}
	}
}

func TestPolicyOptionMinUsed(t *testing.T) {
	testcases := []struct {
		description string
		scenario    func() error
		expectErr   bool
	}{
		{
			"not allow to evict a cache item with min used smaller than setting",
			func() error {
				pool := backend.Adapter(gomap.New(), codec.Gob{})
				pool.IncrRef("123")
				pool.DecrRef("123")
				rr := RR(MinimalUsed(3))
				_, err := rr.Evict(pool)
				return err
			},
			true,
		},
		{
			"allow to evict a referenced cache with min used equal or greater than setting",
			func() error {
				pool := backend.Adapter(gomap.New(), codec.Gob{})
				pool.IncrRef("123")
				pool.IncrRef("123")
				pool.IncrRef("123")
				pool.DecrRef("123")
				pool.DecrRef("123")
				pool.DecrRef("123")
				rr := RR(MinimalUsed(3))
				_, err := rr.Evict(pool)
				return err
			},
			false,
		},
	}

	for idx, tc := range testcases {
		err := tc.scenario()
		if err != nil && !tc.expectErr {
			t.Errorf("[#Case%d]: %s, expect no error, but get %v", idx, tc.description, err)
		}
		if err == nil && tc.expectErr {
			t.Errorf("[#Case%d]: %s, expect an error, but get no errors", idx, tc.description)
		}
	}
}

func TestPolicyOptionMinLiveTime(t *testing.T) {
	testcases := []struct {
		description string
		scenario    func() error
		expectErr   bool
	}{
		{
			"not allow to evict a cache item with live time smaller than setting",
			func() error {
				pool := backend.Adapter(gomap.New(), codec.Gob{})
				pool.Put("123", 456)
				rr := RR(MinLiveTime(time.Second))
				_, err := rr.Evict(pool)
				return err
			},
			true,
		},
		{
			"allow to evict a referenced cache with live time equal or greater than setting",
			func() error {
				pool := backend.Adapter(gomap.New(), codec.Gob{})
				pool.Put("123", 456)
				time.Sleep(time.Millisecond)
				rr := RR(MinLiveTime(time.Millisecond))
				_, err := rr.Evict(pool)
				return err
			},
			false,
		},
	}

	for idx, tc := range testcases {
		err := tc.scenario()
		if err != nil && !tc.expectErr {
			t.Errorf("[#Case%d]: %s, expect no error, but get %v", idx, tc.description, err)
		}
		if err == nil && tc.expectErr {
			t.Errorf("[#Case%d]: %s, expect an error, but get no errors", idx, tc.description)
		}
	}
}

func TestPolicyOptionLastUsed(t *testing.T) {
	testcases := []struct {
		description string
		scenario    func() error
		expectErr   bool
	}{
		{
			"not allow to evict a recently used cache item",
			func() error {
				pool := backend.Adapter(gomap.New(), codec.Gob{})
				pool.IncrRef("123")
				pool.DecrRef("123")
				rr := RR(LastUsed(time.Second))
				_, err := rr.Evict(pool)
				return err
			},
			true,
		},
		{
			"allow to evict a recently used cache item",
			func() error {
				pool := backend.Adapter(gomap.New(), codec.Gob{})
				pool.IncrRef("123")
				pool.DecrRef("123")
				time.Sleep(time.Millisecond)
				rr := RR(LastUsed(time.Millisecond))
				_, err := rr.Evict(pool)
				return err
			},
			false,
		},
	}

	for idx, tc := range testcases {
		err := tc.scenario()
		if err != nil && !tc.expectErr {
			t.Errorf("[#Case%d]: %s, expect no error, but get %v", idx, tc.description, err)
		}
		if err == nil && tc.expectErr {
			t.Errorf("[#Case%d]: %s, expect an error, but get no errors", idx, tc.description)
		}
	}
}
