package player

import (
	"github.com/foxbot/gavalink"
)

type EventHandler struct {
	handler []gavalink.EventHandler
}

func NewEventHandler() *EventHandler {
	return &EventHandler{
		handler: make([]gavalink.EventHandler, 0),
	}
}

func (h *EventHandler) OnTrackEnd(player *gavalink.Player, track string, reason string) error {
	for _, h := range h.handler {
		h.OnTrackEnd(player, track, reason)
	}

	return nil
}

func (h *EventHandler) OnTrackException(player *gavalink.Player, track string, reason string) error {
	for _, h := range h.handler {
		h.OnTrackException(player, track, reason)
	}

	return nil
}

func (h *EventHandler) OnTrackStuck(player *gavalink.Player, track string, threshold int) error {
	for _, h := range h.handler {
		h.OnTrackStuck(player, track, threshold)
	}

	return nil
}

func (h *EventHandler) AddHandler(handler gavalink.EventHandler) {
	h.handler = append(h.handler, handler)
}
