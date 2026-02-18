package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CampaignInvitation 活动邀请码模型
type CampaignInvitation struct {
	ID             uuid.UUID      `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	CampaignID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"campaignId"`
	Code           string         `gorm:"type:varchar(30);not null;uniqueIndex" json:"code"`
	MaxUses       int            `gorm:"type:int;not null;default:1" json:"maxUses"` // 最大使用次数
	UseCount       int            `gorm:"type:int;not null;default:0" json:"useCount"`
	ExpiresAt      *time.Time     `gorm:"type:timestamp;null" json:"expiresAt"` // 过期时间
	Status         string         `gorm:"type:varchar(20);not null;default:'active';index" json:"status"`
	CreatedAt      time.Time      `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt      time.Time      `gorm:"not null;default:now()" json:"updatedAt"`
	DeletedAt      *time.Time     `json:"deletedAt"`

	// 关联
	Campaign *Campaign `gorm:"foreignKey:CampaignID" json:"campaign,omitempty"`
}

// TableName 指定表名
func (CampaignInvitation) TableName() string {
	return "campaign_invitations"
}

// BeforeCreate GORM Hook
func (ci *CampaignInvitation) BeforeCreate(tx *gorm.DB) error {
	if ci.ID == uuid.Nil {
		ci.ID = uuid.New()
	}
	return nil
}

// IsActive 检查邀请码是否有效
func (ci *CampaignInvitation) IsActive() bool {
	if ci.Status != "active" {
		return false
	}
	// 检查是否过期
	if ci.ExpiresAt != nil && ci.ExpiresAt.Before(time.Now()) {
		return false
	}
	// 检查是否超过使用次数
	if ci.MaxUses > 0 && ci.UseCount >= ci.MaxUses {
		return false
	}
	return true
}

// IncrementUse 增加使用次数
func (ci *CampaignInvitation) IncrementUse() {
	ci.UseCount++
}

// Deactivate 停用邀请码
func (ci *CampaignInvitation) Deactivate() {
	ci.Status = "inactive"
}
