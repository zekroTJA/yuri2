package api

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/zekroTJA/timedmap"

	"github.com/zekroTJA/discordgo"

	"github.com/zekroTJA/yuri2/internal/api/auth"
	"github.com/zekroTJA/yuri2/internal/config"
	"github.com/zekroTJA/yuri2/internal/database"
	"github.com/zekroTJA/yuri2/internal/logger"
	"github.com/zekroTJA/yuri2/internal/player"
	"github.com/zekroTJA/yuri2/internal/static"
	"github.com/zekroTJA/yuri2/pkg/discordoauth"
	"github.com/zekroTJA/yuri2/pkg/wsmgr"
)

const (
	limitsCleanupInterval = 5 * time.Minute
	limitsLifetime        = 1 * time.Hour

	wsLimit = 1000 * time.Millisecond
	wsBurst = 10

	restLimit = 1 * time.Second
	restBurst = 10
)

// API maintains the HTTP web server, REST
// API and WS API.
type API struct {
	cfg *config.Main

	qualifiedAddress string
	trackCache       map[string]*soundTrack

	db      database.Middleware
	session *discordgo.Session
	player  *player.Player

	teardownChan chan os.Signal

	server *http.Server
	ws     *wsmgr.WebSocketManager
	mux    *http.ServeMux
	auth   *auth.Auth

	discordAuthFE  *discordoauth.DiscordOAuth
	discordAuthAPI *discordoauth.DiscordOAuth

	limits *timedmap.TimedMap
}

// NewAPI initializes a new API and registers handlers
// for web server, REST API and WS API endpoints.
func NewAPI(cfg *config.Main, db database.Middleware, session *discordgo.Session, player *player.Player, teardownChan chan os.Signal) *API {
	// init API object
	api := &API{
		cfg:          cfg,
		db:           db,
		session:      session,
		player:       player,
		teardownChan: teardownChan,

		trackCache: make(map[string]*soundTrack),
		limits:     timedmap.New(limitsCleanupInterval),
	}

	api.qualifiedAddress = cfg.API.PublicAddress
	if !strings.HasPrefix(api.qualifiedAddress, "http") {
		protocol := "http"
		if cfg.API.TLS != nil && cfg.API.TLS.Enable {
			protocol += "s"
		}
		api.qualifiedAddress = fmt.Sprintf("%s://%s", protocol, api.qualifiedAddress)
	}

	// Initialize URL path mux
	api.mux = http.NewServeMux()

	// Initialize HTTP server
	api.server = &http.Server{
		Handler: api.mux,
		Addr:    api.cfg.API.Address,
	}

	// Initialize web socket manager
	api.ws = wsmgr.New()

	// Initialize Auth manager
	api.auth = auth.NewAuth(db, static.TokenHashRounds, static.TokenLifetime)

	// Create Discord OAuth Router for API token
	// request
	api.discordAuthAPI = discordoauth.NewDiscordOAuth(
		api.cfg.API.ClientID,
		api.cfg.API.ClientSecret,
		api.qualifiedAddress+"/token/authorize",
		errResponseWrapper,
		api.restGetTokenHandler)

	// Create Discord OAuth Router for user
	// interface login
	api.discordAuthFE = discordoauth.NewDiscordOAuth(
		api.cfg.API.ClientID,
		api.cfg.API.ClientSecret,
		api.qualifiedAddress+"/login/authorize",
		errPageResponse,
		api.successfullAuthHandler)

	// register HTTP handlers
	api.registerHTTPHandlers()
	// register WS handlers
	api.registerWSHandlers()

	return api
}

// registerHTTPHandlers registers HTTP request handlers
// for specific endpoint paths to the ServeMux.
func (api *API) registerHTTPHandlers() {
	// Static file server
	api.mux.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir("./web/static"))))

	// MAIN HANDLER
	api.mux.HandleFunc("/", api.indexPageHandler)

	// ADMIN HANDLER
	api.mux.HandleFunc("/admin", api.adminPageHandler)

	// WS UPGRADE
	api.mux.HandleFunc("/ws", api.wsUpgradeHandler)

	// GET /login
	api.mux.HandleFunc("/login", api.discordAuthFE.HandlerInit)

	// GET /login/authorize
	api.mux.HandleFunc("/login/authorize", api.discordAuthFE.HandlerCallback)

	// GET /logout
	api.mux.HandleFunc("/logout", api.logoutHandler)

	/////////////
	// REST API

	// GET /token
	api.mux.HandleFunc("/token", api.discordAuthAPI.HandlerInit)

	// GET /token/authorize
	api.mux.HandleFunc("/token/authorize", api.discordAuthAPI.HandlerCallback)

	// GET /api/localsounds
	api.mux.HandleFunc("/api/localsounds", api.restGetLocalSounds)

	// GET /api/logs/:GUILDID
	api.mux.HandleFunc("/api/logs/", api.restGetLogs)

	// GET /api/stats/:GUILDID
	api.mux.HandleFunc("/api/stats/", api.restGetStats)

	// GET /api/favorites
	api.mux.HandleFunc("/api/favorites", api.restGetFavorites)

	// POST /api/favorites/:SOUND
	// DELETE /api/favorites/:SOUND
	api.mux.HandleFunc("/api/favorites/", api.restPostDeleteFavorites)

	// GET /api/admin/stats
	api.mux.HandleFunc("/api/admin/stats", api.restGetAdminStats)

	// GET /api/admin/soundstats
	api.mux.HandleFunc("/api/admin/soundstats", api.restGetAdminSoundStats)

	// POST /api/admin/restart
	api.mux.HandleFunc("/api/admin/restart", api.restPostAdminRestart)

	// POST /api/admin/refetch
	api.mux.HandleFunc("/api/admin/refetch", api.restPostAdminRefetch)
}

// registerWSHandlers registers WS handlers
// for specific WS commands to the WS manager.
func (api *API) registerWSHandlers() {
	// ERROR HANDLER
	api.ws.OnError(func(m string, e error) {
		logger.Error("WS :: %s: %s", m, e.Error())
	})

	// Event: INIT
	api.ws.On("INIT", api.wsInitHandler)

	// Event: JOIN
	api.ws.On("JOIN", api.wsJoinHandler)

	// Event: LEAVE
	api.ws.On("LEAVE", api.wsLeaveHandler)

	// Event: PLAY
	api.ws.On("PLAY", api.wsPlayHandler)

	// Event: RANDOM
	api.ws.On("RANDOM", api.wsRandomHandler)

	// Event: VOLUME
	api.ws.On("VOLUME", api.wsVolumeHandler)

	// Event: STOP
	api.ws.On("STOP", api.wsStopHandler)
}

// StartBlocking starts the HTTP server
// wther in TLS or non-TLS mode, depending
// on configuration and blocks the current
// go routine waiting for incomming requests.
func (api *API) StartBlocking() error {
	var err error

	if api.cfg.API.TLS != nil && api.cfg.API.TLS.Enable {
		err = api.server.ListenAndServeTLS(api.cfg.API.TLS.CertFile, api.cfg.API.TLS.KeyFile)
	} else {
		err = api.server.ListenAndServe()
	}

	return err
}

// Close cleanly shuts down the server.
// This will not panic if the api instance
// is nil because of failed initialization.
func (api *API) Close() {
	if api == nil {
		return
	}
	api.server.Close()
}
