package logger

import (
	"github.com/op/go-logging"
)

const mainLoggerName = "main"

var log = logging.MustGetLogger(mainLoggerName)

// Setup sets configuration for logger
func Setup(format string, level int) {
	formatter := logging.MustStringFormatter(format)
	logging.SetFormatter(formatter)
	logging.SetLevel(logging.Level(level), mainLoggerName)
}

// SetLogLevel sets the log level for the current logger
func SetLogLevel(logLevel int) {
	logging.SetLevel(logging.Level(logLevel), mainLoggerName)
}

// Debug prints a (formatted) debug message
func Debug(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

// Info prints a (formatted) info message
func Info(format string, args ...interface{}) {
	log.Infof(format, args...)
}

// Warning prints a (formatted) warning message
func Warning(format string, args ...interface{}) {
	log.Warningf(format, args...)
}

// Error prints a (formatted) error message
func Error(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

// Fatal prints a (formatted) fatal message
// followed by an exit call with code 1.
func Fatal(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}
