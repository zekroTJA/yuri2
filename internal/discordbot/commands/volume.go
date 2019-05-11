package commands

import (
	"fmt"
	"strconv"
	"time"

	"github.com/zekroTJA/yuri2/internal/discordbot"
	"github.com/zekroTJA/yuri2/internal/player"
	"github.com/zekroTJA/yuri2/internal/static"

	"github.com/zekroTJA/yuri2/internal/database"
	"github.com/zekroTJA/yuri2/pkg/discordgocmds"
)

// Volume provides command functionalities
// for the volume command
type Volume struct {
	DB      database.Middleware
	Player  *player.Player
	PermLvl int
}

// GetInvokes returns the invokes
// for this command.
func (c *Volume) GetInvokes() []string {
	return []string{"vol", "volume"}
}

// GetDescription returns the description
// for this command
func (c *Volume) GetDescription() string {
	return "Set the volume for the guild and the current player"
}

// GetHelp returns the help text for
// this command.
func (c *Volume) GetHelp() string {
	return "`vol` - display current volume\n" +
		"`vol <int>` - set volume"
}

// GetGroup returns the group of
// the command
func (c *Volume) GetGroup() string {
	return static.CommandGroupPlayer
}

// GetPermission returns the minimum
// required required permission level
// to execute this command.
func (c *Volume) GetPermission() int {
	return c.PermLvl
}

// Exec is the actual function which will
// be executed when the command was invoked.
func (c *Volume) Exec(args *discordgocmds.CommandArgs) error {
	if len(args.Args) < 1 {
		val, err := c.DB.GetGuildVolume(args.Guild.ID)
		if err != nil {
			return err
		}

		msg, err := discordbot.SendEmbedMessage(args.Session, args.Channel.ID,
			fmt.Sprintf("Current volume is **`%d`**.", val), "", 0)
		msg.DeleteAfter(6 * time.Second)
		return err
	}

	iVal, err := strconv.Atoi(args.Args[0])
	if err != nil || iVal > 1000 || iVal < 1 {
		msg, err := discordbot.SendEmbedMessage(args.Session, args.Channel.ID,
			"Volume must be a valid number in range of [1, 1000].", "", 0)
		msg.DeleteAfter(8 * time.Second)
		return err
	}

	if err = c.Player.SetVolume(args.Guild.ID, args.User.ID, iVal); err != nil {
		return err
	}

	msg, err := discordbot.SendEmbedMessage(args.Session, args.Channel.ID,
		fmt.Sprintf("Set guilds players volume to **`%d`**.", iVal), "", 0)
	msg.DeleteAfter(6 * time.Second)

	return err
}
