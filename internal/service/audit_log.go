package service

import (
	"context"
	"fmt"

	"github.com/buiminhduc234/audit-log-api/internal/api/dto"
	"github.com/buiminhduc234/audit-log-api/internal/domain"
	"github.com/buiminhduc234/audit-log-api/internal/repository"
	"github.com/buiminhduc234/audit-log-api/internal/service/queue"
)

type AuditLogService struct {
	repo   repository.Repository
	sqsSvc *queue.SQSService
}

func NewAuditLogService(repo repository.Repository, sqsSvc *queue.SQSService) *AuditLogService {
	return &AuditLogService{
		repo:   repo,
		sqsSvc: sqsSvc,
	}
}

func (s *AuditLogService) Create(ctx context.Context, log dto.CreateAuditLogRequest) error {
	auditLog := &domain.AuditLog{
		TenantID:     log.TenantID,
		UserID:       log.UserID,
		SessionID:    log.SessionID,
		IPAddress:    log.IPAddress,
		UserAgent:    log.UserAgent,
		Action:       log.Action,
		ResourceType: log.ResourceType,
		ResourceID:   log.ResourceID,
		Severity:     log.Severity,
		Message:      log.Message,
		BeforeState:  log.BeforeState,
		AfterState:   log.AfterState,
		Metadata:     log.Metadata,
		Timestamp:    log.Timestamp,
	}
	// Store in PostgreSQL
	if err := s.repo.AuditLog().Create(ctx, auditLog); err != nil {
		return fmt.Errorf("failed to store log in PostgreSQL: %w", err)
	}

	// Send message to SQS for asynchronous indexing
	if err := s.sqsSvc.SendIndexMessage(ctx, auditLog); err != nil {
		// Log the error but don't fail the request
		fmt.Printf("failed to send index message to SQS: %v\n", err)
	}

	return nil
}

func (s *AuditLogService) BulkCreate(ctx context.Context, logs []dto.CreateAuditLogRequest) error {
	auditLogs := make([]domain.AuditLog, len(logs))
	for i := range logs {
		auditLog := domain.AuditLog{
			TenantID:     logs[i].TenantID,
			UserID:       logs[i].UserID,
			SessionID:    logs[i].SessionID,
			IPAddress:    logs[i].IPAddress,
			UserAgent:    logs[i].UserAgent,
			Action:       logs[i].Action,
			ResourceType: logs[i].ResourceType,
			ResourceID:   logs[i].ResourceID,
			Severity:     logs[i].Severity,
			Message:      logs[i].Message,
			BeforeState:  logs[i].BeforeState,
			AfterState:   logs[i].AfterState,
			Metadata:     logs[i].Metadata,
			Timestamp:    logs[i].Timestamp,
		}
		auditLogs[i] = auditLog
	}

	// Store in PostgreSQL
	if err := s.repo.AuditLog().BulkCreate(ctx, auditLogs); err != nil {
		return fmt.Errorf("failed to bulk store logs in PostgreSQL: %w", err)
	}

	// Send message to SQS for asynchronous bulk indexing
	if err := s.sqsSvc.SendBulkIndexMessage(ctx, auditLogs); err != nil {
		// Log the error but don't fail the request
		fmt.Printf("failed to send bulk index message to SQS: %v\n", err)
	}

	return nil
}

func (s *AuditLogService) GetByID(ctx context.Context, id string) (*dto.AuditLogResponse, error) {
	log, err := s.repo.AuditLog().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &dto.AuditLogResponse{
		ID:           log.ID,
		TenantID:     log.TenantID,
		UserID:       log.UserID,
		SessionID:    log.SessionID,
		IPAddress:    log.IPAddress,
		UserAgent:    log.UserAgent,
		Action:       log.Action,
		ResourceType: log.ResourceType,
		ResourceID:   log.ResourceID,
		Severity:     log.Severity,
		Message:      log.Message,
		BeforeState:  log.BeforeState,
		AfterState:   log.AfterState,
		Metadata:     log.Metadata,
		Timestamp:    log.Timestamp,
		CreatedAt:    log.CreatedAt,
		UpdatedAt:    log.UpdatedAt,
	}, nil
}

func (s *AuditLogService) List(ctx context.Context, filter *domain.AuditLogFilter) ([]domain.AuditLog, error) {
	// Set default values for pagination
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 {
		filter.PageSize = 10
	}

	// Convert page and page size to limit and offset
	filter.Limit = filter.PageSize
	filter.Offset = (filter.Page - 1) * filter.PageSize

	// Use OpenSearch for searching if there are search criteria
	if s.hasSearchCriteria(filter) {
		return s.repo.OpenSearch().Search(ctx, filter)
	}

	// Otherwise, use PostgreSQL for simple listing
	return s.repo.AuditLog().List(ctx, *filter)
}

func (s *AuditLogService) Delete(ctx context.Context, id string) error {
	// First, get the log to get its tenant ID
	log, err := s.repo.AuditLog().GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get log for deletion: %w", err)
	}
	if log == nil {
		return fmt.Errorf("log not found")
	}

	// Delete from PostgreSQL
	if err := s.repo.AuditLog().Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete log from PostgreSQL: %w", err)
	}

	// Send message to SQS for asynchronous deletion from OpenSearch
	if err := s.sqsSvc.SendDeleteMessage(ctx, log.TenantID, id); err != nil {
		// Log the error but don't fail the request
		fmt.Printf("failed to send delete message to SQS: %v\n", err)
	}

	return nil
}

func (s *AuditLogService) DeleteOlderThan(ctx context.Context, days int) error {
	// Delete from PostgreSQL
	if err := s.repo.AuditLog().DeleteOlderThan(ctx, days); err != nil {
		return fmt.Errorf("failed to delete old logs from PostgreSQL: %w", err)
	}

	// Note: For OpenSearch, we use Index Lifecycle Management (ILM)
	// to automatically manage old indices. This is configured at
	// the OpenSearch cluster level.

	return nil
}

func (s *AuditLogService) GetStats(ctx context.Context, filter *domain.AuditLogFilter) (*domain.AuditLogStats, error) {
	// Use OpenSearch for aggregations if available, otherwise fall back to PostgreSQL
	logs, err := s.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	stats := &domain.AuditLogStats{
		TotalLogs:      int64(len(logs)),
		ActionCounts:   make(map[domain.ActionType]int64),
		SeverityCounts: make(map[domain.SeverityLevel]int64),
		ResourceCounts: make(map[string]int64),
	}

	for _, log := range logs {
		// Count by action
		stats.ActionCounts[domain.ActionType(log.Action)]++

		// Count by severity
		stats.SeverityCounts[domain.SeverityLevel(log.Severity)]++

		// Count by resource
		if log.ResourceType != "" {
			stats.ResourceCounts[log.ResourceType]++
		}
	}

	return stats, nil
}

func (s *AuditLogService) Search(ctx context.Context, filter *domain.AuditLogFilter) ([]domain.AuditLog, error) {
	// Always use OpenSearch for search operations
	return s.repo.OpenSearch().Search(ctx, filter)
}

// hasSearchCriteria checks if the filter contains search criteria that would benefit from OpenSearch
func (s *AuditLogService) hasSearchCriteria(filter *domain.AuditLogFilter) bool {
	return filter.UserID != "" ||
		filter.Action != "" ||
		filter.ResourceType != "" ||
		filter.Severity != "" ||
		!filter.StartTime.IsZero() ||
		!filter.EndTime.IsZero()
}
