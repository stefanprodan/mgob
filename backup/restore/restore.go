package restore

import (
	"github.com/stefanprodan/mgob/config"
)

// Responsible to restore backup from one source
// using mongorestore
// Responsibilities
// - Download backup from one source
// - Restore backup using mongorestore
// - Testing restoring using queries defined by plan
func Restore(plan plan.Config, archive string) error {
	restore := fmt.Sprintf("mongorestore --archive=%v --gzip --host %v --port %v ",
		archive, plan.Target.Host, plan.Target.Port)
	if plan.Target.Database != "" {
		dump += fmt.Sprintf("--db %v ", plan.Target.Database)
	}

	output, err := sh.Command("/bin/sh", "-c", dump).SetTimeout(time.Duration(plan.Scheduler.Timeout) * time.Minute).CombinedOutput()
	if err != nil {
		ex := ""
		if len(output) > 0 {
			ex = strings.Replace(string(output), "\n", " ", -1)
		}
		return errors.Wrapf(err, "mongodump log %v", ex)
	}
	return nil
}
