package api

import (
	"net/http"
	"time"
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
