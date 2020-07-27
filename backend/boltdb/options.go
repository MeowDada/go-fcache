package boltdb

import (
	"os"

	bolt "go.etcd.io/bbolt"
)

// Options configures boltdb instance.
type Options struct {
	Path    string
	Mode    os.FileMode
	Bucket  string
	Options *bolt.Options
}
