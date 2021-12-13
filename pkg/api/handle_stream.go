package api

import (
	"encoding/json"
	"errors"
	"log"
	"msw-open-music/pkg/database"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func (api *API) CheckGetFileStream(w http.ResponseWriter, r *http.Request) error {
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

// /get_file_stream?id=1&config=ffmpeg_config_name
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
	startArgs := []string{"-threads", strconv.FormatInt(api.APIConfig.FfmpegThreads, 10), "-i", path}
	endArgs := []string{"-vn", "-f", "ogg", "-"}
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
	ID         int64  `json:"id"`
	ConfigName string `json:"config_name"`
}

type PrepareFileStreamDirectResponse struct {
	File *database.File `json:"file"`
}

// /prepare_file_stream_direct?id=1&config=ffmpeg_config_name
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
	if !exists {
		// lock the object
		api.Tmpfs.Lock(objPath)

		args := strings.Split(ffmpegConfig.Args, " ")
		startArgs := []string{"-threads", strconv.FormatInt(api.APIConfig.FfmpegThreads, 10), "-i", srcPath}
		endArgs := []string{"-vn", "-y", objPath}
		ffmpegArgs := append(startArgs, args...)
		ffmpegArgs = append(ffmpegArgs, endArgs...)
		cmd := exec.Command("ffmpeg", ffmpegArgs...)
		err = cmd.Run()
		if err != nil {
			api.HandleError(w, r, err)
			return
		}

		api.Tmpfs.Record(objPath)
		api.Tmpfs.Unlock(objPath)

	}

	fileInfo, err := os.Stat(objPath)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	file.Filesize = fileInfo.Size()

	prepareFileStreamDirectResponse := &PrepareFileStreamDirectResponse{
		File: file,
	}
	json.NewEncoder(w).Encode(prepareFileStreamDirectResponse)
}

// /get_file_stream_direct?id=1&config=ffmpeg_config_name
// return converted file with http.ServeFile method
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
