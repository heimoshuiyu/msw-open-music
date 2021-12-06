package api

import (
	"errors"
	"log"
	"net/http"
)

func (api *API) CheckLimit(w http.ResponseWriter, r *http.Request, limit int64) error {
	if limit <= 0 || limit > 10 {
		log.Println("[api] [Warning] Limit error", limit)
		err := errors.New(`"limit" can't be zero or more than 10`)
		api.HandleError(w, r, err)
		return err
	}
	return nil
}
