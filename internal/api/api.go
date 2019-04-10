package api

import (
	"fmt"
	"net/http"

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

type API struct {
	cfg *config.API

	qualifiedAddress string
	trackCache       map[string]*soundTrack

	db      database.Middleware
	session *discordgo.Session
	player  *player.Player

	server *http.Server
	ws     *wsmgr.WebSocketManager
	mux    *http.ServeMux
	auth   *auth.Auth

	discordAuthFE  *discordoauth.DiscordOAuth
	discordAuthAPI *discordoauth.DiscordOAuth
}

func NewAPI(cfg *config.API, db database.Middleware, session *discordgo.Session, player *player.Player) *API {
	api := &API{
		cfg:     cfg,
		db:      db,
		session: session,
		player:  player,

		trackCache: make(map[string]*soundTrack),
	}

	protocol := "http"
	address := api.cfg.Address
	if api.cfg.TLS != nil && api.cfg.TLS.Enable {
		protocol = "https"
	}
	if address[0] == ':' {
		address = "localhost" + address
	}
	api.qualifiedAddress = fmt.Sprintf("%s://%s", protocol, address)

	api.mux = http.NewServeMux()

	api.server = &http.Server{
		Handler: api.mux,
		Addr:    api.cfg.Address,
	}

	api.ws = wsmgr.New()

	api.auth = auth.NewAuth(db, static.TokenHashRounds, static.TokenLifetime)

	api.discordAuthAPI = discordoauth.NewDiscordOAuth(
		api.cfg.ClientID,
		api.cfg.ClientSecret,
		api.qualifiedAddress+"/token/authorize",
		errResponseWrapper,
		api.restGetTokenHandler)

	api.discordAuthFE = discordoauth.NewDiscordOAuth(
		api.cfg.ClientID,
		api.cfg.ClientSecret,
		api.qualifiedAddress+"/login/authorize",
		errPageResponse,
		api.successfullAuthHandler)

	api.registerHTTPHandlers()
	api.registerWSHandlers()

	return api
}

func (api *API) registerHTTPHandlers() {
	// Static file server
	api.mux.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir("./web/static"))))

	// MAIN HANDLER
	api.mux.HandleFunc("/", api.indexPageHandler)

	// WS UPGRADE
	api.mux.HandleFunc("/ws", api.wsUpgradeHandler)

	/////////////
	// REST API

	// GET /token
	api.mux.HandleFunc("/token", api.discordAuthAPI.HandlerInit)

	// GET /token/authorize
	api.mux.HandleFunc("/token/authorize", api.discordAuthAPI.HandlerCallback)

	// GET /login
	api.mux.HandleFunc("/login", api.discordAuthFE.HandlerInit)

	// GET /login/authorize
	api.mux.HandleFunc("/login/authorize", api.discordAuthFE.HandlerCallback)

	// GET /api/localsounds
	api.mux.HandleFunc("/api/localsounds", api.restGetLocalSounds)

	// GET /api/logs/:GUILDID
	api.mux.HandleFunc("/api/logs/", api.restGetLogs)

	// GET /api/stats/:GUILDID
	api.mux.HandleFunc("/api/stats/", api.restGetStats)
}

func (api *API) registerWSHandlers() {
	// ERROR HANDLER
	api.ws.OnError(func(m string, e error) {
		logger.Error("WS :: %s: %s", m, e.Error())
	})

	// Event: INIT
	api.ws.On("INIT", api.wsInitHandler)
	// Event: PLAY
	api.ws.On("PLAY", api.wsPlayHandler)
}

func (api *API) StartBlocking() error {
	var err error

	if api.cfg.TLS != nil && api.cfg.TLS.Enable {
		err = api.server.ListenAndServeTLS(api.cfg.TLS.CertFile, api.cfg.TLS.KeyFile)
	} else {
		err = api.server.ListenAndServe()
	}

	return err
}

func (api *API) Close() {
	if api == nil {
		return
	}
	api.server.Close()
}
