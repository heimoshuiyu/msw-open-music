package api

import (
	"database/sql"
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

func (api *API) LoginAsAnonymous(w http.ResponseWriter, r *http.Request) {
	user, err := api.Db.LoginAsAnonymous()
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	session, _ := api.store.Get(r, api.defaultSessionName)

	// save session
	session.Values["userId"] = user.ID
	err = session.Save(r, w)
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

func (api *API) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var user *database.User
	var err error
	session, _ := api.store.Get(r, api.defaultSessionName)
	log.Println("Session:", session.Values)

	// Get method will login current or anonymous user
	if r.Method == "GET" {

		// if user already logged in
		if userId, ok := session.Values["userId"]; ok {
			user, err = api.Db.GetUserById(userId.(int64))
			if err != nil {
				if err != sql.ErrNoRows {
					api.HandleError(w, r, err)
					return
				}
				log.Println("User not found")
				// login as anonymous user
				api.LoginAsAnonymous(w, r)
				return
			}
			log.Println("User already logged in:", user)

		} else {
			// login as anonymous user
			log.Println("Login as anonymous user")
			api.LoginAsAnonymous(w, r)
			return
		}

	} else {

		var request LoginRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			api.HandleError(w, r, err)
			return
		}

		log.Println("Login as user", request.Username)

		user, err = api.Db.Login(request.Username, request.Password)
		if err != nil {
			api.HandleError(w, r, err)
			return
		}
	}

	// save session
	session.Values["userId"] = user.ID
	err = session.Save(r, w)
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

func (api *API) CheckAdmin(w http.ResponseWriter, r *http.Request) error {
	session, _ := api.store.Get(r, api.defaultSessionName)
	userId, ok := session.Values["userId"]
	if !ok {
		return ErrNotLoggedIn
	}

	user, err := api.Db.GetUserById(userId.(int64))
	if err != nil {
		return err
	}

	if user.Role != database.RoleAdmin {
		return ErrNotAdmin
	}

	return nil
}

func (api *API) CheckNotAnonymous(w http.ResponseWriter, r *http.Request) error {
	session, _ := api.store.Get(r, api.defaultSessionName)
	userId, ok := session.Values["userId"]
	if !ok {
		return ErrNotLoggedIn
	}

	user, err := api.Db.GetUserById(userId.(int64))
	if err != nil {
		return err
	}

	if user.Role == database.RoleAnonymous {
		return ErrAnonymous
	}

	return nil
}

func (api *API) GetUserID(w http.ResponseWriter, r *http.Request) (int64, error) {
	session, _ := api.store.Get(r, api.defaultSessionName)
	userId, ok := session.Values["userId"]
	if !ok {
		return 0, ErrNotLoggedIn
	}

	return userId.(int64), nil
}

type GetUsersResponse struct {
	Users []*database.User `json:"users"`
}

func (api *API) HandleGetUsers(w http.ResponseWriter, r *http.Request) {
	err := api.CheckAdmin(w, r)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	users, err := api.Db.GetUsers()
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	ret := &GetUsersResponse{
		Users: users,
	}

	err = json.NewEncoder(w).Encode(ret)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
}

type UpdateUserActiveRequest struct {
	ID     int64 `json:"id"`
	Active bool  `json:"active"`
}

func (api *API) HandleUpdateUserActive(w http.ResponseWriter, r *http.Request) {
	err := api.CheckAdmin(w, r)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	req := &UpdateUserActiveRequest{}
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	err = api.Db.UpdateUserActive(req.ID, req.Active)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
	api.HandleOK(w, r)
}
