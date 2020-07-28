package backend

// Mock implements store interface. It provides mock functions for testing.
type Mock struct {
	PutHandler   func(k, v []byte) error
	GetHandler   func(k []byte) ([]byte, error)
	RmHandler    func(k []byte) error
	IterHandler  func(func(k, v []byte) error) error
	CloseHandler func() error
}

// Put implements store interface.
func (m Mock) Put(k, v []byte) error { return m.PutHandler(k, v) }

// Get implements store interface.
func (m Mock) Get(k []byte) (v []byte, e error) { return m.GetHandler(k) }

// Remove implements store interface.
func (m Mock) Remove(k []byte) error { return m.RmHandler(k) }

// Iter implements store interface.
func (m Mock) Iter(iterCb func(k, v []byte) error) error { return m.IterHandler(iterCb) }

// Close implements store interface.
func (m Mock) Close() error { return m.CloseHandler() }
