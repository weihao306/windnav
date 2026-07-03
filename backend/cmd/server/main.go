package main

import (
	"log"

	"windnav/backend/internal/config"
	"windnav/backend/internal/database"
	"windnav/backend/internal/httpapi"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	db, err := database.Open(cfg)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		log.Fatalf("migrate database: %v", err)
	}
	if err := database.Seed(db, cfg); err != nil {
		log.Fatalf("seed database: %v", err)
	}

	router := httpapi.NewRouter(cfg, db)
	if err := router.Run(cfg.HTTPAddr); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
