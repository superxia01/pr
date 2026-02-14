package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// InvitationCode 角色邀请码模型
// 注意：此表现在仅用于组织/实体的固定邀请码（服务商、商家、达人）
// 用户的固定邀请码存储在 users.fixed_invitation_code 字段
type InvitationCode struct {
	ID               uuid.UUID  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Code             string     `gorm:"type:varchar(30);not null;uniqueIndex" json:"code"`
	Type             string     `gorm:"type:varchar(50);not null" json:"type"`                        // 邀请码类型（SERVICE_PROVIDER/MERCHANT/CREATOR）
	TargetRole       string     `gorm:"type:varchar(50);not null" json:"targetRole"`                  // 目标角色
	GeneratorID      string     `gorm:"type:varchar(255);not null" json:"generatorId"`               // 生成者ID（与 users.id 一致）
	GeneratorType    string     `gorm:"type:varchar(50);not null" json:"generatorType"`               // 生成者类型
	OrganizationID   *uuid.UUID `gorm:"type:uuid" json:"organizationId"`                            // 组织ID（可选）
	OrganizationType *string    `gorm:"type:varchar(50)" json:"organizationType"`                     // 组织类型（可选）
	IsActive         bool       `gorm:"not null;default:true" json:"isActive"`                        // 是否激活
	UseCount          int        `gorm:"not null;default:0" json:"useCount"`                           // 已使用次数
	CreatedAt        time.Time  `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt        time.Time  `gorm:"not null;default:now()" json:"updatedAt"`
}

// TableName 指定表名
func (InvitationCode) TableName() string {
	return "invitation_codes"
}

// BeforeCreate GORM Hook
func (ic *InvitationCode) BeforeCreate(tx *gorm.DB) error {
	if ic.ID == uuid.Nil {
		ic.ID = uuid.New()
	}
	return nil
}

// CanBeUsed 检查邀请码是否可用
// 注意：固定邀请码通常不受次数限制
func (ic *InvitationCode) CanBeUsed() bool {
	if !ic.IsActive {
		return false
	}
	return true
}

// IncrementUse 增加使用次数
func (ic *InvitationCode) IncrementUse() {
	ic.UseCount++
}

// Deactivate 停用邀请码
func (ic *InvitationCode) Deactivate() {
	ic.IsActive = false
}
