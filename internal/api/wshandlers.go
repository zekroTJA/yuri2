package api

import (
	"fmt"

	"github.com/zekroTJA/yuri2/pkg/wsmgr"
)

type wsInitData struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
}

func (api *API) wsInitHandler(e *wsmgr.Event) {
	data := new(wsInitData)

	err := e.ParseDataTo(data)
	if err != nil {
		wsSendError(e.Sender, fmt.Sprintf("failed parsing data: %s", err.Error()))
		return
	}

	ok, _, err := api.auth.CheckAndRefersh(data.UserID, data.Token)
	if err != nil {
		wsSendError(e.Sender, fmt.Sprintf("failed checking auth: %s", err.Error()))
		return
	}

	if !ok {
		wsSendError(e.Sender, "unauthorized")
		e.Sender.Close()
		return
	}

	e.Sender.SetIdent(data.UserID)
}

func (api *API) wsPlayHandler(e *wsmgr.Event) {
	wsCheckInitilized(e.Sender)
}
