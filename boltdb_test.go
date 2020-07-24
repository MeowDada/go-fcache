package fcache

import (
	"os"
	"testing"
)

func TestBoltDB(t *testing.T) {

	// Case1: Open boltDB failed.
	_, err := BoltDB("/dev/null")
	if err == nil {
		t.Fatal("Should open db failed")
	}

	// Case2: Open a valid boltDB.
	db, err := BoltDB("bolt.db")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		db.Close()
		os.Remove("bolt.db")
	}()

	if err := db.Put("100", 100); err != nil {
		t.Error(err)
	}
	if err := db.Put("200", 200); err != nil {
		t.Error(err)
	}
	if err := db.Put("300", 300); err != nil {
		t.Error(err)
	}
	if err := db.Put("400", 400); err != nil {
		t.Error(err)
	}
	if err := db.Put("500", 500); err != nil {
		t.Error(err)
	}
	if err := db.IncrRef("100", "200", "abc", "edf"); err != nil {
		t.Error(err)
	}
	if err := db.Put("abc", 1024); err != nil {
		t.Error(err)
	}
	if err := db.DecrRef("qwg", "bcs", "100", "200", "200"); err != nil {
		t.Error(err)
	}
	if _, err := db.Get("100"); err != nil {
		t.Error(err)
	}

	err = db.(*boltDB).iter(func(k, v []byte) error {
		var item Item
		item.Parse(v)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	// Case3: Get a unexist key.
	_, err = db.Get("nop")
	if err != ErrCacheMiss {
		t.Errorf("expect %v, but get %v", ErrCacheMiss, err)
	}

	// Case4: Put a duplicate key.
	err = db.Put("kwc", 100)
	if err != nil {
		t.Fatal(err)
	}

	err = db.Put("kwc", 500)
	if err != ErrDupKey {
		t.Fatalf("expect %v, but get %v", ErrDupKey, err)
	}

}

func TestBoltDBIncrRef(t *testing.T) {
	dbPath := "bolt.db"
	db, err := BoltDB(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		db.Close()
		os.Remove(dbPath)
	}()

	err = db.Close()
	if err != nil {
		t.Fatal(err)
	}

	err = db.IncrRef("123")
	if err == nil {
		t.Error("expect non nil error, but get nil error")
	}
}

func TestBoltDBDecrRef(t *testing.T) {
	dbPath := "bolt.db"
	db, err := BoltDB(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		db.Close()
		os.Remove(dbPath)
	}()

	err = db.IncrRef("123")
	if err != nil {
		t.Fatal(err)
	}

	err = db.Close()
	if err != nil {
		t.Fatal(err)
	}

	err = db.DecrRef("123")
	if err == nil {
		t.Error("expect non nil error, but get nil error")
	}
}
