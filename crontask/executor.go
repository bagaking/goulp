package crontask

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/bagaking/goulp/wlog"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type (
	Crontab string

	executor struct {
		cron   *cron.Cron
		parser cron.Parser

		logger *wlog.Log

		mu   sync.Mutex
		jobs map[string]*Job
	}
)

var (
	ErrJobNotFound   = errors.New("job cannot be found")
	ErrRegDuplicated = errors.New("duplicated Job key are not allowed")

	ErrRegValidateFailed = errors.New("register validate failed")
)

func NewExecutor(parser cron.Parser, entry *logrus.Entry) *executor {
	options := []cron.Option{
		cron.WithParser(parser),
	}
	if entry == nil {
		entry = wlog.Common("cron_executor").Entry
	}
	logger := &wlog.Log{Entry: entry}
	options = append(options, cron.WithLogger(cron.PrintfLogger(logger)))
	return &executor{
		cron:   cron.New(options...),
		parser: parser,
		jobs:   make(map[string]*Job),
		logger: logger,
	}
}

// Parser - returns the copy of parser
func (e *executor) Parser() cron.Parser {
	return e.parser
}

func (e *executor) JobKeys() []string {
	keys := make([]string, len(e.jobs))
	for k := range e.jobs {
		keys = append(keys, k)
	}
	return keys
}

func (e *executor) Start() {
	e.logger.Infof("executor started, jobs= %v", e.JobKeys())
	e.cron.Start()
}

func (e *executor) Stop() context.Context {
	e.logger.Infof("executor stopped, jobs= %v", e.JobKeys())
	return e.cron.Stop()
}

func (e *executor) StopAndWait() {
	<-e.Stop().Done()
}

func (e *executor) TriggerByJobKey(ctx context.Context, jobKey string) error {
	job, ok := e.jobs[jobKey]
	if !ok {
		return fmt.Errorf("key= %s, %w", jobKey, ErrJobNotFound)
	}
	return job.Punch(ctx)
}

func (e *executor) Register(ctx context.Context, jobMeta JobMeta) error {
	job := &Job{
		executor: e,
		JobMeta:  jobMeta,
	}
	ctx = job.MakeCtx(ctx)
	logger := wlog.ByCtx(ctx, "register")
	if job.RegValid != nil {
		if err := job.RegValid(ctx); err != nil {
			return err
		}
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	if _, ok := e.jobs[job.Key]; ok {
		return ErrRegDuplicated
	}

	// todo: reg to persist Job
	e.jobs[job.Key] = job
	_, _ = e.cron.AddJob(job.Crontab, job)

	logger.Infof("register job success, job= %s", job.Key)
	return nil
}
