package commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/zekroTJA/yuri2/internal/discordbot"
	"github.com/zekroTJA/yuri2/internal/static"

	"github.com/zekroTJA/yuri2/internal/database"
	"github.com/zekroTJA/yuri2/pkg/discordgocmds"
)

const (
	deleteTimeoutStats = 3 * time.Minute
)

// Stats provides command functionalities
// for the stats command
type Stats struct {
	DB      database.Middleware
	PermLvl int
}

// GetInvokes returns the invokes
// for this command.
func (c *Stats) GetInvokes() []string {
	return []string{"stats", "stat"}
}

// GetDescription returns the description
// for this command
func (c *Stats) GetDescription() string {
	return "Display most played sounds on this guild"
}

// GetHelp returns the help text for
// this command.
func (c *Stats) GetHelp() string {
	return "`stats (<limit>)` - display most played sounds on this guild"
}

// GetGroup returns the group of
// the command
func (c *Stats) GetGroup() string {
	return static.CommandGroupPlayer
}

// GetPermission returns the minimum
// required required permission level
// to execute this command.
func (c *Stats) GetPermission() int {
	return c.PermLvl
}

// Exec is the actual function which will
// be executed when the command was invoked.
func (c *Stats) Exec(args *discordgocmds.CommandArgs) error {
	var err error
	limit := defLimit

	if len(args.Args) > 0 {
		limit, err = strconv.Atoi(args.Args[0])
		if err != nil || limit < 1 || limit > 500 {
			msg, err := discordbot.SendEmbedError(args.Session, args.Channel.ID,
				"Please enter a valid number for limit in range [1, 500].", "")
			msg.DeleteAfter(6 * time.Second)
			return err
		}
	}

	stats, err := c.DB.GetSoundStats(args.Guild.ID, limit)
	if err != nil {
		return err
	}

	allCount, err := c.DB.GetLogLen(args.Guild.ID)
	if err != nil {
		return err
	}

	if len(stats) < 1 {
		msg, err := discordbot.SendEmbedError(args.Session, args.Channel.ID,
			"Stats are empty.", "")
		msg.DeleteAfter(6 * time.Second)
		return err
	}

	strList := make([]string, len(stats))
	for i, stat := range stats {
		strList[i] = fmt.Sprintf("`%d` - **%s** - `%d`",
			i+1, stat.Sound, stat.Count)
	}

	if limit > 30 {
		lmsg, err := discordbot.NewListMessage(args.Session, args.Channel.ID,
			"Sound Stats", fmt.Sprintf("Total plays: **`%d`**\nTop %d played sounds on this guild",
				allCount, len(stats)), strList, 30, 0)
		lmsg.DeleteAfter(deleteTimeoutLog)
		return err
	}

	_, err = discordbot.SendEmbedMessage(args.Session, args.Channel.ID,
		fmt.Sprintf("Total plays: **`%d`**\nTop %d sounds played on this guild\n\n%s",
			allCount, len(stats), strings.Join(strList, "\n")), "Sound Stats", 0)

	return err
}
