package utils

import (
	"pr-business/constants"
	"pr-business/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HasPermission 检查用户是否拥有指定权限
// 管理员（SUPER_ADMIN, SERVICE_PROVIDER_ADMIN, MERCHANT_ADMIN）默认拥有所有权限
// 员工（SERVICE_PROVIDER_STAFF, MERCHANT_STAFF）需要检查权限表
func HasPermission(db *gorm.DB, user *models.User, permissionCode string) bool {
	// 管理员默认拥有所有权限
	if IsSuperAdmin(user) || IsServiceProviderAdmin(user) || IsMerchantAdmin(user) {
		return true
	}

	// 服务商员工需要检查权限表
	if IsServiceProviderStaff(user) {
		var staff models.ServiceProviderStaff
		if err := db.Where("user_id::text = ?", user.AuthCenterUserID).First(&staff).Error; err != nil {
			return false
		}
		var perm models.ServiceProviderStaffPermission
		return db.Where("staff_id = ? AND permission_code = ?", staff.ID, permissionCode).
			First(&perm).Error == nil
	}

	// 商家员工需要检查权限表
	if IsMerchantStaff(user) {
		var staff models.MerchantStaff
		if err := db.Where("user_id::text = ?", user.AuthCenterUserID).First(&staff).Error; err != nil {
			return false
		}
		var perm models.MerchantStaffPermission
		return db.Where("staff_id = ? AND permission_code = ?", staff.ID, permissionCode).
			First(&perm).Error == nil
	}

	// 其他角色（BASIC_USER, CREATOR）没有特殊权限
	return false
}

// HasAnyPermission 检查用户是否拥有任一指定权限
func HasAnyPermission(db *gorm.DB, user *models.User, permissionCodes ...string) bool {
	for _, code := range permissionCodes {
		if HasPermission(db, user, code) {
			return true
		}
	}
	return false
}

// HasAllPermissions 检查用户是否拥有所有指定权限
func HasAllPermissions(db *gorm.DB, user *models.User, permissionCodes ...string) bool {
	for _, code := range permissionCodes {
		if !HasPermission(db, user, code) {
			return false
		}
	}
	return true
}

// RequirePermission 权限检查中间件
// 返回一个gin中间件函数，用于检查用户是否拥有指定权限
func RequirePermission(db *gorm.DB, permissionCode string) func(*gin.Context) {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(401, gin.H{"error": "未认证"})
			c.Abort()
			return
		}

		u := user.(*models.User)

		if !HasPermission(db, u, permissionCode) {
			c.JSON(403, gin.H{
				"error": "无权限",
				"requiredPermission": permissionCode,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyPermission 权限检查中间件（拥有任一权限即可）
func RequireAnyPermission(db *gorm.DB, permissionCodes ...string) func(*gin.Context) {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(401, gin.H{"error": "未认证"})
			c.Abort()
			return
		}

		u := user.(*models.User)

		if !HasAnyPermission(db, u, permissionCodes...) {
			c.JSON(403, gin.H{
				"error": "无权限",
				"requiredPermissions": permissionCodes,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetCurrentUserPermissions 获取用户的所有权限列表
// 返回用户拥有的所有权限码
func GetCurrentUserPermissions(db *gorm.DB, user *models.User) []string {
	permissions := make([]string, 0)

	// 管理员拥有所有权限
	if IsSuperAdmin(user) || IsServiceProviderAdmin(user) || IsMerchantAdmin(user) {
		return []string{
			constants.PermissionManageCreators,
			constants.PermissionEditCreatorInfo,
			constants.PermissionDeleteCreator,
			constants.PermissionPublishCampaign,
			constants.PermissionEditCommission,
			constants.PermissionEditCampaignInfo,
			constants.PermissionDeleteCampaign,
			constants.PermissionReviewTask,
			constants.PermissionViewAllTasks,
			constants.PermissionRecharge,
			constants.PermissionApproveWithdrawal,
			constants.PermissionViewFinancialReports,
			constants.PermissionManageMerchant,
			constants.PermissionManageStaff,
			constants.PermissionCreateInvitationCode,
		}
	}

	// 服务商员工：从权限表中获取
	if IsServiceProviderStaff(user) {
		var staff models.ServiceProviderStaff
		if err := db.Where("user_id::text = ?", user.AuthCenterUserID).First(&staff).Error; err == nil {
			var perms []models.ServiceProviderStaffPermission
			db.Where("staff_id = ?", staff.ID).Find(&perms)
			for _, perm := range perms {
				permissions = append(permissions, perm.PermissionCode)
			}
		}
	}

	// 商家员工：从权限表中获取
	if IsMerchantStaff(user) {
		var staff models.MerchantStaff
		if err := db.Where("user_id::text = ?", user.AuthCenterUserID).First(&staff).Error; err == nil {
			var perms []models.MerchantStaffPermission
			db.Where("staff_id = ?", staff.ID).Find(&perms)
			for _, perm := range perms {
				permissions = append(permissions, perm.PermissionCode)
			}
		}
	}

	return permissions
}
