package wlog

import (
	"context"

	"github.com/sirupsen/logrus"
)

var (
	KeyMethod      = "method_"
	KeyFingerPrint = "finger_print_"
)

type Log struct {
	*logrus.Entry
}

func (l Log) Dev() Log {
	return Log{
		getL().WithFields(l.Data),
	}
}

func (l Log) WithFPAppends(fingerPrints ...string) (ret *logrus.Entry) {
	fp, ok := l.Entry.Data[KeyFingerPrint]
	if !ok {
		ret = l.Entry.WithField(KeyFingerPrint, fp)
		ret.Warn("should I become a root method?")
	} else {
		ret = l.Entry.WithField(KeyFingerPrint, append(fp.([]string), fingerPrints...))
	}
	return
}

func Common(methodOrFingerPrints ...string) Log {
	return ByCtx(nil, methodOrFingerPrints...)
}

func ByCtx(ctx context.Context, methodOrFingerPrints ...string) Log {
	l := genEntry(ctx)
	if count := len(methodOrFingerPrints); count == 1 {
		l = l.WithField(KeyMethod, methodOrFingerPrints[0])
	} else if count > 1 {
		l = l.WithField(KeyFingerPrint, methodOrFingerPrints[1:])
	}
	return Log{l}
}
