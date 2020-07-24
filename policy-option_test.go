package fcache

import (
	"testing"
	"time"
)

func TestPolicyOptions(t *testing.T) {
	// Case1: AllowPsudo option
	h := Hashmap()
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
	h = Hashmap()
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
	h = Hashmap()
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
	h = Hashmap()
	h.Put("def", 100)
	time.Sleep(time.Millisecond)
	rr = RR(AllowPsudo(), MinimalLiveTime(time.Second))
	_, err = rr.Emit(h)
	if err != ErrNoEmitableCaches {
		t.Errorf("expect %v, but get %v", ErrNoEmitableCaches, err)
	}
}
