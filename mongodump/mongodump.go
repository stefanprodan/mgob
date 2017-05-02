package mongodump

import (
	"fmt"
	"github.com/codeskyblue/go-sh"
	"github.com/pkg/errors"
	"github.com/stefanprodan/mgob/config"
	"io/ioutil"
	"time"
)

func Run(plan config.Plan, conf *config.AppConfig) error {

	err := sh.Command("mkdir", "-p", fmt.Sprintf("%v/%v", conf.StoragePath, plan.Name)).Run()
	if err != nil {
		return errors.Wrapf(err, "creating dir in %v failed", conf.StoragePath)
	}

	ts := time.Now().UTC().Unix()
	archive := fmt.Sprintf("%v/%v/%v.gz", conf.StoragePath, plan.Name, ts)
	log := fmt.Sprintf("%v/%v/%v.log", conf.StoragePath, plan.Name, ts)

	cmd := fmt.Sprintf("mongodump --archive=%v --gzip --host %v --port %v --db %v ",
		archive, plan.Target.Host, plan.Target.Port, plan.Target.Database)
	if plan.Target.Username != "" && plan.Target.Password != "" {
		cmd += fmt.Sprintf("-u %v -p %v", plan.Target.Username, plan.Target.Password)
	}

	output, err := sh.Command("/bin/sh", "-c", cmd).SetTimeout(time.Duration(plan.Scheduler.Timeout) * time.Minute).CombinedOutput()
	logToFile(log, output)
	if err != nil {
		return errors.Wrapf(err, "mongodump failed, see log %v", log)
	}

	return nil
}

func logToFile(file string, data []byte) error {
	if len(data) > 0 {
		err := ioutil.WriteFile(file, data, 0644)
		if err != nil {
			return errors.Wrapf(err, "writing log %v failed", file)
		}
	}

	return nil
}
