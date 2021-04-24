package wlog

import (
	"context"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

func createTextLogger() *logrus.Logger {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	}
	return logger
}

func createStdoutLogger() *logrus.Logger {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		ForceColors:   true,
		DisableColors: false,
		FullTimestamp: true,
	}
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.TraceLevel)
	return logger
}

func createStderrLogger() *logrus.Logger {
	logger := createStdoutLogger()
	logger.SetOutput(os.Stderr)
	return logger
}

func createDiscardLogger() *logrus.Logger {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		ForceColors:   false,
		DisableColors: true,
		FullTimestamp: false,
	}
	logger.SetOutput(ioutil.Discard)
	return logger
}

func insertFingerPrintToEntry(entry *logrus.Entry, fingerPrints []string) *logrus.Entry {
	n := len(fingerPrints)
	if n == 0 {
		return entry
	}

	if v, exist := entry.Data[KeyMethod]; exist && v != nil && v != "-" {
		if fp, ok := entry.Data[KeyFingerPrint]; ok && fp != nil {
			fparr, ok := fp.([]string)
			if !ok {
				fparr = make([]string, 0, n)
			} else {
				// todo: be thread-safe, but lead lots of copy
				// a thread-safe pool can be considered here
				fparr = append(make([]string, 0, len(fparr)+n), fparr...)
			}
			return entry.WithField(KeyFingerPrint, append(fparr, fingerPrints...))
		}
		return entry.WithField(KeyFingerPrint, fingerPrints)
	}

	return entry.WithField(KeyMethod, fingerPrints[0]).WithField(KeyFingerPrint, fingerPrints[1:])
}

func unboxMFPfromCtx(ctx context.Context, entry *logrus.Entry) *logrus.Entry {
	if ctx == nil {
		return entry
	}
	cacheMFP := ctx.Value(CtxKeyCacheMFP)
	if cacheMFP == nil {
		return entry
	}

	cachMFPArr, ok := cacheMFP.([]string)
	if !ok {
		return entry
	}

	if n := len(cachMFPArr); n != 0 {
		return insertFingerPrintToEntry(entry, cachMFPArr)
	}

	return entry
}

func boxMFPToCtx(ctx context.Context, entry *logrus.Entry) context.Context {
	method, ok := entry.Data[KeyMethod]
	if !ok || method == nil || method == "-" {
		return ctx
	}

	methodStr, ok := method.(string)
	if !ok {
		return ctx
	}

	fingerPrint, ok := entry.Data[KeyFingerPrint]
	if !ok || fingerPrint == nil {
		return context.WithValue(ctx, CtxKeyCacheMFP, []string{methodStr})
	}

	fingerPrintArr, ok := fingerPrint.([]string)
	if !ok {
		return context.WithValue(ctx, CtxKeyCacheMFP, []string{methodStr})
	}

	cache := append(append(make([]string, 0, len(fingerPrintArr)+1), methodStr), fingerPrintArr...)
	return context.WithValue(ctx, CtxKeyCacheMFP, cache)
}
