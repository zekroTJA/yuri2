package api

import "github.com/foxbot/gavalink"

func (api *API) OnTrackEnd(player *gavalink.Player, track string, reason string) error {
	return nil
}

func (api *API) OnTrackException(player *gavalink.Player, track string, reason string) error {
	return nil
}

func (api *API) OnTrackStuck(player *gavalink.Player, track string, threshold int) error {
	return nil
}
