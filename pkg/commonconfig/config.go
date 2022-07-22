package commonconfig

type Config struct {
	APIConfig   APIConfig   `json:"api"`
	TmpfsConfig TmpfsConfig `json:"tmpfs"`
}

type APIConfig struct {
	DatabaseName     string         `json:"database_name"`
	SingleThread     bool           `json:"single_thread,default=true"`
	Addr             string         `json:"addr"`
	FfmpegThreads    int64          `json:"ffmpeg_threads"`
	FfmpegConfigList []FfmpegConfig `json:"ffmpeg_config_list"`
	SECRET           string         `json:"secret"`
}

type FfmpegConfigList struct {
	FfmpegConfigList []FfmpegConfig `json:"ffmpeg_config_list"`
}

type FfmpegConfig struct {
	Name   string `json:"name"`
	Args   string `json:"args"`
	Format string `json:"format"`
}

type TmpfsConfig struct {
	FileLifeTime    int64  `json:"file_life_time"`
	CleanerInternal int64  `json:"cleaner_internal"`
	Root            string `json:"root"`
}

// Constructors for Config

func NewAPIConfig() APIConfig {
	apiConfig := APIConfig{}
	return apiConfig
}

func NewTmpfsConfig() *TmpfsConfig {
	config := &TmpfsConfig{}
	return config
}
