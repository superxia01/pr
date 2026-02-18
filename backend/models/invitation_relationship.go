package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// InvitationRelationship 邀请关系记录模型
// 记录谁邀请了谁、什么角色、绑定哪个组织、邀请时间
type InvitationRelationship struct {
	ID              uuid.UUID  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	InviterID       string     `gorm:"type:varchar(255);not null;index" json:"inviterId"`                 // 邀请人ID
	InviterRole     string     `gorm:"type:varchar(50);not null" json:"inviterRole"`                    // 邀请人角色
	InviteeID       string     `gorm:"type:varchar(255);not null;index" json:"inviteeId"`                // 被邀请人ID
	InviteeRole     string     `gorm:"type:varchar(50);not null" json:"inviteeRole"`                   // 被邀请人角色
	OrganizationID  *uuid.UUID `gorm:"type:uuid" json:"organizationId"`                                     // 绑定的组织ID（可选）
	OrganizationType *string    `gorm:"type:varchar(50)" json:"organizationType"`                          // 绑定的组织类型（可选）
	InvitationCode  string     `gorm:"type:varchar(30);not null" json:"invitationCode"`                  // 使用的邀请码
	InvitedAt       time.Time  `gorm:"not null;default:now()" json:"invitedAt"`                            // 邀请时间（成功使用时）
	Status         string     `gorm:"type:varchar(20);not null;default:'active';check:status IN ('active', 'cancelled')" json:"status"` // 状态
	CreatedAt      time.Time  `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"not null;default:now()" json:"updatedAt"`

	// 关联
	Inviter *User `gorm:"foreignKey:InviterID" json:"inviter,omitempty"`
	Invitee *User `gorm:"foreignKey:InviteeID" json:"invitee,omitempty"`
}

// TableName 指定表名
func (InvitationRelationship) TableName() string {
	return "invitation_relationships"
}

// BeforeCreate GORM Hook
func (ir *InvitationRelationship) BeforeCreate(tx *gorm.DB) error {
	if ir.ID == uuid.Nil {
		ir.ID = uuid.New()
	}
	return nil
}
