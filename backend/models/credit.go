package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OwnerType 账户所有者类型
type OwnerType string

const (
	OwnerTypeOrgMerchant  OwnerType = "ORG_MERCHANT"  // 商家组织
	OwnerTypeOrgProvider  OwnerType = "ORG_PROVIDER"  // 服务商组织
	OwnerTypeUserPersonal OwnerType = "USER_PERSONAL" // 个人用户
)

// CreditAccount 积分账户模型
type CreditAccount struct {
	ID             uuid.UUID  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	OwnerID        uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:idx_owner_type" json:"ownerId"`
	OwnerType      OwnerType  `gorm:"type:varchar(20);not null;uniqueIndex:idx_owner_type;index" json:"ownerType"`
	Balance        int        `gorm:"type:int;not null;default:0;check:balance >= 0" json:"balance"`
	FrozenBalance  int        `gorm:"type:int;not null;default:0;check:frozen_balance >= 0" json:"frozenBalance"`
	CreatedAt      time.Time  `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"not null;default:now()" json:"updatedAt"`
}

// TableName 指定表名
func (CreditAccount) TableName() string {
	return "credit_accounts"
}

// BeforeCreate GORM Hook
func (ca *CreditAccount) BeforeCreate(tx *gorm.DB) error {
	if ca.ID == uuid.Nil {
		ca.ID = uuid.New()
	}
	return nil
}

// TransactionType 积分交易类型模型
type TransactionType struct {
	Code             string        `gorm:"primaryKey;type:varchar(50)" json:"code"`
	Name             string        `gorm:"type:varchar(100);not null" json:"name"`
	Description      string        `gorm:"type:text" json:"description"`
	AccountTypes     string        `gorm:"type:varchar(20)[];not null" json:"accountTypes"` // JSON array
	AmountDirection  string        `gorm:"type:varchar(10);not null;check:amount_direction IN ('positive', 'negative')" json:"amountDirection"`
	IsActive         bool          `gorm:"type:boolean;not null;default:true" json:"isActive"`
	Version          int           `gorm:"type:int;not null;default:1" json:"version"`
	CreatedAt        time.Time     `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt        time.Time     `gorm:"not null;default:now()" json:"updatedAt"`
	DeprecatedAt     *time.Time    `json:"deprecatedAt"`
}

// TableName 指定表名
func (TransactionType) TableName() string {
	return "transaction_types"
}

// CreditTransaction 积分流水模型
type CreditTransaction struct {
	ID                   uuid.UUID  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	AccountID            uuid.UUID  `gorm:"type:uuid;not null;index" json:"accountId"`
	Type                 string     `gorm:"type:varchar(50);not null;index" json:"type"`
	Amount               int        `gorm:"type:int;not null" json:"amount"`
	BalanceBefore        int        `gorm:"type:int;not null" json:"balanceBefore"`
	BalanceAfter         int        `gorm:"type:int;not null" json:"balanceAfter"`
	TransactionGroupID   *uuid.UUID  `gorm:"type:uuid" json:"transactionGroupId"`
	GroupSequence        *int       `json:"groupSequence"`
	RelatedCampaignID    *uuid.UUID  `gorm:"type:uuid" json:"relatedCampaignId"`
	RelatedTaskID        *uuid.UUID  `gorm:"type:uuid" json:"relatedTaskId"`
	Description          string     `gorm:"type:varchar(200)" json:"description"`
	CreatedAt            time.Time  `gorm:"not null;default:now();index" json:"createdAt"`

	// 关联
	Account *CreditAccount `gorm:"foreignKey:AccountID" json:"account,omitempty"`
}

// TableName 指定表名
func (CreditTransaction) TableName() string {
	return "credit_transactions"
}

// BeforeCreate GORM Hook
func (ct *CreditTransaction) BeforeCreate(tx *gorm.DB) error {
	if ct.ID == uuid.Nil {
		ct.ID = uuid.New()
	}
	return nil
}

// TransactionCode 交易代码常量
const (
	TransactionRecharge         = "RECHARGE"          // 商家充值
	TransactionTaskIncome       = "TASK_INCOME"       // 任务收入
	TransactionTaskSubmit       = "TASK_SUBMIT"       // 任务提交
	TransactionStaffReferral    = "STAFF_REFERRAL"    // 员工返佣
	TransactionProviderIncome   = "PROVIDER_INCOME"   // 服务商收入
	TransactionTaskPublish     = "TASK_PUBLISH"      // 发布任务
	TransactionTaskAccept      = "TASK_ACCEPT"       // 接任务扣除
	TransactionTaskReject      = "TASK_REJECT"       // 审核拒绝
	TransactionTaskEscalate    = "TASK_ESCALATE"     // 超时拒绝
	TransactionTaskRefund      = "TASK_REFUND"       // 任务退款
	TransactionWithdraw        = "WITHDRAW"          // 提现
	TransactionWithdrawFreeze  = "WITHDRAW_FREEZE"   // 提现冻结
	TransactionWithdrawRefund  = "WITHDRAW_REFUND"   // 提现拒绝
	TransactionBonusGift       = "BONUS_GIFT"        // 系统赠送
)
