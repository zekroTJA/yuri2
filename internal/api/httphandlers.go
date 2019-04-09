package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/zekroTJA/yuri2/internal/logger"
)

type getTokenResponse struct {
	Token  string    `json:"token"`
	Expire time.Time `json:"expires"`
}

func (api *API) getTokenHandler(w http.ResponseWriter, r *http.Request, userID string) {
	token, expire, err := api.auth.CreateToken(userID)
	if err != nil {
		errResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, &getTokenResponse{token, expire})
}

func (api *API) successfullAuthHandler(w http.ResponseWriter, r *http.Request, userID string) {
	token, _, err := api.auth.CreateToken(userID)
	if err != nil {
		errPageResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Add("Set-Cookie",
		fmt.Sprintf("token=%s; Path=/", token))
	w.Header().Add("Set-Cookie",
		fmt.Sprintf("userid=%s; Path=/", userID))
	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (api *API) indexPageHandler(w http.ResponseWriter, r *http.Request) {
	ok, userID, err := api.checkAuthCookie(w, r)
	if err != nil {
		logger.Error("API :: checkAuthCookie: %s", err.Error())
		errPageResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	if !ok || userID == "" {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	http.ServeFile(w, r, "./web/pages/index.html")
}

func (api *API) wsUpgradeHandler(w http.ResponseWriter, r *http.Request) {
	_, err := api.ws.NewConn(w, r, nil)
	if err != nil {
		logger.Error("API :: wsUpgradeHandler: %s", err.Error())
		errResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}
