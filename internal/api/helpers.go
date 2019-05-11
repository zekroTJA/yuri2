package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/zekroTJA/ratelimit"

	"github.com/zekroTJA/yuri2/pkg/wsmgr"
)

// a wsErrorType is the Enum-like
// colelction of web socket error
// codes.
type wsErrorType int

const (
	wsErrBadCommandArgs wsErrorType = iota
	wsErrUnauthorized
	wsErrForbidden
	wsErrInternal
	wsErrBadCommand
	wsErrRateLimitExceed
)

var wsErrTypeStr = []string{
	"bad command args",
	"unatuhorized",
	"forbidden",
	"internal",
	"bad command",
	"rate limit exceed",
}

var stdErrMsgs = map[int]string{
	400: "bad request",
	401: "unauthorized",
	403: "forbidden",
	429: "rate limit exceed",
}

// apIErrorBody contains data for a REST
// API error response
type apiErrorBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// wsErrorData contains data for a WS
// API error event
type wsErrorData struct {
	Code    wsErrorType `json:"code"`
	Type    string      `json:"type"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// errResponse writes a JSON formated error response
// to the passed ResponseWriter with the specified
// code and msg. If msg is empty, a default message
// from stdErrMsgs will be used, if available.
func errResponse(w http.ResponseWriter, code int, msg string) {
	if msg == "" {
		msg = stdErrMsgs[code]
	}

	data := map[string]interface{}{
		"error": apiErrorBody{
			Code:    code,
			Message: msg,
		},
	}

	jsonResponse(w, code, data)
}

// errResponseWrapper is a wrapper function which makes
// errResponse compatible with DiscordOAuth.
func errResponseWrapper(w http.ResponseWriter, r *http.Request, code int, msg string) {
	errResponse(w, code, msg)
}

// parseJSONBody ready a requests body and tries
// it to parse this into the passed object instance
// using a JSON decoder.
func parseJSONBody(body io.ReadCloser, v interface{}) error {
	dec := json.NewDecoder(body)
	return dec.Decode(v)
}

// jsonResponse parses the passed data to a JSON
// byte array which will be written as response to
// the passed responseWriter with the specified
// status code.
func jsonResponse(w http.ResponseWriter, code int, data interface{}) {
	var bData []byte
	var err error

	w.Header().Add("Content-Type", "application/json")

	if data != nil {
		bData, err = json.MarshalIndent(data, "", "  ")
	} else {
		bData = []byte("{}")
	}

	if err != nil {
		return
	}

	w.WriteHeader(code)
	_, err = w.Write(bData)
}

// errPageResponse tries to serve a error HTML page
// located in ./web/pages/errors with the specific
// error status code as name and serves it to
// the passed ResponseWriter.
// msg argument will be ignored and is just passable
// because of compatibility with DiscordOAuth.
func errPageResponse(w http.ResponseWriter, r *http.Request, code int, msg string) {
	pageLoc := fmt.Sprintf("./web/pages/errors/%d.html", code)
	http.ServeFile(w, r, pageLoc)
}

// getCookeValue tries to get a string value from
// a cookie passed from the request. If no cookie
// was passed with the specified name, an empty
// stirng in combination with an nil error will be
// returned. Err is only not nil if the reading
// of the cookie failed for some reason.
func getCookieValue(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil || cookie == nil {
		if err == http.ErrNoCookie {
			err = nil
		}
		return "", err
	}

	return cookie.Value, nil
}

// checkAuthCookie tries to get the userID and token
// from passed cookies and checks their validity.
// The return values are a bool, which is true if the
// check was successful without exceptions, the userID
// and an error, which will be not nil if the cookie
// reading or database access fails.
func (api *API) checkAuthCookie(r *http.Request) (bool, string, error) {
	token, err := getCookieValue(r, "token")
	if err != nil {
		return false, "", err
	}

	userID, err := getCookieValue(r, "userid")
	if err != nil {
		return false, "", err
	}

	ok, _, err := api.auth.CheckAndRefresh(userID, token)

	return ok, userID, err
}

// checkAuthCookie checks the request for an 'Authorization'
// header, tries to parse the value, if existent and checks
// the validity of the token.
// The return values are a bool, which is true if the check
// was successful without exceptions, the userID of the
// authenticated user and an error, which will be not nil
// if something fails during cookie reading and database
// access.
func (api *API) checkAuthHeader(r *http.Request) (bool, string, error) {
	headerVal := r.Header.Get("Authorization")
	if headerVal == "" {
		return false, "", nil
	}

	if !strings.HasPrefix(strings.ToLower(headerVal), "basic ") {
		return false, "", nil
	}

	bRawVal, err := base64.StdEncoding.DecodeString(headerVal[6:])
	if err != nil {
		return false, "", err
	}

	var userID, token string

	if split := strings.Split(string(bRawVal), ":"); len(split) == 2 {
		userID = split[0]
		token = split[1]
	} else {
		return false, "", nil
	}

	ok, _, err := api.auth.CheckAndRefresh(userID, token)

	return ok, userID, err
}

// checkAuthWithResponse first checks for a valid 'Authorization'
// header value using checkAuthHeader. If this was not successful,
// the cookies will be checked for valid authorization values.
// If one of both fails because of an unexpected error, this
// will be written to the ResponseWriter.
// If the authorization fails, this will be written to the
// ResponseWriter as unauthorized error.
// The returned bool is true if the authorization was passed
// successfully and the userID will be returned as well.
func (api *API) checkAuthWithResponse(w http.ResponseWriter, r *http.Request) (bool, string) {
	ok, userID, err := api.checkAuthHeader(r)
	if err != nil {
		if strings.HasPrefix(err.Error(), "illegal base64 data") {
			errResponse(w, http.StatusBadRequest, err.Error())
			return false, ""
		}
		errResponse(w, http.StatusInternalServerError, err.Error())
		return false, ""
	}
	if !ok {
		ok, userID, err = api.checkAuthCookie(r)
		if err != nil {
			errResponse(w, http.StatusInternalServerError, err.Error())
			return false, ""
		}

		if !ok || userID == "" {
			errResponse(w, http.StatusUnauthorized, "")
			return false, ""
		}
	}

	return true, userID
}

// checkMethodWithResponse checks if the used method is
// one of the passed allowed methods. Else, attach "Allow"
// header with allowed request methods and an 405 status.
// Returns false if the check did not pass.
func checkMethodWithResponse(w http.ResponseWriter, r *http.Request, method ...string) bool {
	for _, m := range method {
		if m == r.Method {
			return true
		}
	}

	w.Header().Set("Allow", strings.Join(method, ", "))
	errResponse(w, http.StatusMethodNotAllowed, "method not allowed")
	return false
}

// wsSendError creates a wsErrorData object from passed error
// code and message and sends it to the specified connection.
func wsSendError(wsc *wsmgr.WebSocketConn, code wsErrorType, msg string, data ...interface{}) error {
	e := &wsErrorData{
		Code:    code,
		Type:    wsErrTypeStr[code],
		Message: msg,
	}
	if len(data) > 0 {
		e.Data = data[0]
	}
	return wsc.Out(wsmgr.NewEvent("ERROR", e))
}

// wsCheckInitWithResponse checks if the connection was successfully
// initialized by an `INIT` event, which writes userID and
// the users guilds to the ident property of the connection.
func wsCheckInitWithResponse(wsc *wsmgr.WebSocketConn) *wsIdent {
	ident, ok := wsc.GetIdent().(*wsIdent)
	if !ok || ident == nil {
		wsSendError(wsc, wsErrUnauthorized, "unauthorized")
		wsc.Close()
		return nil
	}

	return ident
}

// GetURLQueryInt tries to get an parsed integer value from
// URL queries by key name. If you pass a number after key,
// an error will be returned when the value is smaller than
// the given value. If you are passing a second number and
// the value is smaller than the first number and larger than
// the second number, an error will be returned.
func GetURLQueryInt(queries url.Values, key string, rng ...int) (bool, int, error) {
	if s := queries.Get(key); s != "" {

		val, err := strconv.Atoi(s)
		if err != nil {
			return false, 0, err
		}

		switch len(rng) {
		case 1:
			if val < rng[0] {
				return false, 0, fmt.Errorf("%s must be a valid number larger than %d", key, rng[0]-1)
			}
		case 2:
			if val < rng[0] || val > rng[1] {
				return false, 0, fmt.Errorf("%s must be a valid number in range [%d, %d]", key, rng[0], rng[1])
			}
		}

		return true, val, nil
	}

	return false, 0, nil
}

// wsCreateLimiter creates a WS limiter with the configured
// parameters for the user ID.
func (api *API) wsCreateLimiter(userID string) *ratelimit.Limiter {
	limiter := ratelimit.NewLimiter(wsLimit, wsBurst)
	api.limits.Set("ws"+userID, limiter, limitsLifetime)
	return limiter
}

// wsCheckLimitWithResponse tries to get a token from the users limiter.
// The success of this will be returned as boolean and a reservation
// object will be returned containing the current ratelimiters status.
// If the limit was exceed, this will be send to the client as ERROR
// event.
func (api *API) wsCheckLimitWithResponse(wsc *wsmgr.WebSocketConn, userID string) (bool, *ratelimit.Reservation) {
	var limiter *ratelimit.Limiter

	limiter, _ = api.limits.GetValue("ws" + userID).(*ratelimit.Limiter)
	if limiter == nil {
		limiter = api.wsCreateLimiter(userID)
	}

	ok, res := limiter.Reserve()
	if !ok {
		wsSendError(wsc, wsErrRateLimitExceed, "Rate limit exceed. Wait reset_time * milliseconds until sending another command.", map[string]int64{
			"reset_time": time.Until(res.Reset.Time).Nanoseconds() / 1000000,
		})
	}

	return ok, res
}

// createLimiter creates a rest limiter with the configured
// parameters for the user ID.
func (api *API) createLimiter(ident string) *ratelimit.Limiter {
	limiter := ratelimit.NewLimiter(restLimit, restBurst)
	api.limits.Set("rest"+ident, limiter, limitsLifetime)
	return limiter
}

// checkLimitWithResponse tries to get a token from the users limiter.
// The success of this will be returned as boolean and a reservation
// object will be returned containing the current ratelimiters status.
// The current limiter status will be set as "X-RateLimit-Limit",
// "X-RateLimit-Remaining" and "X-RateLimit-Reset" header.
// If the limit was exceed, this will be send to the client as 429 error
// response.
func (api *API) checkLimitWithResponse(w http.ResponseWriter, ident string) (bool, *ratelimit.Reservation) {
	var limiter *ratelimit.Limiter

	limiter, _ = api.limits.GetValue("rest" + ident).(*ratelimit.Limiter)
	if limiter == nil {
		limiter = api.createLimiter(ident)
	}

	ok, res := limiter.Reserve()

	var reset int64
	if !res.Reset.IsNil() {
		reset = time.Until(res.Reset.Time).Nanoseconds() / 1000000
	}

	h := w.Header()
	h.Set("X-RateLimit-Limit", fmt.Sprintf("%d", res.Burst))
	h.Set("X-RateLimit-Remaining", fmt.Sprintf("%d", res.Remaining))
	h.Set("X-RateLimit-Reset", fmt.Sprintf("%d", reset))

	if !ok {
		errResponse(w, http.StatusTooManyRequests, "")
	}

	return ok, res
}

// isAdmin checks if the specified userID is
// the owner or in the list of admins defined
// in the config file.
func (api *API) isAdmin(userID string) bool {
	if api.cfg.Discord.OwnerID == userID {
		return true
	}

	for _, id := range api.cfg.API.AdminIDs {
		if id == userID {
			return true
		}
	}

	return false
}
