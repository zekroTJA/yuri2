package commands

import (
	"regexp"
	"strings"
	"time"

	"github.com/zekroTJA/yuri2/internal/discordbot"
	"github.com/zekroTJA/yuri2/internal/player"
	"github.com/zekroTJA/yuri2/internal/static"

	"github.com/zekroTJA/yuri2/internal/database"
	"github.com/zekroTJA/yuri2/pkg/discordgocmds"
)

// List provides command functionalities
// for the list command
type Search struct {
	DB      database.Middleware
	Player  *player.Player
	PermLvl int
}

// GetInvokes returns the invokes
// for this command.
func (c *Search) GetInvokes() []string {
	return []string{"s", "search"}
}

// GetDescription returns the description
// for this command
func (c *Search) GetDescription() string {
	return "Search local sounds by wildcard or regex"
}

// GetHelp returns the help text for
// this command.
func (c *Search) GetHelp() string {
	return "`s <query>` - search sounds by wildcard query\n" +
		"`s rx <rxQuery>` - serch sounds by regex"
}

// GetGroup returns the group of
// the command
func (c *Search) GetGroup() string {
	return static.CommandGroupPlayer
}

// GetPermission returns the minimum
// required required permission level
// to execute this command.
func (c *Search) GetPermission() int {
	return c.PermLvl
}

// Exec is the acual function which will
// be executed when the command was invoked.
func (c *Search) Exec(args *discordgocmds.CommandArgs) error {
	if len(args.Args) < 1 {
		msg, err := discordbot.SendEmbedError(args.Session, args.Channel.ID,
			"Please enter a valid search query.\n*Use `help s` for more information.*", "")
		msg.DeleteAfter(8 * time.Second)
		return err
	}

	byRx := args.Args[0] == "rx"
	var rx *regexp.Regexp
	var err error
	var searchFunc func(string) bool

	if byRx && len(args.Args) < 2 {
		msg, err := discordbot.SendEmbedError(args.Session, args.Channel.ID,
			"Please enter a valid regular expression.\n*Use `help s` for more information.*", "")
		msg.DeleteAfter(8 * time.Second)
		return err
	}

	if byRx {
		if rx, err = regexp.Compile(args.Args[1]); err != nil {
			return err
		}
		searchFunc = c.matchRegExFunc(rx)
	} else {
		searchFunc = c.matchWildcardFunc(args.Args[0])
	}

	sfl, err := c.Player.GetLocalFiles()
	if err != nil {
		return err
	}

	sfl.SortByName()

	strList := make([]string, len(sfl))
	count := 0
	for _, v := range sfl {
		if searchFunc(v.Name) {
			strList[count] = v.Name
			count++
		}
	}

	if count == 0 {
		msg, err := discordbot.SendEmbedError(args.Session, args.Channel.ID,
			"No results matched.", "")
		msg.DeleteAfter(6 * time.Second)
		return err
	}

	if count > 30 {
		msg, err := discordbot.SendEmbedError(args.Session, args.Channel.ID,
			"Too much search results.", "")
		msg.DeleteAfter(6 * time.Second)
		return err
	}

	strList = strList[:count]

	_, err = discordbot.SendEmbedMessage(args.Session, args.Channel.ID,
		strings.Join(strList, "\n"), "Search Results", 0)

	return err
}

func (c *Search) matchWildcardFunc(q string) func(string) bool {
	return func(name string) bool {
		name = strings.ToLower(name)
		if strings.HasPrefix(q, "*") {
			return strings.HasSuffix(name, q[1:])
		}
		if strings.HasSuffix(q, "*") {
			return strings.HasPrefix(name, q[:1])
		}
		return strings.Contains(name, q)
	}
}

func (c *Search) matchRegExFunc(rx *regexp.Regexp) func(string) bool {
	return rx.MatchString
}
