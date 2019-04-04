package inits

import (
	"github.com/zekroTJA/yuri2/internal/database"
	"github.com/zekroTJA/yuri2/internal/discordbot"
	"github.com/zekroTJA/yuri2/internal/discordbot/commands"
	"github.com/zekroTJA/yuri2/internal/logger"
	"github.com/zekroTJA/yuri2/internal/player"
	"github.com/zekroTJA/yuri2/pkg/discordgocmds"
)

// InitDiscordBot initializes the discord bot session and registers
// all commands and handlers.
func InitDiscordBot(token, ownerID, generalPrefix string, dbMiddleware database.Middleware, player *player.Player) *discordbot.Bot {
	handlers := []interface{}{
		player.ReadyHandler,
		player.VoiceServerUpdateHandler,
		player.VoiceStateUpdateHandler,
	}

	cmds := []discordgocmds.Command{
		&commands.Prefix{PermLvl: 5, DB: dbMiddleware},
		&commands.Test{PermLvl: 999, DB: dbMiddleware, Player: player},
	}

	bot, err := discordbot.NewBot(token, ownerID, generalPrefix, dbMiddleware)
	if err != nil {
		logger.Fatal("DBOT :: failed initialization: %s", err.Error())
	}

	bot.RegisterHandler(handlers)
	bot.RegisterCommands(cmds)
	bot.CmdHandler.RegisterDefaultHandler(&commands.Play{DB: dbMiddleware})

	logger.Info("DBOT :: initialized")

	if err := bot.Open(); err != nil {
		logger.Fatal("DBOT :: failed connecting to the discord API: %s", err.Error())
	}

	logger.Info("DBOT :: connection established")

	return bot
}
