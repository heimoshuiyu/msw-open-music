package api

import (
	"encoding/json"
	"log"
	"msw-open-music/pkg/database"
	"net/http"
)

type GetFilesInFolderRequest struct {
	Folder_id int64 `json:"folder_id"`
	Limit     int64 `json:"limit"`
	Offset    int64 `json:"offset"`
}

type GetFilesInFolderResponse struct {
	Files  *[]database.File `json:"files"`
	Folder string           `json:"folder"`
}

func (api *API) HandleGetFilesInFolder(w http.ResponseWriter, r *http.Request) {
	getFilesInFolderRequest := &GetFilesInFolderRequest{
		Folder_id: -1,
	}

	err := json.NewDecoder(r.Body).Decode(getFilesInFolderRequest)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	// check empyt
	if getFilesInFolderRequest.Folder_id < 0 {
		api.HandleErrorString(w, r, `"folder_id" can't be none or negative`)
		return
	}

	files, folder, err := api.Db.GetFilesInFolder(getFilesInFolderRequest.Folder_id, getFilesInFolderRequest.Limit, getFilesInFolderRequest.Offset)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	getFilesInFolderResponse := &GetFilesInFolderResponse{
		Files:  &files,
		Folder: folder,
	}

	log.Println("[api] Get files in folder", getFilesInFolderRequest.Folder_id)

	json.NewEncoder(w).Encode(getFilesInFolderResponse)
}
