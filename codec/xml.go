package codec

import "encoding/xml"

// XML implements
type XML struct {
	prefix string
	indent string
}

// Marshal marshals the input interface into a byte array in
// xml format.
func (x XML) Marshal(v interface{}) (b []byte, e error) {
	return xml.MarshalIndent(v, x.prefix, x.indent)
}

// Unmarshal unmarshals xml binaries into the given interface.
func (x XML) Unmarshal(b []byte, v interface{}) error {
	return xml.Unmarshal(b, v)
}
