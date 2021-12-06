package api

import (
	"encoding/json"
	"log"
	"msw-open-music/pkg/database"
	"net/http"
)

type SearchFoldersRequest struct {
	Foldername string `json:"foldername"`
	Limit      int64  `json:"limit"`
	Offset     int64  `json:"offset"`
}

type SearchFoldersResponse struct {
	Folders []database.Folder `json:"folders"`
}

func (api *API) HandleSearchFolders(w http.ResponseWriter, r *http.Request) {
	searchFoldersRequest := &SearchFoldersRequest{}
	err := json.NewDecoder(r.Body).Decode(searchFoldersRequest)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	// check empty
	if searchFoldersRequest.Foldername == "" {
		api.HandleErrorString(w, r, `"foldername" can't be empty`)
		return
	}
	if api.CheckLimit(w, r, searchFoldersRequest.Limit) != nil {
		return
	}

	searchFoldersResponse := &SearchFoldersResponse{}

	searchFoldersResponse.Folders, err = api.Db.SearchFolders(searchFoldersRequest.Foldername, searchFoldersRequest.Limit, searchFoldersRequest.Offset)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	log.Println("[api] Search folders", searchFoldersRequest.Foldername, searchFoldersRequest.Limit, searchFoldersRequest.Offset)

	json.NewEncoder(w).Encode(searchFoldersResponse)
}
