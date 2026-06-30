package store

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/netbirdio/netbird/management/server/types"
	"github.com/netbirdio/netbird/shared/management/status"
)

func (s *SqlStore) GetCosmosResources(ctx context.Context, accountID string) ([]*types.CosmosResource, error) {
	var resources []*types.CosmosResource
	result := s.db.WithContext(ctx).
		Where(accountIDCondition, accountID).
		Order("name asc").
		Find(&resources)
	if result.Error != nil {
		return nil, status.Errorf(status.Internal, "failed to get resources")
	}
	return resources, nil
}

func (s *SqlStore) GetCosmosResource(ctx context.Context, accountID, resourceID string) (*types.CosmosResource, error) {
	var resource types.CosmosResource
	result := s.db.WithContext(ctx).
		Where(accountAndIDQueryCondition, accountID, resourceID).
		First(&resource)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, status.Errorf(status.NotFound, "resource not found")
	}
	if result.Error != nil {
		return nil, status.Errorf(status.Internal, "failed to get resource")
	}
	return &resource, nil
}

func (s *SqlStore) CreateCosmosResource(ctx context.Context, resource *types.CosmosResource) error {
	if result := s.db.WithContext(ctx).Create(resource); result.Error != nil {
		return status.Errorf(status.Internal, "failed to create resource")
	}
	return nil
}

func (s *SqlStore) SaveCosmosResource(ctx context.Context, resource *types.CosmosResource) error {
	result := s.db.WithContext(ctx).
		Model(&types.CosmosResource{}).
		Where(accountAndIDQueryCondition, resource.AccountID, resource.ID).
		Select("*").
		Updates(resource)
	if result.Error != nil {
		return status.Errorf(status.Internal, "failed to update resource")
	}
	if result.RowsAffected == 0 {
		return status.Errorf(status.NotFound, "resource not found")
	}
	return nil
}

func (s *SqlStore) DeleteCosmosResource(ctx context.Context, accountID, resourceID string) error {
	result := s.db.WithContext(ctx).
		Where(accountAndIDQueryCondition, accountID, resourceID).
		Delete(&types.CosmosResource{})
	if result.Error != nil {
		return status.Errorf(status.Internal, "failed to delete resource")
	}
	if result.RowsAffected == 0 {
		return status.Errorf(status.NotFound, "resource not found")
	}
	return nil
}

func (s *SqlStore) GetCosmosSessions(ctx context.Context, accountID string) ([]*types.CosmosSession, error) {
	var sessions []*types.CosmosSession
	result := s.db.WithContext(ctx).
		Where(accountIDCondition, accountID).
		Order("started_at desc").
		Find(&sessions)
	if result.Error != nil {
		return nil, status.Errorf(status.Internal, "failed to get sessions")
	}
	return sessions, nil
}

func (s *SqlStore) GetCosmosSession(ctx context.Context, accountID, sessionID string) (*types.CosmosSession, error) {
	var session types.CosmosSession
	result := s.db.WithContext(ctx).
		Where(accountAndIDQueryCondition, accountID, sessionID).
		First(&session)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, status.Errorf(status.NotFound, "session not found")
	}
	if result.Error != nil {
		return nil, status.Errorf(status.Internal, "failed to get session")
	}
	return &session, nil
}

func (s *SqlStore) CreateCosmosSession(ctx context.Context, session *types.CosmosSession) error {
	if result := s.db.WithContext(ctx).Create(session); result.Error != nil {
		return status.Errorf(status.Internal, "failed to create session")
	}
	return nil
}

func (s *SqlStore) SaveCosmosSession(ctx context.Context, session *types.CosmosSession) error {
	result := s.db.WithContext(ctx).
		Model(&types.CosmosSession{}).
		Where(accountAndIDQueryCondition, session.AccountID, session.ID).
		Select("*").
		Updates(session)
	if result.Error != nil {
		return status.Errorf(status.Internal, "failed to update session")
	}
	if result.RowsAffected == 0 {
		return status.Errorf(status.NotFound, "session not found")
	}
	return nil
}

func (s *SqlStore) CreateCosmosAuditEvent(ctx context.Context, event *types.CosmosAuditEvent) error {
	if result := s.db.WithContext(ctx).Create(event); result.Error != nil {
		return status.Errorf(status.Internal, "failed to create audit event")
	}
	return nil
}

func (s *SqlStore) GetCosmosAuditEvents(ctx context.Context, accountID string) ([]*types.CosmosAuditEvent, error) {
	var events []*types.CosmosAuditEvent
	result := s.db.WithContext(ctx).
		Where(accountIDCondition, accountID).
		Order("timestamp desc").
		Limit(500).
		Find(&events)
	if result.Error != nil {
		return nil, status.Errorf(status.Internal, "failed to get audit events")
	}
	return events, nil
}
