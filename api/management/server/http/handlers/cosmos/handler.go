package cosmos

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/xid"

	"github.com/netbirdio/netbird/management/server/account"
	nbcontext "github.com/netbirdio/netbird/management/server/context"
	"github.com/netbirdio/netbird/management/server/types"
	"github.com/netbirdio/netbird/shared/management/http/util"
)

type store interface {
	GetCosmosResources(ctx context.Context, accountID string) ([]*types.CosmosResource, error)
	GetCosmosResource(ctx context.Context, accountID, resourceID string) (*types.CosmosResource, error)
	CreateCosmosResource(ctx context.Context, resource *types.CosmosResource) error
	SaveCosmosResource(ctx context.Context, resource *types.CosmosResource) error
	DeleteCosmosResource(ctx context.Context, accountID, resourceID string) error
	GetCosmosSessions(ctx context.Context, accountID string) ([]*types.CosmosSession, error)
	GetCosmosSession(ctx context.Context, accountID, sessionID string) (*types.CosmosSession, error)
	CreateCosmosSession(ctx context.Context, session *types.CosmosSession) error
	SaveCosmosSession(ctx context.Context, session *types.CosmosSession) error
	CreateCosmosAuditEvent(ctx context.Context, event *types.CosmosAuditEvent) error
	GetCosmosAuditEvents(ctx context.Context, accountID string) ([]*types.CosmosAuditEvent, error)
}

type handler struct {
	store store
}

type resourceRequest struct {
	Name             string                       `json:"name"`
	Description      string                       `json:"description"`
	Protocol         types.CosmosResourceProtocol `json:"protocol"`
	Host             string                       `json:"host"`
	Port             int                          `json:"port"`
	GroupIDs         []string                     `json:"group_ids"`
	Enabled          bool                         `json:"enabled"`
	RecordingEnabled bool                         `json:"recording_enabled"`
}

type resourceGroupsRequest struct {
	GroupIDs []string `json:"group_ids"`
}

type sessionRequest struct {
	ResourceID string `json:"resource_id"`
	ClientIP   string `json:"client_ip"`
}

type auditMeta map[string]any

func AddEndpoints(accountManager account.Manager, router *mux.Router) {
	cosmosStore, ok := accountManager.GetStore().(store)
	h := &handler{store: cosmosStore}

	if !ok {
		router.HandleFunc("/cosmos/resources", unsupported).Methods("GET", "POST", "OPTIONS")
		router.HandleFunc("/cosmos/resources/{resourceId}", unsupported).Methods("GET", "PUT", "DELETE", "OPTIONS")
		router.HandleFunc("/cosmos/sessions", unsupported).Methods("GET", "POST", "OPTIONS")
		router.HandleFunc("/cosmos/sessions/{sessionId}/close", unsupported).Methods("POST", "OPTIONS")
		router.HandleFunc("/cosmos/audit/events", unsupported).Methods("GET", "OPTIONS")
		return
	}

	router.HandleFunc("/cosmos/resources", h.getResources).Methods("GET", "OPTIONS")
	router.HandleFunc("/cosmos/resources", h.createResource).Methods("POST", "OPTIONS")
	router.HandleFunc("/cosmos/resources/{resourceId}", h.getResource).Methods("GET", "OPTIONS")
	router.HandleFunc("/cosmos/resources/{resourceId}", h.updateResource).Methods("PUT", "OPTIONS")
	router.HandleFunc("/cosmos/resources/{resourceId}", h.deleteResource).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/cosmos/resources/{resourceId}/groups", h.getResourceGroups).Methods("GET", "OPTIONS")
	router.HandleFunc("/cosmos/resources/{resourceId}/groups", h.updateResourceGroups).Methods("PUT", "OPTIONS")
	router.HandleFunc("/cosmos/sessions", h.getSessions).Methods("GET", "OPTIONS")
	router.HandleFunc("/cosmos/sessions", h.createSession).Methods("POST", "OPTIONS")
	router.HandleFunc("/cosmos/sessions/{sessionId}/close", h.closeSession).Methods("POST", "OPTIONS")
	router.HandleFunc("/cosmos/audit/events", h.getAuditEvents).Methods("GET", "OPTIONS")
}

func unsupported(w http.ResponseWriter, _ *http.Request) {
	util.WriteErrorResponse("cosmos store is not available", http.StatusNotImplemented, w)
}

func (h *handler) getResources(w http.ResponseWriter, r *http.Request) {
	userAuth, err := nbcontext.GetUserAuthFromContext(r.Context())
	if err != nil {
		util.WriteError(r.Context(), err, w)
		return
	}
	resources, err := h.store.GetCosmosResources(r.Context(), userAuth.AccountId)
	if err != nil {
		util.WriteError(r.Context(), err, w)
		return
	}
	util.WriteJSONObject(r.Context(), w, resources)
}

func (h *handler) getResource(w http.ResponseWriter, r *http.Request) {
	userAuth, err := nbcontext.GetUserAuthFromContext(r.Context())
	if err != nil {
		util.WriteError(r.Context(), err, w)
		return
	}
	resource, err := h.store.GetCosmosResource(r.Context(), userAuth.AccountId, mux.Vars(r)["resourceId"])
	if err != nil {
		util.WriteError(r.Context(), err, w)
		return
	}
	util.WriteJSONObject(r.Context(), w, resource)
}

func (h *handler) createResource(w http.ResponseWriter, r *http.Request) {
	userAuth, err := nbcontext.GetUserAuthFromContext(r.Context())
	if err != nil {
		util.WriteError(r.Context(), err, w)
		return
	}
	req, ok := decodeResourceRequest(w, r)
	if !ok {
		return
	}
	now := time.Now().UTC()
	resource := &types.CosmosResource{
		ID:               xid.New().String(),
		AccountID:        userAuth.AccountId,
		Name:             strings.TrimSpace(req.Name),
		Description:      strings.TrimSpace(req.Description),
		Protocol:         req.Protocol,
		Host:             strings.TrimSpace(req.Host),
		Port:             defaultPort(req.Protocol, req.Port),
		GroupIDs:         strings.Join(normalizeLabels(req.GroupIDs), ","),
		Enabled:          req.Enabled,
		RecordingEnabled: req.RecordingEnabled,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	if err := validateResource(resource); err != "" {
		util.WriteErrorResponse(err, http.StatusBadRequest, w)
		return
	}
	if err := h.store.CreateCosmosResource(r.Context(), resource); err != nil {
		util.WriteError(r.Context(), err, w)
		return
	}
	h.audit(r, userAuth.AccountId, userAuth.UserId, userAuth.Name, userAuth.Email, "resource.created", "resource", resource.ID, resource.Name, auditMeta{
		"protocol": resource.Protocol,
		"host":     resource.Host,
		"port":     strconv.Itoa(resource.Port),
	})
	util.WriteJSONObject(r.Context(), w, resource)
}

func (h *handler) updateResource(w http.ResponseWriter, r *http.Request) {
	userAuth, err := nbcontext.GetUserAuthFromContext(r.Context())
	if err != nil {
		util.WriteError(r.Context(), err, w)
		return
	}
	resource, err := h.store.GetCosmosResource(r.Context(), userAuth.AccountId, mux.Vars(r)["resourceId"])
	if err != nil {
		util.WriteError(r.Context(), err, w)
		return
	}
	req, ok := decodeResourceRequest(w, r)
	if !ok {
		return
	}
	resource.Name = strings.TrimSpace(req.Name)
	resource.Description = strings.TrimSpace(req.Description)
	resource.Protocol = req.Protocol
	resource.Host = strings.TrimSpace(req.Host)
	resource.Port = defaultPort(req.Protocol, req.Port)
	resource.GroupIDs = strings.Join(normalizeLabels(req.GroupIDs), ",")
	resource.Enabled = req.Enabled
	resource.RecordingEnabled = req.RecordingEnabled
	resource.UpdatedAt = time.Now().UTC()
	if err := validateResource(resource); err != "" {
		util.WriteErrorResponse(err, http.StatusBadRequest, w)
		return
	}
	if err := h.store.SaveCosmosResource(r.Context(), resource); err != nil {
		util.WriteError(r.Context(), err, w)
		return
	}
	h.audit(r, userAuth.AccountId, userAuth.UserId, userAuth.Name, userAuth.Email, "resource.updated", "resource", resource.ID, resource.Name, nil)
	util.WriteJSONObject(r.Context(), w, resource)
}

func (h *handler) deleteResource(w http.ResponseWriter, r *http.Request) {
	userAuth, err := nbcontext.GetUserAuthFromContext(r.Context())
	if err != nil {
		util.WriteError(r.Context(), err, w)
		return
	}
	resourceID := mux.Vars(r)["resourceId"]
	resource, _ := h.store.GetCosmosResource(r.Context(), userAuth.AccountId, resourceID)
	if err := h.store.DeleteCosmosResource(r.Context(), userAuth.AccountId, resourceID); err != nil {
		util.WriteError(r.Context(), err, w)
		return
	}
	targetName := resourceID
	if resource != nil {
		targetName = resource.Name
	}
	h.audit(r, userAuth.AccountId, userAuth.UserId, userAuth.Name, userAuth.Email, "resource.deleted", "resource", resourceID, targetName, nil)
	util.WriteJSONObject(r.Context(), w, util.EmptyObject{})
}

func (h *handler) getResourceGroups(w http.ResponseWriter, r *http.Request) {
	userAuth, err := nbcontext.GetUserAuthFromContext(r.Context())
	if err != nil {
		util.WriteError(r.Context(), err, w)
		return
	}
	resource, err := h.store.GetCosmosResource(r.Context(), userAuth.AccountId, mux.Vars(r)["resourceId"])
	if err != nil {
		util.WriteError(r.Context(), err, w)
		return
	}
	groupIDs := strings.Split(resource.GroupIDs, ",")
	if len(groupIDs) == 1 && groupIDs[0] == "" {
		groupIDs = nil
	}
	util.WriteJSONObject(r.Context(), w, map[string]any{"group_ids": groupIDs})
}

func (h *handler) updateResourceGroups(w http.ResponseWriter, r *http.Request) {
	userAuth, err := nbcontext.GetUserAuthFromContext(r.Context())
	if err != nil {
		util.WriteError(r.Context(), err, w)
		return
	}
	resource, err := h.store.GetCosmosResource(r.Context(), userAuth.AccountId, mux.Vars(r)["resourceId"])
	if err != nil {
		util.WriteError(r.Context(), err, w)
		return
	}
	var req resourceGroupsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteErrorResponse("couldn't parse JSON request", http.StatusBadRequest, w)
		return
	}
	resource.GroupIDs = strings.Join(normalizeLabels(req.GroupIDs), ",")
	resource.UpdatedAt = time.Now().UTC()
	if err := h.store.SaveCosmosResource(r.Context(), resource); err != nil {
		util.WriteError(r.Context(), err, w)
		return
	}
	h.audit(r, userAuth.AccountId, userAuth.UserId, userAuth.Name, userAuth.Email, "resource.groups_updated", "resource", resource.ID, resource.Name, auditMeta{"group_ids": resource.GroupIDs})
	util.WriteJSONObject(r.Context(), w, resource)
}

func (h *handler) getSessions(w http.ResponseWriter, r *http.Request) {
	userAuth, err := nbcontext.GetUserAuthFromContext(r.Context())
	if err != nil {
		util.WriteError(r.Context(), err, w)
		return
	}
	sessions, err := h.store.GetCosmosSessions(r.Context(), userAuth.AccountId)
	if err != nil {
		util.WriteError(r.Context(), err, w)
		return
	}
	util.WriteJSONObject(r.Context(), w, sessions)
}

func (h *handler) createSession(w http.ResponseWriter, r *http.Request) {
	userAuth, err := nbcontext.GetUserAuthFromContext(r.Context())
	if err != nil {
		util.WriteError(r.Context(), err, w)
		return
	}
	var req sessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteErrorResponse("couldn't parse JSON request", http.StatusBadRequest, w)
		return
	}
	resource, err := h.store.GetCosmosResource(r.Context(), userAuth.AccountId, req.ResourceID)
	if err != nil {
		util.WriteError(r.Context(), err, w)
		return
	}
	if !resource.Enabled {
		util.WriteErrorResponse("resource is disabled", http.StatusPreconditionFailed, w)
		return
	}
	now := time.Now().UTC()
	clientIP := strings.TrimSpace(req.ClientIP)
	if clientIP == "" {
		clientIP, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	session := &types.CosmosSession{
		ID:           xid.New().String(),
		AccountID:    userAuth.AccountId,
		ResourceID:   resource.ID,
		ResourceName: resource.Name,
		UserID:       userAuth.UserId,
		UserName:     userAuth.Name,
		UserEmail:    userAuth.Email,
		Protocol:     resource.Protocol,
		Status:       types.CosmosSessionActive,
		ClientIP:     clientIP,
		StartedAt:    now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := h.store.CreateCosmosSession(r.Context(), session); err != nil {
		util.WriteError(r.Context(), err, w)
		return
	}
	h.audit(r, userAuth.AccountId, userAuth.UserId, userAuth.Name, userAuth.Email, "session.started", "session", session.ID, resource.Name, auditMeta{
		"resource_id": resource.ID,
		"protocol":    resource.Protocol,
	})
	util.WriteJSONObject(r.Context(), w, session)
}

func (h *handler) closeSession(w http.ResponseWriter, r *http.Request) {
	userAuth, err := nbcontext.GetUserAuthFromContext(r.Context())
	if err != nil {
		util.WriteError(r.Context(), err, w)
		return
	}
	session, err := h.store.GetCosmosSession(r.Context(), userAuth.AccountId, mux.Vars(r)["sessionId"])
	if err != nil {
		util.WriteError(r.Context(), err, w)
		return
	}
	now := time.Now().UTC()
	session.Status = types.CosmosSessionClosed
	session.EndedAt = &now
	session.UpdatedAt = now
	if err := h.store.SaveCosmosSession(r.Context(), session); err != nil {
		util.WriteError(r.Context(), err, w)
		return
	}
	h.audit(r, userAuth.AccountId, userAuth.UserId, userAuth.Name, userAuth.Email, "session.closed", "session", session.ID, session.ResourceName, auditMeta{
		"resource_id": session.ResourceID,
	})
	util.WriteJSONObject(r.Context(), w, session)
}

func (h *handler) getAuditEvents(w http.ResponseWriter, r *http.Request) {
	userAuth, err := nbcontext.GetUserAuthFromContext(r.Context())
	if err != nil {
		util.WriteError(r.Context(), err, w)
		return
	}
	events, err := h.store.GetCosmosAuditEvents(r.Context(), userAuth.AccountId)
	if err != nil {
		util.WriteError(r.Context(), err, w)
		return
	}
	util.WriteJSONObject(r.Context(), w, events)
}

func (h *handler) audit(r *http.Request, accountID, userID, userName, userEmail, action, targetType, targetID, targetName string, meta auditMeta) {
	payload := "{}"
	if meta != nil {
		if b, err := json.Marshal(meta); err == nil {
			payload = string(b)
		}
	}
	now := time.Now().UTC()
	_ = h.store.CreateCosmosAuditEvent(r.Context(), &types.CosmosAuditEvent{
		ID:         xid.New().String(),
		AccountID:  accountID,
		UserID:     userID,
		UserName:   userName,
		UserEmail:  userEmail,
		Action:     action,
		TargetType: targetType,
		TargetID:   targetID,
		TargetName: targetName,
		Timestamp:  now,
		Meta:       payload,
		CreatedAt:  now,
	})
}

func decodeResourceRequest(w http.ResponseWriter, r *http.Request) (resourceRequest, bool) {
	var req resourceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteErrorResponse("couldn't parse JSON request", http.StatusBadRequest, w)
		return req, false
	}
	return req, true
}

func validateResource(resource *types.CosmosResource) string {
	if resource.Name == "" {
		return "name is required"
	}
	if resource.Host == "" {
		return "host is required"
	}
	if resource.Protocol != types.CosmosProtocolSSH && resource.Protocol != types.CosmosProtocolRDP && resource.Protocol != types.CosmosProtocolVNC {
		return "protocol must be ssh, rdp, or vnc"
	}
	if resource.Port <= 0 || resource.Port > 65535 {
		return "port must be between 1 and 65535"
	}
	return ""
}

func defaultPort(protocol types.CosmosResourceProtocol, port int) int {
	if port > 0 {
		return port
	}
	switch protocol {
	case types.CosmosProtocolRDP:
		return 3389
	case types.CosmosProtocolVNC:
		return 5900
	default:
		return 22
	}
}

func normalizeLabels(labels []string) []string {
	normalized := make([]string, 0, len(labels))
	seen := make(map[string]struct{}, len(labels))
	for _, label := range labels {
		label = strings.TrimSpace(label)
		if label == "" {
			continue
		}
		if _, ok := seen[label]; ok {
			continue
		}
		seen[label] = struct{}{}
		normalized = append(normalized, label)
	}
	return normalized
}
