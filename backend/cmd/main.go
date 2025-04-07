package main

import (
	"log"

	"github.com/guncv/Poll-Voting-Website/backend/config"
	"github.com/guncv/Poll-Voting-Website/backend/controller"
	"github.com/guncv/Poll-Voting-Website/backend/db"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	database := db.InitDB(*config)
	cacheService := db.NewRedisCacheService(*config)

	server := controller.NewServer(*config, database, cacheService)

	if err := server.Start(config.ServerAddress); err != nil {
		log.Fatal("cannot start server:", err)
	}
}
