package main

import (
	"log"
	"net/http"

	"github.com/ethancls/kyle-server/internal/api"
	"github.com/ethancls/kyle-server/internal/config"
	"github.com/ethancls/kyle-server/internal/db"
)

func main() {
	cfg := config.Load()

	// Try PostgreSQL first, then SQLite, then in-memory fallback
	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Printf("WARNING: PostgreSQL not available: %v", err)
		database, err = db.ConnectSQLite()
		if err != nil {
			log.Printf("WARNING: SQLite not available: %v (server will start with in-memory fallback)", err)
			database = nil
		}
	}
	if database != nil {
		defer database.Close()
		if err := db.RunMigrations(database, "migrations"); err != nil {
			log.Printf("WARNING: migrations failed: %v", err)
		}
	}

	router := api.SetupRouter(database, cfg)

	log.Printf("Kyle server starting on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatal(err)
	}
}
