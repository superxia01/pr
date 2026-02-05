package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AccountPermission 账户权限关联模型
type AccountPermission struct {
	ID          uuid.UUID  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	AccountID   uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:idx_account_user" json:"accountId"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:idx_account_user" json:"userId"`
	AccountType string     `gorm:"type:varchar(20);not null" json:"accountType"` // PERSONAL, ORG_MERCHANT, ORG_PROVIDER
	CanView     bool       `gorm:"type:boolean;not null;default:true" json:"canView"`
	CanOperate  bool       `gorm:"type:boolean;not null;default:false" json:"canOperate"`
	CreatedAt   time.Time  `gorm:"not null;default:now()" json:"createdAt"`

	// 关联
	Account *CreditAccount `gorm:"foreignKey:AccountID" json:"account,omitempty"`
}

// TableName 指定表名
func (AccountPermission) TableName() string {
	return "account_permissions"
}

// BeforeCreate GORM Hook
func (ap *AccountPermission) BeforeCreate(tx *gorm.DB) error {
	if ap.ID == uuid.Nil {
		ap.ID = uuid.New()
	}
	return nil
}
