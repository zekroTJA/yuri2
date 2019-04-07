package commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/zekroTJA/yuri2/internal/discordbot"
	"github.com/zekroTJA/yuri2/internal/player"
	"github.com/zekroTJA/yuri2/internal/static"

	"github.com/zekroTJA/yuri2/internal/database"
	"github.com/zekroTJA/yuri2/pkg/discordgocmds"
)

const (
	defLimit         = 30
	deleteTimeoutLog = 3 * time.Minute
	timeFormat       = "01/02 - 15:04:05"
)

// Log provides command functionalities
// for the log command
type Log struct {
	DB      database.Middleware
	Player  *player.Player
	PermLvl int
}

// GetInvokes returns the invokes
// for this command.
func (c *Log) GetInvokes() []string {
	return []string{"log", "soundlog"}
}

// GetDescription returns the description
// for this command
func (c *Log) GetDescription() string {
	return "Display sound log"
}

// GetHelp returns the help text for
// this command.
func (c *Log) GetHelp() string {
	return "`log (<limit>)` - display sound log"
}

// GetGroup returns the group of
// the command
func (c *Log) GetGroup() string {
	return static.CommandGroupPlayer
}

// GetPermission returns the minimum
// required required permission level
// to execute this command.
func (c *Log) GetPermission() int {
	return c.PermLvl
}

// Exec is the acual function which will
// be executed when the command was invoked.
func (c *Log) Exec(args *discordgocmds.CommandArgs) error {
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

	sles, err := c.DB.GetLogEntries(args.Guild.ID, 0, limit)
	if err != nil {
		return err
	}

	allCount, err := c.DB.GetLogLen(args.Guild.ID)
	if err != nil {
		return err
	}

	if len(sles) < 1 {
		msg, err := discordbot.SendEmbedError(args.Session, args.Channel.ID,
			"Log is empty.", "")
		msg.DeleteAfter(6 * time.Second)
		return err
	}

	strList := make([]string, len(sles))
	for i, sle := range sles {
		strList[i] = fmt.Sprintf("`[%s]` - **%s** *[%s]* - %s",
			sle.Time.Format(timeFormat), sle.Sound,
			strings.ToUpper(sle.Source)[:1], sle.UserTag)
	}

	if limit > 30 {
		lmsg, err := discordbot.NewListMessage(args.Session, args.Channel.ID,
			"Sounds Log", fmt.Sprintf("Last %d of %d log entries", len(sles), allCount), strList, 30, 0)
		lmsg.DeleteAfter(deleteTimeoutLog)
		return err
	}

	_, err = discordbot.SendEmbedMessage(args.Session, args.Channel.ID,
		fmt.Sprintf("Last %d of %d log entries\n\n%s", len(sles), allCount, strings.Join(strList, "\n")), "Sounds Log", 0)

	return err
}
