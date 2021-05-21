package api

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"msw-open-music/internal/pkg/database"
	"net/http"
	"os"
)

type API struct {
	Db *database.Database
	Server http.Server
	token string
}

type Status struct {
	Status string `json:"status,omitempty"`
}
var ok Status = Status{
	Status: "OK",
}

type WalkRequest struct {
	Token string `json:"token"`
	Root string `json:"root"`
	Pattern []string `json:"pattern"`
}

type ResetRequest struct {
	Token string `json:"token"`
}

type SearchFilesRequest struct {
	Filename string `json:"filename"`
	Limit int64 `json:"limit"`
	Offset int64 `json:"offest"`
}

type SearchFoldersRequest struct {
	Foldername string `json:"foldername"`
	Limit int64 `json:"limit"`
	Offset int64 `json:"offest"`
}

type SearchFilesRespond struct {
	Files []database.File `json:"files"`
}

type SearchFoldersRespond struct {
	Folders []database.Folder `json:"folders"`
}

func (api *API) CheckToken(token string) (error) {
	if token != api.token {
		return errors.New("token not matched")
	}
	return nil
}

func (api *API) HandleTokenNotMatch(w http.ResponseWriter, r *http.Request) {
	api.HandleErrorStringCode(w, r, "token not match", 403)
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
	log.Println("Handle Error", errorString)
	errStatus := &Status{
		Status: errorString,
	}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(errStatus)
}

func (api *API) HandleReset(w http.ResponseWriter, r *http.Request) {
	resetRequest := &ResetRequest{}
	err := json.NewDecoder(r.Body).Decode(resetRequest)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	// check token
	err = api.CheckToken(resetRequest.Token)
	if err != nil {
		api.HandleTokenNotMatch(w, r)
		return
	}

	// reset
	err = api.Db.ResetFiles()
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
	err = api.Db.ResetFolder()
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	api.HandleStatus(w, r, "Database reseted")
}

func (api *API) HandleWalk(w http.ResponseWriter, r *http.Request) {
	walkRequest := &WalkRequest{}
	err := json.NewDecoder(r.Body).Decode(walkRequest)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	// check token match
	err = api.CheckToken(walkRequest.Token)
	if err != nil {
		api.HandleTokenNotMatch(w, r)
		return
	}

	// check root empty
	if walkRequest.Root == "" {
		api.HandleErrorString(w, r, `key "root" can't be empty`)
		return
	}

	// check pattern empty
	if len(walkRequest.Pattern) == 0 {
		api.HandleErrorString(w, r, `"[]pattern" can't be empty`)
		return
	}

	// walk
	err = api.Db.Walk(walkRequest.Root, walkRequest.Pattern)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	api.HandleStatus(w, r, "Database udpated")
}

func (api *API) HandleOK(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(&ok)
}

func (api *API) HandleStatus(w http.ResponseWriter, r *http.Request, status string) {
	s := &Status{
		Status: status,
	}

	json.NewEncoder(w).Encode(s)
}

func (api *API) HandleSearchFiles(w http.ResponseWriter, r *http.Request) {
	searchFilesRequest := &SearchFilesRequest{}
	err := json.NewDecoder(r.Body).Decode(searchFilesRequest)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	// check empty
	if searchFilesRequest.Filename == "" {
		api.HandleErrorString(w, r, `"filename" can't be empty`)
		return
	}
	if searchFilesRequest.Limit == 0 {
		api.HandleErrorString(w, r, `"limit" can't be zero`)
		return
	}

	searchFilesRespond := &SearchFilesRespond{}

	searchFilesRespond.Files, err = api.Db.SearchFiles(searchFilesRequest.Filename, searchFilesRequest.Limit, searchFilesRequest.Offset)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	json.NewEncoder(w).Encode(searchFilesRespond)
}

func (api *API) HandleSearchFolders(w http.ResponseWriter, r *http.Request) {
	searchFoldersRequest := &SearchFoldersRequest{}
	err := json.NewDecoder(r.Body).Decode(searchFoldersRequest)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	// check empty
	if searchFoldersRequest.Foldername == "" {
		api.HandleErrorString(w, r, `"foldername" can't be empty`)
		return
	}
	if searchFoldersRequest.Limit == 0 {
		api.HandleErrorString(w, r, `"limit" can't be zero`)
		return
	}

	searchFoldersRespond := &SearchFoldersRespond{}

	searchFoldersRespond.Folders, err = api.Db.SearchFolders(searchFoldersRequest.Foldername, searchFoldersRequest.Limit, searchFoldersRequest.Offset)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	json.NewEncoder(w).Encode(searchFoldersRespond)
}

type GetFileRequest struct {
	ID int64 `json:"id"`
}

func (api *API) HandleGetFile(w http.ResponseWriter, r *http.Request) {
	getFilesRequest := &GetFileRequest{
		ID: -1,
	}

	err := json.NewDecoder(r.Body).Decode(getFilesRequest)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	// check empty
	if getFilesRequest.ID < 0 {
		api.HandleErrorString(w, r, `"id" can't be none or negative`)
		return
	}

	file, err := api.Db.GetFile(getFilesRequest.ID)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	path, err := file.Path()
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	src, err := os.Open(path)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
	defer src.Close()
	io.Copy(w, src)
}

func NewAPI(dbName string, Addr string) (*API, error) {
	var err error

	db, err := database.NewDatabase(dbName)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	apiMux := http.NewServeMux()

	api := &API{
		Db: db,
		Server: http.Server{
			Addr: Addr,
			Handler: mux,
		},
	}

	// mount api
	apiMux.HandleFunc("/hello", api.HandleOK)
	apiMux.HandleFunc("/get_file", api.HandleGetFile)
	apiMux.HandleFunc("/search_files", api.HandleSearchFiles)
	apiMux.HandleFunc("/search_folders", api.HandleSearchFolders)
	// below needs token
	apiMux.HandleFunc("/walk", api.HandleWalk)
	apiMux.HandleFunc("/reset", api.HandleReset)

	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", apiMux))
	mux.Handle("/web/", http.StripPrefix("/web", http.FileServer(http.Dir("web"))))

	api.token = "pwd"

	return api, nil
}
