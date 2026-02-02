package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CampaignStatus 营销活动状态
type CampaignStatus string

const (
	CampaignStatusDraft  CampaignStatus = "DRAFT"  // 草稿
	CampaignStatusOpen   CampaignStatus = "OPEN"   // 开放中
	CampaignStatusClosed CampaignStatus = "CLOSED" // 已关闭
)

// Campaign 营销活动模型
type Campaign struct {
	ID                  uuid.UUID      `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	MerchantID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"merchantId"`
	ProviderID          *uuid.UUID     `gorm:"type:uuid;index" json:"providerId"`
	Title               string         `gorm:"type:varchar(100);not null" json:"title"`
	Requirements        string         `gorm:"type:text;not null" json:"requirements"`
	Platforms           string         `gorm:"type:jsonb;not null" json:"platforms"`
	TaskAmount          int            `gorm:"type:int;not null" json:"taskAmount"`
	CampaignAmount      int            `gorm:"type:int;not null" json:"campaignAmount"`
	CreatorAmount       *int           `json:"creatorAmount"`
	StaffReferralAmount *int           `json:"staffReferralAmount"`
	ProviderAmount      *int           `json:"providerAmount"`
	Quota               int            `gorm:"type:int;not null" json:"quota"`
	TaskDeadline        time.Time      `gorm:"type:timestamp;not null" json:"taskDeadline"`
	SubmissionDeadline  time.Time      `gorm:"type:timestamp;not null" json:"submissionDeadline"`
	Status              CampaignStatus `gorm:"type:varchar(20);not null;default:'DRAFT';index" json:"status"`
	CreatedAt           time.Time      `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt           time.Time      `gorm:"not null;default:now()" json:"updatedAt"`
	DeletedAt           *time.Time     `json:"deletedAt"`

	// 关联
	Merchant *Merchant    `gorm:"foreignKey:MerchantID" json:"merchant,omitempty"`
	Provider *ServiceProvider `gorm:"foreignKey:ProviderID" json:"provider,omitempty"`
	Tasks    []Task       `gorm:"foreignKey:CampaignID" json:"tasks,omitempty"`
}

// TableName 指定表名
func (Campaign) TableName() string {
	return "campaigns"
}

// BeforeCreate GORM Hook
func (c *Campaign) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// TaskStatus 任务状态
type TaskStatus string

const (
	TaskStatusOpen      TaskStatus = "OPEN"      // 开放中
	TaskStatusAssigned  TaskStatus = "ASSIGNED"  // 已分配
	TaskStatusSubmitted TaskStatus = "SUBMITTED" // 已提交
	TaskStatusApproved  TaskStatus = "APPROVED"  // 已通过
	TaskStatusRejected  TaskStatus = "REJECTED"  // 已拒绝
)

// TaskPriority 任务优先级
type TaskPriority string

const (
	TaskPriorityHigh   TaskPriority = "HIGH"
	TaskPriorityMedium TaskPriority = "MEDIUM"
	TaskPriorityLow    TaskPriority = "LOW"
)

// Task 任务名额模型
type Task struct {
	ID               uuid.UUID    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	CampaignID       uuid.UUID    `gorm:"type:uuid;not null;uniqueIndex:idx_campaign_slot;index" json:"campaignId"`
	TaskSlotNumber   int          `gorm:"type:int;not null;uniqueIndex:idx_campaign_slot" json:"taskSlotNumber"`
	Status           TaskStatus   `gorm:"type:varchar(20);not null;default:'OPEN';index" json:"status"`
	CreatorID        *uuid.UUID   `gorm:"type:uuid;uniqueIndex:idx_campaign_creator;index" json:"creatorId"`
	AssignedAt       *time.Time   `json:"assignedAt"`
	Platform         string       `gorm:"type:varchar(50)" json:"platform"`
	PlatformURL      string       `gorm:"type:varchar(500)" json:"platformUrl"`
	Screenshots      string       `gorm:"type:jsonb" json:"screenshots"`
	SubmittedAt      *time.Time   `json:"submittedAt"`
	Notes            string       `gorm:"type:text" json:"notes"`
	AuditedBy        *uuid.UUID   `gorm:"type:uuid;index" json:"auditedBy"`
	AuditedAt        *time.Time   `json:"auditedAt"`
	AuditNote        string       `gorm:"type:text" json:"auditNote"`
	InviterID        *uuid.UUID   `gorm:"type:uuid;index" json:"inviterId"`
	InviterType      string       `gorm:"type:varchar(50)" json:"inviterType"`
	Priority         TaskPriority `gorm:"type:varchar(10);not null;default:'MEDIUM'" json:"priority"`
	Tags             string       `gorm:"type:text[]" json:"tags"`
	Version          int          `gorm:"type:int;not null;default:0" json:"version"`
	CreatedAt        time.Time    `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt        time.Time    `gorm:"not null;default:now()" json:"updatedAt"`

	// 关联
	Campaign *Campaign `gorm:"foreignKey:CampaignID" json:"campaign,omitempty"`
	Creator  *Creator  `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
	Auditor  *User     `gorm:"foreignKey:AuditedBy" json:"auditor,omitempty"`
	Inviter  *User     `gorm:"foreignKey:InviterID" json:"inviter,omitempty"`
}

// TableName 指定表名
func (Task) TableName() string {
	return "tasks"
}

// BeforeCreate GORM Hook
func (t *Task) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}
