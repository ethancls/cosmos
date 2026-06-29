package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/ethancls/kyle-server/internal/config"
)

type Router struct {
	mux *http.ServeMux
	db  *sql.DB
	cfg *config.Config
}

func SetupRouter(database *sql.DB, cfg *config.Config) http.Handler {
	r := &Router{
		mux: http.NewServeMux(),
		db:  database,
		cfg: cfg,
	}

	r.mux.HandleFunc("GET /health", r.handleHealth)
	r.mux.HandleFunc("GET /api/v1/status", r.handleStatus)

	// Placeholder route groups
	r.mux.HandleFunc("/api/v1/servers/", r.handlePlaceholder("servers"))
	r.mux.HandleFunc("/api/v1/users/", r.handlePlaceholder("users"))
	r.mux.HandleFunc("/api/v1/connections/", r.handlePlaceholder("connections"))
	r.mux.HandleFunc("/api/v1/policies/", r.handlePlaceholder("policies"))
	r.mux.HandleFunc("/api/v1/audit/", r.handlePlaceholder("audit"))

	return r
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func (r *Router) handleHealth(w http.ResponseWriter, req *http.Request) {
	err := r.db.Ping()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{
			"status":   "unhealthy",
			"database": "disconnected",
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (r *Router) handleStatus(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"version": "0.1.0",
		"status":  "running",
	})
}

func (r *Router) handlePlaceholder(domain string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": domain + " endpoint (not yet implemented)",
		})
	}
}
