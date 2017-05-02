package backup

import "time"

type Result struct {
	Name      string
	Plan      string
	Duration  time.Duration
	Size      int64
	Status    int
	Timestamp time.Time
}
