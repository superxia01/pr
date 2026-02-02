package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// InvitationCode 邀请码模型
type InvitationCode struct {
	ID          string         `gorm:"primaryKey;type:varchar(255)" json:"id"`
	Code        string         `gorm:"type:varchar(50);not null;uniqueIndex" json:"code"`
	CodeType    string         `gorm:"type:varchar(50);not null;check:code_type IN ('ADMIN_MASTER', 'SP_ADMIN', 'MERCHANT', 'CREATOR', 'STAFF')" json:"codeType"`
	OwnerID     string         `gorm:"type:varchar(255)" json:"ownerId"`
	OwnerType   string         `gorm:"type:varchar(50);check:owner_type IN ('super_admin', 'service_provider', 'merchant')" json:"ownerType"`
	Status      string         `gorm:"type:varchar(20);not null;default:'active';check:status IN ('active', 'disabled', 'expired', 'used')" json:"status"`
	MaxUses     int            `gorm:"not null;default:1" json:"maxUses"`
	UseCount    int            `gorm:"not null;default:0" json:"useCount"`
	ExpiresAt   *time.Time     `json:"expiresAt"`
	UsedBy      []string       `gorm:"type:jsonb;not null;default:'[]'::jsonb" json:"usedBy"`
	Metadata    string         `gorm:"type:jsonb;default:'{}'::jsonb" json:"metadata"`
	CreatedAt   time.Time      `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"not null;default:now()" json:"updatedAt"`
	DeletedAt   *time.Time     `json:"deletedAt"`
}

// TableName 指定表名
func (InvitationCode) TableName() string {
	return "invitation_codes"
}

// BeforeCreate GORM Hook
func (ic *InvitationCode) BeforeCreate(tx *gorm.DB) error {
	if ic.ID == "" {
		ic.ID = "ic_" + uuid.New().String()
	}
	return nil
}

// CanBeUsed 检查邀请码是否可用
func (ic *InvitationCode) CanBeUsed() bool {
	if ic.Status != "active" {
		return false
	}
	if ic.MaxUses > 0 && ic.UseCount >= ic.MaxUses {
		return false
	}
	if ic.ExpiresAt != nil && ic.ExpiresAt.Before(time.Now()) {
		return false
	}
	return true
}

// IncrementUse 增加使用次数
func (ic *InvitationCode) IncrementUse(userID string) {
	ic.UseCount++
	ic.UsedBy = append(ic.UsedBy, userID)
	if ic.MaxUses > 0 && ic.UseCount >= ic.MaxUses {
		ic.Status = "used"
	}
}
