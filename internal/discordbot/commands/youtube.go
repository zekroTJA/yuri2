package commands

import (
	"strings"
	"time"

	"github.com/zekroTJA/yuri2/internal/discordbot"
	"github.com/zekroTJA/yuri2/internal/player"
	"github.com/zekroTJA/yuri2/internal/static"

	"github.com/zekroTJA/yuri2/internal/database"
	"github.com/zekroTJA/yuri2/pkg/discordgocmds"
)

// YouTube provides command functionalities
// for the yt command
type YouTube struct {
	DB      database.Middleware
	Player  *player.Player
	PermLvl int
}

// GetInvokes returns the invokes
// for this command.
func (c *YouTube) GetInvokes() []string {
	return []string{"yt", "youtube"}
}

// GetDescription returns the description
// for this command
func (c *YouTube) GetDescription() string {
	return "Play something by youtube link or id"
}

// GetHelp returns the help text for
// this command.
func (c *YouTube) GetHelp() string {
	return "`yt <link|id>` - play youtube video"
}

// GetGroup returns the group of
// the command
func (c *YouTube) GetGroup() string {
	return static.CommandGroupPlayer
}

// GetPermission returns the minimum
// required required permission level
// to execute this command.
func (c *YouTube) GetPermission() int {
	return c.PermLvl
}

// Exec is the acual function which will
// be executed when the command was invoked.
func (c *YouTube) Exec(args *discordgocmds.CommandArgs) error {
	if len(args.Args) < 1 {
		msg, err := discordbot.SendEmbedError(args.Session, args.Channel.ID,
			"Please specify a youtube link *(`https://youtube.com/watch?v=:ID` or `https://youtu.be/:ID`)* or the pure `:ID` of the video.", "")
		msg.DeleteAfter(6 * time.Second)
		return err
	}

	ident := args.Args[0]

	if i := strings.Index(ident, "&t="); i > -1 {
		ident = ident[:i]
	}

	cut := 0
	if i := strings.Index(ident, "youtube.com/watch?v="); i > -1 {
		cut = i + len("youtube.com/watch?v=")
	}
	if i := strings.Index(ident, "youtu.be/"); i > -1 {
		cut = i + len("youtu.be/")
	}

	ident = ident[cut:]

	err := c.Player.Play(args.Guild, args.User, ident, player.ResourceYouTube)
	if err == player.ErrNotInVoice {
		msg, err := discordbot.SendEmbedError(args.Session, args.Channel.ID,
			"You need to be in a voice channel to play sounds.", "")
		msg.DeleteAfter(6 * time.Second)
		return err
	}

	return err
}
