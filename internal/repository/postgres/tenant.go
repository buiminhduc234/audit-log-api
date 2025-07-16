package postgres

import (
	"context"

	"gorm.io/gorm"

	"github.com/buiminhduc234/audit-log-api/internal/domain"
)

type TenantRepository struct {
	db *gorm.DB
}

func NewTenantRepository(db *gorm.DB) *TenantRepository {
	return &TenantRepository{db: db}
}

func (r *TenantRepository) Create(ctx context.Context, tenant *domain.Tenant) (*domain.Tenant, error) {
	if err := r.db.WithContext(ctx).Create(tenant).Error; err != nil {
		return nil, err
	}
	return tenant, nil
}

func (r *TenantRepository) GetByID(ctx context.Context, id string) (*domain.Tenant, error) {
	var tenant domain.Tenant
	if err := r.db.WithContext(ctx).First(&tenant, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *TenantRepository) GetByAPIKey(ctx context.Context, apiKey string) (*domain.Tenant, error) {
	var tenant domain.Tenant
	if err := r.db.WithContext(ctx).First(&tenant, "api_key = ?", apiKey).Error; err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *TenantRepository) Update(ctx context.Context, tenant *domain.Tenant) error {
	return r.db.WithContext(ctx).Save(tenant).Error
}

func (r *TenantRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.Tenant{}, "id = ?", id).Error
}

func (r *TenantRepository) List(ctx context.Context) ([]domain.Tenant, error) {
	var tenants []domain.Tenant
	if err := r.db.WithContext(ctx).Find(&tenants).Error; err != nil {
		return nil, err
	}
	return tenants, nil
}
