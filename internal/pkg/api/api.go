package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"msw-open-music/internal/pkg/database"
	"msw-open-music/internal/pkg/tmpfs"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type API struct {
	Db *database.Database
	Server http.Server
	token string
	APIConfig APIConfig
	Tmpfs *tmpfs.Tmpfs
}

type FfmpegConfigList struct {
	FfmpegConfigList []FfmpegConfig `json:"ffmpeg_config_list"`
}

type AddFfmpegConfigRequest struct {
	Token string `json:"token"`
	Name string `json:"name"`
	FfmpegConfig FfmpegConfig `json:"ffmpeg_config"`
}

type FfmpegConfig struct {
	Name string `json:"name"`
	Args string `json:"args"`
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

func (api *API) HandleGetFileInfo(w http.ResponseWriter, r *http.Request) {
	getFileRequest := &GetFileRequest{
		ID: -1,
	}

	err := json.NewDecoder(r.Body).Decode(getFileRequest)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	// check empty
	if getFileRequest.ID < 0 {
		api.HandleErrorString(w, r, `"id" can't be none or negative`)
		return
	}

	file, err := api.Db.GetFile(getFileRequest.ID)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	err = json.NewEncoder(w).Encode(file)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
}

func (api *API) CheckGetFileStream(w http.ResponseWriter, r *http.Request) (error) {
	var err error
	q := r.URL.Query()
	ids := q["id"]
	if len(ids) == 0 {
		err = errors.New(`parameter "id" can't be empty`)
		api.HandleError(w, r, err)
		return err
	}
	_, err = strconv.Atoi(ids[0])
	if err != nil {
		err = errors.New(`parameter "id" should be an integer`)
		api.HandleError(w, r, err)
		return err
	}
	configs := q["config"]
	if len(configs) == 0 {
		err = errors.New(`parameter "config" can't be empty`)
		api.HandleError(w, r, err)
		return err
	}
	return nil
}

func (api *API) GetFfmpegConfig(configName string) (FfmpegConfig, bool) {
	ffmpegConfig := FfmpegConfig{}
	for _, f := range api.APIConfig.FfmpegConfigList {
		if f.Name == configName {
			ffmpegConfig = f
		}
	}
	if ffmpegConfig.Name == "" {
		return ffmpegConfig, false
	}
	return ffmpegConfig, true
}

func (api *API) HandleGetFileStream(w http.ResponseWriter, r *http.Request) {
	err := api.CheckGetFileStream(w, r)
	if err != nil {
		return
	}
	q := r.URL.Query()
	ids := q["id"]
	id, err := strconv.Atoi(ids[0])
	configs := q["config"]
	configName := configs[0]
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

	log.Println("[api] Stream file", path, configName)

	ffmpegConfig, ok := api.GetFfmpegConfig(configName)
	if !ok {
		api.HandleErrorStringCode(w, r, `ffmpeg config not found`, 404)
		return
	}
	args := strings.Split(ffmpegConfig.Args, " ")
	startArgs := []string {"-threads", strconv.FormatInt(api.APIConfig.FfmpegThreads, 10), "-i", path}
	endArgs := []string {"-vn", "-f", "ogg", "-"}
	ffmpegArgs := append(startArgs, args...)
	ffmpegArgs = append(ffmpegArgs, endArgs...)
	cmd := exec.Command("ffmpeg", ffmpegArgs...)
	cmd.Stdout = w
	err = cmd.Run()
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
}

type PrepareFileStreamDirectRequest struct {
	ID int64 `json:"id"`
	ConfigName string `json:"config_name"`
}

type PrepareFileStreamDirectResponse struct {
	Filesize int64 `json:"filesize"`
}

func (api *API) HandlePrepareFileStreamDirect(w http.ResponseWriter, r *http.Request) {
	prepareFileStreamDirectRequst := &PrepareFileStreamDirectRequest{
		ID: -1,
	}
	err := json.NewDecoder(r.Body).Decode(prepareFileStreamDirectRequst)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	// check empty
	if prepareFileStreamDirectRequst.ID < 0 {
		api.HandleErrorString(w, r, `"id" can't be none or negative`)
		return
	}
	if prepareFileStreamDirectRequst.ConfigName == "" {
		api.HandleErrorString(w, r, `"config_name" can't be empty`)
		return
	}

	file, err := api.Db.GetFile(prepareFileStreamDirectRequst.ID)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
	srcPath, err := file.Path()
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	log.Println("[api] Prepare stream direct file", srcPath, prepareFileStreamDirectRequst.ConfigName)
	ffmpegConfig, ok := api.GetFfmpegConfig(prepareFileStreamDirectRequst.ConfigName)
	if !ok {
		api.HandleErrorStringCode(w, r, `ffmpeg config not found`, 404)
		return
	}
	objPath := api.Tmpfs.GetObjFilePath(prepareFileStreamDirectRequst.ID, prepareFileStreamDirectRequst.ConfigName)

	// check obj file exists
	exists := api.Tmpfs.Exits(objPath)
	if exists {
		fileInfo, err := os.Stat(objPath)
		if err != nil {
			api.HandleError(w, r, err)
			return
		}
		prepareFileStreamDirectResponse := &PrepareFileStreamDirectResponse{
			Filesize: fileInfo.Size(),
		}
		json.NewEncoder(w).Encode(prepareFileStreamDirectResponse)
			return
	}

	api.Tmpfs.Record(objPath)
	args := strings.Split(ffmpegConfig.Args, " ")
	startArgs := []string {"-threads", strconv.FormatInt(api.APIConfig.FfmpegThreads, 10), "-i", srcPath}
	endArgs := []string {"-vn", "-y", objPath}
	ffmpegArgs := append(startArgs, args...)
	ffmpegArgs = append(ffmpegArgs, endArgs...)
	cmd := exec.Command("ffmpeg", ffmpegArgs...)
	err = cmd.Run()
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	fileInfo, err := os.Stat(objPath)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
	prepareFileStreamDirectResponse := &PrepareFileStreamDirectResponse{
		Filesize: fileInfo.Size(),
	}
	json.NewEncoder(w).Encode(prepareFileStreamDirectResponse)
}

func (api *API) HandleGetFileStreamDirect(w http.ResponseWriter, r *http.Request) {
	err := api.CheckGetFileStream(w, r)
	if err != nil {
		return
	}
	q := r.URL.Query()
	ids := q["id"]
	id, err := strconv.Atoi(ids[0])
	configs := q["config"]
	configName := configs[0]

	path := api.Tmpfs.GetObjFilePath(int64(id), configName)
	if api.Tmpfs.Exits(path) {
		api.Tmpfs.Record(path)
	}

	log.Println("[api] Get direct cached file", path)

	http.ServeFile(w, r, path)
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
	getFileRequest := &GetFileRequest{
		ID: -1,
	}

	err := json.NewDecoder(r.Body).Decode(getFileRequest)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	// check empty
	if getFileRequest.ID < 0 {
		api.HandleErrorString(w, r, `"id" can't be none or negative`)
		return
	}

	file, err := api.Db.GetFile(getFileRequest.ID)
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

func (api *API) HandleGetFfmpegConfigs(w http.ResponseWriter, r *http.Request) {
	log.Println("[api] Get ffmpeg config list")
	ffmpegConfigList:= &FfmpegConfigList{
		FfmpegConfigList: api.APIConfig.FfmpegConfigList,
	}
	json.NewEncoder(w).Encode(&ffmpegConfigList)
}

func (api *API) HandleAddFfmpegConfig(w http.ResponseWriter, r *http.Request) {
	addFfmpegConfigRequest := AddFfmpegConfigRequest{}
	err := json.NewDecoder(r.Body).Decode(&addFfmpegConfigRequest)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	// check token
	err = api.CheckToken(w, r, addFfmpegConfigRequest.Token)
	if err != nil {
		return
	}

	// check name and args not null
	if addFfmpegConfigRequest.Name == "" {
		api.HandleErrorString(w, r, `"ffmpeg_config.name" can't be empty`)
		return
	}
	if addFfmpegConfigRequest.FfmpegConfig.Args == "" {
		api.HandleErrorString(w, r, `"ffmpeg_config.args" can't be empty`)
		return
	}

	log.Println("[api] Add ffmpeg config")

	api.APIConfig.FfmpegConfigList = append(api.APIConfig.FfmpegConfigList, addFfmpegConfigRequest.FfmpegConfig)

	api.HandleOK(w, r)
}

type FeedbackRequest struct {
	Feedback string `json:"feedback"`
}

func (api *API) HandleFeedback(w http.ResponseWriter, r *http.Request) {
	feedbackRequest := &FeedbackRequest{}
	err :=json.NewDecoder(r.Body).Decode(feedbackRequest)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	// check empty feedback
	if feedbackRequest.Feedback == "" {
		api.HandleErrorString(w, r, `"feedback" can't be empty`)
		return
	}

	log.Println("[api] Feedback", feedbackRequest.Feedback)

	headerBuff := &bytes.Buffer{}
	err = r.Header.Write(headerBuff)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
	header := headerBuff.String()

	err = api.Db.InsertFeedback(time.Now().Unix(), feedbackRequest.Feedback, header)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
	api.HandleOK(w, r)
}

func NewAPIConfig() (APIConfig) {
	apiConfig := APIConfig{}
	return apiConfig
}

type APIConfig struct {
	DatabaseName string `json:"database_name"`
	Addr string `json:"addr"`
	Token string `json:"token"`
	FfmpegThreads int64 `json:"ffmpeg_threads"`
	FfmpegConfigList []FfmpegConfig `json:"ffmpeg_config_list"`
}

type Config struct {
	APIConfig APIConfig `json:"api"`
	TmpfsConfig tmpfs.TmpfsConfig `json:"tmpfs"`
}

func NewAPI(config Config) (*API, error) {
	var err error

	apiConfig := config.APIConfig
	tmpfsConfig := config.TmpfsConfig

	db, err := database.NewDatabase(apiConfig.DatabaseName)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	apiMux := http.NewServeMux()

	api := &API{
		Db: db,
		Server: http.Server{
			Addr: apiConfig.Addr,
			Handler: mux,
		},
		APIConfig: apiConfig,
	}
	api.Tmpfs = tmpfs.NewTmpfs(tmpfsConfig)

	// mount api
	apiMux.HandleFunc("/hello", api.HandleOK)
	apiMux.HandleFunc("/get_file", api.HandleGetFile)
	apiMux.HandleFunc("/get_file_direct", api.HandleGetFileDirect)
	apiMux.HandleFunc("/search_files", api.HandleSearchFiles)
	apiMux.HandleFunc("/search_folders", api.HandleSearchFolders)
	apiMux.HandleFunc("/get_files_in_folder", api.HandleGetFilesInFolder)
	apiMux.HandleFunc("/get_random_files", api.HandleGetRandomFiles)
	apiMux.HandleFunc("/get_file_stream", api.HandleGetFileStream)
	apiMux.HandleFunc("/get_ffmpeg_config_list", api.HandleGetFfmpegConfigs)
	apiMux.HandleFunc("/feedback", api.HandleFeedback)
	apiMux.HandleFunc("/get_file_info", api.HandleGetFileInfo)
	apiMux.HandleFunc("/get_file_stream_direct", api.HandleGetFileStreamDirect)
	apiMux.HandleFunc("/prepare_file_stream_direct", api.HandlePrepareFileStreamDirect)
	// below needs token
	apiMux.HandleFunc("/walk", api.HandleWalk)
	apiMux.HandleFunc("/reset", api.HandleReset)
	apiMux.HandleFunc("/add_ffmpeg_config", api.HandleAddFfmpegConfig)

	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", apiMux))
	mux.Handle("/web/", http.StripPrefix("/web", http.FileServer(http.Dir("web"))))

	api.token = apiConfig.Token

	return api, nil
}
