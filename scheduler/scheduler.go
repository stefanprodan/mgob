package scheduler

import (
	"github.com/Sirupsen/logrus"
	"github.com/dustin/go-humanize"
	"github.com/pkg/errors"
	"github.com/robfig/cron"
	"github.com/stefanprodan/mgob/backup"
	"github.com/stefanprodan/mgob/config"
	"github.com/stefanprodan/mgob/metrics"
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
		logrus.Infof("Plan %v next run at %v", e.Job.(backupJob).name, e.Next)
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
	logrus.Infof("%v backup started", b.plan.Name)
	status := "200"
	t1 := time.Now()

	res, err := backup.Run(b.plan, b.conf.TmpPath, b.conf.StoragePath)
	if err != nil {
		status = "500"
		logrus.Errorf("%v backup failed %v", b.plan.Name, err)
	} else {
		logrus.Infof("%v backup finished in %v archive size %v",
			b.plan.Name, res.Duration, humanize.Bytes(uint64(res.Size)))
	}

	t2 := time.Now()
	b.metrics.Total.WithLabelValues(b.plan.Name, status).Inc()
	b.metrics.Latency.WithLabelValues(b.plan.Name, status).Observe(t2.Sub(t1).Seconds())

	b.stats.Set(&res)
}
