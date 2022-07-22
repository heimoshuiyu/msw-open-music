package api

import (
	"bytes"
	"encoding/json"
	"log"
	"msw-open-music/pkg/database"
	"net/http"
	"time"
)

type FeedbackRequest struct {
	Content string `json:"content"`
}

func (api *API) HandleFeedback(w http.ResponseWriter, r *http.Request) {
	feedbackRequest := &FeedbackRequest{}
	err := json.NewDecoder(r.Body).Decode(feedbackRequest)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	// check empty feedback
	if feedbackRequest.Content == "" {
		api.HandleErrorString(w, r, `"feedback" can't be empty`)
		return
	}

	log.Println("[api] Feedback", feedbackRequest.Content)

	headerBuff := &bytes.Buffer{}
	err = r.Header.Write(headerBuff)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
	header := headerBuff.String()

	userID, err := api.GetUserID(w, r)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	err = api.Db.InsertFeedback(time.Now().Unix(), feedbackRequest.Content, userID, header)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
	api.HandleOK(w, r)
}

type GetFeedbacksResponse struct {
	Feedbacks []*database.Feedback `json:"feedbacks"`
}

func (api *API) HandleGetFeedbacks(w http.ResponseWriter, r *http.Request) {
	feedbacks, err := api.Db.GetFeedbacks()
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	resp := &GetFeedbacksResponse{
		Feedbacks: feedbacks,
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
}

type DeleteFeedbackRequest struct {
	ID int64 `json:"id"`
}

func (api *API) HandleDeleteFeedback(w http.ResponseWriter, r *http.Request) {
	req := &DeleteFeedbackRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	err = api.Db.DeleteFeedback(req.ID)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	api.HandleOK(w, r)
}
