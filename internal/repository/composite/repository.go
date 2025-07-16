package composite

import (
	"github.com/buiminhduc234/audit-log-api/internal/config"
	"github.com/buiminhduc234/audit-log-api/internal/repository"
	"github.com/buiminhduc234/audit-log-api/internal/repository/opensearch"
	"github.com/buiminhduc234/audit-log-api/internal/repository/postgres"
	opensearchclient "github.com/opensearch-project/opensearch-go/v2"
	"gorm.io/gorm"
)

type compositeRepository struct {
	postgresRepo repository.PostgresRepository
	osRepo       repository.OpenSearchRepository
}

func NewCompositeRepository(db *gorm.DB, osClient *opensearchclient.Client, osConfig *config.OpenSearchConfig) repository.Repository {
	return &compositeRepository{
		postgresRepo: postgres.NewPostgresRepository(db),
		osRepo:       opensearch.NewRepository(osClient, osConfig),
	}
}

// Implement PostgresRepository interface
func (r *compositeRepository) AuditLog() repository.AuditLogRepository {
	return r.postgresRepo.AuditLog()
}

func (r *compositeRepository) Tenant() repository.TenantRepository {
	return r.postgresRepo.Tenant()
}

// Implement OpenSearch access
func (r *compositeRepository) OpenSearch() repository.OpenSearchRepository {
	return r.osRepo
}
