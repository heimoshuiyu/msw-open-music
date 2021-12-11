package api

import (
	"encoding/json"
	"log"
	"msw-open-music/pkg/database"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User *database.User `json:"user"`
}

func (api *API) HandleLogin(w http.ResponseWriter, r *http.Request) {
	// Get method will login as anonymous user
	if r.Method == "GET" {
		log.Println("Login as anonymous user")
		user, err := api.Db.LoginAsAnonymous()
		if err != nil {
			api.HandleError(w, r, err)
			return
		}
		resp := &LoginResponse{
			User: user,
		}
		err = json.NewEncoder(w).Encode(resp)
		return
	}

	var request LoginRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	log.Println("Login as user", request.Username)

	user, err := api.Db.Login(request.Username, request.Password)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	resp := &LoginResponse{
		User: user,
	}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     int64  `json:"role"`
}

func (api *API) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var request RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	log.Println("Register user", request.Username)

	err = api.Db.Register(request.Username, request.Password, request.Role)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	api.HandleOK(w, r)
}
