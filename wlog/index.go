package wlog

import (
	"context"
)

// DefaultWLog the default wlog instance
var DefaultWLog *WLog

func init() {
	DefaultWLog = NewWLog(createTextLogger())
}

// SetEntryGetter sets the EntryMaker of default wlog instance
func SetEntryGetter(em EntryMaker) {
	DefaultWLog.SetEntryMaker(em)
}

// Common create with given ctx and fingerprints (by default wlog instance)
func Common(fingerPrints ...string) Log {
	l := DefaultWLog.Common(fingerPrints...)
	return Log{l}
}

// ByCtx create with given ctx and fingerprints (by default wlog instance)
func ByCtx(ctx context.Context, fingerPrints ...string) Log {
	l := DefaultWLog.ByCtx(ctx, fingerPrints...)
	return Log{l}
}
