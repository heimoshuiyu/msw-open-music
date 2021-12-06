package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type FeedbackRequest struct {
	Feedback string `json:"feedback"`
}

func (api *API) HandleFeedback(w http.ResponseWriter, r *http.Request) {
	feedbackRequest := &FeedbackRequest{}
	err := json.NewDecoder(r.Body).Decode(feedbackRequest)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	// check empty feedback
	if feedbackRequest.Feedback == "" {
		api.HandleErrorString(w, r, `"feedback" can't be empty`)
		return
	}

	log.Println("[api] Feedback", feedbackRequest.Feedback)

	headerBuff := &bytes.Buffer{}
	err = r.Header.Write(headerBuff)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
	header := headerBuff.String()

	err = api.Db.InsertFeedback(time.Now().Unix(), feedbackRequest.Feedback, header)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
	api.HandleOK(w, r)
}
