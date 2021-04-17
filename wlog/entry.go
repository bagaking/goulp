package wlog

import (
	"context"

	"github.com/sirupsen/logrus"
)

var genEntry func(ctx context.Context) *logrus.Entry

func init() {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	}
	defaultEntry := logger.WithField(KeyMethod, "-")
	genEntry = func(ctx context.Context) *logrus.Entry {
		return defaultEntry
	}
}

func Set(getter func(ctx context.Context) *logrus.Entry) {
	genEntry = getter
}
