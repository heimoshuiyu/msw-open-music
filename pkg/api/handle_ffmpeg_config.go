package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type FfmpegConfig struct {
	Name string `json:"name"`
	Args string `json:"args"`
}

type FfmpegConfigList struct {
	FfmpegConfigList []FfmpegConfig `json:"ffmpeg_config_list"`
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

func (api *API) HandleGetFfmpegConfigs(w http.ResponseWriter, r *http.Request) {
	log.Println("[api] Get ffmpeg config list")
	ffmpegConfigList := &FfmpegConfigList{
		FfmpegConfigList: api.APIConfig.FfmpegConfigList,
	}
	json.NewEncoder(w).Encode(&ffmpegConfigList)
}

type AddFfmpegConfigRequest struct {
	Token        string       `json:"token"`
	Name         string       `json:"name"`
	FfmpegConfig FfmpegConfig `json:"ffmpeg_config"`
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
