package codec

import (
	"bytes"
	"encoding/gob"
)

// Gob implements codec interface.
type Gob struct{}

// Marshal marshals the input interface into a byte array consist of
// go binaries.
func (g Gob) Marshal(v interface{}) (b []byte, e error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	e = enc.Encode(v)
	return buf.Bytes(), e
}

// Unmarshal unmarshals go binaries into the given interface.
func (g Gob) Unmarshal(b []byte, v interface{}) error {
	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)
	return dec.Decode(v)
}
