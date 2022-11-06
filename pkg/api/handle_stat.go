package api

import (
	"encoding/json"
	"log"
	"msw-open-music/pkg/database"
	"net/http"
	"time"
)

type RecordPlaybackRequest struct {
	Playback database.Playback `json:"playback"`
}

func (api *API) HandleRecordPlayback(w http.ResponseWriter, r *http.Request) {
	recordPlaybackRequest := &RecordPlaybackRequest{}
	err := json.NewDecoder(r.Body).Decode(recordPlaybackRequest)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	recordPlaybackRequest.Playback.Time = time.Now()
	recordPlaybackRequest.Playback.UserID, err = api.GetUserID(w, r)
	if err != nil {
		if err == ErrNotLoggedIn {
			user, err := api.Db.LoginAsAnonymous()
			if err != nil {
				api.HandleError(w, r, err)
				return
			}
			recordPlaybackRequest.Playback.UserID = user.ID
		} else {
			api.HandleError(w, r, err)
			return
		}
	}

	log.Println("[api] Record playback history",
		recordPlaybackRequest.Playback.UserID,
		recordPlaybackRequest.Playback.FileID,
		recordPlaybackRequest.Playback.Duration,
		recordPlaybackRequest.Playback.Method)

	err = api.Db.RecordPlayback(recordPlaybackRequest.Playback)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
	api.HandleOK(w, r)
}
