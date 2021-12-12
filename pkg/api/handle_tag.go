package api

import (
	"encoding/json"
	"log"
	"msw-open-music/pkg/database"
	"net/http"
)

type getTagsResponse struct {
	Tags []database.Tag `json:"tags"`
}

func (api *API) HandleGetTags(w http.ResponseWriter, r *http.Request) {
	tags, err := api.Db.GetTags()
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	log.Println("Successfully got tags")

	resp := &getTagsResponse{Tags: tags}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
}

type InsertTagRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
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

	var req InsertTagRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	tag, err := api.Db.InsertTag(req.Name, req.Description)
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
	tag, err := api.Db.UpdateTag(req)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
	err = json.NewEncoder(w).Encode(tag)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
}
