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
	port = 27016
)

// Responsible to restore backup from one source
// using mongorestore
// Responsibilities
// - Download backup from one source
// - Restore backup using mongorestore
// - Testing restoring using queries defined by plan
func Restore(plan config.Plan, archive string) (string, error) {
	err := startMongoToRestore()
	if err != nil {
		return "", err
	}
	defer shutdownMongo()
	restore := fmt.Sprintf("mongorestore --archive=%v --gzip --host=%v --port=%d ",
		archive, host, port)
	output, err := sh.Command("/bin/sh", "-c", restore).SetTimeout(time.Duration(plan.Scheduler.Timeout) * time.Minute).CombinedOutput()
	if err != nil {
		ex := ""
		if len(output) > 0 {
			ex = strings.Replace(string(output), "\n", " ", -1)
		}
		return "", errors.Wrapf(err, "mongorestore log %v", ex)
	}
	fmt.Printf("%s\n", output)
	checkMsg, err := checkRestore(host, port, plan.Restore.Database, plan.Restore)
	if err != nil {
		return "", err
	}
	return checkMsg, nil
}

func cleanMongo(s *mgo.Session) error {
	names, err := s.DatabaseNames()
	if err != nil {
		return err
	}
	for _, name := range names {
		err = s.DB(name).DropDatabase()
		fmt.Printf("droping database  = %v\n", name)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkRestore(host string, port int, database string, r config.Restore) (string, error) {
	mongoURL := fmt.Sprintf("mongodb://%v:%d", host, port)
	session, err := mgo.Dial(mongoURL)
	if err != nil {
		return "", err
	}
	defer session.Close()
	defer cleanMongo(session)

	session.SetSafe(&mgo.Safe{})
	for _, collec := range r.Collections {
		c := session.DB(database).C(collec.Name)
		countRestored, err2 := c.Find(bson.M{}).Count()
		if err2 != nil {
			return "", err
		}

		if countRestored < collec.Count {
			return "", errors.New(
				fmt.Sprintf("Count in collection {%v}  don'n match", collec.Name))
		}
		fmt.Printf("collection = %v | count = %d | expected count = %d \n", collec.Name, countRestored, collec.Count)
	}
	collNames, err := session.DB(database).CollectionNames()
	if err != nil {
		return "", err
	}
	if len(collNames) != r.CollectionsLength {
		return "", errors.New(
			fmt.Sprintf(
				"Collection length don't match , got %d, expected %d",
				len(collNames),
				r.CollectionsLength))

	}
	return fmt.Sprintf("restore check with success\n Total collections restored = %d",
		len(collNames)), nil
}

func wrapShError(prefix string, output []byte, err error) error {
	ex := ""
	if len(output) > 0 {
		ex = strings.Replace(string(output), "\n", " ", -1)
	}
	return errors.Wrapf(err, fmt.Sprintf("%v %v", prefix, ex))

}

func startMongoToRestore() error {
	mongo := fmt.Sprintf("mongod --fork --logpath /var/log/mongodb.log --port %d", port)
	cmd := exec.Command("sh", "-c", mongo)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return wrapShError("fail to start mongo", stdoutStderr, err)
	}
	fmt.Printf("%s\n", stdoutStderr)
	return nil
}

func shutdownMongo() {
	shutdown := fmt.Sprintf("mongo --port %d --eval \"db.getSiblingDB('admin').shutdownServer()\"", port)
	cmd := exec.Command("sh", "-c", shutdown)
	stdoutStderr, err := cmd.CombinedOutput()
	fmt.Printf("%s\n", stdoutStderr)
	if err != nil {
		fmt.Print(wrapShError("fail to stop mongo", stdoutStderr, err))
	}

}
