package handlers

import (
	"github.com/zekroTJA/discordgo"
	"github.com/zekroTJA/yuri2/internal/logger"
	"github.com/zekroTJA/yuri2/internal/static"
)

type Ready struct {
}

func (h *Ready) Handler(s *discordgo.Session, e *discordgo.Ready) {
	logger.Info("session ready")
	logger.Info("Invite: https://discordapp.com/api/oauth2/authorize?client_id=%s&scope=bot&permissions=%d",
		e.User.ID, static.InvitePermission)
}
