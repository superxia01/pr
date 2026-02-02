package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Creator 达人模型
type Creator struct {
	ID                       uuid.UUID    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	UserID                   uuid.UUID    `gorm:"type:uuid;not null;index;uniqueIndex:idx_user_primary" json:"userId"`
	IsPrimary                bool         `gorm:"type:boolean;not null;default:true;uniqueIndex:idx_user_primary;index" json:"isPrimary"`
	Level                    string       `gorm:"type:varchar(20);not null;default:'UGC';check:level IN ('UGC', 'KOC', 'INF', 'KOL');index" json:"level"`
	FollowersCount           int          `gorm:"type:int;not null;default:0" json:"followersCount"`
	WechatOpenID             string       `gorm:"type:varchar(100)" json:"wechatOpenId"`
	WechatNickname           string       `gorm:"type:varchar(100)" json:"wechatNickname"`
	WechatAvatar             string       `gorm:"type:varchar(500)" json:"wechatAvatar"`
	InviterID                *uuid.UUID   `gorm:"type:uuid;index" json:"inviterId"`
	InviterType              string       `gorm:"type:varchar(50);check:inviter_type IS NULL OR inviter_type IN ('PROVIDER_STAFF', 'PROVIDER_ADMIN', 'OTHER')" json:"inviterType"`
	InviterRelationshipBroken bool         `gorm:"type:boolean;not null;default:false" json:"inviterRelationshipBroken"`
	Status                   string       `gorm:"type:varchar(20);not null;default:'active';check:status IN ('active', 'banned', 'inactive');index" json:"status"`
	CreatedAt                time.Time    `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt                time.Time    `gorm:"not null;default:now()" json:"updatedAt"`
	DeletedAt                *time.Time   `json:"deletedAt"`

	// 关联
	User    *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Inviter *User     `gorm:"foreignKey:InviterID" json:"inviter,omitempty"`
}

// TableName 指定表名
func (Creator) TableName() string {
	return "creators"
}

// BeforeCreate GORM Hook
func (c *Creator) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// CreatorLevel 达人等级枚举
type CreatorLevel string

const (
	CreatorLevelUGC CreatorLevel = "UGC" // 普通用户
	CreatorLevelKOC CreatorLevel = "KOC" // 关键意见消费者
	CreatorLevelINF CreatorLevel = "INF" // 达人
	CreatorLevelKOL CreatorLevel = "KOL" // 关键意见领袖
)

// CreatorStatus 达人状态枚举
type CreatorStatus string

const (
	CreatorStatusActive   CreatorStatus = "active"
	CreatorStatusBanned   CreatorStatus = "banned"
	CreatorStatusInactive CreatorStatus = "inactive"
)
