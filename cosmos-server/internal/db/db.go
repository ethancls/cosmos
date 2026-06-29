package db

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/lib/pq"
)

// Connect opens a PostgreSQL database connection.
// Returns an error if PostgreSQL is unavailable — the caller handles the fallback to in-memory mock data.
func Connect(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("sql open: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db ping: %w", err)
	}
	log.Println("Connected to PostgreSQL")
	return db, nil
}

// RunMigrations runs SQL migration files from the given directory.
func RunMigrations(database *sql.DB, migrationsDir string) error {
	var files []string
	err := filepath.WalkDir(migrationsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.EqualFold(filepath.Ext(path), ".sql") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Migrations directory not found, skipping migrations")
			return nil
		}
		return fmt.Errorf("walk migrations dir: %w", err)
	}
	if len(files) == 0 {
		log.Println("No migration files found")
		return nil
	}
	sort.Strings(files)
	for _, f := range files {
		log.Printf("Running migration: %s", filepath.Base(f))
		content, err := os.ReadFile(f)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", f, err)
		}
		if _, err := database.Exec(string(content)); err != nil {
			return fmt.Errorf("exec migration %s: %w", f, err)
		}
	}
	log.Printf("Applied %d migration(s)", len(files))
	return nil
}
