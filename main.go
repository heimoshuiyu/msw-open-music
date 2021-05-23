package main

import (
	"flag"
	"log"
	"msw-open-music/internal/pkg/api"
)

var DatabaseName string
var Listen string
var Token string

func init() {
	flag.StringVar(&DatabaseName, "db", "/tmp/music.sqlite3", "sqlite3 database file path")
	flag.StringVar(&Listen, "listen", ":8080", "http server listening")
	flag.StringVar(&Token, "token", "mikusavetheworld", "secret token")
}

func main() {
	flag.Parse()
	apiConfig := api.NewAPIConfig()
	apiConfig.FfmpegConfigs["libopus 128k"] = &api.FfmpegConfig{
		Name: "libopus 128k",
		Args: "-c:a libopus -ab 128k",
	}
	apiConfig.FfmpegConfigs["libopus 256k"] = &api.FfmpegConfig{
		Name: "libopus 256k",
		Args: "-c:a libopus -ab 256k",
	}
	apiConfig.DatabaseName = DatabaseName
	apiConfig.Addr = Listen
	apiConfig.Token = Token
	api, err := api.NewAPI(apiConfig)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Started")
	log.Fatal(api.Server.ListenAndServe())
}
