package main

import (
	"log"

	"github.com/guncv/Poll-Voting-Website/backend/api"
	"github.com/guncv/Poll-Voting-Website/backend/db"
	"github.com/guncv/Poll-Voting-Website/backend/util"
)

func main() {
	config, err := util.LoadConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	db := db.InitDB(config)
	server := api.NewServer(config, db)

	if err := server.Start(config.ServerAddress); err != nil {
		log.Fatal("cannot start server:", err)
	}
}
