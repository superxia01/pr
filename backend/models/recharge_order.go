package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RechargeOrderStatus 充值订单状态
type RechargeOrderStatus string

const (
	RechargeOrderStatusPending   RechargeOrderStatus = "pending"   // 待审核
	RechargeOrderStatusApproved  RechargeOrderStatus = "approved"  // 已通过
	RechargeOrderStatusRejected  RechargeOrderStatus = "rejected"  // 已拒绝
	RechargeOrderStatusCompleted RechargeOrderStatus = "completed" // 已完成
)

// RechargeOrder 充值订单模型
// 用于线下充值流程：用户提交订单 → 超管审核 → 充值入账
type RechargeOrder struct {
	ID             uuid.UUID            `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	UserID         string               `gorm:"type:varchar(255);not null;index" json:"userId"`
	AccountID      uuid.UUID            `gorm:"type:uuid;not null;index" json:"accountId"`
	Amount         int                  `gorm:"type:int;not null;check:amount > 0" json:"amount"`
	PaymentMethod  string               `gorm:"type:varchar(20);not null" json:"paymentMethod"` // 支付方式：alipay/wechat/bank
	PaymentProof  string               `gorm:"type:varchar(500)" json:"paymentProof"`              // 支付凭证URL
	Status        RechargeOrderStatus   `gorm:"type:varchar(20);not null;default:'pending';check:status IN ('pending', 'approved', 'rejected', 'completed')" json:"status"`
	RejectionNote string               `gorm:"type:text" json:"rejectionNote"`                 // 拒绝原因
	AuditedBy     *string              `gorm:"type:varchar(255)" json:"auditedBy"`          // 审核人ID
	AuditedAt      *time.Time           `json:"auditedAt"`                                    // 审核时间
	ProcessedAt    *time.Time           `json:"processedAt"`                                   // 完成时间
	CreatedAt      time.Time            `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt      time.Time            `gorm:"not null;default:now()" json:"updatedAt"`

	// 关联
	Account *CreditAccount `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	User    *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Auditor *User          `gorm:"foreignKey:AuditedBy" json:"auditor,omitempty"`
}

// TableName 指定表名
func (RechargeOrder) TableName() string {
	return "recharge_orders"
}

// BeforeCreate GORM Hook
func (ro *RechargeOrder) BeforeCreate(tx *gorm.DB) error {
	if ro.ID == uuid.Nil {
		ro.ID = uuid.New()
	}
	return nil
}

// CanAudit 检查是否可以审核
func (ro *RechargeOrder) CanAudit() bool {
	return ro.Status == RechargeOrderStatusPending
}

// CanProcess 检查是否可以处理完成
func (ro *RechargeOrder) CanProcess() bool {
	return ro.Status == RechargeOrderStatusApproved
}
