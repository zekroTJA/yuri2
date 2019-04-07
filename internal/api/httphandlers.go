package api

import (
	"fmt"
	"net/http"
)

func (api *API) handlerRedirectToLogin(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("https://discordapp.com/api/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=identify",
		api.cfg.ClientID, api.authRedirectURI)
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
