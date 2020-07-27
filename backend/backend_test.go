package backend

import (
	"io"
	"testing"

	"github.com/meowdada/go-fcache/cache"
)

func TestIsNoKeyError(t *testing.T) {
	testcases := []struct {
		description string
		err         error
		result      bool
	}{
		{
			"valid no key error",
			cache.ErrNoSuchKey,
			true,
		},
		{
			"nil error as input",
			nil,
			false,
		},
		{
			"invalid no key error",
			io.EOF,
			false,
		},
	}

	for idx, tc := range testcases {
		r := IsNoKeyError(tc.err)
		if r != tc.result {
			t.Errorf("[#Case%d]: %s, expect %v but get %v", idx, tc.description, tc.result, r)
		}
	}
}
