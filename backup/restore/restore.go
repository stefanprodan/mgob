package restore

import (
	"fmt"
	"os/exec"
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
	err := startMongoToRestore()
	if err != nil {
		return err
	}
	defer shutdownMongo()
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
		return errors.Wrapf(err, "mongorestore log %v", ex)
	}
	return nil
}

func wrapShError(prefix string, output []byte, err error) error {
	ex := ""
	if len(output) > 0 {
		ex = strings.Replace(string(output), "\n", " ", -1)
	}
	return errors.Wrapf(err, fmt.Sprintf("%v %v", prefix, ex))

}

func startMongoToRestore() error {
	mongo := "mongod --fork --logpath /var/log/mongodb.log"
	cmd := exec.Command("sh", "-c", mongo)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return wrapShError("fail to start mongo", stdoutStderr, err)
	}
	fmt.Printf("%s\n", stdoutStderr)
	return nil
}

func shutdownMongo() {
	shutdown := "mongo --eval \"db.getSiblingDB('admin').shutdownServer()\""
	cmd := exec.Command("sh", "-c", shutdown)
	stdoutStderr, err := cmd.CombinedOutput()
	fmt.Printf("%s\n", stdoutStderr)
	if err != nil {
		fmt.Print(wrapShError("fail to stop mongo", stdoutStderr, err))
	}

}
