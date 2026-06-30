package types

import "time"

type CosmosResourceProtocol string

const (
	CosmosProtocolSSH CosmosResourceProtocol = "ssh"
	CosmosProtocolRDP CosmosResourceProtocol = "rdp"
	CosmosProtocolVNC CosmosResourceProtocol = "vnc"
)

type CosmosResource struct {
	ID               string                 `gorm:"primaryKey" json:"id"`
	AccountID        string                 `gorm:"index;not null" json:"account_id"`
	Name             string                 `gorm:"not null" json:"name"`
	Description      string                 `json:"description"`
	Protocol         CosmosResourceProtocol `gorm:"not null" json:"protocol"`
	Host             string                 `gorm:"not null" json:"host"`
	Port             int                    `gorm:"not null" json:"port"`
	GroupIDs         string                 `json:"group_ids"`
	Enabled          bool                   `gorm:"not null;default:true" json:"enabled"`
	RecordingEnabled bool                   `gorm:"not null;default:true" json:"recording_enabled"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
}

type CosmosSessionStatus string

const (
	CosmosSessionPending CosmosSessionStatus = "pending"
	CosmosSessionActive  CosmosSessionStatus = "active"
	CosmosSessionClosed  CosmosSessionStatus = "closed"
	CosmosSessionDenied  CosmosSessionStatus = "denied"
)

type CosmosSession struct {
	ID                string                 `gorm:"primaryKey" json:"id"`
	AccountID         string                 `gorm:"index;not null" json:"account_id"`
	ResourceID        string                 `gorm:"index;not null" json:"resource_id"`
	ResourceName      string                 `json:"resource_name"`
	UserID            string                 `gorm:"index;not null" json:"user_id"`
	UserName          string                 `json:"user_name"`
	UserEmail         string                 `json:"user_email"`
	Protocol          CosmosResourceProtocol `gorm:"not null" json:"protocol"`
	Status            CosmosSessionStatus    `gorm:"not null" json:"status"`
	ClientIP          string                 `json:"client_ip"`
	GuacdConnectionID string                 `json:"guacd_connection_id"`
	RecordingEnabled  bool                   `json:"recording_enabled"`
	RecordingPath     string                 `json:"recording_path"`
	StartedAt         time.Time              `json:"started_at"`
	EndedAt           *time.Time             `json:"ended_at"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
}

type CosmosAuditEvent struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	AccountID   string    `gorm:"index;not null" json:"account_id"`
	UserID      string    `gorm:"index" json:"user_id"`
	UserName    string    `json:"user_name"`
	UserEmail   string    `json:"user_email"`
	Action      string    `gorm:"index;not null" json:"action"`
	TargetType  string    `gorm:"index;not null" json:"target_type"`
	TargetID    string    `gorm:"index" json:"target_id"`
	TargetName  string    `json:"target_name"`
	Timestamp   time.Time `gorm:"index;not null" json:"timestamp"`
	Description string    `json:"description"`
	Meta        string    `json:"meta"`
	CreatedAt   time.Time `json:"created_at"`
}
