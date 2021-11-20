package inits

import (
	"github.com/zekroTJA/yuri2/internal/config"
	"github.com/zekroTJA/yuri2/internal/database"
	"github.com/zekroTJA/yuri2/internal/discordbot"
	"github.com/zekroTJA/yuri2/internal/discordbot/handlers"
	"github.com/zekroTJA/yuri2/internal/discordbot/slashcommands"
	"github.com/zekroTJA/yuri2/internal/logger"
	"github.com/zekroTJA/yuri2/internal/player"
	"github.com/zekroTJA/yuri2/internal/static"
)

// InitDiscordBot initializes the discord bot session and registers
// all commands and handlers.
func InitDiscordBot(cfg *config.Discord, db database.Middleware, player *player.Player) *discordbot.Bot {
	handlers := []interface{}{
		player.ReadyHandler,
		player.VoiceServerUpdateHandler,
		player.VoiceStateUpdateHandler,

		handlers.NewReady(cfg.StatusShuffle).Handler,
	}

	if static.Release != "TRUE" {
		// cmds = append(cmds, &commands.Test{PermLvl: 999, DB: dbMiddleware, Player: player})
	}

	bot, err := discordbot.NewBot(cfg.Token, cfg.Token,
		cfg.GeneralPrefix, db)
	if err != nil {
		logger.Fatal("DBOT :: failed initialization: %s", err.Error())
	}

	bot.RegisterHandler(handlers)
	bot.RegisterCommands(
		&slashcommands.Play{player},
		&slashcommands.Bind{player, db},
	)

	logger.Info("DBOT :: initialized")

	if err := bot.Open(); err != nil {
		logger.Fatal("DBOT :: failed connecting to the discord API: %s", err.Error())
	}

	logger.Info("DBOT :: connection established")

	return bot
}
