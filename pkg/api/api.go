package api

import (
	"msw-open-music/pkg/database"
	"msw-open-music/pkg/tmpfs"
	"net/http"
)

type API struct {
	Db        *database.Database
	Server    http.Server
	token     string
	APIConfig APIConfig
	Tmpfs     *tmpfs.Tmpfs
}

func NewAPIConfig() APIConfig {
	apiConfig := APIConfig{}
	return apiConfig
}

type APIConfig struct {
	DatabaseName     string         `json:"database_name"`
	Addr             string         `json:"addr"`
	Token            string         `json:"token"`
	FfmpegThreads    int64          `json:"ffmpeg_threads"`
	FfmpegConfigList []FfmpegConfig `json:"ffmpeg_config_list"`
}

type Config struct {
	APIConfig   APIConfig         `json:"api"`
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
			Addr:    apiConfig.Addr,
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
	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("web/build"))))

	api.token = apiConfig.Token

	return api, nil
}