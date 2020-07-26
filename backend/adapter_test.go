package backend

import (
	"testing"

	"github.com/meowdada/go-fcache/backend/gomap"
	"github.com/meowdada/go-fcache/cache"
	"github.com/meowdada/go-fcache/codec"
)

func TestAdapterIter(t *testing.T) {
	testcases := []struct {
		description string
		adapter     cache.DB
		iterFn      func(k string, v cache.Item) error
		expectErr   error
	}{
		{
			"valid iter",
			Adapter(mock{
				iter: func(func(k, v []byte) error) error { return nil },
			}, codec.Mock{
				UnmarshalFn: func(b []byte, v interface{}) error { return nil },
			}),
			func(k string, v cache.Item) error { return nil },
			nil,
		},
		{
			"iterFn with error",
			Adapter(mock{
				iter: func(iterFn func(k, v []byte) error) error {
					return iterFn(nil, nil)
				},
			}, codec.Mock{
				UnmarshalFn: func(b []byte, v interface{}) error { return nil },
			}),
			func(k string, v cache.Item) error { return errMock },
			errMock,
		},
	}

	for idx, tc := range testcases {
		ada, iterCb := tc.adapter, tc.iterFn
		err := ada.Iter(iterCb)
		if err != tc.expectErr {
			t.Errorf("[#Case %d] %s: expect %v, but get %v", idx, tc.description, tc.expectErr, err)
		}
	}
}

func TestAdapterPut(t *testing.T) {
	testcases := []struct {
		description string
		scenario    func() error
		expectErr   error
	}{
		{
			"valid put",
			func() error {
				ada := Adapter(gomap.New(), codec.Gob{})
				return ada.Put("123", 456)
			},
			nil,
		},
	}

	for idx, tc := range testcases {
		err := tc.scenario()
		if err != tc.expectErr {
			t.Errorf("[#Case %d] %s: expect %v, but get %v", idx, tc.description, tc.expectErr, err)
		}
	}
}

func TestAdapterGet(t *testing.T) {

}

func TestAdapterRemove(t *testing.T) {

}

func TestAdapterIncrRef(t *testing.T) {

}

func TestAdapterDecrRef(t *testing.T) {

}

func TestAdapterClose(t *testing.T) {

}
