package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WithdrawalStatus 提现状态
type WithdrawalStatus string

const (
	WithdrawalStatusPending   WithdrawalStatus = "pending"   // 待审核
	WithdrawalStatusApproved  WithdrawalStatus = "approved"  // 已通过
	WithdrawalStatusRejected  WithdrawalStatus = "rejected"  // 已拒绝
	WithdrawalStatusCompleted WithdrawalStatus = "completed" // 已完成
)

// WithdrawalMethod 提现方式
type WithdrawalMethod string

const (
	WithdrawalMethodAlipay WithdrawalMethod = "ALIPAY" // 支付宝
	WithdrawalMethodWechat WithdrawalMethod = "WECHAT" // 微信
	WithdrawalMethodBank   WithdrawalMethod = "BANK"   // 银行转账
)

// Withdrawal 提现记录
type Withdrawal struct {
	ID              uuid.UUID         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	AccountID       uuid.UUID         `gorm:"type:uuid;not null;index:idx_withdrawals_account_id" json:"accountId"`
	Amount          int               `gorm:"type:int;not null" json:"amount"`                   // 申请金额（积分）
	Fee             int               `gorm:"type:int;not null;default:0" json:"fee"`             // 手续费（积分）
	ActualAmount    int               `gorm:"type:int;not null" json:"actualAmount"`              // 实际到账金额（积分）
	Method          WithdrawalMethod  `gorm:"type:varchar(20);not null" json:"method"`            // 提现方式
	AccountInfo     string            `gorm:"type:jsonb;not null" json:"accountInfo"`             // 提现账户信息（JSON）
	AccountInfoHash string            `gorm:"type:varchar(64)" json:"accountInfoHash"`            // 账户信息哈希（用于去重）
	Status          WithdrawalStatus  `gorm:"type:varchar(20);not null;default:'pending';index:idx_withdrawals_status" json:"status"`
	AuditNote       string            `gorm:"type:varchar(200)" json:"auditNote"`                 // 审核备注
	AuditedBy       *uuid.UUID        `gorm:"type:uuid" json:"auditedBy"`                        // 审核人ID
	AuditedAt       *time.Time        `gorm:"type:timestamp" json:"auditedAt"`                   // 审核时间
	CompletedAt     *time.Time        `gorm:"type:timestamp" json:"completedAt"`                 // 完成时间
	CreatedAt       time.Time         `gorm:"type:timestamp;not null;default:NOW();index:idx_withdrawals_created_at" json:"createdAt"`
	UpdatedAt       time.Time         `gorm:"type:timestamp;not null;default:NOW()" json:"updatedAt"`

	// 关联
	Account CreditAccount `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	Auditor *User         `gorm:"foreigner:AuditedBy;references:ID" json:"auditor,omitempty"`
}

// BeforeCreate GORM hook
func (w *Withdrawal) BeforeCreate(tx *gorm.DB) error {
	return nil
}

// CanAudit 是否可以审核
func (w *Withdrawal) CanAudit() bool {
	return w.Status == WithdrawalStatusPending
}

// CanProcess 是否可以打款
func (w *Withdrawal) CanProcess() bool {
	return w.Status == WithdrawalStatusApproved
}

// IsPending 是否待审核
func (w *Withdrawal) IsPending() bool {
	return w.Status == WithdrawalStatusPending
}

// IsApproved 是否已通过
func (w *Withdrawal) IsApproved() bool {
	return w.Status == WithdrawalStatusApproved
}

// IsRejected 是否已拒绝
func (w *Withdrawal) IsRejected() bool {
	return w.Status == WithdrawalStatusRejected
}

// IsCompleted 是否已完成
func (w *Withdrawal) IsCompleted() bool {
	return w.Status == WithdrawalStatusCompleted
}
