package discordoauth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/zekroTJA/yuri2/internal/static"
)

// OnErrorFunc is the function to be used to handle errors during
// authentication.
type OnErrorFunc func(w http.ResponseWriter, r *http.Request, status int, msg string)

// OnSuccessFuc is the func to be used to handle the successful
// authentication.
type OnSuccessFuc func(w http.ResponseWriter, r *http.Request, userID string)

// DiscordOAuth provides http handlers for
// authenticating a discord User by your Discord
// OAuth application.
type DiscordOAuth struct {
	clientID     string
	clientSecret string
	redirectURI  string

	onError   OnErrorFunc
	onSuccess OnSuccessFuc
}

type oAuthTokenResponse struct {
	Error       string `json:"error"`
	AccessToken string `json:"access_token"`
}

type getUserMeResponse struct {
	Error string `json:"error"`
	ID    string `json:"id"`
}

// NewDiscordOAuth returns a new instance of DiscordOAuth.
func NewDiscordOAuth(clientID, clientSecret, redirectURI string, onError OnErrorFunc, onSuccess OnSuccessFuc) *DiscordOAuth {
	if onError == nil {
		onError = func(w http.ResponseWriter, r *http.Request, status int, msg string) {}
	}
	if onSuccess == nil {
		onSuccess = func(w http.ResponseWriter, r *http.Request, userID string) {}
	}

	return &DiscordOAuth{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,

		onError:   onError,
		onSuccess: onSuccess,
	}
}

// HandlerInit returns a redirect response to the OAuth Apps
// authentication page.
func (d *DiscordOAuth) HandlerInit(w http.ResponseWriter, r *http.Request) {
	uri := fmt.Sprintf("https://discordapp.com/api/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=identify",
		d.clientID, url.QueryEscape(d.redirectURI))
	w.Header().Set("Location", uri)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// HandlerCallback will be requested by discordapp.com on successful
// app authentication. This handler will check the validity of the passed
// authorization code by getting a bearer token and trying to get self
// user data by requesting them using the bearer token.
// If this fails, onError will be called. Else, onSuccess will be
// called passing the userID of the user authenticated.
func (d *DiscordOAuth) HandlerCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	// 1. Request getting bearer token by app auth code

	data := map[string][]string{
		"client_id":     []string{d.clientID},
		"client_secret": []string{d.clientSecret},
		"grant_type":    []string{"authorization_code"},
		"code":          []string{code},
		"redirect_uri":  []string{d.redirectURI},
		"scope":         []string{"identify"},
	}

	values := url.Values(data)
	req, err := http.NewRequest("POST", static.URLDiscordAPIOAuthToken,
		bytes.NewBuffer([]byte(values.Encode())))
	if err != nil {
		d.onError(w, r, http.StatusInternalServerError, "failed creating request: "+err.Error())
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		d.onError(w, r, http.StatusInternalServerError, "failed executing request: "+err.Error())
		return
	}

	if res.StatusCode >= 300 {
		d.onError(w, r, http.StatusUnauthorized, "")
		return
	}

	resAuthBody := new(oAuthTokenResponse)
	err = parseJSONBody(res.Body, resAuthBody)
	if err != nil {
		d.onError(w, r, http.StatusInternalServerError, "failed parsing Discord API response: "+err.Error())
		return
	}

	if resAuthBody.Error != "" || resAuthBody.AccessToken == "" {
		d.onError(w, r, http.StatusUnauthorized, "")
		return
	}

	// 2. Request getting user ID

	req, err = http.NewRequest("GET", static.URLDiscordGetUserMe, nil)
	if err != nil {
		d.onError(w, r, http.StatusInternalServerError, "failed creating request: "+err.Error())
		return
	}

	req.Header.Set("Authorization", "Bearer "+resAuthBody.AccessToken)

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		d.onError(w, r, http.StatusInternalServerError, "failed executing request: "+err.Error())
		return
	}

	if res.StatusCode >= 300 {
		d.onError(w, r, http.StatusUnauthorized, "")
		return
	}

	resGetMe := new(getUserMeResponse)
	err = parseJSONBody(res.Body, resGetMe)
	if err != nil {
		d.onError(w, r, http.StatusInternalServerError, "failed parsing Discord API response: "+err.Error())
		return
	}

	if resGetMe.Error != "" || resGetMe.ID == "" {
		d.onError(w, r, http.StatusUnauthorized, "")
		return
	}

	d.onSuccess(w, r, resGetMe.ID)
}

func parseJSONBody(body io.ReadCloser, v interface{}) error {
	dec := json.NewDecoder(body)
	return dec.Decode(v)
}
