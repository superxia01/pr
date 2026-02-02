package utils

import (
	"pr-business/models"

	"github.com/gin-gonic/gin"
)

// HasRole 检查用户是否拥有指定角色
func HasRole(user *models.User, role string) bool {
	for _, r := range user.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// HasAnyRole 检查用户是否拥有任一指定角色
func HasAnyRole(user *models.User, roles ...string) bool {
	for _, role := range roles {
		if HasRole(user, role) {
			return true
		}
	}
	return false
}

// HasAllRoles 检查用户是否拥有所有指定角色
func HasAllRoles(user *models.User, roles ...string) bool {
	for _, role := range roles {
		if !HasRole(user, role) {
			return false
		}
	}
	return true
}

// GetCurrentUser 从上下文中获取当前用户
func GetCurrentUser(c *gin.Context) (*models.User, bool) {
	user, exists := c.Get("user")
	if !exists {
		return nil, false
	}
	return user.(*models.User), true
}

// IsSuperAdmin 检查是否是超级管理员
func IsSuperAdmin(user *models.User) bool {
	return HasRole(user, "SUPER_ADMIN")
}

// IsMerchantAdmin 检查是否是商家管理员
func IsMerchantAdmin(user *models.User) bool {
	return HasRole(user, "MERCHANT_ADMIN")
}

// IsServiceProviderAdmin 检查是否是服务商管理员
func IsServiceProviderAdmin(user *models.User) bool {
	return HasRole(user, "SP_ADMIN")
}

// IsCreator 检查是否是达人
func IsCreator(user *models.User) bool {
	return HasRole(user, "CREATOR")
}

// IsMerchantStaff 检查是否是商家员工
func IsMerchantStaff(user *models.User) bool {
	return HasRole(user, "MERCHANT_STAFF")
}

// IsServiceProviderStaff 检查是否是服务商员工
func IsServiceProviderStaff(user *models.User) bool {
	return HasRole(user, "SP_STAFF")
}
