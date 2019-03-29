package discordbot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/discordgocmds"
)

// Bot maintains the Discord bot session,
// incomming events and command executions.
type Bot struct {
	Session    *discordgo.Session
	CmdHandler *discordgocmds.CmdHandler
}

// NewBot creates a new instance of Bot.
//   token         : Discord API Bot token
//   ownerID       : the Discord ID of the bot owners account
//   generalPrefix : the general usable prefix for the bot
//   dbMiddleware  : database middleware to access database connection
func NewBot(token, ownerID, generalPrefix string, dbMiddleware discordgocmds.DatabaseMiddleware) (*Bot, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	cmdHandlerOptions := discordgocmds.CmdHandlerOptions{
		BotOwnerID:           ownerID,
		DefaultColor:         0x039BE5,
		InvokeToLower:        true,
		OwnerPermissionLevel: 5,
		ParseMsgEdit:         true,
		Prefix:               generalPrefix,
		ReactToBots:          false,
	}

	cmdHandler := discordgocmds.New(session, dbMiddleware, &cmdHandlerOptions)

	return &Bot{
		Session:    session,
		CmdHandler: cmdHandler,
	}, nil
}

// RegisterCommands registers a set of commands
// to the CmdHandler.
func (b *Bot) RegisterCommands(cmds []discordgocmds.Command) {
	for _, c := range cmds {
		b.CmdHandler.RegisterCommand(c)
	}
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
