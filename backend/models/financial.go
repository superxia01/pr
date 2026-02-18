package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// JSONMap JSON map类型，支持 JSON 序列化
type JSONMap map[string]interface{}

// Scan 实现 sql.Scanner 接口
func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = map[string]interface{}{}
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

	return json.Unmarshal(bytes, j)
}

// Value 实现 driver.Valuer 接口
func (j JSONMap) Value() (driver.Value, error) {
	if len(j) == 0 {
		return "{}", nil
	}
	return json.Marshal(j)
}

// CashAccount 现金账户（真钱账户）
type CashAccount struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	AccountType string    `gorm:"type:varchar(50);not null" json:"account_type"` // WECHAT, ALIPAY, BANK_TRANSFER, MARKETING, OPERATIONS
	Balance     int       `gorm:"type:int;not null;default:0" json:"balance"`     // 单位：分
	Description string    `gorm:"type:text" json:"description"`
	IsActive    bool      `gorm:"type:boolean;not null;default:true" json:"is_active"`
	Metadata    JSONMap   `gorm:"type:jsonb;default:'{}'" json:"metadata"`
	CreatedAt   time.Time `gorm:"type:timestamp;not null;default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"type:timestamp;not null;default:now()" json:"updated_at"`
}

// TableName 指定表名
func (CashAccount) TableName() string {
	return "cash_accounts"
}

// SystemAccount 系统账户（托管+平台收益）
type SystemAccount struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	AccountType string    `gorm:"type:varchar(50);not null;unique" json:"account_type"` // TICKET_ESCROW, TASK_ESCROW, PLATFORM_REVENUE
	Balance     int       `gorm:"type:int;not null;default:0" json:"balance"`
	Description string    `gorm:"type:text" json:"description"`
	IsActive    bool      `gorm:"type:boolean;not null;default:true" json:"is_active"`
	Metadata    JSONMap   `gorm:"type:jsonb;default:'{}'" json:"metadata"`
	CreatedAt   time.Time `gorm:"type:timestamp;not null;default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"type:timestamp;not null;default:now()" json:"updated_at"`
}

// TableName 指定表名
func (SystemAccount) TableName() string {
	return "system_accounts"
}

// WithdrawalRequest 提现申请（增强版-带冻结机制）
type WithdrawalRequest struct {
	ID              uuid.UUID        `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	AccountID       uuid.UUID        `gorm:"type:uuid;not null" json:"account_id"`
	Amount          int              `gorm:"type:int;not null" json:"amount"`                    // 单位：积分
	YuanAmount      float64          `gorm:"type:decimal(10,2);not null" json:"yuan_amount"` // 单位：元
	Status          WithdrawalStatus `gorm:"type:varchar(50);not null;default:'pending'" json:"status"` // pending, processing, completed, rejected
	CashAccountType string           `gorm:"type:varchar(50)" json:"cash_account_type"`        // 用于审核通过时选择扣款账户
	Description     string           `gorm:"type:text" json:"description"`
	RejectReason    *string          `gorm:"type:text" json:"reject_reason"`
	ReviewedBy      string           `gorm:"type:varchar(255)" json:"reviewed_by"`
	ReviewedAt      *time.Time       `gorm:"type:timestamp" json:"reviewed_at"`
	CompletedAt     *time.Time       `gorm:"type:timestamp" json:"completed_at"`
	CreatedAt       time.Time        `gorm:"type:timestamp;not null;default:now()" json:"created_at"`
}

// TableName 指定表名
func (WithdrawalRequest) TableName() string {
	return "withdrawal_requests_enhanced"
}

// FinancialAuditLog 财务审计日志（简化版）
type FinancialAuditLog struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID       string    `gorm:"type:varchar(255);not null" json:"user_id"`
	Action       string    `gorm:"type:varchar(100);not null" json:"action"`     // WITHDRAWAL_APPROVE, WITHDRAWAL_REJECT, CREDIT_RECHARGE, etc.
	ResourceType string    `gorm:"type:varchar(50);not null" json:"resource_type"` // WITHDRAWAL_REQUEST, PAYMENT_ORDER, CREDIT_ACCOUNT, etc.
	ResourceID   string    `gorm:"type:varchar(255);not null" json:"resource_id"`
	Changes      JSONMap   `gorm:"type:jsonb;default:'{}'" json:"changes"`         // 记录 before/after 状态
	IPAddress    string    `gorm:"type:inet" json:"ip_address"`
	UserAgent    string    `gorm:"type:text" json:"user_agent"`
	CreatedAt    time.Time `gorm:"type:timestamp;not null;default:now()" json:"created_at"`
}

// TableName 指定表名
func (FinancialAuditLog) TableName() string {
	return "financial_audit_logs"
}
