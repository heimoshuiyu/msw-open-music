package api

import (
	"net/http"
	"msw-open-music/pkg/database"
	"encoding/json"
	"time"
)

// review.FileId, review.Content
func (api *API) HandleInsertReview(w http.ResponseWriter, r *http.Request) {
	review := &database.Review{}

	err := json.NewDecoder(r.Body).Decode(review)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	// check not anonymous
	err = api.CheckNotAnonymous(w, r)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	review.UserId, err = api.GetUserID(w, r)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	review.CreatedAt = time.Now().Unix()

	err = api.Db.InsertReview(review)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	api.HandleOK(w, r)
}
