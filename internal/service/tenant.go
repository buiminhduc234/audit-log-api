package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/buiminhduc234/audit-log-api/internal/domain"
	"github.com/buiminhduc234/audit-log-api/internal/repository"
)

type TenantService struct {
	repo repository.Repository
}

func NewTenantService(repo repository.Repository) *TenantService {
	return &TenantService{repo: repo}
}

func (s *TenantService) Create(ctx context.Context, tenant *domain.Tenant) error {
	if tenant.ID == "" {
		tenant.ID = uuid.New().String()
	}

	now := time.Now()
	tenant.CreatedAt = now
	tenant.UpdatedAt = now

	return s.repo.Tenant().Create(ctx, tenant)
}

func (s *TenantService) GetByID(ctx context.Context, id string) (*domain.Tenant, error) {
	return s.repo.Tenant().GetByID(ctx, id)
}

func (s *TenantService) GetByAPIKey(ctx context.Context, apiKey string) (*domain.Tenant, error) {
	return s.repo.Tenant().GetByAPIKey(ctx, apiKey)
}

func (s *TenantService) Update(ctx context.Context, tenant *domain.Tenant) error {
	tenant.UpdatedAt = time.Now()
	return s.repo.Tenant().Update(ctx, tenant)
}

func (s *TenantService) Delete(ctx context.Context, id string) error {
	return s.repo.Tenant().Delete(ctx, id)
}

func (s *TenantService) List(ctx context.Context) ([]domain.Tenant, error) {
	return s.repo.Tenant().List(ctx)
}
