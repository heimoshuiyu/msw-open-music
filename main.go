package main

import (
	"encoding/json"
	"flag"
	"log"
	"msw-open-music/internal/pkg/api"
	"os"
)

var APIConfigFilePath string

func init() {
	flag.StringVar(&APIConfigFilePath, "apiconfig", "api_config.json", "API Config Json file")
}

func main() {
	var err error
	flag.Parse()
	apiConfig := api.NewAPIConfig()

	apiConfigFile, err := os.Open(APIConfigFilePath)
	if err != nil {
		log.Fatal(err)
	}
	err = json.NewDecoder(apiConfigFile).Decode(&apiConfig)
	if err != nil {
		log.Fatal(err)
	}
	apiConfigFile.Close()

	api, err := api.NewAPI(apiConfig)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Starting", apiConfig.DatabaseName, apiConfig.Addr, apiConfig.Token)
	log.Fatal(api.Server.ListenAndServe())
}
