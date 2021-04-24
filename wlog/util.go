package wlog

import (
	"os"
	"io/ioutil"

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
			return entry.WithField(KeyFingerPrint, append(fp.([]string), fingerPrints...))
		}
		return entry.WithField(KeyFingerPrint, fingerPrints)
	}

	return entry.WithField(KeyMethod, fingerPrints[0]).WithField(KeyFingerPrint, fingerPrints[1:])
}