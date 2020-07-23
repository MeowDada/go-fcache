package fcache

import "unsafe"

func str2bytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

func bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
