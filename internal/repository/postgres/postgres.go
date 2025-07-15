package postgres

import (
	"gorm.io/gorm"

	"github.com/buiminhduc234/audit-log-api/internal/repository"
)

type postgresRepository struct {
	db           *gorm.DB
	auditLogRepo repository.AuditLogRepository
	tenantRepo   repository.TenantRepository
}

func NewPostgresRepository(db *gorm.DB) repository.Repository {
	return &postgresRepository{
		db:           db,
		auditLogRepo: NewAuditLogRepository(db),
		tenantRepo:   NewTenantRepository(db),
	}
}

func (r *postgresRepository) AuditLog() repository.AuditLogRepository {
	return r.auditLogRepo
}

func (r *postgresRepository) Tenant() repository.TenantRepository {
	return r.tenantRepo
}
