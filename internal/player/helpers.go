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

func (p *Player) checkPerms(guild *discordgo.Guild, userID string) error {
	if p.blockedRoleName == "" && (p.playRoleName == "" || p.playRoleName == "@everyone") {
		return nil
	}

	memb, err := p.session.State.Member(guild.ID, userID)
	if err == discordgo.ErrStateNotFound {
		memb, err = p.session.GuildMember(guild.ID, userID)
	}
	if err != nil {
		return err
	}

	var isPlayer bool

	for _, gRole := range guild.Roles {
		if gRole.Name != p.blockedRoleName && gRole.Name != p.playRoleName {
			continue
		}

		for _, mrID := range memb.Roles {
			if mrID != gRole.ID {
				continue
			}

			switch gRole.Name {
			case p.blockedRoleName:
				return ErrBlocked
			case p.playRoleName:
				isPlayer = true
			}
		}
	}

	if !isPlayer {
		return ErrNoPermission
	}

	return nil
}

func (p *Player) checkPermsByIDs(guildID, userID string) error {
	if p.blockedRoleName == "" && (p.playRoleName == "" || p.playRoleName == "@everyone") {
		return nil
	}

	guild, err := p.session.Guild(guildID)
	if err != nil {
		return err
	}
	return p.checkPerms(guild, userID)
}
