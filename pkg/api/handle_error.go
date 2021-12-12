package api

import (
	"encoding/json"
	"log"
	"net/http"
	"errors"
)

var (
	ErrNotLoggedIn = errors.New("not logged in")
	ErrNotAdmin = errors.New("not admin")
	ErrEmpty = errors.New("Empty field detected, please fill in all fields")
)

type Error struct {
  Error string `json:"error,omitempty"`
}

func (api *API) HandleError(w http.ResponseWriter, r *http.Request, err error) {
	api.HandleErrorString(w, r, err.Error())
}

func (api *API) HandleErrorCode(w http.ResponseWriter, r *http.Request, err error, code int) {
	api.HandleErrorStringCode(w, r, err.Error(), code)
}

func (api *API) HandleErrorString(w http.ResponseWriter, r *http.Request, errorString string) {
	api.HandleErrorStringCode(w, r, errorString, 500)
}

func (api *API) HandleErrorStringCode(w http.ResponseWriter, r *http.Request, errorString string, code int) {
	log.Println("[api] [Error]", code, errorString)
	errStatus := &Error{
		Error: errorString,
	}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(errStatus)
}
