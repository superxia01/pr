package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Roles 角色数组类型，支持 JSON 序列化
type Roles []string

// Scan 实现 sql.Scanner 接口
func (r *Roles) Scan(value interface{}) error {
	if value == nil {
		*r = []string{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return nil
	}

	return json.Unmarshal(bytes, r)
}

// Value 实现 driver.Valuer 接口
func (r Roles) Value() (driver.Value, error) {
	if len(r) == 0 {
		return "[]", nil
	}
	return json.Marshal(r)
}

// Profile 用户资料类型，支持 JSON 序列化
type Profile map[string]interface{}

// Scan 实现 sql.Scanner 接口
func (p *Profile) Scan(value interface{}) error {
	if value == nil {
		*p = map[string]interface{}{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return nil
	}

	return json.Unmarshal(bytes, p)
}

// Value 实现 driver.Valuer 接口
func (p Profile) Value() (driver.Value, error) {
	if len(p) == 0 {
		return "{}", nil
	}
	return json.Marshal(p)
}

// User 用户模型
type User struct {
	ID               string         `gorm:"primaryKey;type:varchar(255)" json:"id"`
	AuthCenterUserID string         `gorm:"type:uuid;not null;uniqueIndex" json:"authCenterUserId"`
	Nickname         string         `gorm:"type:varchar(50)" json:"nickname"`
	AvatarURL        string         `gorm:"type:varchar(500)" json:"avatarUrl"`
	Profile          Profile        `gorm:"type:jsonb;default:'{}'" json:"profile"`
	Roles            Roles          `gorm:"type:jsonb;default:'[]'" json:"roles"`
	ActiveRole       string         `gorm:"column:active_role;type:varchar(50)" json:"currentRole"`
	LastUsedRole     string         `gorm:"type:varchar(50)" json:"lastUsedRole"`
	InvitedBy        string         `gorm:"type:varchar(255)" json:"invitedBy"`
	InvitationCodeID string         `gorm:"type:uuid" json:"invitationCodeId"`
	Status           string         `gorm:"type:varchar(20);not null;default:'active';check:status IN ('active', 'banned', 'inactive')" json:"status"`
	LastLoginAt      *time.Time     `json:"lastLoginAt"`
	LastLoginIP      string         `gorm:"type:varchar(50)" json:"lastLoginIp"`
	CreatedAt        time.Time      `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt        time.Time      `gorm:"not null;default:now()" json:"updatedAt"`
	DeletedAt        *time.Time     `json:"deletedAt"`
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
