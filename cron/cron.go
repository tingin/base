package cron

import (
	"time"

	"github.com/go-co-op/gocron/v2"
)

type TaskFunc func() error

func Do(task TaskFunc, duration time.Duration) (gocron.Scheduler, gocron.Job, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return s, nil, err
	}

	job, err := s.NewJob(
		gocron.DurationJob(
			duration,
		),
		gocron.NewTask(
			func() { task() },
		),
		gocron.WithStartAt(
			gocron.WithStartImmediately(),
		),
	)
	s.Start()
	return s, job, nil
}
