package scheduler

import (
	"github.com/stefanprodan/mgob/config"
	"github.com/robfig/cron"
	"github.com/stefanprodan/mgob/mongodump"
	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
)

func Start(plans []config.Plan, conf *config.AppConfig) error {

	c := cron.New()

	for _, plan := range plans {

		schedule, err := cron.ParseStandard(plan.Scheduler.Cron)
		if err != nil {
			return errors.Wrapf(err, "Invalid cron %v for plan %v", plan.Scheduler.Cron, plan.Name)
		}

		c.Schedule(schedule, backupJob{plan.Name, plan, conf})
	}

	c.Start()

	for i, e := range c.Entries() {
		logrus.Infof("Plan %v next run on %v", plans[i].Name, e.Next)
	}

	return nil
}

type backupJob struct {
	name string
	plan config.Plan
	conf *config.AppConfig
}

func (b backupJob) Run()  {
	logrus.Infof("Starting job for %v", b.plan.Name)
	err := mongodump.Run(b.plan, b.conf)
	if err != nil {
		logrus.Errorf("Job %v failed %v", b.plan.Name, err)
	}else {
		logrus.Infof("Job finished for %v", b.plan.Name)
	}
}