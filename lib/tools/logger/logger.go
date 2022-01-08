package logger

import (
	"io"
	"log"
)

// const var
const (
	EmergencyLogLevel = "EMERGENCY"
	AlertLogLevel     = "ALERT"
	CriticalLogLevel  = "CRITICAL"
	ErrorLogLevel     = "ERROR"
	WarningLogLevel   = "WARNING"
	NoticeLogLevel    = "NOTICE"
	InfoLogLevel      = "INFO"
	DebugLogLevel     = "DEBUG"
)

var levels = map[string]int8{
	EmergencyLogLevel: 8,
	AlertLogLevel:     7,
	CriticalLogLevel:  6,
	ErrorLogLevel:     5,
	WarningLogLevel:   4,
	NoticeLogLevel:    3,
	InfoLogLevel:      2,
	DebugLogLevel:     1,
}

// Logger interface
type Logger interface {
	// System is unusable.
	Emergency(requestID, format string, v ...interface{})
	// Action must be taken immediately.
	//   * Example: Entire website down, database unavailable, etc.
	//     This should trigger the SMS alerts and wake you up.
	Alert(requestID, format string, v ...interface{})
	// Critical conditions.
	//   * Example: Application component unavailable, unexpected exception.
	Critical(requestID, format string, v ...interface{})
	// Runtime errors that do not require immediate action but should typically be logged and monitored.
	Error(requestID, format string, v ...interface{})
	// Exceptional occurrences that are not errors.
	//   * Example: Use of deprecated APIs, poor use of an API, undesirable things that are not necessarily wrong.
	Warning(requestID, format string, v ...interface{})
	// Normal but significant events.
	Notice(requestID, format string, v ...interface{})
	// Interesting events.
	//   * Example: User logs in, SQL logs.
	Info(requestID, format string, v ...interface{})
	// Detailed debug information.
	Debug(requestID, format string, v ...interface{})
}

type logger struct {
	logLevel string
	logger   *log.Logger
}

func (l *logger) Emergency(requestID, format string, v ...interface{}) {
	l.log(requestID, EmergencyLogLevel, format, v...)
}
func (l *logger) Alert(requestID, format string, v ...interface{}) {
	l.log(requestID, AlertLogLevel, format, v...)
}
func (l *logger) Critical(requestID, format string, v ...interface{}) {
	l.log(requestID, CriticalLogLevel, format, v...)
}
func (l *logger) Error(requestID, format string, v ...interface{}) {
	l.log(requestID, ErrorLogLevel, format, v...)
}
func (l *logger) Warning(requestID, format string, v ...interface{}) {
	l.log(requestID, WarningLogLevel, format, v...)
}
func (l *logger) Notice(requestID, format string, v ...interface{}) {
	l.log(requestID, NoticeLogLevel, format, v...)
}
func (l *logger) Info(requestID, format string, v ...interface{}) {
	l.log(requestID, InfoLogLevel, format, v...)
}
func (l *logger) Debug(requestID, format string, v ...interface{}) {
	l.log(requestID, DebugLogLevel, format, v...)
}

func (l *logger) log(requestID, level, format string, v ...interface{}) {
	if levels[l.logLevel] <= levels[level] {
		l.logger.Printf("["+requestID+"]["+level+"] "+format, v...)
	}
}

// New Public
func New(logLevel, prefix string, out io.Writer) Logger {
	desiredLogLevel := logLevel

	_, ok := levels[logLevel]
	if !ok {
		// don't error just notify the user and pick up a default
		logLevel = WarningLogLevel
	}

	l := logger{
		logLevel: logLevel,
		logger:   log.New(out, prefix, log.LstdFlags),
	}

	if !ok {
		l.Warning("log level %q out of the known domain", desiredLogLevel)
	}

	return &l
}
