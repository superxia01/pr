package services

import "errors"

// 财务相关错误定义
var (
	// ErrWithdrawalNotFound 提现申请不存在
	ErrWithdrawalNotFound = errors.New("提现申请不存在")

	// ErrInvalidWithdrawalStatus 提现状态不正确
	ErrInvalidWithdrawalStatus = errors.New("提现状态不正确")

	// ErrInsufficientBalance 余额不足
	ErrInsufficientBalance = errors.New("余额不足")

	// ErrInsufficientFrozenBalance 冻结余额不足
	ErrInsufficientFrozenBalance = errors.New("冻结余额不足")

	// ErrCashAccountNotFound 现金账户不存在
	ErrCashAccountNotFound = errors.New("现金账户不存在")

	// ErrCashAccountBalanceInsufficient 现金账户余额不足
	ErrCashAccountBalanceInsufficient = errors.New("现金账户余额不足")

	// ErrCreditAccountNotFound 积分账户不存在
	ErrCreditAccountNotFound = errors.New("积分账户不存在")

	// ErrInvalidAmount 无效金额
	ErrInvalidAmount = errors.New("无效金额")

	// ErrInvalidCreditType 无效积分类型
	ErrInvalidCreditType = errors.New("无效积分类型")
)
