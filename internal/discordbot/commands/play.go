package commands

import (
	"fmt"

	"github.com/zekroTJA/yuri2/internal/database"
	"github.com/zekroTJA/yuri2/pkg/discordgocmds"
)

// Play provides command functionalities
// for the prefix command
type Play struct {
	DB database.Middleware
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
	fmt.Printf("%+v\n", args)
	return nil
}
