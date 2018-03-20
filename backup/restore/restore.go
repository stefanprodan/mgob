package restore

import (
	"fmt"
	"strings"
	"time"

	"github.com/codeskyblue/go-sh"
	"github.com/pkg/errors"
	"github.com/stefanprodan/mgob/config"
)

// Responsible to restore backup from one source
// using mongorestore
// Responsibilities
// - Download backup from one source
// - Restore backup using mongorestore
// - Testing restoring using queries defined by plan
func Restore(plan config.Plan, archive string) error {
	restore := fmt.Sprintf("mongorestore --archive=%v --gzip --host %v --port %v ",
		archive, plan.Target.Host, plan.Target.Port)
	if plan.Target.Database != "" {
		restore += fmt.Sprintf("--db %v ", plan.Target.Database)
	}

	output, err := sh.Command("/bin/sh", "-c", restore).SetTimeout(time.Duration(plan.Scheduler.Timeout) * time.Minute).CombinedOutput()
	if err != nil {
		ex := ""
		if len(output) > 0 {
			ex = strings.Replace(string(output), "\n", " ", -1)
		}
		return errors.Wrapf(err, "mongodump log %v", ex)
	}
	return nil
}
