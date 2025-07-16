package repository

import (
	"context"
	"time"

	"github.com/buiminhduc234/audit-log-api/internal/domain"
)

type AuditLogRepository interface {
	Create(ctx context.Context, log *domain.AuditLog) error
	GetByID(ctx context.Context, id string) (*domain.AuditLog, error)
	List(ctx context.Context, filter domain.AuditLogFilter) ([]domain.AuditLog, error)
	Delete(ctx context.Context, id string) error
	DeleteOlderThan(ctx context.Context, days int) error
	BulkCreate(ctx context.Context, logs []domain.AuditLog) error
}

type OpenSearchRepository interface {
	Index(ctx context.Context, log *domain.AuditLog) error
	BulkIndex(ctx context.Context, logs []domain.AuditLog) error
	Search(ctx context.Context, filter *domain.AuditLogFilter) ([]domain.AuditLog, error)
	CreateIndex(ctx context.Context, tenantID string, t time.Time) error
	DeleteIndex(ctx context.Context, tenantID string) error
}

type TenantRepository interface {
	Create(ctx context.Context, tenant *domain.Tenant) (*domain.Tenant, error)
	GetByID(ctx context.Context, id string) (*domain.Tenant, error)
	GetByAPIKey(ctx context.Context, apiKey string) (*domain.Tenant, error)
	Update(ctx context.Context, tenant *domain.Tenant) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]domain.Tenant, error)
}

// PostgresRepository represents the PostgreSQL-specific repository interface
type PostgresRepository interface {
	AuditLog() AuditLogRepository
	Tenant() TenantRepository
}

// Repository represents the main repository interface that combines both PostgreSQL and OpenSearch
type Repository interface {
	PostgresRepository
	OpenSearch() OpenSearchRepository
}
