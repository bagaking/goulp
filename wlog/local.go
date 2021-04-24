package wlog

import (
	"sync"
)

// LocalWLogMethod is used to specify the kinds of logger
type LocalWLogMethod string

const (
	keyLocalMethod = "local_"

	// LDev can used to print debug messages
	LDev LocalWLogMethod = "dev"
	// LInit can used to print init messages
	LInit LocalWLogMethod = "init"
	// LExit can used to print exit messages
	LExit LocalWLogMethod = "exit"
)

var (
	localWLog    *WLog
	onceInitWLog = sync.Once{}

	localDiscard = Log{createDiscardLogger().WithField(keyLocalMethod, "discard")}
)

// DevEnabled can make all dev logger print to ioutil.Discard
var DevEnabled = true

// Log are used to print devOnly Logs, all results will be print to stdout
func (m LocalWLogMethod) Log(fingerPrints ...string) (entry Log) {
	if !DevEnabled {
		return localDiscard
	}

	if localWLog == nil {
		onceInitWLog.Do(func() {
			localWLog = NewWLog(createStdoutLogger())
		})
	}

	return Log{
		localWLog.Common(fingerPrints...).WithField(keyLocalMethod, m),
	}
}
