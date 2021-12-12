package api

import (
	"encoding/json"
	"net/http"
)

type UpdateFoldernameRequest struct {
	ID         int64  `json:"id"`
	Foldername string `json:"foldername"`
}

func (api *API) HandleUpdateFoldername(w http.ResponseWriter, r *http.Request) {
	req := &UpdateFoldernameRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// check is admin
	err = api.CheckAdmin(w, r)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	err = api.Db.UpdateFoldername(req.ID, req.Foldername)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	api.HandleOK(w, r)
}
