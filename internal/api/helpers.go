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

	"github.com/zekroTJA/yuri2/pkg/wsmgr"
)

var stdErrMsgs = map[int]string{
	400: "bad request",
	401: "unauthorized",
	403: "forbidden",
	429: "too many requests",
}

type apiErrorBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

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

func errResponseWrapper(w http.ResponseWriter, r *http.Request, code int, msg string) {
	errResponse(w, code, msg)
}

func parseJSONBody(body io.ReadCloser, v interface{}) error {
	dec := json.NewDecoder(body)
	return dec.Decode(v)
}

func jsonResponse(w http.ResponseWriter, code int, data interface{}) {
	var bData []byte
	var err error

	w.Header().Add("Content-Type", "application/json")

	if data != nil {
		bData, err = json.MarshalIndent(data, "", "  ")
	}

	if err != nil {
		return
	}

	w.WriteHeader(code)
	_, err = w.Write(bData)
}

func errPageResponse(w http.ResponseWriter, r *http.Request, code int, msg string) {
	pageLoc := "./web/pages/400.html"

	switch code {
	case 401:
		pageLoc = "./web/pages/errors/401.html"
	case 404:
		pageLoc = "./web/pages/errors/404.html"
	case 500:
		pageLoc = "./web/pages/errors/500.html"
	}

	http.ServeFile(w, r, pageLoc)
}

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

func (api *API) checkAuthCookie(r *http.Request) (bool, string, error) {
	token, err := getCookieValue(r, "token")
	if err != nil {
		return false, "", err
	}

	userID, err := getCookieValue(r, "userid")
	if err != nil {
		return false, "", err
	}

	ok, _, err := api.auth.CheckAndRefersh(userID, token)

	return ok, userID, err
}

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

	ok, _, err := api.auth.CheckAndRefersh(userID, token)

	return ok, userID, err
}

func (api *API) checkAuthHeaderWithResponse(w http.ResponseWriter, r *http.Request) (bool, string) {
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
		errResponse(w, http.StatusUnauthorized, "")
		return false, ""
	}

	return true, userID
}

func wsSendError(wsc *wsmgr.WebSocketConn, msg string) error {
	return wsc.Out(wsmgr.NewEvent("ERROR", msg))
}

func wsCheckInitilized(wsc *wsmgr.WebSocketConn) string {
	ident, ok := wsc.GetIdent().(string)
	if !ok || ident == "" {
		wsSendError(wsc, "unauthorized")
		wsc.Close()
		return ""
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
