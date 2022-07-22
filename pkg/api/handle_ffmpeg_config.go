package api

import (
	"encoding/json"
	"log"
	"msw-open-music/pkg/commonconfig"
	"net/http"
)

func (api *API) GetFfmpegConfig(configName string) (commonconfig.FfmpegConfig, bool) {
	ffmpegConfig := commonconfig.FfmpegConfig{}
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
	ffmpegConfigList := &commonconfig.FfmpegConfigList{
		FfmpegConfigList: api.APIConfig.FfmpegConfigList,
	}
	json.NewEncoder(w).Encode(&ffmpegConfigList)
}
