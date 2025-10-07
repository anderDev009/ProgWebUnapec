package main

import (
	"log"

	"petmatch/internal/config"
	"petmatch/internal/database"
	"petmatch/internal/router"
)

func main() {
	cfg := config.Load()

	db, err := database.Open(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	server, err := router.New(db, cfg)
	if err != nil {
		log.Fatalf("failed to configure router: %v", err)
	}

	log.Printf("PetMatch API listening on port %s", cfg.HTTPPort)

	if err := server.Run(":" + cfg.HTTPPort); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
