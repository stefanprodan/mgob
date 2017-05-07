package scheduler

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/dustin/go-humanize"
	"github.com/pkg/errors"
	"github.com/robfig/cron"
	"github.com/stefanprodan/mgob/backup"
	"github.com/stefanprodan/mgob/config"
	"github.com/stefanprodan/mgob/metrics"
	"github.com/stefanprodan/mgob/notifier"
	"time"
)

type Scheduler struct {
	Cron    *cron.Cron
	Plans   []config.Plan
	Config  *config.AppConfig
	Stats   *Stats
	metrics *metrics.BackupMetrics
}

func New(plans []config.Plan, conf *config.AppConfig, stats *Stats) *Scheduler {
	s := &Scheduler{
		Cron:    cron.New(),
		Plans:   plans,
		Config:  conf,
		Stats:   stats,
		metrics: metrics.New("mgob", "scheduler"),
	}

	return s
}

func (s *Scheduler) Start() error {
	for _, plan := range s.Plans {
		schedule, err := cron.ParseStandard(plan.Scheduler.Cron)
		if err != nil {
			return errors.Wrapf(err, "Invalid cron %v for plan %v", plan.Scheduler.Cron, plan.Name)
		}
		s.Cron.Schedule(schedule, backupJob{plan.Name, plan, s.Config, s.Stats, s.metrics})
	}
	s.Cron.Start()

	for _, e := range s.Cron.Entries() {
		logrus.WithField("plan", e.Job.(backupJob).name).Infof("Next run at %v", e.Next)
	}
	return nil
}

type backupJob struct {
	name    string
	plan    config.Plan
	conf    *config.AppConfig
	stats   *Stats
	metrics *metrics.BackupMetrics
}

func (b backupJob) Run() {
	logrus.WithField("plan", b.plan.Name).Info("Backup started")
	status := "200"
	t1 := time.Now()

	res, err := backup.Run(b.plan, b.conf.TmpPath, b.conf.StoragePath)
	if err != nil {
		status = "500"
		logrus.WithField("plan", b.plan.Name).Errorf("Backup failed %v", err)
		if err := notifier.SendNotification(fmt.Sprintf("%v backup failed", b.plan.Name),
			err.Error(), true, b.plan); err != nil {
			logrus.WithField("plan", b.plan.Name).Errorf("Notifier failed %v", err)
		}
	} else {
		logrus.WithField("plan", b.plan.Name).Infof("Backup finished in %v archive %v size %v",
			res.Duration, res.Name, humanize.Bytes(uint64(res.Size)))
		if err := notifier.SendNotification(fmt.Sprintf("%v backup finished", b.plan.Name),
			fmt.Sprintf("%v backup finished in %v archive size %v",
				res.Name, res.Duration, humanize.Bytes(uint64(res.Size))),
			false, b.plan); err != nil {
			logrus.WithField("plan", b.plan.Name).Errorf("Notifier failed %v", err)
		}
	}

	t2 := time.Now()
	b.metrics.Total.WithLabelValues(b.plan.Name, status).Inc()
	b.metrics.Latency.WithLabelValues(b.plan.Name, status).Observe(t2.Sub(t1).Seconds())

	b.stats.Set(&res)
}
