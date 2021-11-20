package slashcommands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/yuri2/internal/database"
	"github.com/zekroTJA/yuri2/internal/player"
	"github.com/zekrotja/ken"
)

type Bind struct {
	Player *player.Player
	Db     database.Middleware
}

var _ ken.Command = (*Bind)(nil)

func (c *Bind) Name() string {
	return "bind"
}

func (c *Bind) Description() string {
	return "Bind a sound (or random) to the fast trigger."
}

func (c *Bind) Version() string {
	return "1.0.0"
}

func (c *Bind) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Bind) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "name",
			Description: "The name of the sound to be bindet (or `r`/`random`).",
		},
	}
}

func (c *Bind) Run(ctx *ken.Ctx) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	nameV, ok := ctx.Options().GetByNameOptional("name")

	if ok {
		val, err := c.Db.GetFastTrigger(ctx.User().ID)
		if err != nil {
			return err
		}

		if val == "" {
			val = "random"
		}

		err = ctx.FollowUpEmbed(&discordgo.MessageEmbed{
			Description: fmt.Sprintf("Currently, fast trigger is bound to **`%s`**.", val),
		}).Error
	} else {
		name := nameV.StringValue()

		if name == "r" || name == "random" {
			if err := c.Db.SetFastTrigger(ctx.User().ID, ""); err != nil {
				return err
			}
			err = ctx.FollowUpEmbed(&discordgo.MessageEmbed{
				Description: "Bound fast trigger to **`random`**.",
			}).Error
			return
		}

		if _, ok := c.Player.GetLocalSoundPath(name); !ok {
			err = ctx.FollowUpError(fmt.Sprintf("Could not fetch any local sound by identifier `%s`.", name), "").Error
			return
		}

		if err = c.Db.SetFastTrigger(ctx.User().ID, name); err != nil {
			return
		}

		err = ctx.FollowUpEmbed(&discordgo.MessageEmbed{
			Description: fmt.Sprintf("Bound fast trigger to **`%s`**.", name),
		}).Error
	}

	return
}
