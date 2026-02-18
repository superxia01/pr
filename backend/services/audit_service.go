package services

import (
	"fmt"
	"time"
	"pr-business/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuditService 财务审计日志服务（简化版）
type AuditService struct {
	db *gorm.DB
}

func NewAuditService(db *gorm.DB) *AuditService {
	return &AuditService{
		db: db,
	}
}

// LogFinancialOperation 记录关键财务操作
func (s *AuditService) LogFinancialOperation(
	userID string,
	action string,
	resourceType string,
	resourceID string,
	changes map[string]interface{},
	ipAddress string,
	userAgent string,
) error {
	log := models.FinancialAuditLog{
		ID:           uuid.New(),
		UserID:       userID,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Changes:      changes,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		CreatedAt:    time.Now(),
	}

	if err := s.db.Create(&log).Error; err != nil {
		return fmt.Errorf("创建审计日志失败: %w", err)
	}

	return nil
}

// QueryAuditLogs 查询审计日志
func (s *AuditService) QueryAuditLogs(
	userID string,
	action string,
	resourceType string,
	resourceID string,
	limit int,
	offset int,
) ([]models.FinancialAuditLog, error) {
	query := s.db.Model(&models.FinancialAuditLog{}).Where("user_id = ?", userID)

	if action != "" {
		query = query.Where("action = ?", action)
	}

	if resourceType != "" {
		query = query.Where("resource_type = ?", resourceType)
	}

	if resourceID != "" {
		query = query.Where("resource_id = ?", resourceID)
	}

	query = query.Order("created_at DESC").Limit(limit).Offset(offset)

	var logs []models.FinancialAuditLog
	if err := query.Find(&logs).Error; err != nil {
		return nil, fmt.Errorf("查询审计日志失败: %w", err)
	}

	return logs, nil
}
