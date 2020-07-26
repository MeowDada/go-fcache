package codec

import "encoding/json"

// JSON implements codec interface.
type JSON struct {
	prefix string
	indent string
}

// Marshal marshals the input interface into a byte array in
// json format.
func (j JSON) Marshal(v interface{}) (b []byte, e error) {
	return json.MarshalIndent(v, j.prefix, j.indent)
}

// Unmarshal unmarshals json binaries into the given interface.
func (j JSON) Unmarshal(b []byte, v interface{}) error {
	return json.Unmarshal(b, v)
}
