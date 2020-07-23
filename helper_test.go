package fcache

import (
	"testing"
)

func TestBytesStringConvertion(t *testing.T) {
	str := "silver wing"
	converted := bytes2str(str2bytes(str))
	if str != converted {
		t.Errorf("expect %v, but get %v", str, converted)
	}
}
