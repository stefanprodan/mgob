package restore

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/codeskyblue/go-sh"
	"github.com/pkg/errors"
	"github.com/stefanprodan/mgob/config"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	host = "127.0.0.1"
	port = 27017
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
		archive, host, port)
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
	fmt.Printf("%s\n", output)
	err = checkRestore(host, port)
	if err != nil {
		return err
	}
	return nil
}

func checkRestore(host string, port int) error {
	mongoURL := fmt.Sprintf("mongodb://%v:%d", host, port)
	session, err := mgo.Dial(mongoURL)
	if err != nil {
		return err
	}
	defer session.Close()

	session.SetSafe(&mgo.Safe{})
	c := session.DB("garden").C("parameters")

	count, err2 := c.Find(bson.M{}).Count()
	if err2 != nil {
		return err
	}
	fmt.Printf("total  parameters count = %d\n", count)
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
