package api

import (
	"encoding/json"
	"io"
	"net/http"

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

func (api *API) checkAuthCookie(w http.ResponseWriter, r *http.Request) (bool, string, error) {
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
