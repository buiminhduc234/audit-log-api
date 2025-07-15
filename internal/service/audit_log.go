package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/buiminhduc234/audit-log-api/internal/domain"
	"github.com/buiminhduc234/audit-log-api/internal/repository"
)

type AuditLogService struct {
	repo repository.Repository
}

func NewAuditLogService(repo repository.Repository) *AuditLogService {
	return &AuditLogService{repo: repo}
}

func (s *AuditLogService) Create(ctx context.Context, log *domain.AuditLog) error {
	if log.ID == "" {
		log.ID = uuid.New().String()
	}

	now := time.Now()
	log.CreatedAt = now
	log.UpdatedAt = now

	return s.repo.AuditLog().Create(ctx, log)
}

func (s *AuditLogService) BulkCreate(ctx context.Context, logs []domain.AuditLog) error {
	now := time.Now()
	for i := range logs {
		if logs[i].ID == "" {
			logs[i].ID = uuid.New().String()
		}
		logs[i].CreatedAt = now
		logs[i].UpdatedAt = now
	}

	return s.repo.AuditLog().BulkCreate(ctx, logs)
}

func (s *AuditLogService) GetByID(ctx context.Context, id string) (*domain.AuditLog, error) {
	return s.repo.AuditLog().GetByID(ctx, id)
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

	return s.repo.AuditLog().List(ctx, *filter)
}

func (s *AuditLogService) Delete(ctx context.Context, id string) error {
	return s.repo.AuditLog().Delete(ctx, id)
}

func (s *AuditLogService) DeleteOlderThan(ctx context.Context, days int) error {
	return s.repo.AuditLog().DeleteOlderThan(ctx, days)
}

func (s *AuditLogService) GetStats(ctx context.Context, filter *domain.AuditLogFilter) (*domain.AuditLogStats, error) {
	logs, err := s.repo.AuditLog().List(ctx, *filter)
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
	// For now, we'll use the List method with the search query
	// In a real implementation, this would use full-text search capabilities
	// like Elasticsearch or PostgreSQL's full-text search
	return s.repo.AuditLog().List(ctx, *filter)
}
