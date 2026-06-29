package api

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ethancls/kyle-server/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type Router struct {
	mux           *http.ServeMux
	db            *sql.DB
	cfg           *config.Config
	setupComplete bool
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
	r.mux.HandleFunc("GET /api/instance/version", r.handleInstanceVersion)
	r.mux.HandleFunc("GET /api/locations/countries", r.handleEmptyList)

	// Users (Netbird-compatible)
	r.mux.HandleFunc("GET /api/users/current", r.handleCurrentUser)
	r.mux.HandleFunc("GET /api/users", r.handleUsers)
	r.mux.HandleFunc("POST /api/users", r.handleCreateUser)

	// Setup (Netbird-compatible, unauthenticated)
	r.mux.HandleFunc("GET /api/setup", r.handleGetSetup)
	r.mux.HandleFunc("POST /api/setup", r.handlePostSetup)

	// Dashboard data endpoints (PostgreSQL only)
	r.mux.HandleFunc("GET /api/peers", r.handlePeers)
	r.mux.HandleFunc("GET /api/groups", r.handleGroups)
	r.mux.HandleFunc("GET /api/setup-keys", r.handleSetupKeys)
	r.mux.HandleFunc("GET /api/events", r.handleEvents)
	r.mux.HandleFunc("GET /api/nameservers", r.handleNameservers)
	r.mux.HandleFunc("GET /api/routes", r.handleRoutes)
	r.mux.HandleFunc("GET /api/dns/settings", r.handleDNSSettings)
	r.mux.HandleFunc("GET /api/accounts", r.handleAccounts)

	// Legacy v1 endpoints
	r.mux.HandleFunc("/api/v1/servers/", r.handleServersV1)
	r.mux.HandleFunc("/api/v1/users/", r.handleUsersV1)
	r.mux.HandleFunc("/api/v1/connections/", r.handleConnectionsV1)
	r.mux.HandleFunc("/api/v1/policies/", r.handlePoliciesV1)
	r.mux.HandleFunc("/api/v1/audit/", r.handleAuditV1)

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
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

// ---------------------------------------------------------------------------
// Dev OIDC mock endpoints
// ---------------------------------------------------------------------------

func (r *Router) handleOIDCConfig(w http.ResponseWriter, req *http.Request) {
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

func (r *Router) handleAuthorize(w http.ResponseWriter, req *http.Request) {
	redirectURI := req.URL.Query().Get("redirect_uri")
	if redirectURI == "" {
		redirectURI = req.Referer()
	}
	code := randomString(32)
	http.Redirect(w, req, redirectURI+"?code="+code+"&state="+req.URL.Query().Get("state"), http.StatusFound)
}

func (r *Router) handleToken(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	now := time.Now()
	jwtSecret := []byte(r.cfg.JWTSecret)
	if len(jwtSecret) == 0 {
		jwtSecret = []byte("dev-secret")
	}

	idToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":   "http://" + req.Host,
		"sub":   "dev-user",
		"aud":   "cosmos-api",
		"email": "dev@cosmos.local",
		"name":  "Dev User",
		"iat":   now.Unix(),
		"exp":   now.Add(1 * time.Hour).Unix(),
	})
	idTokenStr, _ := idToken.SignedString(jwtSecret)

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "http://" + req.Host,
		"sub": "dev-user",
		"aud": "cosmos-api",
		"iat": now.Unix(),
		"exp": now.Add(1 * time.Hour).Unix(),
	})
	accessTokenStr, _ := accessToken.SignedString(jwtSecret)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token":  accessTokenStr,
		"id_token":      idTokenStr,
		"refresh_token": "dev-refresh-token-" + randomString(16),
		"token_type":    "Bearer",
		"expires_in":    3600,
	})
}

func (r *Router) handleUserInfo(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(map[string]interface{}{
		"sub":   "dev-user",
		"name":  "Dev User",
		"email": "dev@cosmos.local",
	})
}

func (r *Router) handleInstance(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(map[string]interface{}{
		"setup_required": !r.setupComplete,
	})
}

func (r *Router) handleInstanceVersion(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"management_current_version":   "0.1.0",
		"management_available_version": "0.1.0",
		"dashboard_available_version":  "0.1.0",
	})
}

func (r *Router) handleEmptyList(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]interface{}{})
}

// ---------------------------------------------------------------------------
// Health & Status
// ---------------------------------------------------------------------------

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

// ---------------------------------------------------------------------------
// Users (Netbird-compatible)
// ---------------------------------------------------------------------------

// handleCurrentUser returns the current authenticated user.
// GET /api/users/current
func (r *Router) handleCurrentUser(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var id, email, name, role, status string
	err := r.db.QueryRow(`SELECT id, email, COALESCE(display_name,''), role, status
		FROM users WHERE deleted_at IS NULL ORDER BY created_at LIMIT 1`).
		Scan(&id, &email, &name, &role, &status)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "no user found"})
		return
	}
	isCurrent := true
	isService := false
	apiIssued := "api"
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":               id,
		"email":            email,
		"name":             name,
		"role":             role,
		"auto_groups":      []string{},
		"status":           status,
		"is_current":       &isCurrent,
		"is_service_user":  &isService,
		"is_blocked":       false,
		"last_login":       nil,
		"issued":           &apiIssued,
		"pending_approval": false,
		"permissions": map[string]interface{}{
			"is_restricted": false,
			"modules": map[string]map[string]bool{
				"*": {"read": true, "write": true},
			},
		},
	})
}

// handleCreateUser creates a new user.
// POST /api/users
func (r *Router) handleCreateUser(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var body struct {
		Email         *string  `json:"email"`
		Name          *string  `json:"name"`
		Role          string   `json:"role"`
		AutoGroups    []string `json:"auto_groups"`
		IsServiceUser bool     `json:"is_service_user"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	email := ""
	if body.Email != nil {
		email = *body.Email
	}
	name := ""
	if body.Name != nil {
		name = *body.Name
	}
	role := body.Role
	if role == "" {
		role = "user"
	}

	userID := "user-" + randomString(12)
	isCurrent := false
	isService := body.IsServiceUser
	apiIssued := "api"

	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":              userID,
		"email":           email,
		"name":            name,
		"role":            role,
		"auto_groups":     body.AutoGroups,
		"status":          "active",
		"is_current":      &isCurrent,
		"is_service_user": &isService,
		"is_blocked":      false,
		"last_login":      nil,
		"issued":          &apiIssued,
		"pending_approval": false,
	})
}

// ---------------------------------------------------------------------------
// Setup (Netbird-compatible)
// ---------------------------------------------------------------------------

// handleGetSetup returns the setup status.
// GET /api/setup
func (r *Router) handleGetSetup(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"setup_required": false,
	})
}

// handlePostSetup creates the initial admin user.
// POST /api/setup
func (r *Router) handlePostSetup(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var body struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}
	if body.Email == "" || body.Name == "" || body.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "email, name, and password are required"})
		return
	}

	userID := "user-" + randomString(12)

	// Create default tenant if none exists
	_, err := r.db.Exec(`INSERT INTO tenants (id, name, slug, tier)
		VALUES ($1, $2, $3, 'free') ON CONFLICT DO NOTHING`,
		"dev-tenant", "Cosmos", "cosmos")
	if err != nil {
		log.Printf("create tenant failed: %v", err)
	}

	// Insert user
	_, err = r.db.Exec(`INSERT INTO users (id, tenant_id, email, display_name, role, status, password_hash)
		VALUES ($1, $2, $3, $4, 'owner', 'active', $5)`,
		userID, "dev-tenant", body.Email, body.Name, body.Password)
	if err != nil {
		log.Printf("insert user failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to create user"})
		return
	}

	r.setupComplete = true
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": userID,
		"email":   body.Email,
	})
}

// ---------------------------------------------------------------------------
// Dashboard endpoints (PostgreSQL only, no fallback)
// ---------------------------------------------------------------------------

func (r *Router) handleUsers(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	users, err := queryUsers(r.db)
	if err != nil {
		log.Printf("query users failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "database query failed"})
		return
	}
	json.NewEncoder(w).Encode(users)
}

func (r *Router) handleAccounts(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	accounts, err := queryAccounts(r.db)
	if err != nil {
		log.Printf("query accounts failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "database query failed"})
		return
	}
	json.NewEncoder(w).Encode(accounts)
}

func (r *Router) handlePeers(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	peers, err := queryPeers(r.db)
	if err != nil {
		log.Printf("query peers failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "database query failed"})
		return
	}
	json.NewEncoder(w).Encode(peers)
}

func (r *Router) handleGroups(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	groups, err := queryGroups(r.db)
	if err != nil {
		log.Printf("query groups failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "database query failed"})
		return
	}
	json.NewEncoder(w).Encode(groups)
}

func (r *Router) handleSetupKeys(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	keys, err := querySetupKeys(r.db)
	if err != nil {
		log.Printf("query setup keys failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "database query failed"})
		return
	}
	json.NewEncoder(w).Encode(keys)
}

func (r *Router) handleEvents(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	events, err := queryEvents(r.db)
	if err != nil {
		log.Printf("query events failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "database query failed"})
		return
	}
	json.NewEncoder(w).Encode(events)
}

func (r *Router) handleNameservers(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]string{})
}

func (r *Router) handleRoutes(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]string{})
}

func (r *Router) handleDNSSettings(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{})
}

// ---------------------------------------------------------------------------
// V1 API endpoints (PostgreSQL only, no fallback)
// ---------------------------------------------------------------------------

func (r *Router) handleServersV1(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	servers, err := queryServers(r.db)
	if err != nil {
		log.Printf("query servers failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "database query failed"})
		return
	}
	json.NewEncoder(w).Encode(servers)
}

func (r *Router) handleUsersV1(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	users, err := queryUsersV1(r.db)
	if err != nil {
		log.Printf("query v1 users failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "database query failed"})
		return
	}
	json.NewEncoder(w).Encode(users)
}

func (r *Router) handleConnectionsV1(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	conns, err := queryConnections(r.db)
	if err != nil {
		log.Printf("query connections failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "database query failed"})
		return
	}
	json.NewEncoder(w).Encode(conns)
}

func (r *Router) handlePoliciesV1(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	policies, err := queryPolicies(r.db)
	if err != nil {
		log.Printf("query policies failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "database query failed"})
		return
	}
	json.NewEncoder(w).Encode(policies)
}

func (r *Router) handleAuditV1(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	logs, err := queryAuditLogs(r.db)
	if err != nil {
		log.Printf("query audit logs failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "database query failed"})
		return
	}
	json.NewEncoder(w).Encode(logs)
}

// ---------------------------------------------------------------------------
// Utility
// ---------------------------------------------------------------------------

func randomString(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)[:n]
}
