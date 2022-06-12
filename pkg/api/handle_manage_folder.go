package api

import (
	"encoding/json"
	"log"
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

	log.Println("[api] Reset foldername folderID", req.ID)

	err = api.Db.ResetFoldername(req.ID)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	api.HandleOK(w, r)
}

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

	log.Println("[api] Update foldername folderID", req.ID, req.Foldername)

	err = api.Db.UpdateFoldername(req.ID, req.Foldername)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	api.HandleOK(w, r)
}
