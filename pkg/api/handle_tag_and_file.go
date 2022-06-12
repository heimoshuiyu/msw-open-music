package api

import (
	"encoding/json"
	"log"
	"msw-open-music/pkg/database"
	"net/http"
)

type PutTagOnFileRequest struct {
	TagID  int64 `json:"tag_id"`
	FileID int64 `json:"file_id"`
}

func (api *API) HandlePutTagOnFile(w http.ResponseWriter, r *http.Request) {
	// check if the user is admin
	err := api.CheckAdmin(w, r)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	req := &PutTagOnFileRequest{}
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	userID, err := api.GetUserID(w, r)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	// check empty
	if req.TagID == 0 || req.FileID == 0 {
		api.HandleError(w, r, ErrEmpty)
		return
	}

	log.Println("[api] Put tag on file request:", req, "userID:", userID)

	api.Db.PutTagOnFile(req.TagID, req.FileID, userID)

	api.HandleOK(w, r)
}

type GetTagsOnFileRequest struct {
	ID int64 `json:"id"`
}

type GetTagsOnFileResponse struct {
	Tags []*database.Tag `json:"tags"`
}

func (api *API) HandleGetTagsOnFile(w http.ResponseWriter, r *http.Request) {
	req := &GetTagsOnFileRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
	
	log.Println("[api] Get tags on file request:", req)

	tags, err := api.Db.GetTagsOnFile(req.ID)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	resp := &GetTagsOnFileResponse{
		Tags: tags,
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
}

type DeleteTagOnFileRequest struct {
	TagID  int64 `json:"tag_id"`
	FileID int64 `json:"file_id"`
}

func (api *API) HandleDeleteTagOnFile(w http.ResponseWriter, r *http.Request) {
	// check if the user is admin
	err := api.CheckAdmin(w, r)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	req := &DeleteTagOnFileRequest{}
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	// check empty
	if req.TagID == 0 || req.FileID == 0 {
		api.HandleError(w, r, ErrEmpty)
		return
	}

	log.Println("[api] Delete tag on file request:", req)

	err = api.Db.DeleteTagOnFile(req.TagID, req.FileID)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	api.HandleOK(w, r)
}
