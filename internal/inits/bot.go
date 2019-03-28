package inits

import (
	"github.com/zekroTJA/discordgocmds"
	"github.com/zekroTJA/yuri2/internal/discordbot"
)

// InitDiscordBot initializes the discord bot session and registers
// all commands and handlers.
func InitDiscordBot(token, ownerID, generalPrefix string, dbMiddleware discordgocmds.DatabaseMiddleware) *discordbot.Bot {
	handlers := []interface{}{}

	commands := []discordgocmds.Command{}

	bot := discordbot.NewBot(token, ownerID, generalPrefix, dbMiddleware)
	bot.RegisterHandler(handlers)
	bot.RegisterCommands(commands)

	return bot
}
