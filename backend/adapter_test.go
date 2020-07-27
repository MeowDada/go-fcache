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
		{
			"valid put with dummy",
			func() error {
				ada := Adapter(gomap.New(), codec.Gob{})
				err := ada.IncrRef("123")
				if err != nil {
					return err
				}
				return ada.Put("123", 456)
			},
			nil,
		},
		{
			"invalid put with duplicate key",
			func() error {
				ada := Adapter(gomap.New(), codec.Gob{})
				err := ada.Put("123", 456)
				if err != nil {
					t.Error(err)
				}
				return ada.Put("123", 456)
			},
			ErrDupKey,
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
	testcases := []struct {
		description string
		scenario    func() error
		expectErr   error
	}{
		{
			"missing key",
			func() error {
				ada := Adapter(gomap.New(), codec.Gob{})
				_, err := ada.Get("123")
				return err
			},
			cache.ErrNoSuchKey,
		},
		{
			"valid get",
			func() error {
				ada := Adapter(gomap.New(), codec.Gob{})
				err := ada.Put("123", 456)
				if err != nil {
					return err
				}

				_, err = ada.Get("123")
				return err
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

func TestAdapterRemove(t *testing.T) {
	testcases := []struct {
		description string
		scenario    func() error
		expectErr   error
	}{
		{
			"remove missing key",
			func() error {
				ada := Adapter(gomap.New(), codec.Gob{})
				return ada.Remove("123")
			},
			nil,
		},
		{
			"valid remove",
			func() error {
				ada := Adapter(gomap.New(), codec.Gob{})
				err := ada.Put("123", 456)
				if err != nil {
					return err
				}

				err = ada.Remove("123")
				if err != nil {
					return err
				}

				_, err = ada.Get("123")
				return err
			},
			cache.ErrNoSuchKey,
		},
	}

	for idx, tc := range testcases {
		err := tc.scenario()
		if err != tc.expectErr {
			t.Errorf("[#Case %d] %s: expect %v, but get %v", idx, tc.description, tc.expectErr, err)
		}
	}
}

func TestAdapterIncrRef(t *testing.T) {
	testcases := []struct {
		description string
		scenario    func() error
		expectErr   error
	}{
		{
			"valid incr reference",
			func() error {
				ada := Adapter(gomap.New(), codec.Gob{})
				return ada.IncrRef("123", "456", "789")
			},
			nil,
		},
		{
			"error when inner put dummy",
			func() error {
				ada := Adapter(mock{
					put: func(k, v []byte) error { return errMock },
					get: func(k []byte) ([]byte, error) { return nil, cache.ErrNoSuchKey },
				}, codec.Gob{})

				return ada.IncrRef("123", "456", "789")
			},
			errMock,
		},
		{
			"error when inner put real",
			func() error {
				ada := Adapter(mock{
					put: func(k, v []byte) error { return errMock },
					get: func(k []byte) ([]byte, error) {
						item := cache.New(123, "123", 123)
						codec := codec.Gob{}
						return codec.Marshal(item)
					},
				}, codec.Gob{})
				return ada.IncrRef("123")
			},
			errMock,
		},
	}

	for idx, tc := range testcases {
		err := tc.scenario()
		if err != tc.expectErr {
			t.Errorf("[#Case %d] %s: expect %v, but get %v", idx, tc.description, tc.expectErr, err)
		}
	}
}

func TestAdapterDecrRef(t *testing.T) {
	testcases := []struct {
		description string
		scenario    func() error
		expectErr   error
	}{
		{
			"ignore no such key",
			func() error {
				ada := Adapter(gomap.New(), codec.Gob{})
				return ada.DecrRef("123", "456", "789")
			},
			nil,
		},
		{
			"error when inner get",
			func() error {
				ada := Adapter(mock{
					get: func(k []byte) ([]byte, error) { return nil, errMock },
				}, codec.Gob{})

				return ada.DecrRef("123", "456", "789")
			},
			errMock,
		},
		{
			"valid decr",
			func() error {
				ada := Adapter(gomap.New(), codec.Gob{})
				err := ada.IncrRef("123")
				if err != nil {
					return err
				}
				return ada.DecrRef("123")
			},
			nil,
		},
		{
			"error when decr inner put",
			func() error {
				m := gomap.New()
				ada := Adapter(mock{
					put: func(k, v []byte) error { return errMock },
					get: m.Get,
				}, codec.Gob{})

				err := ada.IncrRef("123")
				if err != nil {
					return err
				}
				return ada.DecrRef("123")
			},
			errMock,
		},
	}

	for idx, tc := range testcases {
		err := tc.scenario()
		if err != tc.expectErr {
			t.Errorf("[#Case %d] %s: expect %v, but get %v", idx, tc.description, tc.expectErr, err)
		}
	}
}

func TestAdapterClose(t *testing.T) {
	ada := Adapter(gomap.New(), codec.Gob{})
	err := ada.Close()
	if err != nil {
		t.Error(err)
	}
}
