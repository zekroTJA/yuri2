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

// allowedFileTypes specifies local audio file types
// which are allowed to be played
const allowedFileTypes = "mp3 wav ogg aac"

// ResourceType is the enum-like collection
// of sound resources.
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

// resourceNames contains text descriptions for
// the resource types.
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
	// ErrNoPermission is returned if the user has not the
	// specified player role
	ErrNoPermission = errors.New("insufficient permission")
	// ErrBlocked is returned if the user has the specified
	// blocking role
	ErrBlocked = errors.New("blocked of using the player")
)

// Player maintains multiple gavalink players.
type Player struct {
	restURL  string
	wsURL    string
	password string
	fileLocs []string

	playRoleName    string
	blockedRoleName string

	onError      func(t string, err error)
	eventHandler *EventHandlerManager

	localSounds map[string]string

	lastVoiceStates *timedmap.TimedMap
	link            *gavalink.Lavalink
	session         *discordgo.Session
	db              database.Middleware
	selfVoiceStates map[string]*discordgo.VoiceState
}

// NewPlayer creates a new Player.
//   restURL         : the lavalink REST API URL
//   wsURL           : the lavalink WS URL
//   fileLoc         : the pathes to local sound files
//                     (will be merged together and handled as one location)
//   playRoleName    : the Discord role name whos permitted to play sounds
//   blockedRoleName : the Discord role name who is blocked from using the player
//   db              : the database middleware to use
//   onError         : handler func to be used when errors are occuring
func NewPlayer(restURL, wsURL, password string, fileLocs []string, playRoleName, blockedRoleName string,
	db database.Middleware, onError func(t string, err error)) *Player {

	if onError == nil {
		onError = func(t string, err error) {}
	}

	rand.Seed(time.Now().UnixNano())

	return &Player{
		restURL:         restURL,
		wsURL:           wsURL,
		password:        password,
		fileLocs:        fileLocs,
		playRoleName:    playRoleName,
		blockedRoleName: blockedRoleName,
		eventHandler:    NewEventHandlerManager(),
		onError:         onError,
		db:              db,
		localSounds:     make(map[string]string),
		lastVoiceStates: timedmap.New(10 * time.Minute),
		selfVoiceStates: make(map[string]*discordgo.VoiceState),
	}
}

// Init creates a node to the specified lavalink
// connection using the sessions state user ID.
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

// AddEventHandler adds a struct described
// in the EventHandler interface.
func (p *Player) AddEventHandler(handler EventHandler) {
	p.eventHandler.AddHandler(handler)
}

// ReadyHandler is the handler set for the bot sessions
// ready event wich will initialize the player.
func (p *Player) ReadyHandler(s *discordgo.Session, e *discordgo.Ready) {
	if err := p.Init(s); err != nil {
		p.onError("Ready#Init", err)
	}
}

// loadTrack requests the lavalink node to load
// a track by ident and getting the returned track
// ID and information.
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

// FetchLocalSounds updates the local sounds
// list by reading all local sounds in the
// specified location matching the specified
// audio file types.
func (p *Player) FetchLocalSounds() error {
	if len(p.fileLocs) == 0 {
		return errors.New("sound file locations can not be empty")
	}

	p.localSounds = make(map[string]string)

	for _, floc := range p.fileLocs {
		files, err := ioutil.ReadDir(floc)
		if err != nil {
			return err
		}

		loaded := 0
		for _, f := range files {
			if f.IsDir() {
				continue
			}

			nameSplit := strings.Split(f.Name(), ".")
			if len(nameSplit) == 1 || !strings.Contains(allowedFileTypes, nameSplit[1]) {
				continue
			}

			p.localSounds[nameSplit[0]] = fmt.Sprintf("%s/%s", floc, f.Name())
			loaded++
		}

		logger.Info("PLAYER :: loaded %d sounds from %s", loaded, floc)
	}

	return nil
}

// GetLocalFiles returns the list of local files.
// this function does not Refetch the local file
// list because of performance reasons.
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

// GetLocalSoundPath returns the full path to
// a local sound file by its ident.
func (p *Player) GetLocalSoundPath(short string) (string, bool) {
	path, ok := p.localSounds[short]
	return path, ok
}

// PlayRandomSound plays a random local sound.
// If the bot is not connected to a voice channel on the
// guild the comamnd was executed, the bot will join this
// channel and then play the randomly chosen sound.
func (p *Player) PlayRandomSound(guild *discordgo.Guild, user *discordgo.User) error {
	sounds, err := p.GetLocalFiles()
	if err != nil {
		return err
	}

	r := rand.Intn(len(sounds))
	return p.Play(guild, user, sounds[r].Name, ResourceLocal)
}

// JoinVoiceCannel joins the voice channel by passed
// VoiceState.
func (p *Player) JoinVoiceCannel(vs *discordgo.VoiceState) error {
	if err := p.checkPermsByIDs(vs.GuildID, vs.UserID); err != nil {
		return err
	}

	if vs == nil {
		return errors.New("voiceState is nil")
	}

	err := p.session.ChannelVoiceJoinManual(vs.GuildID, vs.ChannelID, false, true)
	if err != nil {
		return err
	}

	p.eventHandler.OnVoiceJoined(vs.GuildID, vs.ChannelID)

	return nil
}

// LeaveVoiceChannel leaves the voice channel the
// bot is connected to on the specified guild.
func (p *Player) LeaveVoiceChannel(guildID, userID string) error {
	if err := p.checkPermsByIDs(guildID, userID); err != nil {
		return err
	}

	err := p.session.ChannelVoiceJoinManual(guildID, "", false, true)
	if err != nil {
		return err
	}

	var channelID string
	if vs, ok := p.selfVoiceStates[guildID]; ok {
		channelID = vs.ChannelID
	}

	p.eventHandler.OnVoiceLeft(guildID, channelID)

	return nil
}

// Play plays the specified sound by passed ident and resource.
// If the bot is not connected to a voice channel on the
// guild the comamnd was executed, the bot will join this
// channel and then play the sound.
func (p *Player) Play(guild *discordgo.Guild, user *discordgo.User, ident string, t ResourceType) error {
	if err := p.checkPerms(guild, user.ID); err != nil {
		return err
	}

	var joined bool
	originalIdent := ident

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

	// Fire payer track start event
	p.eventHandler.OnTrackStart(player, track.Data, originalIdent, t, guild.ID,
		userVS.ChannelID, user.ID, user.String())

	logger.Debug("PLAYER :: playing sound '%s' (resource %s)", ident, resourceNames[t])

	// Actually playing the track
	err = player.Play(track.Data)
	if err != nil {
		return err
	}

	// Add play action to sounds log
	err = p.db.AddLogEntry(&database.SoundLogEntry{
		GuildID: guild.ID,
		UserID:  user.ID,
		UserTag: user.String(),
		Source:  resourceNames[t],
		Sound:   originalIdent,
	})

	if t == ResourceLocal {
		err = p.db.AddSoundStatsCount(guild.ID, originalIdent)
	}

	return err
}

// Stop stops a playing sound.
func (p *Player) Stop(guild *discordgo.Guild, user *discordgo.User) error {
	if err := p.checkPerms(guild, user.ID); err != nil {
		return err
	}

	pl, err := p.link.GetPlayer(guild.ID)
	if err != nil && err.Error() == "Couldn't find a player for that guild" {
		return nil
	}
	if err != nil {
		return err
	}

	return pl.Stop()
}

// SetVolume sets the volume for the current guilds
// player. The volume will also be saved and set up
// as volume for all following players created on
// the guild.
func (p *Player) SetVolume(guildID, userID string, vol int) error {
	if err := p.checkPermsByIDs(guildID, userID); err != nil {
		return err
	}

	pl, err := p.link.GetPlayer(guildID)
	if err != nil {
		return nil
	}

	if err = pl.Volume(vol); err != nil {
		return err
	}

	p.eventHandler.OnVolumeChanged(pl, guildID, vol)

	return p.db.SetGuildVolume(guildID, vol)
}

// GetVolume returns the volume of the player
// for the specified guild.
func (p *Player) GetVolume(guildID string) (int, error) {
	pl, err := p.link.GetPlayer(guildID)
	if err != nil {
		return 0, nil
	}

	return pl.GetVolume(), nil
}

// GetSelfVoiceState returns the current voice state on
// the specified guild. This will return nil if the bot
// is not in any voice channel on this guild.
func (p *Player) GetSelfVoiceState(guildID string) *discordgo.VoiceState {
	return p.selfVoiceStates[guildID]
}
