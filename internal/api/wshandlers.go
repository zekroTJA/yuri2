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

type wsHelloData struct {
	Connected bool          `json:"connected"`
	Vol       int           `json:"vol"`
	VS        *wsVoiceState `json:"voice_state"`
}

type wsVoiceState struct {
	GuildID   string `json:"guild_id"`
	ChannelID string `json:"channel_id"`
}

// Event: INIT
func (api *API) wsInitHandler(e *wsmgr.Event) {
	data := new(wsInitData)

	// Checking if current connection was already initialized
	// by checking if the set ident is not nil.
	if e.Sender.GetIdent() != nil {
		wsSendError(e.Sender, wsErrBadCommand, "your connection is initialized already")
		return
	}

	// Parsing argument data.
	err := e.ParseDataTo(data)
	if err != nil {
		wsSendError(e.Sender, wsErrBadCommandArgs, fmt.Sprintf("failed parsing data: %s", err.Error()))
		return
	}

	// Check if the user is member of a guild the bot is also member of
	guilds := discordbot.GetUsersGuilds(api.session, data.UserID)
	if guilds == nil {
		wsSendError(e.Sender, wsErrForbidden, "you must be a member of a guild the bot is also member of")
		return
	}

	// Check authorization
	ok, _, err := api.auth.CheckAndRefresh(data.UserID, data.Token)
	if err != nil {
		wsSendError(e.Sender, wsErrInternal, fmt.Sprintf("failed checking auth: %s", err.Error()))
		return
	}
	if !ok {
		wsSendError(e.Sender, wsErrUnauthorized, "unauthorized")
		e.Sender.Close()
		return
	}

	// Create and register rate limiter for current connection
	api.wsCreateLimiter(data.UserID)

	// Setting connections Ident
	e.Sender.SetIdent(&wsIdent{
		UserID: data.UserID,
		Guilds: guilds,
	})

	// Get the users voice channel, if the user is connected
	// to a voice channel.
	// Then, the information will be added to the HELLO event payload.
	guild, _ := discordbot.GetUsersGuildInVoice(api.session, data.UserID)
	var svs *discordgo.VoiceState
	var vol int

	if guild != nil {
		svs = api.player.GetSelfVoiceState(guild.ID)
		vol, _ = api.player.GetVolume(guild.ID)
	}

	event := &wsHelloData{
		Vol:       vol,
		Connected: svs != nil,
	}

	if svs != nil {
		event.VS = &wsVoiceState{
			ChannelID: svs.ChannelID,
			GuildID:   svs.GuildID,
		}
	}

	e.Sender.Out(wsmgr.NewEvent("HELLO", event))
}

// Event: JOIN
func (api *API) wsJoinHandler(e *wsmgr.Event) {
	ident := wsCheckInitWithResponse(e.Sender)
	if ident == nil {
		return
	}

	ok, _ := api.wsCheckLimitWithResponse(e.Sender, ident.UserID)
	if !ok {
		return
	}

	_, vs := discordbot.GetUsersGuildInVoice(api.session, ident.UserID)
	if vs == nil {
		wsSendError(e.Sender, wsErrForbidden, "you need to be in a voice channel to perform this command")
		return
	}

	err := api.player.JoinVoiceCannel(vs)
	if err != nil {
		wsSendError(e.Sender, wsErrInternal, fmt.Sprintf("command failed: %s", err.Error()))
	}
}

// Event: LEAVE
func (api *API) wsLeaveHandler(e *wsmgr.Event) {
	ident := wsCheckInitWithResponse(e.Sender)
	if ident == nil {
		return
	}

	ok, _ := api.wsCheckLimitWithResponse(e.Sender, ident.UserID)
	if !ok {
		return
	}

	guild, _ := discordbot.GetUsersGuildInVoice(api.session, ident.UserID)
	if guild == nil {
		wsSendError(e.Sender, wsErrForbidden, "you need to be in a voice channel to perform this command")
		return
	}

	err := api.player.LeaveVoiceChannel(guild.ID, ident.UserID)
	if err != nil {
		wsSendError(e.Sender, wsErrInternal, fmt.Sprintf("command failed: %s", err.Error()))
	}
}

// Event: PLAY
func (api *API) wsPlayHandler(e *wsmgr.Event) {
	ident := wsCheckInitWithResponse(e.Sender)
	if ident == nil {
		return
	}

	ok, _ := api.wsCheckLimitWithResponse(e.Sender, ident.UserID)
	if !ok {
		return
	}

	guild, _ := discordbot.GetUsersGuildInVoice(api.session, ident.UserID)
	if guild == nil {
		wsSendError(e.Sender, wsErrForbidden, "you need to be in a voice channel to perform this command")
		return
	}

	user, err := api.session.User(ident.UserID)
	if err != nil || user == nil {
		wsSendError(e.Sender, wsErrInternal, fmt.Sprintf("faield getting user context: %s", err.Error()))
		return
	}

	data := new(wsPlayData)
	err = e.ParseDataTo(data)
	if err != nil {
		wsSendError(e.Sender, wsErrBadCommandArgs, fmt.Sprintf("failed parsing data: %s", err.Error()))
		return
	}

	if data.Ident == "" {
		wsSendError(e.Sender, wsErrBadCommandArgs, "ident must be a valid string value")
		return
	}

	err = api.player.Play(guild, user, data.Ident, player.ResourceType(data.Source))
	if err != nil {
		wsSendError(e.Sender, wsErrInternal, fmt.Sprintf("command failed: %s", err.Error()))
	}
}

// Event: RANDOM
func (api *API) wsRandomHandler(e *wsmgr.Event) {
	ident := wsCheckInitWithResponse(e.Sender)
	if ident == nil {
		return
	}

	ok, _ := api.wsCheckLimitWithResponse(e.Sender, ident.UserID)
	if !ok {
		return
	}

	guild, _ := discordbot.GetUsersGuildInVoice(api.session, ident.UserID)
	if guild == nil {
		wsSendError(e.Sender, wsErrForbidden, "you need to be in a voice channel to perform this command")
		return
	}

	user, err := api.session.User(ident.UserID)
	if err != nil || user == nil {
		wsSendError(e.Sender, wsErrInternal, fmt.Sprintf("faield getting user context: %s", err.Error()))
		return
	}

	err = api.player.PlayRandomSound(guild, user)
	if err != nil {
		wsSendError(e.Sender, wsErrBadCommandArgs, fmt.Sprintf("command failed: %s", err.Error()))
	}
}

// Event: VOLUME
func (api *API) wsVolumeHandler(e *wsmgr.Event) {
	ident := wsCheckInitWithResponse(e.Sender)
	if ident == nil {
		return
	}

	ok, _ := api.wsCheckLimitWithResponse(e.Sender, ident.UserID)
	if !ok {
		return
	}

	vol, ok := e.Data.(float64)
	if !ok {
		wsSendError(e.Sender, wsErrBadCommandArgs, "invalid command data format")
		return
	}

	if vol < 0 || vol > 1000 {
		wsSendError(e.Sender, wsErrBadCommandArgs, "invalid value: must be in range [0, 1000]")
		return
	}

	guild, _ := discordbot.GetUsersGuildInVoice(api.session, ident.UserID)
	if guild == nil {
		wsSendError(e.Sender, wsErrForbidden, "you need to be in a voice channel to perform this command")
		return
	}

	err := api.player.SetVolume(guild.ID, ident.UserID, int(vol))
	if err != nil {
		wsSendError(e.Sender, wsErrInternal, fmt.Sprintf("command failed: %s", err.Error()))
	}
}

// Event: STOP
func (api *API) wsStopHandler(e *wsmgr.Event) {
	ident := wsCheckInitWithResponse(e.Sender)
	if ident == nil {
		return
	}

	ok, _ := api.wsCheckLimitWithResponse(e.Sender, ident.UserID)
	if !ok {
		return
	}

	guild, _ := discordbot.GetUsersGuildInVoice(api.session, ident.UserID)
	if guild == nil {
		wsSendError(e.Sender, wsErrForbidden, "you need to be in a voice channel to perform this command")
		return
	}

	user, err := api.session.User(ident.UserID)
	if err != nil || user == nil {
		wsSendError(e.Sender, wsErrInternal, fmt.Sprintf("faield getting user context: %s", err.Error()))
		return
	}

	err = api.player.Stop(guild, user)
	if err != nil {
		wsSendError(e.Sender, wsErrBadCommandArgs, fmt.Sprintf("command failed: %s", err.Error()))
	}
}
