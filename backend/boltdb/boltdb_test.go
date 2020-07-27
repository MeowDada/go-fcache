package boltdb

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/meowdada/go-fcache/cache"
)

func TestPut(t *testing.T) {
	testcases := []struct {
		description string
		k           []byte
		v           []byte
		expectErr   bool
	}{
		{
			"valid put",
			[]byte("123"),
			[]byte("456"),
			false,
		},
		{
			"put nil bytes as key",
			nil,
			[]byte("nil"),
			true,
		},
		{
			"put nil value",
			[]byte("nil value"),
			nil,
			false,
		},
		{
			"put nil key-value pair",
			nil,
			nil,
			true,
		},
	}

	db := New(Options{
		Path:    "bolt.db",
		Mode:    0666,
		Bucket:  "cache",
		Options: nil,
	})
	defer os.Remove("bolt.db")
	for idx, tc := range testcases {
		desc := tc.description
		err := db.Put(tc.k, tc.v)
		if err != nil && !tc.expectErr {
			t.Errorf("[Case#%d]%s: expect no error, but get %v", idx, desc, err)
		}
		if err == nil && tc.expectErr {
			t.Errorf("[Case#%d]%s: expect error occurs, but get no error", idx, desc)
		}
	}

	// Invalid Open
	db2 := New(Options{
		Path:    "/dev/null",
		Mode:    0666,
		Bucket:  "cache",
		Options: nil,
	})
	db2.Put([]byte("123"), []byte("456"))
}

func TestGet(t *testing.T) {
	testcases := []struct {
		description string
		scenario    func() error
		expectErr   bool
	}{
		{
			"invalid init",
			func() error {
				db := New(Options{
					Path:    "/dev/null",
					Mode:    0666,
					Bucket:  "cache",
					Options: nil,
				})
				k := []byte("123")
				_, err := db.Get(k)
				return err
			},
			true,
		},
		{
			"invalid bucket",
			func() error {
				db := New(Options{
					Path:    "bolt.db",
					Mode:    0666,
					Bucket:  "",
					Options: nil,
				})
				defer os.Remove("bolt.db")

				k := []byte("123")
				_, err := db.Get(k)
				return err
			},
			true,
		},
		{
			"valid get",
			func() error {
				db := New(Options{
					Path:    "bolt.db",
					Mode:    0666,
					Bucket:  "cache",
					Options: nil,
				})
				defer os.Remove("bolt.db")

				k, v := []byte("123"), []byte("456")
				err := db.Put(k, v)
				if err != nil {
					return err
				}

				r, err := db.Get(k)
				if err != nil {
					return err
				}
				if bytes.Compare(r, v) != 0 {
					return fmt.Errorf("expect %v, but get %v", v, r)
				}
				return nil
			},
			false,
		},
		{
			"cache miss",
			func() error {
				db := New(Options{
					Path:    "bolt.db",
					Mode:    0666,
					Bucket:  "cache",
					Options: nil,
				})
				defer os.Remove("bolt.db")

				k := []byte("123")
				_, err := db.Get(k)
				if err != nil {
					return err
				}
				return nil
			},
			true,
		},
	}

	for idx, tc := range testcases {
		desc := tc.description
		err := tc.scenario()
		if err != nil && !tc.expectErr {
			t.Errorf("[Case#%d] %s: expect no errors, but get %v", idx, desc, err)
		}
		if err == nil && tc.expectErr {
			t.Errorf("[Case#%d] %s: expect error occurs, but get no errors", idx, desc)
		}
	}
}

func TestRemove(t *testing.T) {
	testcases := []struct {
		description string
		scenario    func() error
		expectErr   bool
	}{
		{
			"invalid init",
			func() error {
				db := New(Options{
					Path:    "/dev/null",
					Mode:    0666,
					Bucket:  "cache",
					Options: nil,
				})
				return db.Remove(nil)
			},
			true,
		},
		{
			"valid remove",
			func() error {
				db := New(Options{
					Path:    "bolt.db",
					Mode:    0666,
					Bucket:  "cache",
					Options: nil,
				})
				defer os.Remove("bolt.db")

				k, v := []byte("123"), []byte("456")
				err := db.Put(k, v)
				if err != nil {
					return err
				}

				err = db.Remove(k)
				if err != nil {
					return err
				}

				_, err = db.Get(k)
				if err == nil {
					return fmt.Errorf("expect missing key, but get no error")
				}
				return nil
			},
			false,
		},
		{
			"cache miss",
			func() error {
				db := New(Options{
					Path:    "bolt.db",
					Mode:    0666,
					Bucket:  "cache",
					Options: nil,
				})
				defer os.Remove("bolt.db")

				k := []byte("123")
				return db.Remove(k)
			},
			false,
		},
	}

	for idx, tc := range testcases {
		desc := tc.description
		err := tc.scenario()
		if err != nil && !tc.expectErr {
			t.Errorf("[Case#%d] %s: expect no errors, but get %v", idx, desc, err)
		}
		if err == nil && tc.expectErr {
			t.Errorf("[Case#%d] %s: expect error occurs, but get no errors", idx, desc)
		}
	}
}

func TestIter(t *testing.T) {
	testcases := []struct {
		description string
		scenario    func() error
		expectErr   bool
	}{
		{
			"invalid init",
			func() error {
				db := New(Options{
					Path:    "/dev/null",
					Mode:    0666,
					Bucket:  "cache",
					Options: nil,
				})
				return db.Iter(nil)
			},
			true,
		},
		{
			"Iterate all elements",
			func() error {
				db := New(Options{
					Path:    "bolt.db",
					Mode:    0666,
					Bucket:  "cache",
					Options: nil,
				})
				defer os.Remove("bolt.db")
				list := map[string][]byte{
					"happy":    []byte("new year"),
					"oh":       []byte("my god"),
					"yakusoku": []byte("mamoru"),
				}
				for k, v := range list {
					err := db.Put([]byte(k), v)
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
				return db.Iter(iterCb)
			},
			false,
		},
	}

	for idx, tc := range testcases {
		err := tc.scenario()
		if err != nil && !tc.expectErr {
			t.Errorf("[Case#%d] %s: expect no errors, but get %v", idx, tc.description, err)
		}
		if err == nil && tc.expectErr {
			t.Errorf("[Case#%d] %s: expect error occurs, but get no error", idx, tc.description)
		}
	}
}

func TestClose(t *testing.T) {
	db := New(Options{
		Path:    "/dev/nill",
		Mode:    0666,
		Bucket:  "cache",
		Options: nil,
	})

	err := db.Close()
	if err != nil {
		t.Error(err)
	}

	db2 := New(Options{
		Path:    "bolt.db",
		Mode:    0666,
		Bucket:  "cache",
		Options: nil,
	})
	defer os.Remove("bolt.db")

	err = db2.Put([]byte("123"), []byte("456"))
	if err != nil {
		t.Error(err)
	}

	err = db2.Close()
	if err != nil {
		t.Error(err)
	}
}
