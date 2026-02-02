package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID                string         `gorm:"primaryKey;type:varchar(255)" json:"id"`
	AuthCenterUserID  string         `gorm:"type:uuid;not null;uniqueIndex" json:"authCenterUserId"`
	Nickname          string         `gorm:"type:varchar(50)" json:"nickname"`
	AvatarURL         string         `gorm:"type:varchar(500)" json:"avatarUrl"`
	Profile           string         `gorm:"type:jsonb;default:'{}'::jsonb" json:"profile"`
	Roles             []string       `gorm:"type:jsonb;not null;default:'[]'::jsonb" json:"roles"`
	CurrentRole        string         `gorm:"type:varchar(50)" json:"currentRole"`
	LastUsedRole      string         `gorm:"type:varchar(50)" json:"lastUsedRole"`
	InvitedBy         string         `gorm:"type:varchar(255)" json:"invitedBy"`
	InvitationCodeID  string         `gorm:"type:uuid" json:"invitationCodeId"`
	Status            string         `gorm:"type:varchar(20);not null;default:'active';check:status IN ('active', 'banned', 'inactive')" json:"status"`
	LastLoginAt       *time.Time     `json:"lastLoginAt"`
	LastLoginIP       string         `gorm:"type:varchar(50)" json:"lastLoginIp"`
	CreatedAt         time.Time      `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt         time.Time      `gorm:"not null;default:now()" json:"updatedAt"`
	DeletedAt         *time.Time     `json:"deletedAt"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate GORM Hook
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = "usr_" + uuid.New().String()
	}
	return nil
}
