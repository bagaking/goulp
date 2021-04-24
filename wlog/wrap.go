package wlog

import (
	"github.com/sirupsen/logrus"
)

var (
	// KeyMethod is the key which is used in context to sepcify the method
	KeyMethod      = "method_"
	// KeyFingerPrint is the key which is used in context to sepcify the finger prints
	KeyFingerPrint = "finger_print_"
)

// Log is a wrap of entry
type Log struct {
	*logrus.Entry
}

// Dev returns a Dev Log in which all fields of the origin Logger has been inserted
func (l Log) Dev() Log {
	return Log{
		Entry: LDev.Log().WithFields(l.Data),
	}
}

// WithFPAppends returns a new logger append fingerPrints to the origin logger
func (l Log) WithFPAppends(fingerPrints ...string) Log {
	return Log{ Entry:insertFingerPrintToEntry(l.Entry, fingerPrints) }
}
