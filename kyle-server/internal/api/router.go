package api

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
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

	// Dev OIDC mock endpoints (for local development without identity provider)
	r.mux.HandleFunc("GET /.well-known/openid-configuration", r.handleOIDCConfig)
	r.mux.HandleFunc("GET /authorize", r.handleAuthorize)
	r.mux.HandleFunc("POST /oauth/token", r.handleToken)
	r.mux.HandleFunc("GET /userinfo", r.handleUserInfo)
	r.mux.HandleFunc("GET /api/instance", r.handleInstance)

	// Placeholder API endpoints (return empty arrays for dashboard to render)
	r.mux.HandleFunc("GET /api/users", r.handleList("users"))
	r.mux.HandleFunc("GET /api/peers", r.handleList("peers"))
	r.mux.HandleFunc("GET /api/groups", r.handleList("groups"))
	r.mux.HandleFunc("GET /api/setup-keys", r.handleList("setup-keys"))
	r.mux.HandleFunc("GET /api/nameservers", r.handleList("nameservers"))
	r.mux.HandleFunc("GET /api/routes", r.handleList("routes"))
	r.mux.HandleFunc("GET /api/dns/settings", r.handleDNS("dns/settings"))
	r.mux.HandleFunc("GET /api/events", r.handleList("events"))
	r.mux.HandleFunc("GET /api/accounts", r.handleList("accounts"))
	r.mux.HandleFunc("/api/v1/servers/", r.handleList("servers"))
	r.mux.HandleFunc("/api/v1/users/", r.handleList("users"))
	r.mux.HandleFunc("/api/v1/connections/", r.handleList("connections"))
	r.mux.HandleFunc("/api/v1/policies/", r.handleList("policies"))
	r.mux.HandleFunc("/api/v1/audit/", r.handleList("audit"))

	return r
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	corsHeaders(w)
	if req.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	r.mux.ServeHTTP(w, req)
}

func corsHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

// Dev OIDC: discovery endpoint
func (r *Router) handleOIDCConfig(w http.ResponseWriter, req *http.Request) {
	corsHeaders(w)
	base := "http://" + req.Host
	json.NewEncoder(w).Encode(map[string]interface{}{
		"issuer":                 base,
		"authorization_endpoint": base + "/authorize",
		"token_endpoint":         base + "/oauth/token",
		"userinfo_endpoint":      base + "/userinfo",
		"jwks_uri":               base + "/.well-known/jwks.json",
		"scopes_supported":       []string{"openid", "profile", "email"},
		"response_types_supported": []string{"code", "token", "id_token"},
	})
}

// Dev OIDC: authorization redirect
func (r *Router) handleAuthorize(w http.ResponseWriter, req *http.Request) {
	corsHeaders(w)
	redirectURI := req.URL.Query().Get("redirect_uri")
	if redirectURI == "" {
		redirectURI = req.Referer()
	}
	code := randomString(32)
	http.Redirect(w, req, redirectURI+"?code="+code+"&state="+req.URL.Query().Get("state"), http.StatusFound)
}

// Dev OIDC: token exchange
func (r *Router) handleToken(w http.ResponseWriter, req *http.Request) {
	corsHeaders(w)
	if req.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token":  "dev-access-token-" + randomString(16),
		"id_token":      "dev-id-token-" + randomString(16),
		"refresh_token": "dev-refresh-token-" + randomString(16),
		"token_type":    "Bearer",
		"expires_in":    3600,
	})
}

// Dev OIDC: user info
func (r *Router) handleUserInfo(w http.ResponseWriter, req *http.Request) {
	corsHeaders(w)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"sub":   "dev-user",
		"name":  "Dev User",
		"email": "dev@kyle.local",
	})
}

// Instance info (called by dashboard on startup)
func (r *Router) handleInstance(w http.ResponseWriter, req *http.Request) {
	corsHeaders(w)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"name":    "Kyle Dev",
		"version": "0.1.0",
		"setup":   true,
	})
}

func randomString(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)[:n]
}

func (r *Router) handleHealth(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	dbStatus := "connected"
	if r.db == nil {
		dbStatus = "unavailable"
	} else if err := r.db.Ping(); err != nil {
		dbStatus = "disconnected"
	}
	status := "ok"
	httpStatus := http.StatusOK
	if dbStatus != "connected" {
		status = "degraded"
		httpStatus = http.StatusOK // still 200, just degraded
	}
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(map[string]string{
		"status":   status,
		"database": dbStatus,
	})
}

func (r *Router) handleStatus(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"version": "0.1.0",
		"status":  "running",
	})
}

func (r *Router) handleList(_ string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("[]"))
	}
}

func (r *Router) handleDNS(_ string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}
}
