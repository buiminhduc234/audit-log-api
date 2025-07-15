package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/buiminhduc234/audit-log-api/internal/domain"
	"github.com/buiminhduc234/audit-log-api/internal/utils"
)

type AuditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

func (r *AuditLogRepository) Create(ctx context.Context, log *domain.AuditLog) error {
	if log.ID == "" {
		log.ID = uuid.New().String()
	}

	tenantID, err := utils.GetTenantIDFromContext(ctx)
	if err != nil {
		return err
	}
	log.TenantID = tenantID

	return r.db.WithContext(ctx).Create(log).Error
}

func (r *AuditLogRepository) GetByID(ctx context.Context, id string) (*domain.AuditLog, error) {
	var log domain.AuditLog

	db, err := getTenantScope(r.db, ctx)
	if err != nil {
		return nil, err
	}

	if err := db.First(&log, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *AuditLogRepository) List(ctx context.Context, filter domain.AuditLogFilter) ([]domain.AuditLog, error) {
	var logs []domain.AuditLog

	db, err := getTenantScope(r.db, ctx)
	if err != nil {
		return nil, err
	}

	// Apply additional filters
	if filter.UserID != "" {
		db = db.Where("user_id = ?", filter.UserID)
	}
	if filter.Action != "" {
		db = db.Where("action = ?", filter.Action)
	}
	if filter.ResourceType != "" {
		db = db.Where("resource_type = ?", filter.ResourceType)
	}
	if filter.ResourceID != "" {
		db = db.Where("resource_id = ?", filter.ResourceID)
	}
	if filter.Severity != "" {
		db = db.Where("severity = ?", filter.Severity)
	}
	if !filter.StartTime.IsZero() {
		db = db.Where("timestamp >= ?", filter.StartTime)
	}
	if !filter.EndTime.IsZero() {
		db = db.Where("timestamp <= ?", filter.EndTime)
	}

	// Apply pagination
	if filter.Limit > 0 {
		db = db.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		db = db.Offset(filter.Offset)
	}

	// Apply sorting
	db = db.Order("timestamp DESC")

	if err := db.Find(&logs).Error; err != nil {
		return nil, err
	}

	return logs, nil
}

func (r *AuditLogRepository) Delete(ctx context.Context, id string) error {
	db, err := getTenantScope(r.db, ctx)
	if err != nil {
		return err
	}

	return db.Delete(&domain.AuditLog{}, "id = ?", id).Error
}

func (r *AuditLogRepository) DeleteOlderThan(ctx context.Context, days int) error {
	db, err := getTenantScope(r.db, ctx)
	if err != nil {
		return err
	}

	return db.Where("timestamp < NOW() - INTERVAL '? days'", days).
		Delete(&domain.AuditLog{}).
		Error
}

func (r *AuditLogRepository) BulkCreate(ctx context.Context, logs []domain.AuditLog) error {
	tenantID, err := utils.GetTenantIDFromContext(ctx)
	if err != nil {
		return err
	}

	// Generate UUIDs for logs without IDs
	for i := range logs {
		if logs[i].ID == "" {
			logs[i].ID = uuid.New().String()
		}
		logs[i].TenantID = tenantID
	}

	return r.db.WithContext(ctx).CreateInBatches(logs, 100).Error
}
