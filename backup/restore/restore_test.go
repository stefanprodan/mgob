// +build unit

package restore_test

import (
	"testing"

	"github.com/stefanprodan/mgob/backup/restore"
	"github.com/stefanprodan/mgob/config"
)

func assertError(t *testing.T, err error) {
	t.Log(err)
	if err == nil {
		t.Error(err)
	}
}

func assertNoError(t *testing.T, err error) {
	t.Log(err)
	if err != nil {
		t.Error(err)
	}
}

func TestMongoRestoreReturnErrorOnInvalidArchive(t *testing.T) {
	sched := config.Scheduler{Timeout: 60}
	plan := config.Plan{
		Target:    config.Target{},
		Scheduler: sched,
	}
	_, err := restore.Restore(plan, "invalid")
	assertError(t, err)
}

func setUp(host string, port int, collCount int) config.Plan {
	target := config.Target{
		Host: host,
		Port: port,
	}
	collections := []config.Collection{
		{
			Name:  "parameters",
			Count: collCount,
		},
	}
	restore := config.Restore{
		Database:          "garden",
		Collections:       collections,
		CollectionsLength: 1,
	}
	sched := config.Scheduler{Timeout: 60}
	plan := config.Plan{
		Target:    target,
		Scheduler: sched,
		Restore:   restore,
	}
	return plan
}

func TestMongoRestoreWithSuccess(t *testing.T) {
	collCount := 5
	plan := setUp("localhost", 27017, collCount)
	_, err := restore.Restore(plan, "/tmp/dump_test.gz")
	assertNoError(t, err)
}

func TestShouldGetErrorOnInvalidCount(t *testing.T) {
	collCount := 10
	plan := setUp("localhost", 27017, collCount)
	_, err := restore.Restore(plan, "/tmp/dump_test.gz")
	assertError(t, err)
}
