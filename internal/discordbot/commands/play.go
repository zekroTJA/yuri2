package commands

import (
	"time"

	"github.com/zekroTJA/yuri2/internal/discordbot"
	"github.com/zekroTJA/yuri2/internal/player"

	"github.com/zekroTJA/yuri2/internal/database"
	"github.com/zekroTJA/yuri2/pkg/discordgocmds"
)

// Play provides command functionalities
// for the play command
type Play struct {
	DB     database.Middleware
	Player *player.Player
}

// GetInvokes returns the invokes
// for this command.
func (c *Play) GetInvokes() []string {
	return []string{}
}

// GetDescription returns the description
// for this command
func (c *Play) GetDescription() string {
	return "Default sound player command"
}

// GetHelp returns the help text for
// this command.
func (c *Play) GetHelp() string {
	return ""
}

// GetGroup returns the group of
// the command
func (c *Play) GetGroup() string {
	return discordgocmds.GroupGeneral
}

// GetPermission returns the minimum
// required required permission level
// to execute this command.
func (c *Play) GetPermission() int {
	return 0
}

// Exec is the acual function which will
// be executed when the command was invoked.
func (c *Play) Exec(args *discordgocmds.CommandArgs) error {
	err := c.Player.Play(args.Guild, args.User, args.Args[0], player.ResourceLocal)
	if err == player.ErrNotFound {
		return nil
	}
	if err == player.ErrNotInVoice {
		msg, err := discordbot.SendEmbedError(args.Session, args.Channel.ID,
			"You need to be in a voice channel to play sounds.", "")
		msg.DeleteAfter(6 * time.Second)
		return err
	}
	return err
}
