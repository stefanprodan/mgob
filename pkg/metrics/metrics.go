package metrics

import "github.com/prometheus/client_golang/prometheus"

type BackupMetrics struct {
	Total   *prometheus.CounterVec
	Size    *prometheus.GaugeVec
	Latency *prometheus.SummaryVec
}

func New(namespace string, subsystem string) *BackupMetrics {
	prom := &BackupMetrics{}

	prom.Total = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "backup_total",
			Help:      "The total number of backups.",
		},
		[]string{"plan", "status"},
	)

	prom.Size = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "backup_size",
			Help:      "The size of backup.",
		},
		[]string{"plan", "status"},
	)

	prom.Latency = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "backup_latency",
			Help:      "Backup duration in seconds.",
		},
		[]string{"plan", "status"},
	)

	prometheus.MustRegister(prom.Total)
	prometheus.MustRegister(prom.Size)
	prometheus.MustRegister(prom.Latency)

	return prom
}
