package commands

import (
	"fmt"

	"github.com/zekroTJA/yuri2/internal/database"
	"github.com/zekroTJA/yuri2/internal/player"
	"github.com/zekroTJA/yuri2/pkg/discordgocmds"
)

// Test provides command functionalities
// for the test command
type Test struct {
	PermLvl int
	DB      database.Middleware
	Player  *player.Player
}

// GetInvokes returns the invokes
// for this command.
func (c *Test) GetInvokes() []string {
	return []string{"test"}
}

// GetDescription returns the description
// for this command
func (c *Test) GetDescription() string {
	return ""
}

// GetHelp returns the help text for
// this command.
func (c *Test) GetHelp() string {
	return ""
}

// GetGroup returns the group of
// the command
func (c *Test) GetGroup() string {
	return discordgocmds.GroupAdmin
}

// GetPermission returns the minimum
// required required permission level
// to execute this command.
func (c *Test) GetPermission() int {
	return c.PermLvl
}

// Exec is the acual function which will
// be executed when the command was invoked.
func (c *Test) Exec(args *discordgocmds.CommandArgs) error {
	m, err := args.Session.GuildMember("526196711962705925", "123123273489234")
	fmt.Print(m, err)
	return err
}
