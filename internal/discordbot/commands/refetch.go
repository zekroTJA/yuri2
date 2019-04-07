package commands

import (
	"time"

	"github.com/zekroTJA/yuri2/internal/discordbot"
	"github.com/zekroTJA/yuri2/internal/player"
	"github.com/zekroTJA/yuri2/internal/static"

	"github.com/zekroTJA/yuri2/internal/database"
	"github.com/zekroTJA/yuri2/pkg/discordgocmds"
)

// Refetch provides command functionalities
// for the refetch command
type Refetch struct {
	DB      database.Middleware
	Player  *player.Player
	PermLvl int
}

// GetInvokes returns the invokes
// for this command.
func (c *Refetch) GetInvokes() []string {
	return []string{"refetch", "reload"}
}

// GetDescription returns the description
// for this command
func (c *Refetch) GetDescription() string {
	return "Refetch local sounds"
}

// GetHelp returns the help text for
// this command.
func (c *Refetch) GetHelp() string {
	return "`refetch`- refetch local sounds"
}

// GetGroup returns the group of
// the command
func (c *Refetch) GetGroup() string {
	return static.CommandGroupPlayer
}

// GetPermission returns the minimum
// required required permission level
// to execute this command.
func (c *Refetch) GetPermission() int {
	return c.PermLvl
}

// Exec is the acual function which will
// be executed when the command was invoked.
func (c *Refetch) Exec(args *discordgocmds.CommandArgs) error {
	err := c.Player.FetchLocalSounds()
	if err != nil {
		return err
	}

	msg, err := discordbot.SendEmbedError(args.Session, args.Channel.ID,
		"Refetched local sounds.", "")
	msg.DeleteAfter(8 * time.Second)

	return err
}
