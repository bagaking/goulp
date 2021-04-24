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

// ByCtxAndCache returns the entry and cache it in the context
func ByCtxAndCache(ctx context.Context, fingerPrints ...string) (Log, context.Context) {
	l, ctx := DefaultWLog.ByCtxAndCache(ctx, fingerPrints...)
	return Log{l}, ctx
}

// ByCtxAndRemoveCache returns the entry and remove the cache of log entry in the context
// method and finger print will be transferred to ctx, thus the mfp works in future
func ByCtxAndRemoveCache(ctx context.Context, fingerPrints ...string) (Log, context.Context) {
	l, ctx := DefaultWLog.ByCtxAndRemoveCache(ctx, fingerPrints...)
	return Log{l}, ctx
}
