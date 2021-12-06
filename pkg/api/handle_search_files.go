package api

import (
	"encoding/json"
	"log"
	"msw-open-music/pkg/database"
	"net/http"
)

type SearchFilesRequest struct {
	Filename string `json:"filename"`
	Limit    int64  `json:"limit"`
	Offset   int64  `json:"offset"`
}

type SearchFilesResponse struct {
	Files []database.File `json:"files"`
}

func (api *API) HandleSearchFiles(w http.ResponseWriter, r *http.Request) {
	searchFilesRequest := &SearchFilesRequest{}
	err := json.NewDecoder(r.Body).Decode(searchFilesRequest)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	// check empty
	if searchFilesRequest.Filename == "" {
		api.HandleErrorString(w, r, `"filename" can't be empty`)
		return
	}
	if api.CheckLimit(w, r, searchFilesRequest.Limit) != nil {
		return
	}

	searchFilesResponse := &SearchFilesResponse{}

	searchFilesResponse.Files, err = api.Db.SearchFiles(searchFilesRequest.Filename, searchFilesRequest.Limit, searchFilesRequest.Offset)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	log.Println("[api] Search files", searchFilesRequest.Filename, searchFilesRequest.Limit, searchFilesRequest.Offset)

	json.NewEncoder(w).Encode(searchFilesResponse)
}
