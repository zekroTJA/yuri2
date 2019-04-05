package player

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"

	"github.com/zekroTJA/yuri2/internal/database"

	"github.com/zekroTJA/yuri2/internal/logger"

	"github.com/zekroTJA/timedmap"

	"github.com/foxbot/gavalink"
	"github.com/zekroTJA/discordgo"
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

var resourceNames = []string{"local", "youtube", "http"}

const (
	voiceStateLifetime      = 3 * time.Hour
	autoQuitDuration        = 6 * time.Second
	fastMuteTriggerDuration = 250 * time.Millisecond
)

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
	db              database.Middleware
	selfVoiceStates map[string]*discordgo.VoiceState
}

func NewPlayer(restURL, wsURL, password, fileLoc string, db database.Middleware, handler gavalink.EventHandler, onError func(t string, err error)) *Player {
	if onError == nil {
		onError = func(t string, err error) {}
	}

	if handler == nil {
		handler = new(gavalink.DummyEventHandler)
	}

	rand.Seed(time.Now().UnixNano())

	return &Player{
		restURL:         restURL,
		wsURL:           wsURL,
		password:        password,
		fileLoc:         fileLoc,
		eventHandler:    handler,
		onError:         onError,
		db:              db,
		localSounds:     make(map[string]string),
		lastVoiceStates: timedmap.New(10 * time.Minute),
		selfVoiceStates: make(map[string]*discordgo.VoiceState),
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

func (p *Player) GetLocalFiles() (SoundFileList, error) {
	sounds := make([]*SoundFile, len(p.localSounds))

	i := 0
	for name, path := range p.localSounds {
		sf, err := NewSoundFile(name, path)
		if err != nil {
			return nil, err
		}
		sounds[i] = sf
		i++
	}

	return sounds, nil
}

func (p *Player) GetLocalSoundPath(short string) (string, bool) {
	path, ok := p.localSounds[short]
	return path, ok
}

func (p *Player) PlayRandomSound(guild *discordgo.Guild, user *discordgo.User) error {
	sounds, err := p.GetLocalFiles()
	if err != nil {
		return err
	}

	r := rand.Intn(len(sounds))
	return p.Play(guild, user, sounds[r].Name, ResourceLocal)
}

func (p *Player) JoinVoiceCannel(vs *discordgo.VoiceState) error {
	if vs == nil {
		return errors.New("voiceState is nil")
	}

	return p.session.ChannelVoiceJoinManual(vs.GuildID, vs.ChannelID, false, true)
}

func (p *Player) QuitVoiceChannel(guildID string) error {
	return p.session.ChannelVoiceJoinManual(guildID, "", false, true)
}

func (p *Player) Play(guild *discordgo.Guild, user *discordgo.User, ident string, t ResourceType) error {
	var joined bool

	switch t {
	case ResourceLocal:
		// Check if sound file actually exists
		var ok bool
		if ident, ok = p.GetLocalSoundPath(ident); !ok {
			return ErrNotFound
		}
	}

	// Check if executing user is in a voice channel
	userVS, err := p.getUsersVoiceState(guild.ID, user.ID)
	if err != nil {
		return err
	}

	if userVS == nil {
		return ErrNotInVoice
	}

	// Get player node
	node, err := p.link.BestNode()
	if err != nil {
		return err
	}

	// Check if the bot is in a voice channel.
	// if not, or if it is in another voice channel
	// than the executor, join their channel.
	selfVS := p.selfVoiceStates[guild.ID]
	if selfVS == nil || selfVS.ChannelID != userVS.ChannelID {
		if err = p.JoinVoiceCannel(userVS); err != nil {
			return err
		}
		joined = true
	}

	// Load track
	track, err := p.loadTrack(node, ident)
	if err != nil {
		return err
	}

	if track == nil {
		return ErrNotFound
	}

	// Getting guild player.
	// If the bot joined a new channel, this
	// will be repeatet 10 times with a timeout
	// delay of 100 milliseconds because it
	// can take up some time until a new player
	// wa created after establishing a voice
	// connection.
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

	logger.Debug("PLAYER :: playing sound '%s' (resource %s)", ident, resourceNames[t])

	// Actually playing the track
	return player.Play(track.Data)
}

func (p *Player) Stop(guild *discordgo.Guild, user *discordgo.User) error {
	pl, err := p.link.GetPlayer(guild.ID)
	if err != nil && err.Error() == "Couldn't find a player for that guild" {
		return nil
	}
	if err != nil {
		return err
	}

	return pl.Stop()
}
