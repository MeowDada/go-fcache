package ioutil

import "unsafe"

// Str2Bytes converts a string to byte array with better performance.
func Str2Bytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

// Bytes2Str converts a byte array to a string with better performance.
func Bytes2Str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
