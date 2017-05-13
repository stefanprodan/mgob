package scheduler

import (
	"sync"

	"github.com/stefanprodan/mgob/backup"
	"github.com/stefanprodan/mgob/config"
)

type Stats struct {
	sync.Mutex
	Items map[string]*backup.Result
}

func NewStats(plans []config.Plan) *Stats {
	m := make(map[string]*backup.Result)
	for _, plan := range plans {
		m[plan.Name] = &backup.Result{
			Plan: plan.Name,
		}
	}

	return &Stats{
		Items: m,
	}
}

func (s *Stats) Set(res *backup.Result) {
	s.Lock()
	s.Items[res.Plan] = res
	s.Unlock()
}

func (s *Stats) GetAll() []backup.Result {
	s.Lock()
	list := make([]backup.Result, 0)
	for _, r := range s.Items {
		list = append(list, *r)
	}
	s.Unlock()

	return list
}
