package player

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/zekroTJA/slms/pkg/timedmap"

	"github.com/bwmarrin/discordgo"
	"github.com/foxbot/gavalink"
)

const allowedFileTypes = "mp3 wav ogg aac"

type ResourceType int

const (
	// ResourceLocal identifies a local audio file
	ResourceLocal ResourceType = iota
	// ResourceYouTube identifies a youtube video
	// resource ID
	ResourceYouTube
	// ResourceHTTP identifies a raw HTTP resource
	// file URL
	ResourceHTTP
)

const voiceStateLifetime = 3 * time.Hour

var (
	// ErrNotFound will be returned wenn no resource could
	// be fetched for the used identifier
	ErrNotFound = errors.New("no resource found for this identifier")
	// ErrNotInVoice is returned if the executor is not
	// in a voice channel on this guild
	ErrNotInVoice = errors.New("user is not in voice channel")
)

type Player struct {
	restURL  string
	wsURL    string
	password string
	fileLoc  string

	onError      func(t string, err error)
	eventHandler gavalink.EventHandler

	localSounds map[string]string

	lastVoiceStates *timedmap.TimedMap
	link            *gavalink.Lavalink
	session         *discordgo.Session
	selfVoiceState  *discordgo.VoiceState
}

func NewPlayer(restURL, wsURL, password, fileLoc string, handler gavalink.EventHandler, onError func(t string, err error)) *Player {
	if onError == nil {
		onError = func(t string, err error) {}
	}

	if handler == nil {
		handler = new(gavalink.DummyEventHandler)
	}

	return &Player{
		restURL:         restURL,
		wsURL:           wsURL,
		password:        password,
		fileLoc:         fileLoc,
		eventHandler:    handler,
		onError:         onError,
		localSounds:     make(map[string]string),
		lastVoiceStates: timedmap.New(10 * time.Minute),
	}
}

func (p *Player) Init(session *discordgo.Session) error {
	p.link = gavalink.NewLavalink("1", session.State.User.ID)
	p.session = session

	err := p.link.AddNodes(gavalink.NodeConfig{
		REST:      p.restURL,
		WebSocket: p.wsURL,
		Password:  p.password,
	})

	if err != nil {
		return err
	}

	return p.FetchLocalSounds()
}

func (p *Player) ReadyHandler(s *discordgo.Session, e *discordgo.Ready) {
	if err := p.Init(s); err != nil {
		p.onError("Ready#Init", err)
	}
}

func (p *Player) loadTrack(node *gavalink.Node, ident string) (*gavalink.Track, error) {
	tracks, err := node.LoadTracks(ident)
	if err != nil {
		return nil, err
	}

	var track *gavalink.Track

	if len(tracks.Tracks) > 0 {
		track = &tracks.Tracks[0]
	}

	return track, err
}

func (p *Player) FetchLocalSounds() error {
	files, err := ioutil.ReadDir(p.fileLoc)
	if err != nil {
		return err
	}

	p.localSounds = make(map[string]string)
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		nameSplit := strings.Split(f.Name(), ".")
		if len(nameSplit) == 1 || !strings.Contains(allowedFileTypes, nameSplit[1]) {
			continue
		}

		p.localSounds[nameSplit[0]] = fmt.Sprintf("%s/%s", p.fileLoc, f.Name())
	}

	return nil
}

func (p *Player) JoinVoiceCannel(vs *discordgo.VoiceState) error {
	if vs == nil {
		return errors.New("voiceState is nil")
	}

	// _, err := p.session.ChannelVoiceJoin(vs.GuildID, vs.ChannelID, false, false)
	// return err
	return p.session.ChannelVoiceJoinManual(vs.GuildID, vs.ChannelID, false, true)
}

func (p *Player) QuitVoiceChannel(guildID string) error {
	// This method of "quitting" a voice channel is also kind
	// of "unconventional", but as long as I do not have a
	// solution for the issue below, I need to solve the
	// problem like this.
	// https://github.com/foxbot/gavalink/issues/3

	ch, err := p.session.GuildChannelCreate(guildID, "tmp", discordgo.ChannelTypeGuildVoice)
	if err != nil {
		return err
	}

	if err = p.session.ChannelVoiceJoinManual(guildID, ch.ID, false, true); err != nil {
		return err
	}

	defer p.link.GetPlayer(guildID)

	_, err = p.session.ChannelDelete(ch.ID)
	return err
}

func (p *Player) Play(guild *discordgo.Guild, user *discordgo.User, ident string, t ResourceType) error {
	var joined bool

	node, err := p.link.BestNode()
	if err != nil {
		return err
	}

	var selfVS, userVS *discordgo.VoiceState

	for _, vs := range guild.VoiceStates {
		switch vs.UserID {
		case p.session.State.User.ID:
			selfVS = vs
		case user.ID:
			userVS = vs
		}
	}

	if userVS == nil {
		return ErrNotInVoice
	}

	if selfVS == nil || selfVS.ChannelID != userVS.ChannelID {
		if err = p.JoinVoiceCannel(userVS); err != nil {
			return err
		}
		joined = true
	}

	switch t {
	case ResourceLocal:
		var ok bool
		if ident, ok = p.localSounds[ident]; !ok {
			return ErrNotFound
		}
	}

	track, err := p.loadTrack(node, ident)
	if err != nil {
		return err
	}

	if track == nil {
		return ErrNotFound
	}

	player, err := p.link.GetPlayer(guild.ID)
	if err != nil {
		// Yes, this is quite dirty but also
		// kind of really effective.
		for i := 0; i < 10 && joined && err.Error() == "Couldn't find a player for that guild"; i++ {
			player, err = p.link.GetPlayer(guild.ID)
			if err == nil {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
		if err != nil {
			return err
		}
	}

	return player.Play(track.Data)
}
