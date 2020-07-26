package codec

// Mock implements codec interface.
type Mock struct {
	MarshalFn   func(v interface{}) ([]byte, error)
	UnmarshalFn func(b []byte, v interface{}) error
}

// Marshal marshals given interface into binaries.
func (m Mock) Marshal(v interface{}) ([]byte, error) {
	return m.MarshalFn(v)
}

// Unmarshal unmarshals given binaries into interface.
func (m Mock) Unmarshal(b []byte, v interface{}) error {
	return m.UnmarshalFn(b, v)
}
