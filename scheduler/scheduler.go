package scheduler

import (
	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/robfig/cron"
	"github.com/stefanprodan/mgob/config"
	"github.com/stefanprodan/mgob/backup"
)

type Scheduler struct {
	Cron   *cron.Cron
	Plans  []config.Plan
	Config *config.AppConfig
}

func New(plans []config.Plan, conf *config.AppConfig) *Scheduler {
	s := &Scheduler{
		Cron:   cron.New(),
		Plans:  plans,
		Config: conf,
	}
	return s
}

func (s *Scheduler) Start() error {

	for _, plan := range s.Plans {

		schedule, err := cron.ParseStandard(plan.Scheduler.Cron)
		if err != nil {
			return errors.Wrapf(err, "Invalid cron %v for plan %v", plan.Scheduler.Cron, plan.Name)
		}

		s.Cron.Schedule(schedule, backupJob{plan.Name, plan, s.Config})
	}

	s.Cron.Start()

	for _, e := range s.Cron.Entries() {
		logrus.Infof("Plan %v next run on %v", e.Job.(backupJob).name, e.Next)
	}

	return nil
}

type backupJob struct {
	name string
	plan config.Plan
	conf *config.AppConfig
}

func (b backupJob) Run() {
	logrus.Infof("Starting job for %v", b.plan.Name)
	err := backup.Run(b.plan, b.conf.TmpPath, b.conf.StoragePath)
	if err != nil {
		logrus.Errorf("Job %v failed %v", b.plan.Name, err)
	} else {
		logrus.Infof("Job finished for %v", b.plan.Name)
	}
}
