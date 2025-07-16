package dto

import (
	"encoding/json"
	"time"
)

type CreateTenantResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AuditLogResponse struct {
	ID           string          `json:"id"`
	TenantID     string          `json:"tenant_id"`
	UserID       string          `json:"user_id"`
	SessionID    string          `json:"session_id"`
	IPAddress    string          `json:"ip_address"`
	UserAgent    string          `json:"user_agent"`
	Action       string          `json:"action"`
	ResourceType string          `json:"resource_type"`
	ResourceID   string          `json:"resource_id"`
	Severity     string          `json:"severity"`
	Message      string          `json:"message"`
	BeforeState  json.RawMessage `json:"before_state,omitempty"`
	AfterState   json.RawMessage `json:"after_state,omitempty"`
	Metadata     json.RawMessage `json:"metadata,omitempty"`
	Timestamp    time.Time       `json:"timestamp"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}
