package inits

import (
	"github.com/zekroTJA/yuri2/internal/database"
	"github.com/zekroTJA/yuri2/internal/logger"
)

// InitDatabase initializes the database middleware and
// and tries to connect to the database.
// If the connect failes, this function will exit the
// process with a fatal log output.
func InitDatabase(middleware database.Middleware, params ...interface{}) {
	if err := middleware.Connect(params...); err != nil {
		logger.Fatal("DATABASE :: Failed connecting: %s", err.Error())
	}
	logger.Info("DATABASE :: initialized")
}
