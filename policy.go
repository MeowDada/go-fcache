package fcache

// Policy is a cache replacement algorithm which able to emit a cache item.
type Policy interface {
	Emit(db DB) (Item, error)
}
