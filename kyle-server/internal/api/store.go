package api

import (
	"database/sql"
	"time"
)

// ---------------------------------------------------------------------------
// Type definitions matching Netbird API response shapes
// ---------------------------------------------------------------------------

type apiUser struct {
	ID            string   `json:"id"`
	Email         string   `json:"email"`
	Name          string   `json:"name"`
	Role          string   `json:"role"`
	AutoGroups    []string `json:"auto_groups"`
	Status        string   `json:"status"`
	IsServiceUser bool     `json:"is_service_user"`
	IsBlocked     bool     `json:"is_blocked"`
	NonDeletable  bool     `json:"non_deletable"`
	LastLogin     string   `json:"last_login"`
	Issued        string   `json:"issued"`
}

type apiAccount struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	Tier      string `json:"tier"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

type apiPeer struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
	OS       string `json:"os"`
	Status   string `json:"status"`
	Version  string `json:"version"`
}

type apiGroup struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type apiSetupKey struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Key       string `json:"key"`
	ExpiresAt string `json:"expires_at"`
	CreatedAt string `json:"created_at"`
}

type apiRoute struct {
	ID          string `json:"id"`
	Network     string `json:"network"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

type apiNameserver struct {
	ID  string `json:"id"`
	IP  string `json:"ip"`
	Port int   `json:"port"`
}

type apiEvent struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Action    string `json:"action"`
	CreatedAt string `json:"created_at"`
}

type apiServer struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Hostname string `json:"hostname"`
	Status   string `json:"status"`
}

type apiConnection struct {
	ID        string `json:"id"`
	ServerID  string `json:"server_id"`
	UserID    string `json:"user_id"`
	Protocol  string `json:"protocol"`
	Status    string `json:"status"`
	StartedAt string `json:"started_at"`
}

type apiPolicy struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	Action      string `json:"action"`
}

type apiAuditLog struct {
	ID        string `json:"id"`
	EventType string `json:"event_type"`
	Action    string `json:"action"`
	Severity  string `json:"severity"`
	CreatedAt string `json:"created_at"`
}

// ---------------------------------------------------------------------------
// PostgreSQL query functions
// ---------------------------------------------------------------------------

func queryUsers(db *sql.DB) ([]apiUser, error) {
	rows, err := db.Query(`SELECT id, email, COALESCE(display_name, ''), role,
		status, COALESCE(last_login_at, TIMESTAMPTZ 'epoch'), created_at
		FROM users WHERE deleted_at IS NULL ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []apiUser
	for rows.Next() {
		var u apiUser
		var lastLogin time.Time
		var createdAt time.Time
		if err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.Role,
			&u.Status, &lastLogin, &createdAt); err != nil {
			return nil, err
		}
		u.AutoGroups = []string{}
		u.IsBlocked = u.Status == "suspended"
		u.Issued = "api"
		if !lastLogin.IsZero() {
			u.LastLogin = lastLogin.Format(time.RFC3339)
		}
		_ = createdAt
		users = append(users, u)
	}
	if users == nil {
		users = []apiUser{}
	}
	return users, rows.Err()
}

func queryAccounts(db *sql.DB) ([]apiAccount, error) {
	rows, err := db.Query(`SELECT id, name, slug, tier, status, created_at
		FROM tenants WHERE deleted_at IS NULL ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []apiAccount
	for rows.Next() {
		var a apiAccount
		var createdAt time.Time
		if err := rows.Scan(&a.ID, &a.Name, &a.Slug, &a.Tier, &a.Status, &createdAt); err != nil {
			return nil, err
		}
		a.CreatedAt = createdAt.Format(time.RFC3339)
		accounts = append(accounts, a)
	}
	if accounts == nil {
		accounts = []apiAccount{}
	}
	return accounts, rows.Err()
}

func queryPeers(db *sql.DB) ([]apiPeer, error) {
	rows, err := db.Query(`SELECT id, name, COALESCE(hostname, ''),
		COALESCE(host(ip_address), ''), COALESCE(os, ''), status, COALESCE(agent_version, '')
		FROM servers WHERE deleted_at IS NULL ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var peers []apiPeer
	for rows.Next() {
		var p apiPeer
		if err := rows.Scan(&p.ID, &p.Name, &p.Hostname, &p.IP, &p.OS, &p.Status, &p.Version); err != nil {
			return nil, err
		}
		peers = append(peers, p)
	}
	if peers == nil {
		peers = []apiPeer{}
	}
	return peers, rows.Err()
}

func queryGroups(db *sql.DB) ([]apiGroup, error) {
	rows, err := db.Query(`SELECT id, name FROM user_groups ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []apiGroup
	for rows.Next() {
		var g apiGroup
		if err := rows.Scan(&g.ID, &g.Name); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	if groups == nil {
		groups = []apiGroup{}
	}
	return groups, rows.Err()
}

func querySetupKeys(db *sql.DB) ([]apiSetupKey, error) {
	rows, err := db.Query(`SELECT id, name, public_key, created_at
		FROM server_access_keys ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []apiSetupKey
	for rows.Next() {
		var k apiSetupKey
		var createdAt time.Time
		if err := rows.Scan(&k.ID, &k.Name, &k.Key, &createdAt); err != nil {
			return nil, err
		}
		k.CreatedAt = createdAt.Format(time.RFC3339)
		keys = append(keys, k)
	}
	if keys == nil {
		keys = []apiSetupKey{}
	}
	return keys, rows.Err()
}

func queryEvents(db *sql.DB) ([]apiEvent, error) {
	rows, err := db.Query(`SELECT id, event_type, event_action, created_at
		FROM audit_logs ORDER BY created_at DESC LIMIT 50`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []apiEvent
	for rows.Next() {
		var e apiEvent
		var createdAt time.Time
		if err := rows.Scan(&e.ID, &e.Type, &e.Action, &createdAt); err != nil {
			return nil, err
		}
		e.CreatedAt = createdAt.Format(time.RFC3339)
		events = append(events, e)
	}
	if events == nil {
		events = []apiEvent{}
	}
	return events, rows.Err()
}

func queryServers(db *sql.DB) ([]apiServer, error) {
	rows, err := db.Query(`SELECT id, name, COALESCE(hostname, ''), status
		FROM servers WHERE deleted_at IS NULL ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var servers []apiServer
	for rows.Next() {
		var s apiServer
		if err := rows.Scan(&s.ID, &s.Name, &s.Hostname, &s.Status); err != nil {
			return nil, err
		}
		servers = append(servers, s)
	}
	if servers == nil {
		servers = []apiServer{}
	}
	return servers, rows.Err()
}

func queryUsersV1(db *sql.DB) ([]apiUser, error) {
	return queryUsers(db)
}

func queryConnections(db *sql.DB) ([]apiConnection, error) {
	rows, err := db.Query(`SELECT id, server_id, user_id, protocol,
		status, started_at FROM connections ORDER BY started_at DESC LIMIT 50`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conns []apiConnection
	for rows.Next() {
		var c apiConnection
		var startedAt time.Time
		if err := rows.Scan(&c.ID, &c.ServerID, &c.UserID, &c.Protocol, &c.Status, &startedAt); err != nil {
			return nil, err
		}
		c.StartedAt = startedAt.Format(time.RFC3339)
		conns = append(conns, c)
	}
	if conns == nil {
		conns = []apiConnection{}
	}
	return conns, rows.Err()
}

func queryPolicies(db *sql.DB) ([]apiPolicy, error) {
	rows, err := db.Query(`SELECT id, name, COALESCE(description, ''),
		status, action FROM policies ORDER BY priority DESC, name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var policies []apiPolicy
	for rows.Next() {
		var p apiPolicy
		var status string
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &status, &p.Action); err != nil {
			return nil, err
		}
		p.Enabled = status == "enabled"
		policies = append(policies, p)
	}
	if policies == nil {
		policies = []apiPolicy{}
	}
	return policies, rows.Err()
}

func queryAuditLogs(db *sql.DB) ([]apiAuditLog, error) {
	rows, err := db.Query(`SELECT id, event_type, event_action, severity, created_at
		FROM audit_logs ORDER BY created_at DESC LIMIT 50`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []apiAuditLog
	for rows.Next() {
		var l apiAuditLog
		var createdAt time.Time
		if err := rows.Scan(&l.ID, &l.EventType, &l.Action, &l.Severity, &createdAt); err != nil {
			return nil, err
		}
		l.CreatedAt = createdAt.Format(time.RFC3339)
		logs = append(logs, l)
	}
	if logs == nil {
		logs = []apiAuditLog{}
	}
	return logs, rows.Err()
}

// ---------------------------------------------------------------------------
// Mock data (in-memory) for development without PostgreSQL
// ---------------------------------------------------------------------------

func mockUsers() []apiUser {
	return []apiUser{
		{
			ID:         "dev-user",
			Email:      "dev@kyle.local",
			Name:       "Dev User",
			Role:       "owner",
			AutoGroups: []string{},
			Status:     "active",
			IsBlocked:  false,
			Issued:     "api",
			LastLogin:  time.Now().UTC().Format(time.RFC3339),
		},
	}
}

func mockAccounts() []apiAccount {
	return []apiAccount{
		{
			ID:        "dev-account",
			Name:      "Kyle Dev",
			Slug:      "kyle-dev",
			Tier:      "free",
			Status:    "active",
			CreatedAt: time.Now().UTC().Format(time.RFC3339),
		},
	}
}

func mockPeers() []apiPeer {
	return []apiPeer{
		{
			ID:       "dev-peer-1",
			Name:     "dev-server",
			Hostname: "dev.local",
			IP:       "10.0.0.1",
			OS:       "linux",
			Status:   "offline",
			Version:  "0.1.0",
		},
	}
}

func mockGroups() []apiGroup {
	return []apiGroup{
		{ID: "dev-group-all", Name: "All"},
	}
}

func mockSetupKeys() []apiSetupKey {
	return []apiSetupKey{}
}

func mockNameservers() []apiNameserver {
	return []apiNameserver{}
}

func mockRoutes() []apiRoute {
	return []apiRoute{}
}

func mockDNSSettings() map[string]interface{} {
	return map[string]interface{}{
		"dns_domain":        "kyle.local",
		"custom_zones":      []interface{}{},
		"nameserver_groups": []interface{}{},
	}
}

func mockEvents() []apiEvent {
	return []apiEvent{}
}

func mockServers() []apiServer {
	return []apiServer{
		{
			ID:       "dev-server-1",
			Name:     "dev-server",
			Hostname: "dev.local",
			Status:   "offline",
		},
	}
}

func mockUsersV1() []apiUser {
	return mockUsers()
}

func mockConnections() []apiConnection {
	return []apiConnection{}
}

func mockPolicies() []apiPolicy {
	return []apiPolicy{}
}

func mockAuditLogs() []apiAuditLog {
	return []apiAuditLog{}
}
