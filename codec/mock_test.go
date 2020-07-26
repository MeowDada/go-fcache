package codec

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMOCK(t *testing.T) {
	mock := Mock{
		MarshalFn:   JSON{}.Marshal,
		UnmarshalFn: JSON{}.Unmarshal,
	}
	t1 := T{100, "123", []byte("456")}
	data, err := mock.Marshal(t1)
	if err != nil {
		t.Fatal(err)
	}

	var t2 T
	err = mock.Unmarshal(data, &t2)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(t1, t2) {
		t.Errorf("expect %v, but get %v", t1, t2)
	}
}
