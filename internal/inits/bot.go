package inits

import (
	"github.com/zekroTJA/yuri2/internal/config"
	"github.com/zekroTJA/yuri2/internal/database"
	"github.com/zekroTJA/yuri2/internal/discordbot"
	"github.com/zekroTJA/yuri2/internal/discordbot/commands"
	"github.com/zekroTJA/yuri2/internal/discordbot/handlers"
	"github.com/zekroTJA/yuri2/internal/logger"
	"github.com/zekroTJA/yuri2/internal/player"
	"github.com/zekroTJA/yuri2/internal/static"
	"github.com/zekroTJA/yuri2/pkg/discordgocmds"
)

// InitDiscordBot initializes the discord bot session and registers
// all commands and handlers.
func InitDiscordBot(cfg *config.Discord, dbMiddleware database.Middleware, player *player.Player) *discordbot.Bot {
	handlers := []interface{}{
		player.ReadyHandler,
		player.VoiceServerUpdateHandler,
		player.VoiceStateUpdateHandler,

		handlers.NewReady(cfg.StatusShuffle).Handler,
	}

	cmds := []discordgocmds.Command{
		&commands.Prefix{PermLvl: 5, DB: dbMiddleware},
		&commands.Bind{PermLvl: 0, DB: dbMiddleware, Player: player},

		&commands.List{PermLvl: 0, Player: player},
		&commands.Search{PermLvl: 0, Player: player},
		&commands.Log{PermLvl: 0, DB: dbMiddleware},
		&commands.Stats{PermLvl: 0, DB: dbMiddleware},

		&commands.Random{PermLvl: 0, Player: player},
		&commands.Stop{PermLvl: 0, Player: player},
		&commands.YouTube{PermLvl: 0, Player: player},
		&commands.Volume{PermLvl: 0, DB: dbMiddleware, Player: player},
		&commands.Join{PermLvl: 0, Player: player},
		&commands.Leave{PermLvl: 0, Player: player},
		&commands.Refetch{PermLvl: 0, Player: player},
	}

	if static.Release != "TRUE" {
		cmds = append(cmds, &commands.Test{PermLvl: 999, DB: dbMiddleware, Player: player})
	}

	bot, err := discordbot.NewBot(cfg.Token, cfg.Token,
		cfg.GeneralPrefix, dbMiddleware)
	if err != nil {
		logger.Fatal("DBOT :: failed initialization: %s", err.Error())
	}

	bot.RegisterHandler(handlers)
	bot.RegisterCommands(cmds)
	bot.CmdHandler.RegisterDefaultHandler(&commands.Play{Player: player})

	logger.Info("DBOT :: initialized")

	if err := bot.Open(); err != nil {
		logger.Fatal("DBOT :: failed connecting to the discord API: %s", err.Error())
	}

	logger.Info("DBOT :: connection established")

	return bot
}
