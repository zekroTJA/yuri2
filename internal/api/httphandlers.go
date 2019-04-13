package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/zekroTJA/yuri2/internal/discordbot"
	"github.com/zekroTJA/yuri2/internal/logger"
)

type getTokenResponse struct {
	Token  string    `json:"token"`
	UserID string    `json:"user_id"`
	Expire time.Time `json:"expires"`
}

type listResponse struct {
	N       int         `json:"n"`
	Results interface{} `json:"results"`
}

// -----------------------------------------------
// --- REST API HANDLERS

// GET /token
func (api *API) restGetTokenHandler(w http.ResponseWriter, r *http.Request, userID string) {
	guilds := discordbot.GetUsersGuilds(api.session, userID)
	if guilds == nil {
		errResponse(w, http.StatusForbidden, "you must be a member of a guild the bot is also member of")
		return
	}

	token, expire, err := api.auth.CreateToken(userID)
	if err != nil {
		errResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, &getTokenResponse{
		Token:  token,
		UserID: userID,
		Expire: expire,
	})
}

// GET /api/localsounds
func (api *API) restGetLocalSounds(w http.ResponseWriter, r *http.Request) {
	if ok, _ := api.checkAuthWithResponse(w, r); !ok {
		return
	}

	queries := r.URL.Query()

	sort := queries.Get("sort")

	_, from, err := GetURLQueryInt(queries, "from", 0)
	if err != nil {
		errResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	okLimit, limit, err := GetURLQueryInt(queries, "limit", 1)
	if err != nil {
		errResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	soundList, err := api.player.GetLocalFiles()
	if err != nil {
		errResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	switch strings.ToUpper(sort) {
	case "NAME":
		soundList.SortByName()
	case "DATE":
		soundList.SortByDate()
	}

	if okLimit {
		if limit > len(soundList)-from {
			limit = len(soundList) - from
		}
		soundList = soundList[from : from+limit]
	} else {
		soundList = soundList[from:]
	}

	jsonResponse(w, http.StatusOK, &listResponse{
		N:       len(soundList),
		Results: soundList,
	})
}

// GET /api/logs/:GUILDID
func (api *API) restGetLogs(w http.ResponseWriter, r *http.Request) {
	ok, userID := api.checkAuthWithResponse(w, r)
	if !ok {
		return
	}

	gidInd := strings.LastIndex(r.URL.Path, "/")
	if gidInd == -1 || gidInd == len(r.URL.Path)-1 {
		errResponse(w, http.StatusBadRequest, "GUILDID must be a valid snowflake value")
		return
	}

	guildID := r.URL.Path[gidInd+1:]

	queries := r.URL.Query()

	_, from, err := GetURLQueryInt(queries, "from", 0)
	if err != nil {
		errResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	okLimit, limit, err := GetURLQueryInt(queries, "limit", 1)
	if err != nil {
		errResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if !okLimit {
		limit = 1000
	}

	if _, err = api.session.GuildMember(guildID, userID); err != nil {
		errResponse(w, http.StatusForbidden, "you must be a member of this guild")
		return
	}

	logs, err := api.db.GetLogEntries(guildID, from, limit)
	if err != nil {
		errResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, &listResponse{
		N:       len(logs),
		Results: logs,
	})
}

// GET /api/stats/:GUILDID
func (api *API) restGetStats(w http.ResponseWriter, r *http.Request) {
	ok, userID := api.checkAuthWithResponse(w, r)
	if !ok {
		return
	}

	gidInd := strings.LastIndex(r.URL.Path, "/")
	if gidInd == -1 || gidInd == len(r.URL.Path)-1 {
		errResponse(w, http.StatusBadRequest, "GUILDID must be a valid snowflake value")
		return
	}

	guildID := r.URL.Path[gidInd+1:]

	queries := r.URL.Query()

	okLimit, limit, err := GetURLQueryInt(queries, "limit", 1)
	if err != nil {
		errResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if !okLimit {
		limit = 1000
	}

	if _, err = api.session.GuildMember(guildID, userID); err != nil {
		errResponse(w, http.StatusForbidden, "you must be a member of this guild")
		return
	}

	stats, err := api.db.GetSoundStats(guildID, limit)
	if err != nil {
		errResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, &listResponse{
		N:       len(stats),
		Results: stats,
	})
}

// -----------------------------------------------
// --- FE HANDLERS

func (api *API) successfullAuthHandler(w http.ResponseWriter, r *http.Request, userID string) {
	guilds := discordbot.GetUsersGuilds(api.session, userID)
	if guilds == nil {
		errPageResponse(w, r, http.StatusForbidden, "")
		return
	}

	token, _, err := api.auth.CreateToken(userID)
	if err != nil {
		errPageResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Add("Set-Cookie",
		fmt.Sprintf("token=%s; Max-Age=2147483647; Path=/", token))
	w.Header().Add("Set-Cookie",
		fmt.Sprintf("userid=%s; Max-Age=2147483647; Path=/", userID))
	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (api *API) indexPageHandler(w http.ResponseWriter, r *http.Request) {
	ok, userID, err := api.checkAuthCookie(r)
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

	guilds := discordbot.GetUsersGuilds(api.session, userID)
	if guilds == nil {
		errPageResponse(w, r, http.StatusForbidden, "")
		return
	}

	http.ServeFile(w, r, "./web/pages/index.html")
}

func (api *API) logoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Set-Cookie", "token=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/;")
	w.Header().Add("Set-Cookie", "userid=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/;")
	http.ServeFile(w, r, "./web/pages/logout.html")
}

func (api *API) wsUpgradeHandler(w http.ResponseWriter, r *http.Request) {
	_, err := api.ws.NewConn(w, r, nil)
	if err != nil {
		logger.Error("API :: wsUpgradeHandler: %s", err.Error())
		errResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}
