package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ServiceProvider 服务商模型
type ServiceProvider struct {
	ID          uuid.UUID              `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	AdminID     *string                `gorm:"type:varchar(255);uniqueIndex" json:"adminId"`
	UserID      string                 `gorm:"type:varchar(255);not null;index" json:"userId"`
	Name        string                 `gorm:"type:varchar(100);not null;index" json:"name"`
	Description string                 `gorm:"type:text" json:"description"`
	Status      string                 `gorm:"type:varchar(20);not null;default:'active';check:status IN ('active', 'suspended', 'inactive');index" json:"status"`
	CreatedAt   time.Time              `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt   time.Time              `gorm:"not null;default:now()" json:"updatedAt"`
	DeletedAt   *time.Time             `json:"deletedAt"`

	// 关联
	Admin   *User                    `gorm:"foreignKey:AdminID" json:"admin,omitempty"`
	User    *User                    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Staff   []ServiceProviderStaff   `gorm:"foreignKey:ProviderID" json:"staff,omitempty"`
	Merchants []Merchant             `gorm:"foreignKey:ProviderID" json:"merchants,omitempty"`
}

// TableName 指定表名
func (ServiceProvider) TableName() string {
	return "service_providers"
}

// BeforeCreate GORM Hook
func (sp *ServiceProvider) BeforeCreate(tx *gorm.DB) error {
	if sp.ID == uuid.Nil {
		sp.ID = uuid.New()
	}
	return nil
}

// ServiceProviderStaff 服务商员工模型
type ServiceProviderStaff struct {
	ID         uuid.UUID                    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	UserID     string                       `gorm:"type:varchar(255);not null;uniqueIndex" json:"userId"`
	ProviderID uuid.UUID                    `gorm:"type:uuid;not null;index" json:"providerId"`
	Title      string                       `gorm:"type:varchar(50)" json:"title"`
	Status     string                       `gorm:"type:varchar(20);not null;default:'active';check:status IN ('active', 'inactive');index" json:"status"`
	CreatedAt  time.Time                    `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt  time.Time                    `gorm:"not null;default:now()" json:"updatedAt"`

	// 关联
	User        *User                          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Provider    *ServiceProvider               `gorm:"foreignKey:ProviderID" json:"provider,omitempty"`
	Permissions []ServiceProviderStaffPermission `gorm:"foreignKey:StaffID" json:"permissions,omitempty"`
}

// TableName 指定表名
func (ServiceProviderStaff) TableName() string {
	return "service_provider_staff"
}

// BeforeCreate GORM Hook
func (sps *ServiceProviderStaff) BeforeCreate(tx *gorm.DB) error {
	if sps.ID == uuid.Nil {
		sps.ID = uuid.New()
	}
	return nil
}

// ServiceProviderStaffPermission 服务商员工权限模型
type ServiceProviderStaffPermission struct {
	ID             uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	StaffID        uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_provider_staff_permission" json:"staffId"`
	PermissionCode string    `gorm:"type:varchar(50);not null;uniqueIndex:idx_provider_staff_permission;index" json:"permissionCode"`
	GrantedAt      time.Time `gorm:"not null;default:now()" json:"grantedAt"`
	GrantedBy      uuid.UUID `gorm:"type:uuid;not null" json:"grantedBy"`

	// 关联
	Staff   *ServiceProviderStaff `gorm:"foreignKey:StaffID" json:"staff,omitempty"`
	Grantor *ServiceProviderStaff `gorm:"foreignKey:GrantedBy" json:"grantor,omitempty"`
}

// TableName 指定表名
func (ServiceProviderStaffPermission) TableName() string {
	return "provider_staff_permissions"
}

// BeforeCreate GORM Hook
func (pspp *ServiceProviderStaffPermission) BeforeCreate(tx *gorm.DB) error {
	if pspp.ID == uuid.Nil {
		pspp.ID = uuid.New()
	}
	return nil
}
