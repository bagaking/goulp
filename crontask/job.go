package crontask

import (
	"context"

	"github.com/bagaking/goulp/wlog"
)

type (
	Job struct {
		ID int64

		JobMeta

		executor *executor
	}
)

func (j *Job) RunWithCtx(ctx context.Context) error {
	wlog.ByCtx(ctx).Tracef("run Job %s", j.Key)
	return j.Punch(ctx)
}

func (j *Job) MakeCtx(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if j.CtxMaker != nil {
		ctx = j.CtxMaker(ctx)
	}

	// save the logEntry of executor into the context
	ctx = context.WithValue(ctx, wlog.CtxKeyCacheEntry, j.executor.logger.Entry)
	return ctx
}

func (j *Job) Run() {
	ctx := j.MakeCtx(context.Background())
	logger, ctx := wlog.ByCtxAndCache(ctx, "job", j.Key)
	if err := j.RunWithCtx(ctx); err != nil {
		logger.Warnf("run Job failed, %s", err)
	}
}
