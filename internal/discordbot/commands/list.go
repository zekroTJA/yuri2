package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/zekroTJA/yuri2/internal/discordbot"

	"github.com/zekroTJA/yuri2/internal/player"
	"github.com/zekroTJA/yuri2/internal/static"

	"github.com/zekroTJA/yuri2/internal/database"
	"github.com/zekroTJA/yuri2/pkg/discordgocmds"
)

const deleteTimeoutList = 5 * time.Minute

// List provides command functionalities
// for the list command
type List struct {
	DB      database.Middleware
	Player  *player.Player
	PermLvl int
}

// GetInvokes returns the invokes
// for this command.
func (c *List) GetInvokes() []string {
	return []string{"ls", "list"}
}

// GetDescription returns the description
// for this command
func (c *List) GetDescription() string {
	return "Play a random, local file"
}

// GetHelp returns the help text for
// this command.
func (c *List) GetHelp() string {
	return "`r` - play a random, local file"
}

// GetGroup returns the group of
// the command
func (c *List) GetGroup() string {
	return static.CommandGroupPlayer
}

// GetPermission returns the minimum
// required required permission level
// to execute this command.
func (c *List) GetPermission() int {
	return c.PermLvl
}

// Exec is the acual function which will
// be executed when the command was invoked.
func (c *List) Exec(args *discordgocmds.CommandArgs) error {
	sfl, err := c.Player.GetLocalFiles()
	if err != nil {
		return err
	}

	if len(args.Args) > 0 && strings.ToLower(args.Args[0]) == "s" {
		sfl.SortByDate()
	} else {
		sfl.SortByName()
	}

	strList := make([]string, len(sfl))
	for i, v := range sfl {
		strList[i] = v.Name
	}

	msg, err := discordbot.NewListMessage(args.Session, args.Channel.ID,
		"Sound List", fmt.Sprintf("**%d Sounds**", len(sfl)), strList, 30, 0)
	msg.DeleteAfter(deleteTimeoutList)

	return err
}
