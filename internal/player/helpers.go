package player

import (
	"github.com/zekroTJA/discordgo"
)

func (p *Player) getMemberCountInVoiceChan(guild *discordgo.Guild, chanID string) int {
	var c int
	for _, vs := range guild.VoiceStates {
		if vs.ChannelID == chanID {
			c++
		}
	}

	return c
}

func (p *Player) updateSelfVS(newVS *discordgo.VoiceState) {
	if newVS != nil && newVS.UserID == p.session.State.User.ID {
		cVS, ok := p.selfVoiceStates[newVS.GuildID]
		if !ok {
			p.selfVoiceStates[newVS.GuildID] = newVS
			return
		}
		cVS.ChannelID = newVS.ChannelID
	}
}

func (p *Player) getUsersVoiceState(guildID, userID string) (*discordgo.VoiceState, error) {
	guild, err := p.session.Guild(guildID)
	if err != nil {
		return nil, err
	}

	for _, vs := range guild.VoiceStates {
		if vs.UserID == userID {
			return vs, nil
		}
	}

	return nil, nil
}
