package main

import (
	"log"
	"net/http"

	"github.com/ethancls/cosmos-server/internal/api"
	"github.com/ethancls/cosmos-server/internal/config"
	"github.com/ethancls/cosmos-server/internal/db"
)

func main() {
	cfg := config.Load()

	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer database.Close()
	if err := db.RunMigrations(database, "migrations"); err != nil {
		log.Printf("WARNING: migrations failed: %v", err)
	}

	router := api.SetupRouter(database, cfg)

	log.Printf("Cosmos server starting on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatal(err)
	}
}
