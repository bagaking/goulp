package wlog

import (
	"context"

	"github.com/sirupsen/logrus"
)

var (
	// CtxKeyCacheEntry is the key to cache log entry into a context
	CtxKeyCacheEntry = struct{ CtxKeyCacheEntry struct{} }{}

	// CtxKeyCacheMFP is the key to cache method and finger print into a context
	CtxKeyCacheMFP = struct{ CtxKeyCacheMFP struct{} }{}
)

type (
	// EntryMaker is the function which specified how wlog create a logger associated with a context
	EntryMaker func(ctx context.Context) *logrus.Entry

	// WLog is the sandbox of logger
	WLog struct {
		entryMaker   EntryMaker
		defaultEntry *logrus.Entry
	}
)

// SetEntryMaker can update the EntryMaker of a wlog instance
func (w *WLog) SetEntryMaker(em EntryMaker) *WLog {
	w.entryMaker = em
	return w
}

// ByCtx Create a new log entry associated with the given ctx
func (w *WLog) ByCtx(ctx context.Context, fingerPrints ...string) *logrus.Entry {
	entry := w.makeEntry(ctx)
	if entry == nil {
		return nil // todo: do nothing?
	}
	return insertFingerPrintToEntry(entry, fingerPrints)
}

// ByCtxAndCache returns the entry and cache it in the context
func (w *WLog) ByCtxAndCache(ctx context.Context, fingerPrints ...string) (*logrus.Entry, context.Context) {
	entry := w.ByCtx(ctx, fingerPrints...)
	if entry == nil {
		return nil, ctx
	}
	return entry, context.WithValue(ctx, CtxKeyCacheEntry, entry)
}

// ByCtxAndRemoveCache returns the entry and remove the cache of log entry in the context
func (w *WLog) ByCtxAndRemoveCache(ctx context.Context, fingerPrints ...string) (*logrus.Entry, context.Context) {
	entry := w.ByCtx(ctx, fingerPrints...)
	if entry == nil {
		return nil, ctx
	}
	ctxClear := context.WithValue(ctx, CtxKeyCacheEntry, nil)
	return entry, boxMFPToCtx(ctxClear, entry)
}

// Common Create a new log entry with empty context
func (w *WLog) Common(fingerPrints ...string) *logrus.Entry {
	return w.ByCtx(nil, fingerPrints...)
}

func (w *WLog) makeEntry(ctx context.Context) *logrus.Entry { // todo:
	if ctx != nil {
		if l := ctx.Value(CtxKeyCacheEntry); l != nil {
			return l.(*logrus.Entry)
		}
	}

	if w.entryMaker != nil {
		return w.entryMaker(ctx)
	}

	return unboxMFPFromCtx(ctx, w.defaultEntry)
}

// NewWLog create a new WLog instance
// argument can be EntryMaker, *logrus.Logger or nil
func NewWLog(entryMakerOrLogger interface{}) *WLog {
	if entryMakerOrLogger == nil {
		panic("invalid arguments: entryMakerOrLogger must be given")
	}

	if em, ok := entryMakerOrLogger.(EntryMaker); ok {
		return &WLog{
			entryMaker: em,
		}
	}

	if entry, ok := entryMakerOrLogger.(*logrus.Entry); ok {
		return &WLog{
			defaultEntry: entry,
		}
	}

	if logger, ok := entryMakerOrLogger.(*logrus.Logger); ok {
		return &WLog{
			defaultEntry: logger.WithField(KeyMethod, "-"),
		}
	}

	panic("invalid arguments: type error")
}
