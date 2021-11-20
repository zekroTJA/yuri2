package discordgocmds

import "github.com/bwmarrin/discordgo"

// CommandArgs will be passed to the
// command Exec function and contains the
// Channel, User, Guild, Message, Session
// and Command Handler Object pointers and
// the list of command arguments
type CommandArgs struct {
	Channel    *discordgo.Channel
	User       *discordgo.User
	Guild      *discordgo.Guild
	Message    *discordgo.Message
	Args       []string
	Session    *discordgo.Session
	CmdHandler *CmdHandler
}
