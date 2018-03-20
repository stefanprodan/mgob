package restore_test

import (
	"github.com/stefanprodan/mgob/backup/restore"
	"github.com/stefanprodan/mgob/config"
	"testing"
)

func assertError(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func TestMungoRestoreReturnErrorOnInvalidArchive(t *testing.T) {
	target := config.Target{
		Host: "localost",
		Port: 27017,
	}
	sched := config.Scheduler{Timeout: 60}
	plan := config.Plan{
		Target:    target,
		Scheduler: sched,
	}
	err := restore.Restore(plan, "invalid")
	assertError(t, err)
}
