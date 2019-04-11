package api

import (
	"fmt"

	"github.com/zekroTJA/discordgo"
	"github.com/zekroTJA/yuri2/internal/player"

	"github.com/zekroTJA/yuri2/internal/discordbot"
	"github.com/zekroTJA/yuri2/pkg/wsmgr"
)

type wsInitData struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
}

type wsPlayData struct {
	Ident  string `json:"ident"`
	Source int    `json:"source"`
}

type wsIdent struct {
	UserID string
	Guilds []*discordgo.Guild
}

// Event: INIT
func (api *API) wsInitHandler(e *wsmgr.Event) {
	data := new(wsInitData)

	err := e.ParseDataTo(data)
	if err != nil {
		wsSendError(e.Sender, fmt.Sprintf("failed parsing data: %s", err.Error()))
		return
	}

	guilds := discordbot.GetUsersGuilds(api.session, data.UserID)
	if guilds == nil {
		wsSendError(e.Sender, "forbidden: you must be a member of a guild the bot is also member of")
		return
	}

	ok, _, err := api.auth.CheckAndRefersh(data.UserID, data.Token)
	if err != nil {
		wsSendError(e.Sender, fmt.Sprintf("failed checking auth: %s", err.Error()))
		return
	}

	if !ok {
		wsSendError(e.Sender, "unauthorized")
		e.Sender.Close()
		return
	}

	e.Sender.SetIdent(&wsIdent{
		UserID: data.UserID,
		Guilds: guilds,
	})
}

// Event: JOIN
func (api *API) wsJoinHandler(e *wsmgr.Event) {
	ident := wsCheckInitilized(e.Sender)
	if ident == nil {
		return
	}

	_, vs := discordbot.GetUsersGuildInVoice(api.session, ident.UserID)
	if vs == nil {
		wsSendError(e.Sender, "you need to be in a voice channel to perform this command")
		return
	}

	err := api.player.JoinVoiceCannel(vs)
	if err != nil {
		wsSendError(e.Sender, fmt.Sprintf("command failed: %s", err.Error()))
	}
}

// Event: LEAVE
func (api *API) wsLeaveHandler(e *wsmgr.Event) {
	ident := wsCheckInitilized(e.Sender)
	if ident == nil {
		return
	}

	guild, _ := discordbot.GetUsersGuildInVoice(api.session, ident.UserID)
	if guild == nil {
		wsSendError(e.Sender, "you need to be in a voice channel to perform this command")
		return
	}

	err := api.player.LeaveVoiceChannel(guild.ID)
	if err != nil {
		wsSendError(e.Sender, fmt.Sprintf("command failed: %s", err.Error()))
	}
}

// Event: PLAY
func (api *API) wsPlayHandler(e *wsmgr.Event) {
	ident := wsCheckInitilized(e.Sender)
	if ident == nil {
		return
	}

	guild, _ := discordbot.GetUsersGuildInVoice(api.session, ident.UserID)
	if guild == nil {
		wsSendError(e.Sender, "you need to be in a voice channel to perform this command")
		return
	}

	user, err := api.session.User(ident.UserID)
	if err != nil || user == nil {
		wsSendError(e.Sender, fmt.Sprintf("faield getting user context: %s", err.Error()))
		return
	}

	data := new(wsPlayData)
	err = e.ParseDataTo(data)
	if err != nil {
		wsSendError(e.Sender, fmt.Sprintf("failed parsing data: %s", err.Error()))
		return
	}

	if data.Ident == "" {
		wsSendError(e.Sender, "invalid arguments: ident must be a valid string value")
		return
	}

	err = api.player.Play(guild, user, data.Ident, player.ResourceType(data.Source))
	if err != nil {
		wsSendError(e.Sender, fmt.Sprintf("command failed: %s", err.Error()))
	}
}

// Event: RANDOM
func (api *API) wsRandomHandler(e *wsmgr.Event) {
	ident := wsCheckInitilized(e.Sender)
	if ident == nil {
		return
	}

	guild, _ := discordbot.GetUsersGuildInVoice(api.session, ident.UserID)
	if guild == nil {
		wsSendError(e.Sender, "you need to be in a voice channel to perform this command")
		return
	}

	user, err := api.session.User(ident.UserID)
	if err != nil || user == nil {
		wsSendError(e.Sender, fmt.Sprintf("faield getting user context: %s", err.Error()))
		return
	}

	err = api.player.PlayRandomSound(guild, user)
	if err != nil {
		wsSendError(e.Sender, fmt.Sprintf("command failed: %s", err.Error()))
	}
}

// Event: VOLUME
func (api *API) wsVolumeHandler(e *wsmgr.Event) {
	ident := wsCheckInitilized(e.Sender)
	if ident == nil {
		return
	}

	vol, ok := e.Data.(int)
	if !ok {
		wsSendError(e.Sender, "invalid command data format")
		return
	}

	if vol < 0 || vol > 1000 {
		wsSendError(e.Sender, "invalid value: must be in range [0, 1000]")
		return
	}

	guild, _ := discordbot.GetUsersGuildInVoice(api.session, ident.UserID)
	if guild == nil {
		wsSendError(e.Sender, "you need to be in a voice channel to perform this command")
		return
	}

	err := api.player.SetVolume(guild.ID, vol)
	if err != nil {
		wsSendError(e.Sender, fmt.Sprintf("command failed: %s", err.Error()))
	}
}
