package ioutil

import (
	"bytes"
	"strings"
	"testing"
)

func TestStr2Bytes(t *testing.T) {
	// Case1: valid string to bytes.
	str := "bang dream"
	lhs := Str2Bytes(str)
	rhs := []byte(str)
	if bytes.Compare(lhs, rhs) != 0 {
		t.Errorf("expect %v, but get %v", rhs, lhs)
	}

	// Case2: empty string to bytes
	str = ""
	lhs = Str2Bytes(str)
	rhs = []byte(nil)
	if bytes.Compare(lhs, rhs) != 0 {
		t.Errorf("expect %v, but get %v", rhs, lhs)
	}
}

func TestBytes2Str(t *testing.T) {
	// Case1: valid bytes to string
	b := []byte{98, 97, 110, 103, 32, 100, 114, 101, 97, 109}
	lhs := Bytes2Str(b)
	rhs := string(b)
	if strings.Compare(lhs, rhs) != 0 {
		t.Errorf("expect %v, but get %v", rhs, lhs)
	}

	// Case2: emptry bytes to string
	b = nil
	lhs = Bytes2Str(b)
	rhs = string(b)
	if strings.Compare(lhs, rhs) != 0 {
		t.Errorf("expect %v, but get %v", rhs, lhs)
	}
}
