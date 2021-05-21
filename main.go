package main

import (
	"log"
	"msw-open-music/internal/pkg/api"
)

func main() {
	api, err := api.NewAPI("/tmp/test.sqlite3", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Started")
	log.Fatal(api.Server.ListenAndServe())
}
