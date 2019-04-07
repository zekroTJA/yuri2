package api

import (
	"fmt"
	"net/http"

	"github.com/zekroTJA/discordgo"
	"github.com/zekroTJA/yuri2/internal/database"

	"github.com/zekroTJA/yuri2/internal/config"
)

type API struct {
	cfg *config.API

	authRedirectURI string

	db      *database.Middleware
	session *discordgo.Session

	server *http.Server
	mux    *http.ServeMux
}

func NewAPI(cfg *config.API, db *database.Middleware, session *discordgo.Session) *API {
	api := &API{
		cfg:     cfg,
		db:      db,
		session: session,
	}

	protocol := "http"
	address := api.cfg.Address
	if api.cfg.TLS != nil && api.cfg.TLS.Enable {
		protocol = "https"
	}
	if address[0] == ':' {
		address = "localhost" + address
	}
	api.authRedirectURI = fmt.Sprintf("%s://%s", protocol, address)

	api.mux = http.NewServeMux()

	api.server = &http.Server{
		Handler: api.mux,
		Addr:    api.cfg.Address,
	}

	api.registerHandlers()

	return api
}

func (api *API) registerHandlers() {
	// GET /
	api.mux.HandleFunc("/", api.handlerRedirectToLogin)
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
