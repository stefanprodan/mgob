package backup

import "time"

type Result struct {
	Name      string        `json:"name"`
	Plan      string        `json:"plan"`
	Duration  time.Duration `json:"duration"`
	Size      int64         `json:"size"`
	Status    int           `json:"status"`
	Timestamp time.Time     `json:"timestamp"`
	Log       string        `json:"log"`
}
