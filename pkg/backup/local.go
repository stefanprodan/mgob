package backup

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/codeskyblue/go-sh"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/stefanprodan/mgob/pkg/config"
)

func dump(plan config.Plan, tmpPath string, ts time.Time) (string, string, error) {
	archive := fmt.Sprintf("%v/%v-%v.gz", tmpPath, plan.Name, ts.Unix())
	mlog := fmt.Sprintf("%v/%v-%v.log", tmpPath, plan.Name, ts.Unix())
	dump := fmt.Sprintf("mongodump --archive=%v --gzip ", archive)

	if plan.Target.Uri != "" {
		// using uri (New in version 3.4.6)
		// host/port/username/password are incompatible with uri
		// https://docs.mongodb.com/manual/reference/program/mongodump/#cmdoption-mongodump-uri
		dump += fmt.Sprintf("--uri %v ", plan.Target.Uri)
	} else {
		// use older host/port
		dump += fmt.Sprintf("--host %v --port %v ", plan.Target.Host, plan.Target.Port)

		if plan.Target.Username != "" && plan.Target.Password != "" {
			dump += fmt.Sprintf("-u %v -p %v ", plan.Target.Username, plan.Target.Password)
		}
	}

	if plan.Target.Database != "" {
		dump += fmt.Sprintf("--db %v ", plan.Target.Database)
	}

	if plan.Target.Params != "" {
		dump += fmt.Sprintf("%v", plan.Target.Params)
	}

	// TODO: mask password
	log.Debugf("dump cmd: %v", dump)
	output, err := sh.Command("/bin/sh", "-c", dump).SetTimeout(time.Duration(plan.Scheduler.Timeout) * time.Minute).CombinedOutput()
	if err != nil {
		ex := ""
		if len(output) > 0 {
			ex = strings.Replace(string(output), "\n", " ", -1)
		}
		return "", "", errors.Wrapf(err, "mongodump log %v", ex)
	}
	logToFile(mlog, output)

	return archive, mlog, nil
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

func applyRetention(path string, retention int) error {
	gz := fmt.Sprintf("cd %v && rm -f $(ls -1t *.gz | tail -n +%v)", path, retention+1)
	err := sh.Command("/bin/sh", "-c", gz).Run()
	if err != nil {
		return errors.Wrapf(err, "removing old gz files from %v failed", path)
	}

	log.Debug("apply retention")
	log := fmt.Sprintf("cd %v && rm -f $(ls -1t *.log | tail -n +%v)", path, retention+1)
	err = sh.Command("/bin/sh", "-c", log).Run()
	if err != nil {
		return errors.Wrapf(err, "removing old log files from %v failed", path)
	}

	return nil
}

// TmpCleanup remove files older than one day
func TmpCleanup(path string) error {
	rm := fmt.Sprintf("find %v -not -name \"mgob.db\" -mtime +%v -type f -delete", path, 1)
	err := sh.Command("/bin/sh", "-c", rm).Run()
	if err != nil {
		return errors.Wrapf(err, "%v cleanup failed", path)
	}

	return nil
}
