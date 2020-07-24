package fcache

import "testing"

func TestMockDB(t *testing.T) {
	m := mockDB{
		iter: func(iterCb func(k string, v Item) error) error {
			return errMockErr
		},
		put: func(key string, size int64) error {
			return errMockErr
		},
		get: func(key string) (Item, error) {
			return Item{}, errMockErr
		},
		rm: func(key string) error {
			return errMockErr
		},
		incrRef: func(keys ...string) error {
			return errMockErr
		},
		decrRef: func(keys ...string) error {
			return errMockErr
		},
		close: func() error {
			return errMockErr
		},
	}
	if err := m.Iter(nil); err == nil {
		t.Error("expect error, but get no error")
	}
	if err := m.Put("123", 456); err == nil {
		t.Error("expect error, but get no error")
	}
	if _, err := m.Get("456"); err == nil {
		t.Error("expect error, but get no error")
	}
	if err := m.Remove("456"); err == nil {
		t.Error("expect error, but get no error")
	}
	if err := m.IncrRef("123"); err == nil {
		t.Error("expect error, but get no error")
	}
	if err := m.DecrRef(); err == nil {
		t.Error("expect error, but get no error")
	}
	if err := m.Close(); err == nil {
		t.Error("expect error, but get no error")
	}
}
