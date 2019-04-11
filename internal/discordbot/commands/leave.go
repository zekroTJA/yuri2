package commands

import (
	"github.com/zekroTJA/yuri2/internal/player"
	"github.com/zekroTJA/yuri2/internal/static"

	"github.com/zekroTJA/yuri2/pkg/discordgocmds"
)

// Leave provides command functionalities
// for the leave command
type Leave struct {
	Player  *player.Player
	PermLvl int
}

// GetInvokes returns the invokes
// for this command.
func (c *Leave) GetInvokes() []string {
	return []string{"leave", "quit"}
}

// GetDescription returns the description
// for this command
func (c *Leave) GetDescription() string {
	return "Quit the currently connected voice channel"
}

// GetHelp returns the help text for
// this command.
func (c *Leave) GetHelp() string {
	return "`quit` - quit voice channel"
}

// GetGroup returns the group of
// the command
func (c *Leave) GetGroup() string {
	return static.CommandGroupPlayer
}

// GetPermission returns the minimum
// required required permission level
// to execute this command.
func (c *Leave) GetPermission() int {
	return c.PermLvl
}

// Exec is the actual function which will
// be executed when the command was invoked.
func (c *Leave) Exec(args *discordgocmds.CommandArgs) error {
	return c.Player.LeaveVoiceChannel(args.Guild.ID)
}
