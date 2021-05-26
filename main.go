package main

import (
	"encoding/json"
	"flag"
	"log"
	"msw-open-music/internal/pkg/api"
	"msw-open-music/internal/pkg/tmpfs"
	"os"
)

var ConfigFilePath string

func init() {
	flag.StringVar(&ConfigFilePath, "config", "config.json", "backend config file path")
}

type Config struct {
	APIConfig api.APIConfig `json:"api"`
	TmpfsConfig tmpfs.TmpfsConfig `json:"tmpfs"`
}

func main() {
	var err error
	flag.Parse()

	config := Config{}
	configFile, err := os.Open(ConfigFilePath)
	if err != nil {
		log.Fatal(err)
	}
	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		log.Fatal(err)
	}
	configFile.Close()

	api, err := api.NewAPI(config.APIConfig, config.TmpfsConfig)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Starting",
		config.APIConfig.DatabaseName,
		config.APIConfig.Addr,
		config.APIConfig.Token,
	)
	log.Fatal(api.Server.ListenAndServe())
}
