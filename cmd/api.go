package main

import (
	"log"

	"faceit/appdata"
	"faceit/database"
	"faceit/server"
)

func main() {
	appdata.InitApp("api")
	err := database.EnsureDbSchema()
	if err != nil {
		log.Fatal(err)
	}

	server.Run()
}
