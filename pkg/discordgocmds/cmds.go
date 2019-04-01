package discordgocmds

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// CmdHandler is the main controller of the command
// handler and contains the bot session,  the
// registered commands and registers the message
// handler checking messages for commands.
type CmdHandler struct {
	discordSession         *discordgo.Session
	options                *CmdHandlerOptions
	permHandler            PermissionHandler
	databaseMiddleware     DatabaseMiddleware
	defaultHandler         Command
	registeredCmds         map[string]Command
	registeredCmdInstances []Command
	logger                 *logger
}

// New creates a new instance of CmdHandler by passing
// the discordgo session, the database middleware instance
// and the command handler options as argument.
func New(session *discordgo.Session, dbMiddleware DatabaseMiddleware, options *CmdHandlerOptions) *CmdHandler {
	c := &CmdHandler{
		discordSession:         session,
		options:                options,
		databaseMiddleware:     dbMiddleware,
		registeredCmds:         make(map[string]Command),
		registeredCmdInstances: make([]Command, 0),
		logger:                 newLogger(),
	}

	c.discordSession.AddHandler(c.messageHandler)
	c.discordSession.AddHandler(c.readyHandler)
	c.RegisterCommand(new(CmdHelp))

	return c
}

// RegisterCommand registers a Command class in the
// command handler and will be available for execution.
func (c *CmdHandler) RegisterCommand(cmd Command) {
	c.registeredCmdInstances = append(c.registeredCmdInstances, cmd)
	for _, invoke := range cmd.GetInvokes() {
		c.registeredCmds[invoke] = cmd
	}
}

// RegisterDefaultHandler registers a command handler
// whill be fired when no other command invoke matches.
func (c *CmdHandler) RegisterDefaultHandler(handler Command) {
	c.defaultHandler = handler
}

//////// private functions ////////

func (c *CmdHandler) sendEmbedError(chanID, body, title string) (*discordgo.Message, error) {
	emb := &discordgo.MessageEmbed{
		Color:       cErrorColor,
		Description: body,
		Title:       title,
	}
	return c.discordSession.ChannelMessageSendEmbed(chanID, emb)
}

//////// discordgo event handlers ////////

func (c *CmdHandler) messageHandler(s *discordgo.Session, e *discordgo.MessageCreate) {
	if e.Message.Author.ID == s.State.User.ID {
		return
	}
	if !c.options.ReactToBots && e.Message.Author.Bot {
		return
	}
	channel, err := s.Channel(e.ChannelID)
	if err != nil {
		c.logger.e.Printf("Failed getting discord channel from ID (%s): %s", e.ChannelID, err.Error())
		return
	}
	if channel.Type != discordgo.ChannelTypeGuildText {
		return
	}
	guildPrefix, err := c.databaseMiddleware.GetGuildPrefix(e.GuildID)
	if err != nil {
		c.logger.e.Printf("Failed fetching guild prefix from database: %s", err.Error())
	}

	var pre string
	if strings.HasPrefix(e.Message.Content, c.options.Prefix) {
		pre = c.options.Prefix
	} else if guildPrefix != "" && strings.HasPrefix(e.Message.Content, guildPrefix) {
		pre = guildPrefix
	} else {
		return
	}

	contSplit := strings.Fields(e.Message.Content)
	invoke := contSplit[0][len(pre):]
	if c.options.InvokeToLower {
		invoke = strings.ToLower(invoke)
	}

	cmdInstance, ok := c.registeredCmds[invoke]

	guild, _ := s.Guild(e.GuildID)
	cmdArgs := &CommandArgs{
		Args:       contSplit[1:],
		Channel:    channel,
		CmdHandler: c,
		Guild:      guild,
		Message:    e.Message,
		Session:    s,
		User:       e.Author,
	}

	if c.options.DeleteCmdMessages {
		s.ChannelMessageDelete(e.ChannelID, e.ID)
	}

	if !ok {
		cmdInstance = c.defaultHandler
		args := make([]string, len(cmdArgs.Args)+1)
		args[0] = invoke
		for i, arg := range cmdArgs.Args {
			args[i+1] = arg
		}
		cmdArgs.Args = args
	}

	hasPerm, err := c.permHandler.CheckUserPermission(cmdArgs, s, cmdInstance)
	if err != nil {
		c.sendEmbedError(channel.ID, fmt.Sprintf("Failed getting permission von database: ```\n%s\n```", err.Error()), "Permission Error")
		return
	}
	if !hasPerm {
		c.sendEmbedError(channel.ID, "You are not permitted to use this command!", "Missing permission")
		return
	}
	err = cmdInstance.Exec(cmdArgs)
	if err != nil {
		c.sendEmbedError(channel.ID, fmt.Sprintf("Failed executing command: ```\n%s\n```", err.Error()), "Command execution failed")
	}

}

func (c *CmdHandler) readyHandler(s *discordgo.Session, e *discordgo.Ready) {
	if c.databaseMiddleware == nil {
		panic("Database middleware must be registered")
	}
	if c.permHandler == nil {
		c.permHandler = NewDefaultPermissionHandler(c.databaseMiddleware)
	}
}
