package utils

import (
	"fmt"
	"strings"
)

// GenerateFixedInvitationCode 生成固定邀请码（基于ID后缀）
// 规则：
// - ADMIN-MASTER: 超级管理员主邀请码（固定）
// - SP-{id后8位}: 服务商管理员邀请码
// - MERCHANT-{provider_id后8位}-{id后8位}: 商家邀请码
// - CREATOR-{id后8位}: 达人邀请码
// - STAFF-{owner_id后8位}-{type}: 员工邀请码
func GenerateFixedInvitationCode(codeType string, id string, extraID ...string) string {
	// 提取ID后8位（去掉前缀）
	suffix := extractIDSuffix(id, 8)

	switch codeType {
	case "ADMIN_MASTER":
		// 超级管理员主邀请码是固定的
		return "ADMIN-MASTER"

	case "SP_ADMIN":
		// 服务商管理员邀请码
		return fmt.Sprintf("SP-%s", suffix)

	case "MERCHANT":
		// 商家邀请码，需要provider_id
		if len(extraID) > 0 {
			providerSuffix := extractIDSuffix(extraID[0], 8)
			return fmt.Sprintf("MERCHANT-%s-%s", providerSuffix, suffix)
		}
		return fmt.Sprintf("MERCHANT-%s", suffix)

	case "CREATOR":
		// 达人邀请码
		return fmt.Sprintf("CREATOR-%s", suffix)

	case "STAFF":
		// 员工邀请码，需要owner_id和staff_type
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

// ValidateInvitationCodeFormat 验证邀请码格式
func ValidateInvitationCodeFormat(code string) (codeType string, valid bool) {
	if code == "ADMIN-MASTER" {
		return "ADMIN_MASTER", true
	}

	if strings.HasPrefix(code, "SP-") {
		return "SP_ADMIN", true
	}

	if strings.HasPrefix(code, "MERCHANT-") {
		return "MERCHANT", true
	}

	if strings.HasPrefix(code, "CREATOR-") {
		return "CREATOR", true
	}

	if strings.HasPrefix(code, "STAFF-") {
		return "STAFF", true
	}

	return "", false
}

// ParseMerchantInvitationCode 解析商家邀请码，返回provider_id后缀
// 例如: "MERCHANT-abc12345-def67890" -> "abc12345"
func ParseMerchantInvitationCode(code string) (providerSuffix string, valid bool) {
	parts := strings.Split(code, "-")
	if len(parts) != 3 || parts[0] != "MERCHANT" {
		return "", false
	}
	return parts[1], true
}
