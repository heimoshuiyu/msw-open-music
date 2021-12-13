package api

import (
	"encoding/json"
	"errors"
	"msw-open-music/pkg/database"
	"net/http"
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

type GetReviewsOnFileRequest struct {
	ID int64 `json:"id"`
}

type GetReviewsOnFileResponse struct {
	Reviews []*database.Review `json:"reviews"`
}

func (api *API) HandleGetReviewsOnFile(w http.ResponseWriter, r *http.Request) {
	req := &GetReviewsOnFileRequest{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	reviews, err := api.Db.GetReviewsOnFile(req.ID)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	resp := &GetReviewsOnFileResponse{
		Reviews: reviews,
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
}

type GetReviewRequest struct {
	ID int64 `json:"id"`
}

type GetReviewResponse struct {
	Review *database.Review `json:"review"`
}

func (api *API) HandleGetReview(w http.ResponseWriter, r *http.Request) {
	req := &GetReviewRequest{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	review, err := api.Db.GetReview(req.ID)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	ret := &GetReviewResponse{
		Review: review,
	}

	err = json.NewEncoder(w).Encode(ret)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
}

func (api *API) CheckUserCanModifyReview(w http.ResponseWriter, r *http.Request, reviewID int64) error {
	review, err := api.Db.GetReview(reviewID)
	if err != nil {
		return err
	}

	err = api.CheckNotAnonymous(w, r)
	if err != nil {
		return err
	}

	err = api.CheckAdmin(w, r)
	if err != nil {
		userID, err := api.GetUserID(w, r)
		if err != nil {
			return err
		}

		if review.UserId != userID {
			return errors.New("you are not allowed to modify this review")
		}
	}

	return nil
}

func (api *API) HandleUpdateReview(w http.ResponseWriter, r *http.Request) {
	req := &database.Review{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	err = api.CheckUserCanModifyReview(w, r, req.ID)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	req.UpdatedAt = time.Now().Unix()

	err = api.Db.UpdateReview(req)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	api.HandleOK(w, r)
}

type DeleteReviewRequest struct {
	ID int64 `json:"id"`
}

func (api *API) HandleDeleteReview(w http.ResponseWriter, r *http.Request) {
	req := &DeleteReviewRequest{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	err = api.CheckUserCanModifyReview(w, r, req.ID)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	err = api.Db.DeleteReview(req.ID)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	api.HandleOK(w, r)
}

type GetReviewsByUserRequest struct {
	ID int64 `json:"id"`
}

func (api *API) HandleGetReviewsByUser(w http.ResponseWriter, r *http.Request) {
	req := &GetReviewsByUserRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	reviews, err := api.Db.GetReviewsByUser(req.ID)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	err = json.NewEncoder(w).Encode(reviews)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
}
