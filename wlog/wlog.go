package wlog

import (
	"context"

	"github.com/sirupsen/logrus"
)

// CtxKey is the key where a log entry set in a context
var CtxKey = struct{}{}

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

// Common Create a new log entry with empty context
func (w *WLog) Common(fingerPrints ...string) *logrus.Entry {
	return w.ByCtx(nil, fingerPrints...)
}

func (w *WLog) makeEntry(ctx context.Context) *logrus.Entry {
	if ctx != nil {
		if l := ctx.Value(CtxKey); l != nil {
			return l.(*logrus.Entry)
		}
	}

	if w.entryMaker != nil {
		return w.entryMaker(ctx)
	}

	return w.defaultEntry
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

