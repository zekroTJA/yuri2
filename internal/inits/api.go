package inits

import (
	"github.com/zekroTJA/discordgo"
	"github.com/zekroTJA/yuri2/internal/api"
	"github.com/zekroTJA/yuri2/internal/config"
	"github.com/zekroTJA/yuri2/internal/database"
	"github.com/zekroTJA/yuri2/internal/logger"
	"github.com/zekroTJA/yuri2/internal/player"
)

// InitAPI initializes the HTTP and WS API exposure.
func InitAPI(cfg *config.API, db database.Middleware, s *discordgo.Session, player *player.Player) *api.API {
	if !cfg.Enable {
		return nil
	}

	api := api.NewAPI(cfg, db, s)

	player.AddEventHandler(api)

	logger.Info("API :: initialized")

	go func() {
		err := api.StartBlocking()
		if err != nil {
			logger.Fatal("API :: failed exposing API: %s", err.Error())
		}
		logger.Info("API :: running and exposed on address '%s'", cfg.Address)
	}()

	return api
}
