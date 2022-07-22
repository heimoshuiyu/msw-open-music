package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type DeleteFileRequest struct {
	ID int64 `json:"id"`
}

func (api *API) HandleDeleteFile(w http.ResponseWriter, r *http.Request) {
	req := &DeleteFileRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	log.Println("[api] delete file", req.ID)

	err = api.Db.DeleteFile(req.ID)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	api.HandleOK(w, r)
}

type UpdateFilenameRequest struct {
	ID       int64  `json:"id"`
	Filename string `json:"filename"`
}

func (api *API) HandleUpdateFilename(w http.ResponseWriter, r *http.Request) {
	req := &UpdateFilenameRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	log.Println("[api] update filename", req.ID, req.Filename)

	err = api.Db.UpdateFilename(req.ID, req.Filename)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	api.HandleOK(w, r)
}

type ResetFilenameRequest struct {
	ID int64 `json:"id"`
}

func (api *API) HandleResetFilename(w http.ResponseWriter, r *http.Request) {
	req := &ResetFilenameRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	log.Println("[api] reset filename", req.ID)

	err = api.Db.ResetFilename(req.ID)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	api.HandleOK(w, r)
}
