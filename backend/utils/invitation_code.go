package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"
)

// ========== 固定邀请码（人邀请人） ==========

// RoleCode 角色代码映射
type RoleCode struct {
	Role        string // 完整角色名
	Code        string // 邀请码后缀代码
	Label       string // 显示名称
}

// 角色代码映射表
var RoleCodeMapping = []RoleCode{
	{Role: "SERVICE_PROVIDER_ADMIN", Code: "SP-ADMIN", Label: "服务商管理员"},
	{Role: "SERVICE_PROVIDER_STAFF", Code: "SP-STAFF", Label: "服务商员工"},
	{Role: "MERCHANT_ADMIN", Code: "M-ADMIN", Label: "商家管理员"},
	{Role: "MERCHANT_STAFF", Code: "M-STAFF", Label: "商家员工"},
	{Role: "CREATOR", Code: "CREATOR", Label: "达人"},
}

// GetRoleCodeByRole 根据角色名获取角色代码
func GetRoleCodeByRole(role string) string {
	for _, rc := range RoleCodeMapping {
		if rc.Role == role {
			return rc.Code
		}
	}
	return ""
}

// GetRoleByCode 根据角色代码获取角色名
func GetRoleByCode(code string) string {
	for _, rc := range RoleCodeMapping {
		if rc.Code == code {
			return rc.Role
		}
	}
	return ""
}

// RequiresOrganizationBinding 检查角色是否需要组织绑定
// 需要组织绑定的角色：服务商管理员、服务商员工、商家管理员、商家员工
func RequiresOrganizationBinding(role string) bool {
	switch role {
	case "SERVICE_PROVIDER_ADMIN", "SERVICE_PROVIDER_STAFF", "MERCHANT_ADMIN", "MERCHANT_STAFF":
		return true
	case "CREATOR":
		return false
	default:
		return false
	}
}

// GetOrganizationTypeByRole 根据角色获取组织类型
func GetOrganizationTypeByRole(role string) string {
	switch role {
	case "SERVICE_PROVIDER_ADMIN", "SERVICE_PROVIDER_STAFF":
		return "service_provider"
	case "MERCHANT_ADMIN", "MERCHANT_STAFF":
		return "merchant"
	default:
		return ""
	}
}

// GenerateUserFixedInvitationCode 生成用户固定邀请码（短码版本）
// 生成8位随机短码，完整信息存储在数据库映射表中
// 短码格式: 8位大写字母+数字（排除易混淆字符 0/O、1/I/l）
//   例如: ABC123XY
func GenerateUserFixedInvitationCode(userID string, targetRole string, organizationID ...string) string {
	// 直接生成8位随机短码
	return generateRandomString(8)
}

// ParseUserFixedInvitationCode 解析用户固定邀请码
// 返回: (邀请人基础码, 组织ID, 目标角色, 是否有效)
// 需要组织绑定: "INV-785aa30f-provider_abc123-SP-ADMIN" -> ("INV-785aa30f", "provider_abc123", "SERVICE_PROVIDER_ADMIN", true)
// 不需要组织绑定: "INV-785aa30f-CREATOR" -> ("INV-785aa30f", "", "CREATOR", true)
func ParseUserFixedInvitationCode(code string) (string, string, string, bool) {
	parts := strings.Split(code, "-")
	if len(parts) < 3 || parts[0] != "INV" {
		return "", "", "", false
	}

	targetRole := GetRoleByCode(parts[len(parts)-1])
	if targetRole == "" {
		return "", "", "", false
	}

	inviterBase := fmt.Sprintf("INV-%s", parts[1])

	// 判断是否需要组织绑定
	if RequiresOrganizationBinding(targetRole) {
		// 格式: INV-{inviter}-{organization}-{role}
		if len(parts) != 4 {
			return "", "", "", false
		}
		organizationID := parts[2]
		return inviterBase, organizationID, targetRole, true
	}

	// 格式: INV-{inviter}-{role}
	if len(parts) != 3 {
		return "", "", "", false
	}
	return inviterBase, "", targetRole, true
}

// extractIDSuffix 从ID中提取后n位
// 例如: "usr_abc123def456" -> "def456" (后6位)
func extractIDSuffix(id string, n int) string {
	// 去掉前缀（如果有）
	parts := strings.Split(id, "_")
	var baseID string
	if len(parts) > 1 {
		baseID = parts[len(parts)-1]
	} else {
		baseID = id
	}

	// 取后n位
	if len(baseID) <= n {
		return baseID
	}
	return baseID[len(baseID)-n:]
}

// ========== 活动邀请码（营销活动） ==========

// GenerateCampaignInvitationCode 生成活动邀请码
// 格式: CAMP-{活动ID后6位}-{随机4位}
// 例如: CAMP-a1b2c3-X9Y7
func GenerateCampaignInvitationCode(campaignID string) string {
	suffix := extractIDSuffix(campaignID, 6)
	randomPart := generateRandomString(4)
	return fmt.Sprintf("CAMP-%s-%s", suffix, randomPart)
}

// ParseCampaignInvitationCode 解析活动邀请码
// 返回: (活动ID后6位, 是否有效)
func ParseCampaignInvitationCode(code string) (string, bool) {
	parts := strings.Split(code, "-")
	if len(parts) != 3 || parts[0] != "CAMP" {
		return "", false
	}
	return parts[1], true
}

// generateRandomString 生成随机字符串
func generateRandomString(n int) string {
	const chars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // 排除易混淆字符
	code := make([]byte, n)
	for i := 0; i < n; i++ {
		randInt, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			// 如果随机数生成失败，使用时间戳作为后备
			return fmt.Sprintf("%0*d", n, time.Now().UnixNano()%int64(pow10(n)))
		}
		code[i] = chars[randInt.Int64()]
	}
	return string(code)
}

func pow10(n int) int64 {
	result := int64(1)
	for i := 0; i < n; i++ {
		result *= 10
	}
	return result
}

// ========== 旧版兼容（保留用于其他地方） ==========

// GenerateFixedInvitationCode 旧版固定邀请码生成（兼容保留）
// 新代码请使用 GenerateUserFixedInvitationCode
func GenerateFixedInvitationCode(codeType string, id string, extraID ...string) string {
	// 提取ID后8位（去掉前缀）
	suffix := extractIDSuffix(id, 8)

	switch codeType {
	case "USER":
		return fmt.Sprintf("INV-%s", suffix)
	case "ADMIN_MASTER":
		return "ADMIN-MASTER"
	case "SP_ADMIN":
		return fmt.Sprintf("SP-%s", suffix)
	case "MERCHANT":
		if len(extraID) > 0 {
			providerSuffix := extractIDSuffix(extraID[0], 8)
			return fmt.Sprintf("MERCHANT-%s-%s", providerSuffix, suffix)
		}
		return fmt.Sprintf("MERCHANT-%s", suffix)
	case "CREATOR":
		return fmt.Sprintf("CREATOR-%s", suffix)
	case "STAFF":
		staffType := "STAFF"
		if len(extraID) > 1 {
			staffType = extraID[1]
		}
		if len(extraID) > 0 {
			ownerSuffix := extractIDSuffix(extraID[0], 8)
			return fmt.Sprintf("STAFF-%s-%s", ownerSuffix, staffType)
		}
		return fmt.Sprintf("STAFF-%s-%s", suffix, staffType)
	default:
		return fmt.Sprintf("UNKNOWN-%s", suffix)
	}
}

// GenerateInvitationCode 生成6位随机邀请码（任务邀请）
func GenerateInvitationCode() string {
	return generateRandomString(6)
}

// ExtractUserIDFromFixedCode 从固定邀请码中提取用户ID
// 例如: "INV-785aa30f-SP-ADMIN" -> "usr_...785aa30f"
func ExtractUserIDFromFixedCode(code string) string {
	parts := strings.Split(code, "-")
	if len(parts) < 2 || parts[0] != "INV" {
		return ""
	}
	// 返回后8位，用于查询用户
	return parts[1]
}

// StringPtr 返回字符串指针
func StringPtr(s string) *string {
	return &s
}
