package api

import (
	"errors"
	"log"
	"net/http"
)

func (api *API) CheckToken(w http.ResponseWriter, r *http.Request, token string) error {
	if token != api.token {
		err := errors.New("token not matched")
		log.Println("[api] [Warning] Token not matched", token)
		api.HandleErrorCode(w, r, err, 403)
		return err
	}
	log.Println("[api] Token passed")
	return nil
}
