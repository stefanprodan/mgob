package db

import(
	"github.com/pkg/errors"
	"github.com/boltdb/bolt"
	"time"
	"encoding/json"
)

type Status struct {
	Plan string `json:"plan"`
	NextRun time.Time `json:"next_run"`
	LastRun *time.Time `json:"last_run,omitempty"`
	LastRunStatus string `json:"last_run_status,omitempty"`
	LastRunLog string `json:"last_run_log,omitempty"`
}

type StatusStore struct {
	*Store
	bucket []byte
}

func NewStatusStore(store *Store) (*StatusStore, error)  {
	bucket := []byte("scheduler_status")

	err := store.NewBucket(bucket)
	if err != nil{
		return nil, errors.Wrap(err,"Status store bucket init failed")
	}

	return &StatusStore{store, bucket}, nil
}

func (db *StatusStore) Put(status *Status) error {

	buf, err := json.Marshal(status)
	if err != nil {
		return errors.Wrap(err, "Status store json marshal failed")
	}

	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(db.bucket)
		err = b.Put([]byte(status.Plan), buf)
		return err
	})
}

func (db *StatusStore) GetAll() ([]*Status, error){
	stats := make([]*Status, 0)

	err:= db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db.bucket))

		return b.ForEach(func(k, v []byte) error {
			var status Status
			err := json.Unmarshal(v, &status)
			if err != nil {
				return errors.Wrap(err, "Status store json unmarshal failed")
			}
			stats = append(stats, &status)
			return nil
		})
	})

	if err != nil{
		return nil, err
	}

	return  stats, nil
}
