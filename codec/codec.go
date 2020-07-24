package codec

// Codec is a encoder and decoder which is able
// to encode an interface into binaries, and decode
// byte array into an interface.
type Codec interface {
	// Marshal marshals the input interface into a byte array.
	Marshal(v interface{}) (b []byte, e error)

	// Unmarshal unmarshals a byte array into the given interface.
	Unmarshal(b []byte, v interface{}) error
}
