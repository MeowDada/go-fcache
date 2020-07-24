package fcache

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	retry "github.com/avast/retry-go"
	"github.com/google/go-cmp/cmp"
)

func TestManager(t *testing.T) {
	db, err := BoltDB("bolt.db")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		db.Close()
		os.Remove("bolt.db")
	}()

	type pair struct {
		path string
		size int64
	}

	pairs := []pair{
		{"a1", 40},
		{"a2", 20},
		{"a3", 30},
		{"a4", 10},
		{"a5", 0},
		{"b1", 100},
		{"b2", 50},
		{"c1", 50},
		{"c2", 50},
	}
	defer func(pairs []pair) {
		for _, pair := range pairs {
			os.Remove(pair.path)
		}
	}(pairs)

	for _, pair := range pairs {
		if err := ioutil.WriteFile(pair.path, nil, 0644); err != nil {
			t.Fatal(err)
		}
	}

	mgr := New(Options{
		Capacity:     200,
		Backend:      db,
		CachePolicy:  LRU(),
		RetryOptions: nil,
	})

	for _, pair := range pairs {
		if strings.Contains(pair.path, "a") {
			mgr.Register(pair.path)
		}
	}

	for _, pair := range pairs {
		err = mgr.Set(pair.path, pair.size)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestManagerCap(t *testing.T) {
	cap := int64(1024)
	mgr := New(Options{
		Capacity:     cap,
		Backend:      Hashmap(),
		CachePolicy:  LRU(),
		RetryOptions: nil,
	})

	if mgr.Cap() != cap {
		t.Errorf("expect %d, but get %d", cap, mgr.Cap())
	}
}

func TestManagerPut(t *testing.T) {
	cap := int64(1000)
	mgr := New(Options{
		Capacity:    cap,
		Backend:     Hashmap(),
		CachePolicy: LRU(),
		RetryOptions: []retry.Option{
			retry.Attempts(10),
			retry.MaxDelay(10 * time.Millisecond),
			retry.LastErrorOnly(true),
		},
	})
	err := mgr.Set("123", cap+1)
	if err != ErrCacheTooLarge {
		t.Errorf("expect %v, but get %v", ErrCacheTooLarge, err)
	}

	err = mgr.Set("456", 100)
	if err != nil {
		t.Fatal(err)
	}

	_, err = mgr.Get("456")
	if err != nil {
		t.Fatal(err)
	}

	mgr.Register("456")
	err = mgr.Set("789", 950)
	if err == nil {
		t.Error("expect impossible to fit the cache, but get no error")
	}
}

func TestManagerPutFailed(t *testing.T) {
	mock := &mockDB{
		put: func(path string, size int64) error {
			return errMockErr
		},
	}
	cap := int64(1000)
	mgr := New(Options{
		Capacity:    cap,
		Backend:     mock,
		CachePolicy: LRU(AllowPsudo()),
		RetryOptions: []retry.Option{
			retry.Attempts(10),
			retry.MaxDelay(10 * time.Millisecond),
			retry.LastErrorOnly(true),
		},
	})
	err := mgr.Set("123", 456)
	if err != errMockErr {
		t.Fatalf("expect %v, but get %v", errMockErr, err)
	}

	mock.put = func(string, int64) error { return nil }
	mock.rm = func(string) error { return errMockErr }
	mock.iter = func(iterCb func(k string, v Item) error) error {
		return nil
	}
	err = mgr.Set("123", 800)
	if err != nil {
		t.Fatal(err)
	}
	err = mgr.Set("789", 200)
	if err != nil {
		t.Fatal(err)
	}

	err = mgr.Set("456", 900)
	if err != ErrNoEmitableCaches {
		t.Fatalf("expect %v, but get %v", ErrNoEmitableCaches, err)
	}

	mock.iter = func(iterCb func(k string, v Item) error) error {
		return iterCb("123", NewDummyItem("123"))
	}
	mock.rm = func(string) error {
		return errMockErr
	}

	err = mgr.Set("456", 900)
	if err != errMockErr {
		t.Fatalf("expect %v, but get %v", errMockErr, err)
	}
}

func TestManagerRegister(t *testing.T) {
	// Mock data.
	pairs := []struct {
		path string
		size int64
	}{
		{"tmp/a", 500},
		{"tmp/e", 600},
		{"tmp/i", 700},
		{"tmp/o", 800},
		{"tmp/u", 900},
	}

	err := os.Mkdir("tmp", 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("tmp")

	for _, p := range pairs {
		err := ioutil.WriteFile(p.path, nil, 0644)
		if err != nil {
			t.Fatal(err)
		}
	}
	defer func() {
		for _, p := range pairs {
			os.Remove(p.path)
		}
	}()

	mgr := New(Options{
		Capacity:     1000,
		Backend:      Hashmap(),
		CachePolicy:  LRU(),
		RetryOptions: nil,
	})
	for _, p := range pairs {
		mgr.Register(p.path)
		mgr.Unregister(p.path)
		err := mgr.Set(p.path, p.size)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestManagerOnce(t *testing.T) {
	target := NewItem("456", 60)
	createFn := func(fn1 func(int64) error, fn2 func(string, int64) error) (Item, error) {
		size := target.Size
		if err := fn1(size); err != nil {
			return Item{}, err
		}

		if err := fn2(target.Path, size); err != nil {
			return Item{}, err
		}
		return target, nil
	}

	mgr := New(Options{
		Capacity:    100,
		Backend:     Hashmap(),
		CachePolicy: LRU(),
		RetryOptions: []retry.Option{
			retry.LastErrorOnly(true),
		},
	})
	err := mgr.Set("123", 40)
	if err != nil {
		t.Fatal(err)
	}
	item1, err := mgr.Once("123", createFn)
	if err != nil {
		t.Fatal(err)
	}

	item2, err := mgr.Get("123")
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(item1, item2) {
		t.Errorf("expect %v, but get %v", item2, item1)
	}

	item3, err := mgr.Once("456", createFn)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(item3, target) {
		t.Errorf("expect %v, but get %v", target, item3)
	}
}
