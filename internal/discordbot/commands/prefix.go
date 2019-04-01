package commands

import (
	"fmt"
	"time"

	"github.com/zekroTJA/yuri2/internal/database"
	"github.com/zekroTJA/yuri2/internal/discordbot"
	"github.com/zekroTJA/yuri2/internal/static"
	"github.com/zekroTJA/yuri2/pkg/discordgocmds"
)

// Prefix provides command functionalities
// for the prefix command
type Prefix struct {
	PermLvl int
	DB      database.Middleware
}

// GetInvokes returns the invokes
// for this command.
func (c *Prefix) GetInvokes() []string {
	return []string{"prefix", "pre"}
}

// GetDescription returns the description
// for this command
func (c *Prefix) GetDescription() string {
	return "Set custom prefix for this guild"
}

// GetHelp returns the help text for
// this command.
func (c *Prefix) GetHelp() string {
	return "`prefix <newPrefix>` - set new prefix"
}

// GetGroup returns the group of
// the command
func (c *Prefix) GetGroup() string {
	return discordgocmds.GroupAdmin
}

// GetPermission returns the minimum
// required required permission level
// to execute this command.
func (c *Prefix) GetPermission() int {
	return c.PermLvl
}

// Exec is the acual function which will
// be executed when the command was invoked.
func (c *Prefix) Exec(args *discordgocmds.CommandArgs) error {
	if len(args.Args) == 0 {
		m, err := discordbot.SendEmbedError(args.Session, args.Channel.ID,
			"Please enter a prefix as argument.\nSee `help prefix` for further information.",
			"Argument Error")
		m.DeleteAfter(6 * time.Second)
		return err
	}

	prefix := args.Args[0]
	err := c.DB.SetGuildPrefix(args.Guild.ID, args.Args[0])
	if err != nil {
		return err
	}

	msg, err := discordbot.SendEmbedMessage(args.Session, args.Channel.ID,
		fmt.Sprintf("Succesfully set `%s` as prefix for this guild.", prefix), "", static.ColorGreen)
	msg.DeleteAfter(5 * time.Second)

	return err
}
