package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Merchant 商家模型
type Merchant struct {
	ID          uuid.UUID      `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	AdminID     string         `gorm:"type:varchar(255);uniqueIndex" json:"adminId"`
	ProviderID  uuid.UUID      `gorm:"type:uuid;not null;index" json:"providerId"`
	UserID      string         `gorm:"type:varchar(255);not null;index" json:"userId"`
	Name        string         `gorm:"type:varchar(100);not null;index" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Industry    string         `gorm:"type:varchar(50)" json:"industry"`
	Status      string         `gorm:"type:varchar(20);not null;default:'active';check:status IN ('active', 'suspended', 'inactive');index" json:"status"`
	CreatedAt   time.Time      `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"not null;default:now()" json:"updatedAt"`
	DeletedAt   *time.Time     `json:"deletedAt"`

	// 关联
	Admin       *User          `gorm:"foreignKey:AdminID" json:"admin,omitempty"`
	Provider    *ServiceProvider `gorm:"foreignKey:ProviderID" json:"provider,omitempty"`
	User        *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Staff       []MerchantStaff `gorm:"foreignKey:MerchantID" json:"staff,omitempty"`
}

// TableName 指定表名
func (Merchant) TableName() string {
	return "merchants"
}

// BeforeCreate GORM Hook
func (m *Merchant) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

// MerchantStaff 商家员工模型
type MerchantStaff struct {
	ID         uuid.UUID    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	UserID     string       `gorm:"type:varchar(255);not null;uniqueIndex" json:"userId"`
	MerchantID uuid.UUID    `gorm:"type:uuid;not null;index" json:"merchantId"`
	Title      string       `gorm:"type:varchar(50)" json:"title"`
	Status     string       `gorm:"type:varchar(20);not null;default:'active';check:status IN ('active', 'inactive');index" json:"status"`
	CreatedAt  time.Time    `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt  time.Time    `gorm:"not null;default:now()" json:"updatedAt"`

	// 关联
	User      *User                       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Merchant  *Merchant                   `gorm:"foreignKey:MerchantID" json:"merchant,omitempty"`
	Permissions []MerchantStaffPermission `gorm:"foreignKey:StaffID" json:"permissions,omitempty"`
}

// TableName 指定表名
func (MerchantStaff) TableName() string {
	return "merchant_staff"
}

// BeforeCreate GORM Hook
func (ms *MerchantStaff) BeforeCreate(tx *gorm.DB) error {
	if ms.ID == uuid.Nil {
		ms.ID = uuid.New()
	}
	return nil
}

// MerchantStaffPermission 商家员工权限模型
type MerchantStaffPermission struct {
	ID             uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	StaffID        uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_staff_permission" json:"staffId"`
	PermissionCode string    `gorm:"type:varchar(50);not null;uniqueIndex:idx_staff_permission;index" json:"permissionCode"`
	GrantedAt      time.Time `gorm:"not null;default:now()" json:"grantedAt"`
	GrantedBy      uuid.UUID `gorm:"type:uuid;not null" json:"grantedBy"`

	// 关联
	Staff *MerchantStaff `gorm:"foreignKey:StaffID" json:"staff,omitempty"`
	Grantor *MerchantStaff `gorm:"foreignKey:GrantedBy" json:"grantor,omitempty"`
}

// TableName 指定表名
func (MerchantStaffPermission) TableName() string {
	return "merchant_staff_permissions"
}

// BeforeCreate GORM Hook
func (msp *MerchantStaffPermission) BeforeCreate(tx *gorm.DB) error {
	if msp.ID == uuid.Nil {
		msp.ID = uuid.New()
	}
	return nil
}
