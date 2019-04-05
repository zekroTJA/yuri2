package commands

import (
	"fmt"
	"strings"

	"github.com/zekroTJA/discordgo"

	"github.com/zekroTJA/yuri2/internal/player"
	"github.com/zekroTJA/yuri2/internal/static"

	"github.com/zekroTJA/yuri2/internal/database"
	"github.com/zekroTJA/yuri2/pkg/discordgocmds"
)

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
	// sfl, err := c.Player.GetLocalFiles()
	// if err != nil {
	// 	return err
	// }

	sfl := player.SoundFileList(make([]*player.SoundFile, 550))

	for i := range sfl {
		sfl[i] = &player.SoundFile{
			Name: "testsound",
		}
	}

	if len(args.Args) > 0 && strings.ToLower(args.Args[0]) == "s" {
		sfl.SortByDate()
	} else {
		sfl.SortByName()
	}

	sites := 1
	switch {
	case len(sfl) > 60:
		sites = 3
	case len(sfl) > 40:
		sites = 2
	}

	perSite := int(len(sfl)/sites) + 1

	emb := &discordgo.MessageEmbed{
		Color: static.ColorDefault,
		Provider: &discordgo.MessageEmbedProvider{
			Name: "test provider",
		},
		Title:  "Sound List",
		Fields: make([]*discordgo.MessageEmbedField, sites),
	}

	fmt.Println(len(sfl), sites)

	strList := make([][]string, sites)
	for i := range strList {
		strList[i] = make([]string, perSite)
	}

	strListC := 0
	site := 0
	for i, sound := range sfl {
		skip := (site + 1) * perSite
		if i >= skip {
			strListC = 0
			site++
		}
		strList[site][strListC] = sound.Name
		strListC++
	}

	for i, l := range strList {
		emb.Fields[i] = &discordgo.MessageEmbedField{
			Inline: true,
			Name:   "-",
			Value:  strings.Join(l, "\n"),
		}
	}

	_, err := args.Session.ChannelMessageSendEmbed(args.Channel.ID, emb)

	return err
}
