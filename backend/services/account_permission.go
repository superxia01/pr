package services

import (
	"pr-business/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AccountPermissionService 账户权限服务
type AccountPermissionService struct {
	db *gorm.DB
}

// NewAccountPermissionService 创建账户权限服务
func NewAccountPermissionService(db *gorm.DB) *AccountPermissionService {
	return &AccountPermissionService{db: db}
}

// GrantPersonalAccountPermission 为用户授予个人账户权限
func (s *AccountPermissionService) GrantPersonalAccountPermission(userID uuid.UUID, accountID uuid.UUID) error {
	permission := models.AccountPermission{
		AccountID:   accountID,
		UserID:      userID,
		AccountType: "PERSONAL",
		CanView:     true,
		CanOperate:  true,
	}
	return s.db.Create(&permission).Error
}

// GrantMerchantAccountPermission 为用户授予商家账户权限
func (s *AccountPermissionService) GrantMerchantAccountPermission(userID uuid.UUID, accountID uuid.UUID, canOperate bool) error {
	permission := models.AccountPermission{
		AccountID:   accountID,
		UserID:      userID,
		AccountType: "ORG_MERCHANT",
		CanView:     true,
		CanOperate:  canOperate,
	}
	return s.db.Create(&permission).Error
}

// GrantProviderAccountPermission 为用户授予服务商账户权限
func (s *AccountPermissionService) GrantProviderAccountPermission(userID uuid.UUID, accountID uuid.UUID, canOperate bool) error {
	permission := models.AccountPermission{
		AccountID:   accountID,
		UserID:      userID,
		AccountType: "ORG_PROVIDER",
		CanView:     true,
		CanOperate:  canOperate,
	}
	return s.db.Create(&permission).Error
}

// RevokeAccountPermission 撤销账户权限
func (s *AccountPermissionService) RevokeAccountPermission(accountID, userID uuid.UUID) error {
	return s.db.Where("account_id = ? AND user_id = ?", accountID, userID).Delete(&models.AccountPermission{}).Error
}

// UpdateAccountPermission 更新账户权限
func (s *AccountPermissionService) UpdateAccountPermission(accountID, userID uuid.UUID, canView, canOperate bool) error {
	return s.db.Model(&models.AccountPermission{}).
		Where("account_id = ? AND user_id = ?", accountID, userID).
		Updates(map[string]interface{}{
			"can_view":    canView,
			"can_operate": canOperate,
		}).Error
}

// GetUserAccountPermissions 获取用户的所有账户权限
func (s *AccountPermissionService) GetUserAccountPermissions(userID uuid.UUID) ([]models.AccountPermission, error) {
	var permissions []models.AccountPermission
	err := s.db.Where("user_id = ?", userID).Find(&permissions).Error
	return permissions, err
}

// GetAccountUsers 获取有权限访问某个账户的所有用户
func (s *AccountPermissionService) GetAccountUsers(accountID uuid.UUID) ([]models.AccountPermission, error) {
	var permissions []models.AccountPermission
	err := s.db.Where("account_id = ?", accountID).Preload("User").Find(&permissions).Error
	return permissions, err
}

// CheckAccountPermission 检查用户是否有账户权限
func (s *AccountPermissionService) CheckAccountPermission(accountID, userID uuid.UUID, requiredPermission string) (bool, error) {
	var permission models.AccountPermission
	err := s.db.Where("account_id = ? AND user_id = ?", accountID, userID).First(&permission).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}

	if requiredPermission == "view" {
		return permission.CanView, nil
	} else if requiredPermission == "operate" {
		return permission.CanOperate, nil
	}
	return false, nil
}
