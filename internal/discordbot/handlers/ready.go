package handlers

import (
	"time"

	"github.com/zekroTJA/discordgo"
	"github.com/zekroTJA/yuri2/internal/config"
	"github.com/zekroTJA/yuri2/internal/logger"
	"github.com/zekroTJA/yuri2/internal/static"
)

// Ready handles a Discord WS ready event.
type Ready struct {
	status []string
	delay  time.Duration

	s *discordgo.Session

	currInd int
	t       *time.Ticker
}

// NewReady cteates a new instance of Ready.
func NewReady(cfg *config.StatusShuffle) *Ready {
	if len(cfg.Status) == 0 {
		logger.Fatal("DBOT :: failed initializing ready handler: status must contain at least one string")
	}
	delay, err := time.ParseDuration(cfg.Delay)
	if err != nil {
		logger.Fatal("DBOT :: failed initializing ready handler: %s", err.Error())
	}

	return &Ready{
		status: cfg.Status,
		delay:  delay,
	}
}

// Handler is the actual handler function
// executes when the event was fired.
func (h *Ready) Handler(s *discordgo.Session, e *discordgo.Ready) {
	logger.Info("session ready")
	logger.Info("Invite: https://discordapp.com/api/oauth2/authorize?client_id=%s&scope=bot&permissions=%d",
		e.User.ID, static.InvitePermission)

	h.s = s
	h.updateStatus(h.status[0])

	if len(h.status) > 1 {
		h.currInd = 1
		h.t = time.NewTicker(h.delay)
		go h.timerLoopBlocking()
	}
}

// updateStatus updates the presence status
// of the Discord bots client.
func (h *Ready) updateStatus(status string) error {
	return h.s.UpdateStatus(0, status)
}

// timerLoopBlocking starts a timer loop
// blocking the current go routine.
// Everytime the loop evaluates, the
// status will shuffle through the defined
// status messages in the configuration.
func (h *Ready) timerLoopBlocking() {
	var err error
	for {
		<-h.t.C // # no ad intended

		if h.currInd >= len(h.status) {
			h.currInd = 0
		}

		err = h.updateStatus(h.status[h.currInd])
		if err != nil {
			logger.Error("DBOT :: READY HANDLER :: updateStatus failed: %s", err.Error())
			h.t.Stop()
		}

		h.currInd++
	}
}
