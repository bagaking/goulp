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

func insertFingerPrintToEntry(entry *logrus.Entry, fingerPrints FingerPrints) *logrus.Entry {
	n := len(fingerPrints)
	if n == 0 {
		return entry
	}

	if v, exist := entry.Data[KeyMethod]; exist && v != nil && v != "-" {
		fp := entry.Data[KeyFingerPrint]
		return entry.WithField(KeyFingerPrint, mustCombineFingerPrint(fp, fingerPrints))
	}

	return entry.WithField(KeyMethod, fingerPrints[0]).WithField(KeyFingerPrint, fingerPrints[1:])
}

func unboxMFPFromCtx(ctx context.Context, entry *logrus.Entry) *logrus.Entry {
	if ctx == nil {
		return entry
	}

	cacheMFP := mustCombineFingerPrint(ctx.Value(CtxKeyCacheMFP), nil)
	if cacheMFP == nil {
		return entry
	}

	if n := len(cacheMFP); n != 0 {
		return insertFingerPrintToEntry(entry, cacheMFP)
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

	head := FingerPrints{methodStr}

	fingerPrint, ok := entry.Data[KeyFingerPrint]
	if !ok || fingerPrint == nil {
		return context.WithValue(ctx, CtxKeyCacheMFP, head)
	}

	fingerPrintArr, ok := fingerPrint.(FingerPrints)
	if !ok {
		return context.WithValue(ctx, CtxKeyCacheMFP, head)
	}

	return context.WithValue(ctx, CtxKeyCacheMFP, mustCombineFingerPrint(head, fingerPrintArr))
}
