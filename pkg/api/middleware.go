package api

import (
	"errors"
	"net/http"
)

func (api *API) PermissionMiddleware(next http.Handler) http.Handler {
	// 0 anonymous user
	// 1 admin
	// 2 normal user
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get permission of URL
		permission, ok := api.APIConfig.Permission[r.URL.Path]
		// 0 means no permission required
		if !ok || permission == 0 {
			next.ServeHTTP(w, r)
			return
		}

		// ger user permission level
		userLevel := api.GetUserLevel(r)

		// admin has root (highest) permission level 1
		if userLevel == 1 {
			next.ServeHTTP(w, r)
			return
		}

		// anonymous userLevel 0 don't have any permission
		// check permission level for other users
		if userLevel == 0 || userLevel > permission {
			api.HandleError(w, r, errors.New("No enougth permission"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (api *API) GetUserLevel(r *http.Request) int64 {
	session, _ := api.store.Get(r, api.defaultSessionName)
	userId, ok := session.Values["userId"]
	if !ok {
		// not logined user is considered anonymous user
		return 0
	}

	user, err := api.Db.GetUserById(userId.(int64))
	if err != nil {
		return 0
	}

	return user.Role
}
