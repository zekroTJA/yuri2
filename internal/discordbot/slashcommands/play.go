package slashcommands

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/yuri2/internal/discordbot"
	"github.com/zekroTJA/yuri2/internal/player"
	"github.com/zekrotja/ken"
)

type Play struct {
	Player *player.Player
}

var _ ken.Command = (*Play)(nil)

func (c *Play) Name() string {
	return "play"
}

func (c *Play) Description() string {
	return "Play a sound."
}

func (c *Play) Version() string {
	return "1.0.0"
}

func (c *Play) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Play) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "name",
			Description: "The name of the sound to be played (or empty if random).",
		},
	}
}

func (c *Play) Run(ctx *ken.Ctx) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	name, ok := ctx.Options().GetByNameOptional("name")

	guild, err := discordbot.Guild(ctx.Session, ctx.Event.GuildID)
	if err != nil {
		return
	}

	if ok {
		err = c.Player.Play(guild, ctx.User(), name.StringValue(), player.ResourceLocal)

	} else {
		err = c.Player.PlayRandomSound(guild, ctx.User())
	}

	if err == player.ErrNotFound {
		return nil
	}

	if err == player.ErrNotInVoice {
		err = ctx.FollowUpError("You need to be in a voice channel to play sounds.", "").Error
		return
	}

	err = ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: "Sound is being played.",
	}).DeleteAfter(5 * time.Second).Error

	return
}
