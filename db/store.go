package db

import(
	"github.com/pkg/errors"
	"github.com/boltdb/bolt"
	"time"
)

type Store struct {
	*bolt.DB
}

// Open creates or opens a bolt db at the specified path.
func Open(path string) (*Store, error) {
	config := &bolt.Options{Timeout: 1 * time.Second}
	d, err := bolt.Open(path, 0600, config)
	if err != nil {
		return nil, errors.Wrapf(err,"Opening store %s failed", path)
	}

	return &Store{d}, nil
}

func (db *Store) NewBucket(name []byte) error {
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(name)
		if err != nil {
			return err
		}
		return nil
	})
}

func (db *Store) DeleteBucket(name []byte) error {
	return db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket(name)
	})
}
