package player

import (
	"time"

	"github.com/foxbot/gavalink"
	"github.com/zekroTJA/discordgo"
	"github.com/zekroTJA/yuri2/internal/logger"
)

// VoiceServerUpdateHandler is the players handler for
// the VoiceServerUpdate event.
func (p *Player) VoiceServerUpdateHandler(s *discordgo.Session, e *discordgo.VoiceServerUpdate) {
	vsu := gavalink.VoiceServerUpdate{
		Endpoint: e.Endpoint,
		GuildID:  e.GuildID,
		Token:    e.Token,
	}

	player, err := p.link.GetPlayer(e.GuildID)
	if err == nil {
		err = player.Forward(s.State.SessionID, vsu)
		if err != nil {
			p.onError("VoiceServerUpdate#GetPlayer", err)
		}
		logger.Debug("PLAYER :: using player for guild %s", e.GuildID)
		return
	}

	node, err := p.link.BestNode()
	if err != nil {
		p.onError("VoiceServerUpdate#BestNode", err)
		return
	}

	player, err = node.CreatePlayer(e.GuildID, s.State.SessionID, vsu, p.eventHandler)
	if err != nil {
		p.onError("VoiceServerUpdate#CreatePlayer", err)
	}

	vol, err := p.db.GetGuildVolume(e.GuildID)
	if err != nil {
		p.onError("VoiceServerUpdate#GetVolume", err)
	}

	if err = player.Volume(vol); err != nil {
		p.onError("VoiceServerUpdate#SetVolume", err)
	}

	logger.Debug("PLAYER :: created player for guild %s", e.GuildID)
}

// VoiceStateUpdateHandler is the players event handler for
// the VoiceStateUpdate event.
func (p *Player) VoiceStateUpdateHandler(s *discordgo.Session, e *discordgo.VoiceStateUpdate) {
	var oldVS *discordgo.VoiceState
	if p.lastVoiceStates.Contains(e.UserID) {
		oldVS, _ = p.lastVoiceStates.GetValue(e.UserID).(*discordgo.VoiceState)
	}

	newVS := e.VoiceState
	defer p.lastVoiceStates.Set(e.UserID, newVS, voiceStateLifetime)

	switch {

	// User left the voice channel
	case oldVS != nil && oldVS.ChannelID != "" && newVS.ChannelID == "":
		p.onVoiceChannelLeft(oldVS, newVS)

	// User moved to another channel
	case oldVS != nil && oldVS.ChannelID != "" && newVS.ChannelID != "" && oldVS.ChannelID != newVS.ChannelID:
		p.onVoiceChannelChange(oldVS, newVS)

	// User joins a channel
	case (oldVS == nil || oldVS.ChannelID == "") && newVS.ChannelID != "":
		p.onVoiceChannelJoined(oldVS, newVS)

	// User self muted
	case oldVS != nil && !oldVS.SelfMute && newVS.SelfMute:
		p.onUserSelfMuted(oldVS, newVS)
	}
}

func (p *Player) onVoiceChannelLeft(oldVS, newVS *discordgo.VoiceState) {
	if newVS.UserID == p.session.State.User.ID {
		delete(p.selfVoiceStates, newVS.GuildID)
	}

	p.autoLeftEmptyVoice(oldVS, newVS)
}

func (p *Player) onVoiceChannelJoined(oldVS, newVS *discordgo.VoiceState) {
	p.updateSelfVS(newVS)
	// if oldVS.UserID == p.session.State.User.ID || oldVS.ChannelID
}

func (p *Player) onVoiceChannelChange(oldVS, newVS *discordgo.VoiceState) {
	p.updateSelfVS(newVS)

	p.autoLeftEmptyVoice(oldVS, newVS)
	// if oldVS.UserID == p.session.State.User.ID || oldVS.ChannelID
}

func (p *Player) onUserSelfMuted(oldVS, newVS *discordgo.VoiceState) {
	p.fastMuteTrigger(oldVS, newVS)
}

func (p *Player) autoLeftEmptyVoice(oldVS, newVS *discordgo.VoiceState) {
	cVS, ok := p.selfVoiceStates[newVS.GuildID]

	if ok && oldVS.UserID != p.session.State.User.ID && oldVS.ChannelID == cVS.ChannelID {

		guild, err := p.session.Guild(newVS.GuildID)
		if err != nil {
			p.onError("autoLeftEmptyVoice#getGuild", err)
			return
		}

		if p.getMemberCountInVoiceChan(guild, cVS.ChannelID) <= 1 {
			time.AfterFunc(autoQuitDuration, func() {
				if p.getMemberCountInVoiceChan(guild, cVS.ChannelID) <= 1 {
					if err = p.LeaveVoiceChannel(cVS.GuildID, cVS.UserID); err != nil {
						p.onError("autoLeftEmptyVoice#QuitVoice", err)
					}
				}
			})
		}
	}
}

func (p *Player) fastMuteTrigger(oldVS, newVS *discordgo.VoiceState) {
	time.AfterFunc(fastMuteTriggerDuration, func() {
		vs, err := p.getUsersVoiceState(newVS.GuildID, newVS.UserID)
		if err != nil {
			p.onError("fastMuteTrigger#getVS", err)
			return
		}
		if vs != nil && !vs.SelfMute {
			val, err := p.db.GetFastTrigger(newVS.UserID)
			if err != nil {
				p.onError("fastMuteTrigger#getValueFromDB", err)
				return
			}

			user, err := p.session.User(newVS.UserID)
			if err != nil {
				p.onError("fastMuteTrigger#getUser", err)
				return
			}

			guild, err := p.session.Guild(newVS.GuildID)
			if err != nil {
				p.onError("fastMuteTrigger#getGuild", err)
				return
			}

			if val == "" {
				err = p.PlayRandomSound(guild, user)
			} else {
				err = p.Play(guild, user, val, ResourceLocal)
			}

			p.onError("fastMuteTrigger#playSound", err)
		}
	})
}
