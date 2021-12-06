package api

import (
	"encoding/json"
	"net/http"
)

type Status struct {
	Status string `json:"status,omitempty"`
}

func (api *API) HandleStatus(w http.ResponseWriter, r *http.Request, status string) {
	s := &Status{
		Status: status,
	}

	json.NewEncoder(w).Encode(s)
}

var ok Status = Status{
	Status: "OK",
}

func (api *API) HandleOK(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(&ok)
}
