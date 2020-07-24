package fcache

type mockDB struct {
	iter    func(iterCb func(k string, v Item) error) error
	put     func(key string, size int64) error
	get     func(key string) (Item, error)
	rm      func(key string) error
	incrRef func(keys ...string) error
	decrRef func(keys ...string) error
	close   func() error
}

func (m *mockDB) Iter(iterCb func(k string, v Item) error) error { return m.iter(iterCb) }
func (m *mockDB) Put(key string, size int64) error               { return m.put(key, size) }
func (m *mockDB) Get(key string) (Item, error)                   { return m.get(key) }
func (m *mockDB) Remove(key string) error                        { return m.rm(key) }
func (m *mockDB) IncrRef(keys ...string) error                   { return m.incrRef(keys...) }
func (m *mockDB) DecrRef(keys ...string) error                   { return m.decrRef(keys...) }
func (m *mockDB) Close() error                                   { return m.close() }
