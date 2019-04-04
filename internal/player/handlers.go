package player

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/foxbot/gavalink"
	"github.com/zekroTJA/yuri2/internal/logger"
)

func (p *Player) VoiceServerUpdateHandler(s *discordgo.Session, e *discordgo.VoiceServerUpdate) {
	vsu := gavalink.VoiceServerUpdate{
		Endpoint: e.Endpoint,
		GuildID:  e.GuildID,
		Token:    e.Token,
	}

	fmt.Printf("%+v\n", e)

	if player, err := p.link.GetPlayer(e.GuildID); err == nil {
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

	_, err = node.CreatePlayer(e.GuildID, s.State.SessionID, vsu, p.eventHandler)
	if err != nil {
		p.onError("VoiceServerUpdate#CreatePlayer", err)
	}
	logger.Debug("PLAYER :: created player for guild %s", e.GuildID)
}

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
	}
}

func (p *Player) onVoiceChannelLeft(oldVS, newVS *discordgo.VoiceState) {
	if newVS.UserID == p.session.State.User.ID {
		p.selfVoiceState = nil
		fmt.Println("self left")
	}
	// if oldVS.UserID == p.session.State.User.ID || oldVS.ChannelID
}

func (p *Player) onVoiceChannelJoined(oldVS, newVS *discordgo.VoiceState) {
	if newVS.UserID == p.session.State.User.ID {
		fmt.Println("self joined")

		p.selfVoiceState = newVS
	}
	// if oldVS.UserID == p.session.State.User.ID || oldVS.ChannelID
}

func (p *Player) onVoiceChannelChange(oldVS, newVS *discordgo.VoiceState) {
	if newVS.UserID == p.session.State.User.ID {
		fmt.Println("self switched")

		p.selfVoiceState = newVS
	}
	// if oldVS.UserID == p.session.State.User.ID || oldVS.ChannelID
}
