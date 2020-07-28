package cache

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	key, size := "123", int64(456)
	item := New(10, key, size)
	if item.Key != key {
		t.Errorf("expect %v, but get %v", key, item.Key)
	}
	if item.Ref != 0 {
		t.Errorf("expect %v, but get %v", 0, item.Ref)
	}
	if item.Size != size {
		t.Errorf("expect %v, but get %v", size, item.Size)
	}
	if item.Real != true {
		t.Errorf("expect %v, but get %v", true, item.Real)
	}
	if item.Used != 0 {
		t.Errorf("expect %v, but get %v", 0, item.Used)
	}
}

func TestDummy(t *testing.T) {
	key, size := "123", int64(0)
	item := Dummy(10, key)
	if item.Key != key {
		t.Errorf("expect %v, but get %v", key, item.Key)
	}
	if item.Ref != 0 {
		t.Errorf("expect %v, but get %v", 0, item.Ref)
	}
	if item.Size != size {
		t.Errorf("expect %v, but get %v", size, item.Size)
	}
	if item.Real != false {
		t.Errorf("expect %v, but get %v", false, item.Real)
	}
	if item.Used != 0 {
		t.Errorf("expect %v, but get %v", 0, item.Used)
	}
}

func TestSetSize(t *testing.T) {
	size := int64(100)
	item := Dummy(10, "123")
	item.SetSize(size)
	if item.Size != size {
		t.Errorf("expect %v, but get %v", size, item.Size)
	}
}

func TestIncrRef(t *testing.T) {
	item := Dummy(10, "123")
	item.IncrRef()
	if item.Ref != 1 {
		t.Errorf("expect %v, but get %v", 1, item.Ref)
	}
}

func TestDecrRef(t *testing.T) {
	item := Dummy(10, "123")
	item.DecrRef()
	if item.Ref != -1 {
		t.Errorf("expect %v, but get %v", -1, item.Ref)
	}
}

func TestIncrUsed(t *testing.T) {
	item := Dummy(10, "123")
	item.IncrUsed()
	if item.Used != 1 {
		t.Errorf("expect %v, but get %v", 1, item.Used)
	}
}

func TestReference(t *testing.T) {
	item := Dummy(10, "123")
	if item.Reference() != 0 {
		t.Errorf("expect %v, but get %v", 0, item.Reference())
	}
}

func TestUsedCount(t *testing.T) {
	item := Dummy(10, "123")
	if item.UsedCount() != 0 {
		t.Errorf("expect %v, but get %v", 0, item.UsedCount())
	}
}

func TestCTime(t *testing.T) {
	item := Dummy(10, "123")
	ctime := time.Time{}
	if item.CTime().Nanosecond() != ctime.Nanosecond() {
		t.Errorf("expect %v, but get %v", ctime.Nanosecond(), item.CTime().Nanosecond())
	}
}

func TestATime(t *testing.T) {
	item := Dummy(10, "123")
	ctime := time.Time{}
	if item.ATime().Nanosecond() != ctime.Nanosecond() {
		t.Errorf("expect %v, but get %v", ctime.Nanosecond(), item.ATime().Nanosecond())
	}
}

func TestUpdateCreatedAt(t *testing.T) {
	item := Dummy(10, "123")
	ctime := time.Now()
	item.UpdateCreatedAt()
	if item.CTime().Sub(ctime).Seconds() > float64(1.0) {
		t.Errorf("expect item.Ctime close to %v, but get %v", ctime, item.CTime())
	}
}

func TestUpdateLastUsed(t *testing.T) {
	item := Dummy(10, "123")
	atime := time.Now()
	item.UpdateLastUsed()
	if item.ATime().Sub(atime).Seconds() > float64(1.0) {
		t.Errorf("expect item.Atime close to %v, but get %v", atime, item.ATime())
	}
}

func TestRemove(t *testing.T) {
	item := Dummy(10, "123")
	err := item.Remove()
	if err != nil {
		t.Fatal(err)
	}

	item2 := New(789, "123", 456)
	err = item2.Remove()
	if err == nil {
		t.Errorf("expect error raise, but get no error")
	}
}

func TestIsReal(t *testing.T) {
	item := Dummy(10, "123")
	if item.IsReal() != false {
		t.Errorf("expect %v, but get %v", false, item.IsReal())
	}
}

func TestSetReal(t *testing.T) {
	item := Dummy(10, "123")
	item.SetReal()
	if item.IsReal() != true {
		t.Errorf("expect %v, but get %v", true, item.IsReal())
	}
}

func TestIsZero(t *testing.T) {
	item := Item{}
	if item.IsZero() != true {
		t.Errorf("expect %v, but get %v", true, item.IsZero())
	}

	item.SetReal()
	if item.IsZero() != false {
		t.Errorf("expect %v, but get %v", false, item.IsZero())
	}
}
