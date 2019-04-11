package commands

import (
	"time"

	"github.com/zekroTJA/yuri2/internal/discordbot"
	"github.com/zekroTJA/yuri2/internal/player"
	"github.com/zekroTJA/yuri2/internal/static"

	"github.com/zekroTJA/yuri2/pkg/discordgocmds"
)

// Join provides command functionalities
// for the leave command
type Join struct {
	Player  *player.Player
	PermLvl int
}

// GetInvokes returns the invokes
// for this command.
func (c *Join) GetInvokes() []string {
	return []string{"join"}
}

// GetDescription returns the description
// for this command
func (c *Join) GetDescription() string {
	return "Join your voice channel"
}

// GetHelp returns the help text for
// this command.
func (c *Join) GetHelp() string {
	return "`join` - join your voice channel"
}

// GetGroup returns the group of
// the command
func (c *Join) GetGroup() string {
	return static.CommandGroupPlayer
}

// GetPermission returns the minimum
// required required permission level
// to execute this command.
func (c *Join) GetPermission() int {
	return c.PermLvl
}

// Exec is the actual function which will
// be executed when the command was invoked.
func (c *Join) Exec(args *discordgocmds.CommandArgs) error {
	for _, vs := range args.Guild.VoiceStates {
		if vs.UserID == args.User.ID {
			return c.Player.JoinVoiceCannel(vs)
		}
	}

	msg, err := discordbot.SendEmbedError(args.Session, args.Channel.ID,
		"You need to be in a voice channel that I can join.", "")
	msg.DeleteAfter(6 * time.Second)
	return err
}
