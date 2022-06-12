package api

import (
	"encoding/json"
	"log"
	"msw-open-music/pkg/database"
	"net/http"
)

type getTagsResponse struct {
	Tags []*database.Tag `json:"tags"`
}

func (api *API) HandleGetTags(w http.ResponseWriter, r *http.Request) {
	tags, err := api.Db.GetTags()
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	log.Println("[api] Successfully got tags")

	resp := &getTagsResponse{Tags: tags}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
}

type InsertTagResponse struct {
	Tag *database.Tag `json:"tag"`
}

func (api *API) HandleInsertTag(w http.ResponseWriter, r *http.Request) {
	// check if user is admin
	err := api.CheckAdmin(w, r)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	req := &database.Tag{}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	req.CreatedByUserId, err = api.GetUserID(w, r)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	tagID, err := api.Db.InsertTag(req)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	tag, err := api.Db.GetTag(tagID)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	resp := &InsertTagResponse{Tag: tag}
	
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
}

type GetTagInfoRequest struct {
	ID int64 `json:"id"`
}

type GetTagInfoResponse struct {
	Tag *database.Tag `json:"tag"`
}

func (api *API) HandleGetTagInfo(w http.ResponseWriter, r *http.Request) {
	var req GetTagInfoRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	tag, err := api.Db.GetTag(req.ID)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	resp := &GetTagInfoResponse{Tag: tag}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
}

func (api *API) HandleUpdateTag(w http.ResponseWriter, r *http.Request) {
	// check if user is admin
	err := api.CheckAdmin(w, r)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	req := &database.Tag{}
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	err = api.Db.UpdateTag(req)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	api.HandleOK(w, r)
}

type DeleteTagRequest struct {
	ID int64 `json:"id"`
}

func (api *API) HandleDeleteTag(w http.ResponseWriter, r *http.Request) {
	// check if user is admin
	err := api.CheckAdmin(w, r)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	req := &DeleteTagRequest{}
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	err = api.Db.DeleteTag(req.ID)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	log.Println("[api] Successfully deleted tag and its references", req.ID)

	api.HandleOK(w, r)
}

