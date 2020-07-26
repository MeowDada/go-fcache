package gomap

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/meowdada/go-fcache/cache"
)

func TestPut(t *testing.T) {
	testcases := []struct {
		description string
		k           []byte
		v           []byte
		expectErr   error
	}{
		{
			"valid put",
			[]byte("123"),
			[]byte("456"),
			nil,
		},
		{
			"put nil bytes as key",
			nil,
			[]byte("nil"),
			nil,
		},
		{
			"put nil value",
			[]byte("nil value"),
			nil,
			nil,
		},
		{
			"put nil key-value pair",
			nil,
			nil,
			nil,
		},
	}

	m := New()
	for idx, tc := range testcases {
		desc := tc.description
		err := m.Put(tc.k, tc.v)
		if err != tc.expectErr {
			t.Errorf("[Case#%d]%s: expect %v, but get %v", idx, desc, tc.expectErr, err)
		}
	}
}

func TestGet(t *testing.T) {
	testcases := []struct {
		description string
		m           *Map
		k           []byte
		expectV     []byte
		expectErr   error
	}{
		{
			"valid get",
			&Map{
				ma: map[string][]byte{
					"bang": []byte("dream"),
				},
			},
			[]byte("bang"),
			[]byte("dream"),
			nil,
		},
		{
			"no such key",
			New(),
			[]byte("hello"),
			nil,
			cache.ErrNoSuchKey,
		},
		{
			"get nil key",
			&Map{
				ma: map[string][]byte{"": nil},
			},
			nil,
			nil,
			nil,
		},
	}

	for idx, tc := range testcases {
		desc := tc.description
		v, err := tc.m.Get(tc.k)
		if err != tc.expectErr {
			t.Errorf("[Case#%d] %s: expect %v, but get %v", idx, desc, tc.expectErr, err)
		}
		if bytes.Compare(v, tc.expectV) != 0 {
			t.Errorf("[Case#%d] %s: expect %v, but get %v", idx, desc, tc.expectV, v)
		}
	}
}

func TestRemove(t *testing.T) {
	testcases := []struct {
		description string
		m           *Map
		k           []byte
		expectErr   error
	}{
		{
			"valid remove",
			&Map{ma: map[string][]byte{"haha": []byte("1233")}},
			[]byte("haha"),
			nil,
		},
		{
			"remove nil key",
			&Map{ma: map[string][]byte{"": []byte("1234")}},
			nil,
			nil,
		},
		{
			"remove a unexist key",
			New(),
			[]byte("hello"),
			nil,
		},
	}

	for idx, tc := range testcases {
		desc := tc.description
		err := tc.m.Remove(tc.k)
		if err != tc.expectErr {
			t.Errorf("[Case#%d] %s: expect %v, but get %v", idx, desc, tc.expectErr, err)
		}

		_, err = tc.m.Get(tc.k)
		if err != cache.ErrNoSuchKey {
			t.Errorf("[Case#%d] %s: expect key %s to be deleted, but it does not", idx, desc, tc.k)
		}
	}
}

func TestIter(t *testing.T) {

	mockErr := fmt.Errorf("mock error")

	testcases := []struct {
		description string
		scenario    func() error
		expectErr   error
	}{
		{
			"Iterate all elements",
			func() error {
				m := New()
				list := map[string][]byte{
					"happy":    []byte("new year"),
					"oh":       []byte("my god"),
					"yakusoku": []byte("mamoru"),
				}
				for k, v := range list {
					err := m.Put([]byte(k), v)
					if err != nil {
						return err
					}
				}
				iterCb := func(k, v []byte) error {
					_, ok := list[string(k)]
					if !ok {
						return cache.ErrNoSuchKey
					}
					return nil
				}
				return m.Iter(iterCb)
			},
			nil,
		},
		{
			"Iterate with error",
			func() error {
				m := New()
				list := map[string][]byte{
					"happy":    []byte("new year"),
					"oh":       []byte("my god"),
					"yakusoku": []byte("mamoru"),
				}
				for k, v := range list {
					err := m.Put([]byte(k), v)
					if err != nil {
						return err
					}
				}
				iterCb := func(k, v []byte) error {
					return mockErr
				}
				return m.Iter(iterCb)
			},
			mockErr,
		},
	}

	for idx, tc := range testcases {
		err := tc.scenario()
		if err != tc.expectErr {
			t.Errorf("[Case#%d] %s: expect %v, but get %v", idx, tc.description, tc.expectErr, err)
		}
	}
}

func TestClose(t *testing.T) {
	m := New()
	err := m.Close()
	if err != nil {
		t.Fatal(err)
	}
}
