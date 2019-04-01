package player

import (
	"github.com/bwmarrin/discordgo"
	"github.com/foxbot/gavalink"
)

type Player struct {
	resURL   string
	wsURL    string
	password string
	link     *gavalink.Lavalink
}

func NewPlayer(restURL, wsURL, password string) *Player {
	return &Player{
		resURL:   restURL,
		wsURL:    wsURL,
		password: password,
	}
}

func (p *Player) Init() {

}

func (p *Player) VoiceServerUpdateHandler(s *discordgo.Session, e *discordgo.VoiceServerUpdate) {
	// if p, err := lavalink.
}
