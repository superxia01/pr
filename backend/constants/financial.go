package constants

// 积分类型常量（统一为积分）
const (
	CreditTypeCredit = "CREDIT" // 统一的积分类型
)

// 提现状态常量
const (
	WithdrawalStatusPending    = "PENDING"
	WithdrawalStatusProcessing = "PROCESSING"
	WithdrawalStatusCompleted  = "COMPLETED"
	WithdrawalStatusRejected   = "REJECTED"
)

// 现金账户类型常量
const (
	CashAccountTypeWeChat       = "WECHAT"
	CashAccountTypeAlipay      = "ALIPAY"
	CashAccountTypeBankTransfer = "BANK_TRANSFER"
	CashAccountTypeMarketing    = "MARKETING"
	CashAccountTypeOperations  = "OPERATIONS"
)

// 系统账户类型常量
const (
	SystemAccountTypeTicketEscrow   = "TICKET_ESCROW"
	SystemAccountTypeTaskEscrow     = "TASK_ESCROW"
	SystemAccountTypePlatformRevenue = "PLATFORM_REVENUE" // 统一的平台收益账户
)

// 审计操作类型常量
const (
	AuditActionWithdrawalRequest   = "WITHDRAWAL_REQUEST"
	AuditActionWithdrawalApprove   = "WITHDRAWAL_APPROVE"
	AuditActionWithdrawalReject    = "WITHDRAWAL_REJECT"
	AuditActionCreditRecharge      = "CREDIT_RECHARGE"
	AuditActionCampaignPublish     = "CAMPAIGN_PUBLISH"
	AuditActionCampaignTaskCreate  = "CAMPAIGN_TASK_CREATE"
	AuditActionSystemAdjust       = "SYSTEM_ADJUST"
)

// 审计资源类型常量
const (
	AuditResourceWithdrawalRequest = "WITHDRAWAL_REQUEST"
	AuditResourcePaymentOrder      = "PAYMENT_ORDER"
	AuditResourceCreditAccount     = "CREDIT_ACCOUNT"
	AuditResourceCashAccount       = "CASH_ACCOUNT"
	AuditResourceSystemAccount     = "SYSTEM_ACCOUNT"
	AuditResourceCampaign          = "CAMPAIGN"
)
