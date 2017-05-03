package backup

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/codeskyblue/go-sh"
	"github.com/pkg/errors"
	"github.com/stefanprodan/mgob/config"
	"os"
	"path/filepath"
	"time"
)

func Run(plan config.Plan, tmpPath string, storagePath string) (Result, error) {
	t1 := time.Now()
	planDir := fmt.Sprintf("%v/%v", storagePath, plan.Name)

	archive, log, err := dump(plan, tmpPath, t1.UTC())
	res := Result{
		Plan:      plan.Name,
		Timestamp: t1.UTC(),
		Status:    500,
	}
	_, res.Name = filepath.Split(archive)

	if err != nil {
		return res, err
	}

	err = sh.Command("mkdir", "-p", planDir).Run()
	if err != nil {
		return res, errors.Wrapf(err, "creating dir %v in %v failed", plan.Name, storagePath)
	}

	fi, err := os.Stat(archive)
	if err != nil {
		return res, errors.Wrapf(err, "stat file %v failed", archive)
	}
	res.Size = fi.Size()

	err = sh.Command("mv", archive, planDir).Run()
	if err != nil {
		return res, errors.Wrapf(err, "moving file from %v to %v failed", archive, planDir)
	}

	err = sh.Command("mv", log, planDir).Run()
	if err != nil {
		return res, errors.Wrapf(err, "moving file from %v to %v failed", log, planDir)
	}

	if plan.Scheduler.Retention > 0 {
		err = applyRetention(planDir, plan.Scheduler.Retention)
		if err != nil {
			return res, errors.Wrap(err, "retention job failed")
		}
	}

	if plan.S3 != nil {
		file := filepath.Join(planDir, res.Name)
		s3output, err := s3Upload(file, plan)
		if err != nil {
			return res, err
		} else {
			logrus.Infof("S3 upload finished %v", s3output)
		}
	}

	t2 := time.Now()
	res.Status = 200
	res.Duration = t2.Sub(t1)
	return res, nil
}
