package services

import (
	"fmt"
	"time"
	"pr-business/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SystemAccountService 系统账户服务（托管+平台收益）
type SystemAccountService struct {
	db               *gorm.DB
	validatorService *ValidatorService
}

func NewSystemAccountService(db *gorm.DB, validatorService *ValidatorService) *SystemAccountService {
	return &SystemAccountService{
		db:               db,
		validatorService: validatorService,
	}
}

// GetEscrowBalance 获取托管账户总余额（任务托管+平台收益）
func (s *SystemAccountService) GetEscrowBalance(accountType string) (int, error) {
	var accounts []models.SystemAccount
	err := s.db.Where("account_type = ? AND is_active = ?", accountType, true).Find(&accounts).Error
	if err != nil {
		return 0, fmt.Errorf("查询系统账户失败: %w", err)
	}

	// 计算总余额
	totalBalance := 0
	for _, account := range accounts {
		totalBalance += account.Balance
	}

	return totalBalance, nil
}

// UpdateSystemBalance 更新系统账户余额
func (s *SystemAccountService) UpdateSystemBalance(accountID string, amount int, description string, tx *gorm.DB) error {
	id, err := uuid.Parse(accountID)
	if err != nil {
		return fmt.Errorf("无效的账户ID: %w", err)
	}

	account := models.SystemAccount{ID: id}
	if err := tx.First(&account).Error; err != nil {
		return err
	}

	account.Balance += amount
	account.UpdatedAt = time.Now()

	if err := tx.Save(&account).Error; err != nil {
		return err
	}

	return nil
}
