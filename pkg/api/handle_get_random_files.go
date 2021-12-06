package api

import (
	"encoding/json"
	"log"
	"msw-open-music/pkg/database"
	"net/http"
)

type GetRandomFilesResponse struct {
	Files *[]database.File `json:"files"`
}

func (api *API) HandleGetRandomFiles(w http.ResponseWriter, r *http.Request) {
	files, err := api.Db.GetRandomFiles(10)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
	getRandomFilesResponse := &GetRandomFilesResponse{
		Files: &files,
	}
	log.Println("[api] Get random files")
	json.NewEncoder(w).Encode(getRandomFilesResponse)
}
