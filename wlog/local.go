package wlog

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
	localDiscard = Log{createDiscardLogger().WithField(keyLocalMethod, "discard")}
)

// DevEnabled can make all dev logger print to ioutil.Discard
var DevEnabled = true

func init() {
	wlog, err := NewWLog(createStdoutLogger())
	if err != nil {
		panic(err)
	}
	localWLog = wlog
}

// Log are used to print devOnly Logs, all results will be print to stdout
func (m LocalWLogMethod) Log(fingerPrints ...string) (entry Log) {
	if !DevEnabled {
		return localDiscard
	}

	return Log{
		localWLog.Common(fingerPrints...).WithField(keyLocalMethod, m),
	}
}
