package codec

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type T struct {
	A int
	B string
	C []byte
}

func TestGob(t *testing.T) {
	gob := Gob{}
	t1 := T{100, "123", []byte("456")}
	data, err := gob.Marshal(t1)
	if err != nil {
		t.Fatal(err)
	}

	var t2 T
	err = gob.Unmarshal(data, &t2)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(t1, t2) {
		t.Errorf("expect %v, but get %v", t1, t2)
	}
}
