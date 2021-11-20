package discordbot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/yuri2/internal/logger"
	"github.com/zekroTJA/yuri2/pkg/discordgocmds"
	"github.com/zekrotja/ken"
	"github.com/zekrotja/ken/state"
	"github.com/zekrotja/ken/store"
)

// Bot maintains the Discord bot session,
// incomming events and command executions.
type Bot struct {
	Session    *discordgo.Session
	CmdHandler *ken.Ken
}

// NewBot creates a new instance of Bot.
//   token         : Discord API Bot token
//   ownerID       : the Discord ID of the bot owners account
//   generalPrefix : the general usable prefix for the bot
//   dbMiddleware  : database middleware to access database connection
func NewBot(token, ownerID, generalPrefix string, dbMiddleware discordgocmds.DatabaseMiddleware) (b *Bot, err error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return
	}

	session.StateEnabled = true
	session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAllWithoutPrivileged)

	cmdHandler, err := ken.New(session, ken.Options{
		CommandStore:   store.NewDefault(),
		State:          state.NewInternal(),
		OnSystemError:  systemErrorHandler,
		OnCommandError: commandErrorHandler,
	})
	if err != nil {
		return
	}

	b = &Bot{
		Session:    session,
		CmdHandler: cmdHandler,
	}
	return
}

// RegisterCommands registers a set of commands
// to the CmdHandler.
func (b *Bot) RegisterCommands(cmds ...ken.Command) {
	b.CmdHandler.RegisterCommands(cmds...)
}

// RegisterHandler registers a set of event handlers
// to the discordgo Session.
func (b *Bot) RegisterHandler(handler []interface{}) {
	for _, h := range handler {
		b.Session.AddHandler(h)
	}
}

// Open initiates the Bot's connection  to
// the Discord API.
func (b *Bot) Open() error {
	return b.Session.Open()
}

// Close cleanly closes the connection to the Discord
// Web Socket so that the WS does not need to wait
// up to 45 seconds until timeout.
func (b *Bot) Close() {
	if b.Session != nil {
		b.Session.Close()
	}
}

func systemErrorHandler(context string, err error, args ...interface{}) {
	logger.Error("Ken Error [%s]: %s", context, err.Error())
}

func commandErrorHandler(err error, ctx *ken.Ctx) {
	// Is ignored if interaction has already been responded
	ctx.Defer()

	if err == ken.ErrNotDMCapable {
		ctx.FollowUpError("This command can not be used in DMs.", "")
		return
	}

	ctx.FollowUpError(
		fmt.Sprintf("The command execution failed unexpectedly:\n```\n%s\n```", err.Error()),
		"Command execution failed")
}
