package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

// ConnectSQLite opens a SQLite database connection using the data/ directory.
func ConnectSQLite() (*sql.DB, error) {
	if err := os.MkdirAll("data", 0755); err != nil {
		return nil, fmt.Errorf("mkdir data dir: %w", err)
	}
	db, err := sql.Open("sqlite", "file:data/kyle.db?cache=shared&_journal_mode=WAL")
	if err != nil {
		return nil, fmt.Errorf("sqlite open: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("sqlite ping: %w", err)
	}
	log.Println("Connected to SQLite")
	return db, nil
}
