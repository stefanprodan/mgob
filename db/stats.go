package db

import (
	"encoding/json"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

type Status struct {
	Plan          string     `json:"plan"`
	NextRun       time.Time  `json:"next_run"`
	LastRun       *time.Time `json:"last_run,omitempty"`
	LastRunStatus string     `json:"last_run_status,omitempty"`
	LastRunLog    string     `json:"last_run_log,omitempty"`
}

type StatusStore struct {
	*Store
	bucket []byte
}

// NewStatusStore creates bucket if not found
func NewStatusStore(store *Store) (*StatusStore, error) {
	bucket := []byte("scheduler_status")

	err := store.NewBucket(bucket)
	if err != nil {
		return nil, errors.Wrap(err, "Status store bucket init failed")
	}

	return &StatusStore{store, bucket}, nil
}

// Put upserts job status
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

// Sync plans found on disk with db
func (db *StatusStore) Sync(stats []*Status) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(db.bucket)

		dbStats := make([]*Status, 0)

		// get all jobs from db
		err := b.ForEach(func(k, v []byte) error {
			var status Status
			err := json.Unmarshal(v, &status)
			if err != nil {
				return errors.Wrap(err, "Status store json unmarshal failed")
			}
			dbStats = append(dbStats, &status)
			return nil
		})
		if err != nil {
			return err
		}

		for _, newS := range stats {
			found := false
			for _, oldS := range dbStats {
				// update next run for existing job
				if newS.Plan == oldS.Plan {
					oldS.NextRun = newS.NextRun
					buf, err := json.Marshal(oldS)
					if err != nil {
						return errors.Wrapf(err, "Json marshal for %v failed", oldS.Plan)
					}
					err = b.Put([]byte(oldS.Plan), buf)
					if err != nil {
						return errors.Wrapf(err, "Updating %v to store failed", oldS.Plan)
					}
					logrus.WithField("plan", oldS.Plan).Infof("Next run at %v", oldS.NextRun)
					found = true
				}
			}

			// insert new job
			if !found {
				logrus.WithField("plan", newS.Plan).Info("New job found, saving to store")
				buf, err := json.Marshal(newS)
				if err != nil {
					return errors.Wrapf(err, "Json marshal for %v failed", newS.Plan)
				}
				err = b.Put([]byte(newS.Plan), buf)
				if err != nil {
					return errors.Wrapf(err, "Saving %v to store failed", newS.Plan)
				}
				logrus.WithField("plan", newS.Plan).Infof("Next run at %v", newS.NextRun)
			}
		}

		// remove jobs not found on disk
		for _, oldS := range dbStats {
			found := false
			for _, newS := range stats {
				if oldS.Plan == newS.Plan {
					found = true
				}
			}

			if !found {
				logrus.WithField("plan", oldS.Plan).Info("Plan not found on disk, removing from store")
				err = b.Delete([]byte(oldS.Plan))
				if err != nil {
					return errors.Wrapf(err, "Removing %v from store failed", oldS.Plan)
				}
			}
		}

		return nil

	})
}

// GetAll loads all jobs stats from db
func (db *StatusStore) GetAll() ([]*Status, error) {
	stats := make([]*Status, 0)

	err := db.View(func(tx *bolt.Tx) error {
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

	if err != nil {
		return nil, err
	}

	return stats, nil
}
