package commands

import (
	"fmt"
	"time"

	"github.com/zekroTJA/yuri2/internal/discordbot"
	"github.com/zekroTJA/yuri2/internal/player"
	"github.com/zekroTJA/yuri2/internal/static"

	"github.com/zekroTJA/yuri2/internal/database"
	"github.com/zekroTJA/yuri2/pkg/discordgocmds"
)

// Bind provides command functionalities
// for the bind command
type Bind struct {
	DB      database.Middleware
	Player  *player.Player
	PermLvl int
}

// GetInvokes returns the invokes
// for this command.
func (c *Bind) GetInvokes() []string {
	return []string{"bind"}
}

// GetDescription returns the description
// for this command
func (c *Bind) GetDescription() string {
	return "Bind a sound to fast trigger"
}

// GetHelp returns the help text for
// this command.
func (c *Bind) GetHelp() string {
	return "`bind <sound>` - bind a specific sound\n" +
		"`bind r` - bind random sounds"
}

// GetGroup returns the group of
// the command
func (c *Bind) GetGroup() string {
	return static.CommandGroupSettings
}

// GetPermission returns the minimum
// required required permission level
// to execute this command.
func (c *Bind) GetPermission() int {
	return c.PermLvl
}

// Exec is the actual function which will
// be executed when the command was invoked.
func (c *Bind) Exec(args *discordgocmds.CommandArgs) error {
	if len(args.Args) < 1 {
		val, err := c.DB.GetFastTrigger(args.User.ID)
		if err != nil {
			return err
		}

		if val == "" {
			val = "random"
		}

		msg, err := discordbot.SendEmbedMessage(args.Session, args.Channel.ID,
			fmt.Sprintf("Currently, fast trigger is bound to **`%s`**.", val), "", 0)
		msg.DeleteAfter(6 * time.Second)
		return err
	}

	val := args.Args[0]

	if val == "r" {
		if err := c.DB.SetFastTrigger(args.User.ID, ""); err != nil {
			return err
		}

		msg, err := discordbot.SendEmbedMessage(args.Session, args.Channel.ID,
			"Bound fast trigger to **`random`**.", "", 0)
		msg.DeleteAfter(6 * time.Second)
		return err
	}

	if _, ok := c.Player.GetLocalSoundPath(val); !ok {
		msg, err := discordbot.SendEmbedError(args.Session, args.Channel.ID,
			fmt.Sprintf("Could not fetch any local sound by identifier `%s`.", val), "")
		msg.DeleteAfter(8 * time.Second)
		return err
	}

	if err := c.DB.SetFastTrigger(args.User.ID, val); err != nil {
		return err
	}

	msg, err := discordbot.SendEmbedMessage(args.Session, args.Channel.ID,
		fmt.Sprintf("Bound fast trigger to **`%s`**.", val), "", 0)
	msg.DeleteAfter(6 * time.Second)

	return err
}
