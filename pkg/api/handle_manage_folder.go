package api

import (
	"encoding/json"
	"net/http"
)

type ResetFoldernameRequest struct {
	ID int64 `json:"id"`
}

func (api *API) HandleResetFoldername(w http.ResponseWriter, r *http.Request) {
	// check admin
	err := api.CheckAdmin(w, r)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	req := &ResetFoldernameRequest{}
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = api.Db.ResetFoldername(req.ID)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	api.HandleOK(w, r)
}
