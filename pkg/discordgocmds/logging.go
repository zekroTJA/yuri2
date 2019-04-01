package discordgocmds

import (
	"log"
	"os"
)

const (
	logPrefixDebug = "DEBUG | "
	logPrefixInfo  = "INFO  | "
	logPrefixWarn  = "WARN  | "
	logPrefixError = "ERROR | "
	logPrefixFatal = "FATAL | "
)

type logger struct {
	d *log.Logger
	i *log.Logger
	w *log.Logger
	e *log.Logger
	f *log.Logger
}

func newLogger() *logger {
	return &logger{
		d: log.New(os.Stdout, logPrefixDebug, log.LstdFlags|log.Lshortfile),
		i: log.New(os.Stdout, logPrefixInfo, log.LstdFlags),
		w: log.New(os.Stdout, logPrefixWarn, log.LstdFlags),
		e: log.New(os.Stdout, logPrefixError, log.LstdFlags),
		f: log.New(os.Stdout, logPrefixFatal, log.LstdFlags),
	}
}
