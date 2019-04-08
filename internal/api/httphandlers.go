package api

import (
	"fmt"
	"net/http"
)

func (api *API) getTokenHandler(w http.ResponseWriter, r *http.Request, userID string) {
	fmt.Println("USERID:", userID)
}
