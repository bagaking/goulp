package wlog

import (
	"os"

	"github.com/sirupsen/logrus"
)

var KeyLocalMethod = "local_"

var DevEnabled = false

var localLogger *logrus.Entry

func init() {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		ForceColors:   true,
		DisableColors: false,
		FullTimestamp: true,
	}
	localLogger = logger.WithField(KeyLocalMethod, "custom")
	localLogger.Logger.SetOutput(os.Stdout)
}

func getL() (entry *logrus.Entry) {
	return localLogger
}

func LPrint(localMethod string, fmtOrMsg string, args ...interface{}) {
	l := getL().WithField(KeyLocalMethod, localMethod)
	if len(args) > 0 {
		l.Infof(fmtOrMsg, args...)
		return
	}
	l.Info(fmtOrMsg)
}

func LDev(fmtOrMsg string, args ...interface{}) {
	if !DevEnabled {
		return
	}
	LPrint("dev", fmtOrMsg, args...)
}

func LInit(fmtOrMsg string, args ...interface{}) {
	LPrint("init", fmtOrMsg, args...)
}

func LExit(fmtOrMsg string, args ...interface{}) {
	LPrint("exit", fmtOrMsg, args...)
}
