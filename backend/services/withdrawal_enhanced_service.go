package services

import (
	"errors"
	"fmt"
	"time"
	"pr-business/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateWithdrawalRequest 创建提现请求的输入参数
type CreateWithdrawalRequest struct {
	UserID     string
	Amount      int      // 单位：积分
	YuanAmount  float64  // 单位：元
	Description string
}

// WithdrawalEnhancedService 增强提现服务（完整冻结+解冻机制）
type WithdrawalEnhancedService struct {
	db                *gorm.DB
	validatorService  *ValidatorService
	cashAccountService *CashAccountService
	systemAccountService *SystemAccountService
}

func NewWithdrawalEnhancedService(
	db *gorm.DB,
	validatorService *ValidatorService,
	cashAccountService *CashAccountService,
	systemAccountService *SystemAccountService,
) *WithdrawalEnhancedService {
	return &WithdrawalEnhancedService{
		db:                 db,
		validatorService:   validatorService,
		cashAccountService: cashAccountService,
		systemAccountService: systemAccountService,
	}
}

// CreateWithdrawalRequest 创建提现申请（带冻结）
func (s *WithdrawalEnhancedService) CreateWithdrawalRequest(
	req *CreateWithdrawalRequest,
	requestorID string,
) (*models.WithdrawalRequest, error) {
	// 1. 获取用户积分账户
	creditAccount, err := s.getCreditAccountForWithdrawal(req.UserID)
	if err != nil {
		return nil, err
	}

	// 2. 验证余额
	if creditAccount.Balance < req.Amount {
		return nil, fmt.Errorf("积分余额不足。当前: %d, 需要: %d",
			creditAccount.Balance, req.Amount)
	}

	var createdWithdrawal *models.WithdrawalRequest

	// 3. 开启数据库事务（确保原子性）
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 4. 冻结用户积分到 FrozenBalance
		creditAccount.Balance -= req.Amount
		creditAccount.FrozenBalance += req.Amount

		if err := tx.Save(creditAccount).Error; err != nil {
			return err
		}

		// 5. 创建提现申请记录
		withdrawal := models.WithdrawalRequest{
			ID:          uuid.New(),
			AccountID:   creditAccount.ID,
			Amount:      req.Amount,
			YuanAmount:  req.YuanAmount,
			Status:      models.WithdrawalStatusPending,
			Description: req.Description,
		}

		if err := tx.Create(&withdrawal).Error; err != nil {
			return err
		}

		createdWithdrawal = &withdrawal
		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdWithdrawal, nil
}

// ApproveWithdrawalRequest 审核通过（解冻+打款）
func (s *WithdrawalEnhancedService) ApproveWithdrawalRequest(
	id string,
	reviewerID string,
	cashAccountType string,
) (*models.WithdrawalRequest, error) {
	var result *models.WithdrawalRequest

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 1. 获取提现申请
		var req models.WithdrawalRequest
		if err := tx.Where("id = ?", id).First(&req).Error; err != nil {
			return errors.New("提现申请不存在")
		}

		// 2. 检查状态
		if req.Status != models.WithdrawalStatusPending {
			return errors.New("提现申请状态不正确")
		}

		// 3. 获取用户积分账户
		var creditAccount models.CreditAccount
		if err := tx.Where("id = ?", req.AccountID).First(&creditAccount).Error; err != nil {
			return err
		}

		// 4. 验证冻结余额是否足够
		if creditAccount.FrozenBalance < req.Amount {
			return errors.New("冻结积分不足")
		}

		// 5. 解冻并扣除积分（FrozenBalance -> 0，积分被消耗）
		creditAccount.FrozenBalance -= req.Amount
		// 注意：这里不减Balance，因为积分已被消耗（兑换为现金）

		if err := tx.Save(&creditAccount).Error; err != nil {
			return err
		}

		// 6. 根据提现类型获取现金账户
		var cashAccount models.CashAccount
		if err := tx.Where("account_type = ? AND is_active = ?", cashAccountType, true).First(&cashAccount).Error; err != nil {
			return fmt.Errorf("现金账户不存在: %s", cashAccountType)
		}

		// 7. 验证现金余额（将元转为分）
		requiredAmount := int(req.YuanAmount * 100)
		if cashAccount.Balance < requiredAmount {
			return fmt.Errorf("现金账户余额不足。当前: ¥%.2f, 需要: ¥%.2f",
				float64(cashAccount.Balance)/100, req.YuanAmount)
		}

		// 8. 扣除现金
		cashAccount.Balance -= requiredAmount
		cashAccount.UpdatedAt = time.Now()

		if err := tx.Save(&cashAccount).Error; err != nil {
			return err
		}

		// 9. 更新提现状态为已完成
		now := time.Now()
		req.Status = models.WithdrawalStatusCompleted
		req.ReviewedBy = reviewerID
		req.ReviewedAt = &now
		req.CompletedAt = &now
		req.CashAccountType = cashAccountType

		if err := tx.Save(&req).Error; err != nil {
			return err
		}

		result = &req
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// RejectWithdrawalRequest 审核拒绝（退还积分）
func (s *WithdrawalEnhancedService) RejectWithdrawalRequest(
	id string,
	reason string,
	reviewerID string,
) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 1. 获取提现申请
		var req models.WithdrawalRequest
		if err := tx.Where("id = ?", id).First(&req).Error; err != nil {
			return errors.New("提现申请不存在")
		}

		// 2. 检查状态
		if req.Status != models.WithdrawalStatusPending {
			return errors.New("提现申请状态不正确")
		}

		// 3. 获取用户积分账户
		var creditAccount models.CreditAccount
		if err := tx.Where("id = ?", req.AccountID).First(&creditAccount).Error; err != nil {
			return err
		}

		// 4. 验证冻结余额
		if creditAccount.FrozenBalance < req.Amount {
			return errors.New("冻结积分不足")
		}

		// 5. 解冻并退还积分（FrozenBalance -> Balance）
		creditAccount.FrozenBalance -= req.Amount
		creditAccount.Balance += req.Amount

		if err := tx.Save(&creditAccount).Error; err != nil {
			return err
		}

		// 6. 更新状态为拒绝
		now := time.Now()
		req.Status = models.WithdrawalStatusRejected
		req.RejectReason = &reason
		req.ReviewedBy = reviewerID
		req.ReviewedAt = &now

		if err := tx.Save(&req).Error; err != nil {
			return err
		}

		return nil
	})
}

// 辅助方法

// getCreditAccountForWithdrawal 获取用户积分账户
func (s *WithdrawalEnhancedService) getCreditAccountForWithdrawal(userID string) (*models.CreditAccount, error) {
	var account models.CreditAccount
	err := s.db.Where("owner_id = ?", userID).First(&account).Error
	if err != nil {
		return nil, fmt.Errorf("获取积分账户失败: %w", err)
	}
	return &account, nil
}
