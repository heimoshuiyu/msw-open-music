package main

import (
	"encoding/json"
	"flag"
	"log"
	"msw-open-music/pkg/api"
	"msw-open-music/pkg/commonconfig"
	"os"
)

var ConfigFilePath string

func init() {
	flag.StringVar(&ConfigFilePath, "config", "config.json", "backend config file path")
}

func main() {
	var err error
	flag.Parse()

	config := commonconfig.Config{}
	configFile, err := os.Open(ConfigFilePath)
	if err != nil {
		log.Fatal(err)
	}
	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		log.Fatal(err)
	}
	configFile.Close()

	api, err := api.NewAPI(config)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Starting",
		config.APIConfig.DatabaseName,
		config.APIConfig.Addr,
	)
	log.Fatal(api.Server.ListenAndServe())
}
