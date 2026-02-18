package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TaskInvitationCode 任务邀请码模型
type TaskInvitationCode struct {
	ID            uuid.UUID  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Code          string     `gorm:"type:varchar(20);not null;uniqueIndex" json:"code"`
	CampaignID    uuid.UUID  `gorm:"type:uuid;not null;index:idx_task_invitation_codes_campaign" json:"campaignId"`
	GeneratorID   string     `gorm:"type:varchar(255);not null" json:"generatorId"`
	GeneratorType string     `gorm:"type:varchar(50);not null" json:"generatorType"`
	MaxUses       int        `gorm:"type:int;default:100" json:"maxUses"`
	UseCount      int        `gorm:"type:int;default:0" json:"useCount"`
	ExpiresAt     *time.Time `json:"expiresAt"`
	CreatedAt     time.Time  `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt     time.Time  `gorm:"not null;default:now()" json:"updatedAt"`

	// 关联
	Campaign *Campaign `gorm:"foreignKey:CampaignID" json:"campaign,omitempty"`
}

// TableName 指定表名
func (TaskInvitationCode) TableName() string {
	return "task_invitation_codes"
}

// BeforeCreate GORM Hook
func (tic *TaskInvitationCode) BeforeCreate(tx *gorm.DB) error {
	if tic.ID == uuid.Nil {
		tic.ID = uuid.New()
	}
	return nil
}

// IsValid 检查邀请码是否有效
func (tic *TaskInvitationCode) IsValid() bool {
	// 检查过期时间
	if tic.ExpiresAt != nil && time.Now().After(*tic.ExpiresAt) {
		return false
	}
	// 检查使用次数
	if tic.MaxUses > 0 && tic.UseCount >= tic.MaxUses {
		return false
	}
	return true
}

// IncrementUse 增加使用次数
func (tic *TaskInvitationCode) IncrementUse() {
	tic.UseCount++
}

// TaskInvitation 任务邀请记录模型
type TaskInvitation struct {
	ID                uuid.UUID  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	InvitationCodeID  uuid.UUID  `gorm:"type:uuid;not null;index:idx_task_invitations_invitation;uniqueIndex:idx_task_invitations_task" json:"invitationCodeId"`
	CreatorID         *string    `gorm:"type:varchar(255);index:idx_task_invitations_creator" json:"creatorId"`
	TaskID            *uuid.UUID `gorm:"type:uuid;index:idx_task_invitations_task" json:"taskId"`
	Status            string     `gorm:"type:varchar(20);not null;default:'accepted'" json:"status"`
	AcceptedAt        time.Time  `gorm:"not null;default:now()" json:"acceptedAt"`
	CompletedAt       *time.Time `json:"completedAt"`
	CreatedAt         time.Time  `gorm:"not null;default:now()" json:"createdAt"`

	// 关联
	InvitationCode *TaskInvitationCode `gorm:"foreignKey:InvitationCodeID" json:"invitationCode,omitempty"`
	Task           *Task               `gorm:"foreignKey:TaskID" json:"task,omitempty"`
}

// TableName 指定表名
func (TaskInvitation) TableName() string {
	return "task_invitations"
}

// BeforeCreate GORM Hook
func (ti *TaskInvitation) BeforeCreate(tx *gorm.DB) error {
	if ti.ID == uuid.Nil {
		ti.ID = uuid.New()
	}
	return nil
}
