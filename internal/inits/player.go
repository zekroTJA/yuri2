package inits

import (
	"github.com/zekroTJA/yuri2/internal/config"
	"github.com/zekroTJA/yuri2/internal/database"
	"github.com/zekroTJA/yuri2/internal/logger"
	"github.com/zekroTJA/yuri2/internal/player"
)

// InitPlayer initializes the sound player.
func InitPlayer(cfg *config.Lavalink, db database.Middleware) *player.Player {
	errHandler := func(t string, err error) {
		if err == nil {
			return
		}
		logger.Error("PLAYER :: %s :: %s", t, err.Error())
	}

	return player.NewPlayer("http://"+cfg.Address, "ws://"+cfg.Address, cfg.Password, cfg.SoundsLocation, db, errHandler)
}
