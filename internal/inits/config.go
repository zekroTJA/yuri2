package inits

import (
	"os"

	"github.com/zekroTJA/yuri2/internal/config"
	"github.com/zekroTJA/yuri2/internal/logger"
)

// InitConfig initializes the parsing or creating of the config file.
// If something throws an error, the process will exit with a
// fatal error message.
// If the config file was not existent and was created, the process
// will exit with an information text.
func InitConfig(loc string, unmarshaler config.UnmarshalFunc, marshaler config.MarshalIndentFunc, dbConfStruct interface{}) *config.Main {
	c, isNew, err := config.OpenAndParse(loc, unmarshaler, marshaler, dbConfStruct)
	if err != nil {
		logger.Fatal("CONFIG :: Failed opening or parsing: %s", err.Error())
	}
	if isNew {
		logger.Info("CONFIG :: New config file was created. " +
			"Please open the file and enter your configuration, then restart.")
		os.Exit(1)
	}

	logger.Debug("%+v\n", c)
	logger.SetLogLevel(c.Misc.LogLevel)
	logger.Info("CONFIG :: initialized")

	return c
}
