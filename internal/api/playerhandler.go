package api

import (
	"github.com/foxbot/gavalink"
	"github.com/zekroTJA/yuri2/internal/logger"
	"github.com/zekroTJA/yuri2/internal/player"
	"github.com/zekroTJA/yuri2/pkg/wsmgr"
)

type soundTrack struct {
	Ident     string              `json:"ident,omitempty"`
	Source    player.ResourceType `json:"source"`
	GuildID   string              `json:"guild_id,omitempty"`
	ChannelID string              `json:"channel_id,omitempty"`
	UserID    string              `json:"user_id,omitempty"`
	UserTag   string              `json:"user_tag,omitempty"`
}

type wsPlayExceptionData struct {
	Reason string      `json:"reason"`
	Track  *soundTrack `json:"track,omitempty"`
}

type wsPlayStuckData struct {
	Threshold int         `json:"threshold"`
	Track     *soundTrack `json:"track,omitempty"`
}

type wsVolumeChangedData struct {
	Vol     int    `json:"vol"`
	GuildID string `json:"guild_id"`
}

type wsGuildChannelData struct {
	GuildID   string `json:"guild_id"`
	ChannelID string `json:"channel_id"`
}

// OnTrackStart is the handler for PLAYING event
func (api *API) OnTrackStart(player *gavalink.Player, track, ident string,
	source player.ResourceType, guildID, channelID, userID, userTag string) {

	s := &soundTrack{
		Ident:     ident,
		Source:    source,
		GuildID:   guildID,
		ChannelID: channelID,
		UserID:    userID,
		UserTag:   userTag,
	}

	// The saved track ID is shortened by 5 characters because
	// the IDs are exactly the same until the last 3 to 5 chars
	// for some reason.
	api.trackCache[track[:len(track)-5]] = s

	logger.Debug("API :: PLAYER HANDLER :: track start event: %s", ident)

	cond := condFactory(guildID)

	if err := api.ws.BroadcastExclusive(wsmgr.NewEvent("PLAYING", s), cond); err != nil {
		logger.Error("WS :: PLAYER HANDLER :: failed broadcasting: %s", err.Error())
	}
}

// OnTrackEnd is the handler for END event
func (api *API) OnTrackEnd(player *gavalink.Player, track string, reason string) error {
	defer delete(api.trackCache, track)

	logger.Debug("API :: PLAYER HANDLER :: track end event")

	s, ok := api.trackCache[track[:len(track)-5]]
	if ok {
		cond := condFactory(s.GuildID)
		if err := api.ws.BroadcastExclusive(wsmgr.NewEvent("END", s), cond); err != nil {
			logger.Error("WS :: PLAYER HANDLER :: failed broadcasting: %s", err.Error())
		}
	}

	return nil
}

// OnTrackException is the handler for PLAY_ERROR event
func (api *API) OnTrackException(player *gavalink.Player, track string, reason string) error {
	defer delete(api.trackCache, track)

	logger.Debug("API :: PLAYER HANDLER :: track exception: %s", reason)

	s, ok := api.trackCache[track[:len(track)-5]]

	e := &wsPlayExceptionData{
		Reason: reason,
		Track:  s,
	}

	if ok {
		cond := condFactory(s.GuildID)
		if err := api.ws.BroadcastExclusive(wsmgr.NewEvent("PLAY_ERROR", e), cond); err != nil {
			logger.Error("WS :: PLAYER HANDLER :: failed broadcasting: %s", err.Error())
		}
	}

	return nil
}

// OnTrackStuck is the handler for STUCK event
func (api *API) OnTrackStuck(player *gavalink.Player, track string, threshold int) error {
	logger.Debug("API :: PLAYER HANDLER :: track stuck event")

	s, ok := api.trackCache[track[:len(track)-5]]

	e := &wsPlayStuckData{
		Threshold: threshold,
		Track:     s,
	}

	if ok {
		cond := condFactory(s.GuildID)
		if err := api.ws.BroadcastExclusive(wsmgr.NewEvent("STUCK", e), cond); err != nil {
			logger.Error("WS :: PLAYER HANDLER :: failed broadcasting: %s", err.Error())
		}
	}

	return nil
}

// OnVolumeChanged is the handler for VOLUME_CHANGED event
func (api *API) OnVolumeChanged(player *gavalink.Player, guildID string, vol int) {
	logger.Debug("API :: PLAYER HANDLER :: volume changed event")

	e := &wsVolumeChangedData{
		GuildID: guildID,
		Vol:     vol,
	}

	cond := condFactory(guildID)
	if err := api.ws.BroadcastExclusive(wsmgr.NewEvent("VOLUME_CHANGED", e), cond); err != nil {
		logger.Error("WS :: PLAYER HANDLER :: failed broadcasting: %s", err.Error())
	}
}

// OnJVoiceJoined is the handler for JOINED event
func (api *API) OnVoiceJoined(guildID, channelID string) {
	logger.Debug("API :: PLAYER HANDLER :: voice joined event")

	e := &wsGuildChannelData{
		GuildID:   guildID,
		ChannelID: channelID,
	}

	cond := condFactory(guildID)
	if err := api.ws.BroadcastExclusive(wsmgr.NewEvent("JOINED", e), cond); err != nil {
		logger.Error("WS :: PLAYER HANDLER :: failed broadcasting: %s", err.Error())
	}
}

// OnVoiceLeft is the handler for LEFT event
func (api *API) OnVoiceLeft(guildID, channelID string) {
	logger.Debug("API :: PLAYER HANDLER :: voice left event")

	e := &wsGuildChannelData{
		GuildID:   guildID,
		ChannelID: channelID,
	}

	cond := condFactory(guildID)
	if err := api.ws.BroadcastExclusive(wsmgr.NewEvent("LEFT", e), cond); err != nil {
		logger.Error("WS :: PLAYER HANDLER :: failed broadcasting: %s", err.Error())
	}
}

// ------------------------------------

// condFactory creates a filter function for BroadcastExclusive
// to just send events to connections which are on the same
// guild the event was fired from.
func condFactory(guildID string) func(c *wsmgr.WebSocketConn) bool {
	return func(c *wsmgr.WebSocketConn) bool {
		ident, _ := c.GetIdent().(*wsIdent)
		if ident == nil {
			return false
		}

		for _, g := range ident.Guilds {
			if g.ID == guildID {
				return true
			}
		}

		return false
	}
}
