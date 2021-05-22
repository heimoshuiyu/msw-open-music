package api

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"msw-open-music/internal/pkg/database"
	"net/http"
	"os"
	"os/exec"
	"strconv"
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
	Offset int64 `json:"offset"`
}

type SearchFoldersRequest struct {
	Foldername string `json:"foldername"`
	Limit int64 `json:"limit"`
	Offset int64 `json:"offset"`
}

type SearchFilesResponse struct {
	Files []database.File `json:"files"`
}

type SearchFoldersResponse struct {
	Folders []database.Folder `json:"folders"`
}

type GetFilesInFolderRequest struct {
	Folder_id int64 `json:"folder_id"`
	Limit int64 `json:"limit"`
	Offset int64 `json:"offset"`
}

type GetFilesInFolderResponse struct {
	Files *[]database.File `json:"files"`
}

type GetRandomFilesResponse struct {
	Files *[]database.File `json:"files"`
}

func (api *API) HandleGetRandomFiles(w http.ResponseWriter, r *http.Request) {
	files, err := api.Db.GetRandomFiles(10);
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
	getRandomFilesResponse := &GetRandomFilesResponse{
		Files: &files,
	}
	log.Println("[api] Get random files")
	json.NewEncoder(w).Encode(getRandomFilesResponse)
}

func (api *API) HandleGetFilesInFolder(w http.ResponseWriter, r *http.Request) {
	getFilesInFolderRequest := &GetFilesInFolderRequest{
		Folder_id: -1,
	}

	err := json.NewDecoder(r.Body).Decode(getFilesInFolderRequest)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	// check empyt
	if getFilesInFolderRequest.Folder_id < 0 {
		api.HandleErrorString(w, r, `"folder_id" can't be none or negative`)
		return
	}

	files, err := api.Db.GetFilesInFolder(getFilesInFolderRequest.Folder_id, getFilesInFolderRequest.Limit, getFilesInFolderRequest.Offset)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	getFilesInFolderResponse := &GetFilesInFolderResponse{
		Files: &files,
	}

	log.Println("[api] Get files in folder", getFilesInFolderRequest.Folder_id)

	json.NewEncoder(w).Encode(getFilesInFolderResponse)
}

func (api *API) CheckToken(w http.ResponseWriter, r *http.Request, token string) (error) {
	if token != api.token {
		err := errors.New("token not matched")
		log.Println("[api] [Warning] Token not matched", token)
		api.HandleErrorCode(w, r, err, 403)
		return err
	}
	log.Println("[api] Token passed")
	return nil
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
	err = api.CheckToken(w, r, resetRequest.Token)
	if err != nil {
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
	err = api.CheckToken(w, r, walkRequest.Token)
	if err != nil {
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
	if api.CheckLimit(w, r, searchFilesRequest.Limit) != nil {
		return
	}

	searchFilesResponse := &SearchFilesResponse{}

	searchFilesResponse.Files, err = api.Db.SearchFiles(searchFilesRequest.Filename, searchFilesRequest.Limit, searchFilesRequest.Offset)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	log.Println("[api] Search files", searchFilesRequest.Filename, searchFilesRequest.Limit, searchFilesRequest.Offset)

	json.NewEncoder(w).Encode(searchFilesResponse)
}

func (api *API) CheckLimit(w http.ResponseWriter, r *http.Request, limit int64) (error) {
	if limit <= 0 || limit > 10 {
		log.Println("[api] [Warning] Limit error", limit)
		err := errors.New(`"limit" can't be zero or more than 10`)
		api.HandleError(w, r, err)
		return err
	}
	return nil
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
	if api.CheckLimit(w, r, searchFoldersRequest.Limit) != nil {
		return
	}

	searchFoldersResponse := &SearchFoldersResponse{}

	searchFoldersResponse.Folders, err = api.Db.SearchFolders(searchFoldersRequest.Foldername, searchFoldersRequest.Limit, searchFoldersRequest.Offset)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	log.Println("[api] Search folders", searchFoldersRequest.Foldername, searchFoldersRequest.Limit, searchFoldersRequest.Offset)

	json.NewEncoder(w).Encode(searchFoldersResponse)
}

type GetFileRequest struct {
	ID int64 `json:"id"`
}

func (api *API) HandleGetFileStream(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	ids := q["id"]
	if len(ids) == 0 {
		api.HandleErrorString(w, r, `parameter "id" can't be empty`)
		return
	}
	id, err := strconv.Atoi(ids[0])
	if err != nil {
		api.HandleErrorString(w, r, `parameter "id" should be an integer`)
		return
	}
	file, err := api.Db.GetFile(int64(id))
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	path, err := file.Path()
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	log.Println("[api] Stream file", path)

	cmd := exec.Command("ffmpeg",
		"-i", path,
		"-c:a", "libopus",
		"-ab", "128k",
		"-vn",
		"-f", "matroska",
		"-",
	)
	cmd.Stdout = w
	err = cmd.Run()
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
}

func (api *API) HandleGetFileDirect(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	ids := q["id"]
	if len(ids) == 0 {
		api.HandleErrorString(w, r, `parameter "id" can't be empty`)
		return
	}
	id, err := strconv.Atoi(ids[0])
	if err != nil {
		api.HandleErrorString(w, r, `parameter "id" should be an integer`)
		return
	}
	file, err := api.Db.GetFile(int64(id))
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	path, err := file.Path()
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	log.Println("[api] Get direct raw file", path)

	http.ServeFile(w, r, path)
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

	log.Println("[api] Get pipe raw file", path)

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
	apiMux.HandleFunc("/get_file_direct", api.HandleGetFileDirect)
	apiMux.HandleFunc("/search_files", api.HandleSearchFiles)
	apiMux.HandleFunc("/search_folders", api.HandleSearchFolders)
	apiMux.HandleFunc("/get_files_in_folder", api.HandleGetFilesInFolder)
	apiMux.HandleFunc("/get_random_files", api.HandleGetRandomFiles)
	apiMux.HandleFunc("/get_file_stream", api.HandleGetFileStream)
	// below needs token
	apiMux.HandleFunc("/walk", api.HandleWalk)
	apiMux.HandleFunc("/reset", api.HandleReset)

	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", apiMux))
	mux.Handle("/web/", http.StripPrefix("/web", http.FileServer(http.Dir("web"))))

	api.token = "pwd"

	return api, nil
}
