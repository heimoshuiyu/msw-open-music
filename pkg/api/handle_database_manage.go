package api

import (
	"encoding/json"
	"net/http"
)

type WalkRequest struct {
	Root    string   `json:"root"`
	Pattern []string `json:"pattern"`
	TagIDs  []int64  `json:"tag_ids"`
}

func (api *API) HandleReset(w http.ResponseWriter, r *http.Request) {
	var err error
	// check admin
	err = api.CheckAdmin(w, r)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	// reset
	err = api.Db.ResetFiles()
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
	err = api.Db.ResetFolder()
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	api.HandleStatus(w, r, "Database reseted")
}

func (api *API) HandleWalk(w http.ResponseWriter, r *http.Request) {
	walkRequest := &WalkRequest{}
	err := json.NewDecoder(r.Body).Decode(walkRequest)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	// check admin
	err = api.CheckAdmin(w, r)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	// check root empty
	if walkRequest.Root == "" {
		api.HandleErrorString(w, r, `key "root" can't be empty`)
		return
	}

	// check pattern empty
	if len(walkRequest.Pattern) == 0 {
		api.HandleErrorString(w, r, `"[]pattern" can't be empty`)
		return
	}

	// get userID
	userID, err := api.GetUserID(w, r)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	// walk
	err = api.Db.Walk(walkRequest.Root, walkRequest.Pattern, walkRequest.TagIDs, userID)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	api.HandleStatus(w, r, "Database udpated")
}
