package crontask

import (
	"context"
)

type (
	Punch        func(ctx context.Context) error
	RegValidator func(ctx context.Context) error

	CtxMaker func(ctx context.Context) context.Context

	JobMeta struct {
		Key     string
		Crontab string

		// Punch is the real method to be executed
		Punch Punch

		// RegValid is used to validate the service instance it is on
		// during the registration process. If the service instance does
		// not meet certain conditions, the Job will not be registered to
		// this instance.
		RegValid RegValidator

		// CtxMaker creates the context for the timed execution of the
		// Job, if set to nil, then context.Background() is used by default
		CtxMaker CtxMaker

		JobStrategy
	}

	JobStrategy struct {
		ExecuteGroup string // 在哪个 Group 里执行, 这个字段存在时, 只在对应 Group 的 Participant Instance 里才会执行
	}
)
