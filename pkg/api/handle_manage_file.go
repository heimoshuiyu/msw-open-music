package api

import (
	"encoding/json"
	"net/http"
)

type DeleteFileRequest struct {
	ID int64 `json:"id"`
}

func (api *API) HandleDeleteFile(w http.ResponseWriter, r *http.Request) {
	// check admin
	err := api.CheckAdmin(w, r)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	req := &DeleteFileRequest{}
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	err = api.Db.DeleteFile(req.ID)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	api.HandleOK(w, r)
}
