package policy

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/meowdada/go-fcache/backend"
	"github.com/meowdada/go-fcache/backend/gomap"
	"github.com/meowdada/go-fcache/codec"
)

func TestPolicyOptions(t *testing.T) {
	// Case1: AllowPsudo option
	h := backend.Adapter(gomap.New(), codec.Gob{})
	err := h.IncrRef("psudo")
	if err != nil {
		t.Fatal(err)
	}

	err = h.DecrRef("psudo")
	if err != nil {
		t.Fatal(err)
	}

	rr := RR(AllowPsudo())
	_, err = rr.Emit(h)
	if err != nil {
		t.Fatal(err)
	}

	// Case2: AllowReferenced option
	h = backend.Adapter(gomap.New(), codec.Gob{})
	err = h.IncrRef("refed")
	if err != nil {
		t.Fatal(err)
	}

	err = h.IncrRef("refed")
	if err != nil {
		t.Fatal(err)
	}

	rr = RR(AllowPsudo(), AllowReferenced())
	_, err = rr.Emit(h)
	if err != nil {
		t.Fatal(err)
	}

	// Case3: MinimalUsed Option
	h = backend.Adapter(gomap.New(), codec.Gob{})
	h.IncrRef("abc")
	h.IncrRef("abc")
	h.DecrRef("abc")
	h.DecrRef("abc")
	rr = RR(AllowPsudo(), MinimalUsed(3))
	_, err = rr.Emit(h)
	if err != ErrNoEmitableCaches {
		t.Errorf("execpt %v, but get %v", ErrNoEmitableCaches, err)
	}

	// Case4: MinimalLivedTime option
	h = backend.Adapter(gomap.New(), codec.Gob{})
	h.Put("def", 100)
	time.Sleep(time.Millisecond)
	rr = RR(AllowPsudo(), MinimalLiveTime(time.Second))
	_, err = rr.Emit(h)
	if err != ErrNoEmitableCaches {
		t.Errorf("expect %v, but get %v", ErrNoEmitableCaches, err)
	}
}

func TestPolicyAllowReferenced(t *testing.T) {
	db := backend.Adapter(gomap.New(), codec.Gob{})

	err := ioutil.WriteFile("test123", nil, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("test123")

	err = db.Put("test123", 100)
	if err != nil {
		t.Fatal(err)
	}

	err = db.IncrRef("test123")
	if err != nil {
		t.Fatal(err)
	}

	rr := RR()
	v, err := rr.Emit(db)
	if err == nil {
		t.Fatalf("expect evict no cache item, but get %v", v)
	}
}
