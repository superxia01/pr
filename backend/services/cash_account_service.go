package services

import (
	"errors"
	"fmt"
	"time"
	"pr-business/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CashAccountService 现金账户服务
type CashAccountService struct {
	db               *gorm.DB
	validatorService *ValidatorService
}

func NewCashAccountService(db *gorm.DB, validatorService *ValidatorService) *CashAccountService {
	return &CashAccountService{
		db:               db,
		validatorService: validatorService,
	}
}

// GetActiveCashAccounts 获取激活的现金账户
func (s *CashAccountService) GetActiveCashAccounts(accountType string) ([]models.CashAccount, error) {
	var accounts []models.CashAccount
	err := s.db.Where("account_type = ? AND is_active = ?", accountType, true).Find(&accounts).Error
	if err != nil {
		return nil, fmt.Errorf("查询现金账户失败: %w", err)
	}
	return accounts, nil
}

// UpdateBalance 更新现金账户余额
func (s *CashAccountService) UpdateBalance(accountID string, amount int, description string, tx *gorm.DB) error {
	id, err := uuid.Parse(accountID)
	if err != nil {
		return fmt.Errorf("无效的账户ID: %w", err)
	}

	account := models.CashAccount{ID: id}
	if err := tx.First(&account).Error; err != nil {
		return err
	}

	// 验证金额
	if amount < 0 && amount < -1000000 {
		return errors.New("金额超出允许范围")
	}

	account.Balance += amount
	account.UpdatedAt = time.Now()

	if err := tx.Save(&account).Error; err != nil {
		return err
	}

	return nil
}

// CreateCashAccount 创建现金账户
func (s *CashAccountService) CreateCashAccount(accountType string, description string, operatorID string, initialBalance int) (*models.CashAccount, error) {
	account := models.CashAccount{
		AccountType: accountType,
		Balance:     initialBalance,
		Description: description,
		IsActive:    true,
	}

	if err := s.db.Create(&account).Error; err != nil {
		return nil, fmt.Errorf("创建现金账户失败: %w", err)
	}

	return &account, nil
}
