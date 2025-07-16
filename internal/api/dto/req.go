package dto

import (
	"encoding/json"
	"time"
)

type CreateTenantRequest struct {
	Name string `json:"name" binding:"required"`
}

type CreateAuditLogRequest struct {
	TenantID     string          `json:"tenant_id" binding:"required"`
	UserID       string          `json:"user_id"`
	SessionID    string          `json:"session_id"`
	IPAddress    string          `json:"ip_address"`
	UserAgent    string          `json:"user_agent"`
	Action       string          `json:"action" binding:"required"`
	ResourceType string          `json:"resource_type" binding:"required"`
	ResourceID   string          `json:"resource_id" binding:"required"`
	Severity     string          `json:"severity" binding:"required"`
	Message      string          `json:"message" binding:"required"`
	BeforeState  json.RawMessage `json:"before_state"`
	AfterState   json.RawMessage `json:"after_state"`
	Metadata     json.RawMessage `json:"metadata"`
	Timestamp    time.Time       `json:"timestamp" binding:"required"`
}
