package commands

import (
	"time"

	"github.com/zekroTJA/yuri2/internal/discordbot"
	"github.com/zekroTJA/yuri2/internal/player"
	"github.com/zekroTJA/yuri2/internal/static"

	"github.com/zekroTJA/yuri2/internal/database"
	"github.com/zekroTJA/yuri2/pkg/discordgocmds"
)

// Random provides command functionalities
// for the random command
type Random struct {
	DB      database.Middleware
	Player  *player.Player
	PermLvl int
}

// GetInvokes returns the invokes
// for this command.
func (c *Random) GetInvokes() []string {
	return []string{"r"}
}

// GetDescription returns the description
// for this command
func (c *Random) GetDescription() string {
	return "Play a random, local file"
}

// GetHelp returns the help text for
// this command.
func (c *Random) GetHelp() string {
	return "`r` - play a random, local file"
}

// GetGroup returns the group of
// the command
func (c *Random) GetGroup() string {
	return static.CommandGroupPlayer
}

// GetPermission returns the minimum
// required required permission level
// to execute this command.
func (c *Random) GetPermission() int {
	return c.PermLvl
}

// Exec is the acual function which will
// be executed when the command was invoked.
func (c *Random) Exec(args *discordgocmds.CommandArgs) error {
	err := c.Player.PlayRandomSound(args.Guild, args.User)
	if err == player.ErrNotInVoice {
		msg, err := discordbot.SendEmbedError(args.Session, args.Channel.ID,
			"You need to be in a voice channel to play sounds.", "")
		msg.DeleteAfter(6 * time.Second)
		return err
	}
	return err
}
