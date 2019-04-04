package inits

import (
	"github.com/zekroTJA/yuri2/internal/config"
	"github.com/zekroTJA/yuri2/internal/logger"
	"github.com/zekroTJA/yuri2/internal/player"
)

// InitPlayer initializes the sound player.
func InitPlayer(cfg *config.Lavalink) *player.Player {
	errHandler := func(t string, err error) {
		logger.Error("PLAYER :: %s :: %s", t, err.Error())
	}

	return player.NewPlayer(cfg.RESTAddress, cfg.WSAddress, cfg.Password, cfg.SoundsLocation, nil, errHandler)
}
