package commands

import (
	"github.com/zekroTJA/yuri2/internal/player"
	"github.com/zekroTJA/yuri2/internal/static"

	"github.com/zekroTJA/yuri2/internal/database"
	"github.com/zekroTJA/yuri2/pkg/discordgocmds"
)

// Stop provides command functionalities
// for the random command
type Stop struct {
	DB      database.Middleware
	Player  *player.Player
	PermLvl int
}

// GetInvokes returns the invokes
// for this command.
func (c *Stop) GetInvokes() []string {
	return []string{"stop"}
}

// GetDescription returns the description
// for this command
func (c *Stop) GetDescription() string {
	return "Stop a playing sound"
}

// GetHelp returns the help text for
// this command.
func (c *Stop) GetHelp() string {
	return "`stop` - stop a playing sound"
}

// GetGroup returns the group of
// the command
func (c *Stop) GetGroup() string {
	return static.CommandGroupPlayer
}

// GetPermission returns the minimum
// required required permission level
// to execute this command.
func (c *Stop) GetPermission() int {
	return c.PermLvl
}

// Exec is the acual function which will
// be executed when the command was invoked.
func (c *Stop) Exec(args *discordgocmds.CommandArgs) error {
	return c.Player.Stop(args.Guild, args.User)
}
