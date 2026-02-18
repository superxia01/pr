package constants

// 用户角色常量（统一使用大写下划线格式）
const (
	// 超级管理系统管理员
	RoleSuperAdmin = "SUPER_ADMIN"

	// 服务商管理员
	RoleServiceProviderAdmin = "SERVICE_PROVIDER_ADMIN"

	// 服务商员工
	RoleServiceProviderStaff = "SERVICE_PROVIDER_STAFF"

	// 商家管理员
	RoleMerchantAdmin = "MERCHANT_ADMIN"

	// 商家员工
	RoleMerchantStaff = "MERCHANT_STAFF"

	// 达人
	RoleCreator = "CREATOR"

	// 基础用户（新注册用户的默认角色）
	RoleBasicUser = "BASIC_USER"
)

// GetAllRoles 返回所有角色列表
func GetAllRoles() []string {
	return []string{
		RoleSuperAdmin,
		RoleServiceProviderAdmin,
		RoleServiceProviderStaff,
		RoleMerchantAdmin,
		RoleMerchantStaff,
		RoleCreator,
		RoleBasicUser,
	}
}

// IsValidRole 检查角色是否有效
func IsValidRole(role string) bool {
	for _, r := range GetAllRoles() {
		if r == role {
			return true
		}
	}
	return false
}

// GetRoleDisplayName 获取角色显示名称
func GetRoleDisplayName(role string) string {
	switch role {
	case RoleSuperAdmin:
		return "超级管理员"
	case RoleServiceProviderAdmin:
		return "服务商管理员"
	case RoleServiceProviderStaff:
		return "服务商员工"
	case RoleMerchantAdmin:
		return "商家管理员"
	case RoleMerchantStaff:
		return "商家员工"
	case RoleCreator:
		return "达人"
	case RoleBasicUser:
		return "基础用户"
	default:
		return "未知角色"
	}
}

// GetRoleShortName 获取角色短名称（用于代码中）
func GetRoleShortName(role string) string {
	switch role {
	case RoleSuperAdmin:
		return "超管"
	case RoleServiceProviderAdmin:
		return "服务商管理员"
	case RoleServiceProviderStaff:
		return "服务商员工"
	case RoleMerchantAdmin:
		return "商家管理员"
	case RoleMerchantStaff:
		return "商家员工"
	case RoleCreator:
		return "达人"
	case RoleBasicUser:
		return "基础用户"
	default:
		return "未知"
	}
}

// ========== 邀请权限配置 ==========

// InvitableRole 可邀请的角色信息
type InvitableRole struct {
	Role        string // 角色名
	Code        string // 邀请码后缀
	Label       string // 显示名称
}

// GetInvitableRoles 获取用户可以邀请的角色列表（基于用户拥有的所有角色）
// 返回: 可邀请的角色数组（去重后）
func GetInvitableRoles(roles []string) []InvitableRole {
	// 使用 map 去重
	uniqueRoles := make(map[string]InvitableRole)

	// 遍历用户的所有角色，收集可邀请的角色
	for _, role := range roles {
		var invitable []InvitableRole

		switch role {
		case RoleSuperAdmin:
			// 超级管理员可以邀请所有类型
			invitable = []InvitableRole{
				{Role: RoleServiceProviderAdmin, Code: "SP-ADMIN", Label: "服务商管理员"},
				{Role: RoleServiceProviderStaff, Code: "SP-STAFF", Label: "服务商员工"},
				{Role: RoleMerchantAdmin, Code: "M-ADMIN", Label: "商家管理员"},
				{Role: RoleMerchantStaff, Code: "M-STAFF", Label: "商家员工"},
				{Role: RoleCreator, Code: "CREATOR", Label: "达人"},
			}

		case RoleServiceProviderAdmin:
			// 服务商管理员可以邀请：服务商员工、商家管理员、达人
			invitable = []InvitableRole{
				{Role: RoleServiceProviderStaff, Code: "SP-STAFF", Label: "服务商员工"},
				{Role: RoleMerchantAdmin, Code: "M-ADMIN", Label: "商家管理员"},
				{Role: RoleCreator, Code: "CREATOR", Label: "达人"},
			}

		case RoleMerchantAdmin:
			// 商家管理员可以邀请：商家员工
			invitable = []InvitableRole{
				{Role: RoleMerchantStaff, Code: "M-STAFF", Label: "商家员工"},
			}

		case RoleServiceProviderStaff:
			// 服务商员工可以邀请：达人
			invitable = []InvitableRole{
				{Role: RoleCreator, Code: "CREATOR", Label: "达人"},
			}

		case RoleMerchantStaff, RoleCreator:
			// 商家员工和达人不能邀请他人
			invitable = []InvitableRole{}

		default:
			invitable = []InvitableRole{}
		}

		// 添加到 map 中去重
		for _, ir := range invitable {
			uniqueRoles[ir.Role] = ir
		}
	}

	// 将 map 转换为切片返回
	result := make([]InvitableRole, 0, len(uniqueRoles))
	for _, ir := range uniqueRoles {
		result = append(result, ir)
	}

	return result
}

// CanInviteRole 检查用户是否可以邀请指定角色（基于用户拥有的所有角色）
func CanInviteRole(roles []string, targetRole string) bool {
	invitableRoles := GetInvitableRoles(roles)
	for _, ir := range invitableRoles {
		if ir.Role == targetRole {
			return true
		}
	}
	return false
}
