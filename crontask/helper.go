package crontask

import (
	"time"

	"github.com/robfig/cron/v3"
)

func getDuration(crontab Crontab, parser cron.Parser) (time.Duration, error) {
	schedule, err := parser.Parse(string(crontab))
	if err != nil {
		return 0, err
	}

	next := schedule.Next(time.Now())
	duration := schedule.Next(next).Sub(next)
	return duration, nil
}
