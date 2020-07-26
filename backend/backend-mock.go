package backend

// mock implements store interface.
type mock struct {
	put   func(k, v []byte) error
	get   func(k []byte) ([]byte, error)
	rm    func(k []byte) error
	iter  func(func(k, v []byte) error) error
	close func() error
}

func (m mock) Put(k, v []byte) error                     { return m.put(k, v) }
func (m mock) Get(k []byte) (v []byte, e error)          { return m.get(k) }
func (m mock) Remove(k []byte) error                     { return m.rm(k) }
func (m mock) Iter(iterCb func(k, v []byte) error) error { return m.iter(iterCb) }
func (m mock) Close() error                              { return m.close() }
