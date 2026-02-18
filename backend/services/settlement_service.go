package services

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"pr-business/models"
)

// SettlementService 结算服务
// 处理任务审核通过后的积分分配
type SettlementService struct {
	db                *gorm.DB
	permissionService  *AccountPermissionService
	validatorService   *ValidatorService
	cashAccountService *CashAccountService
}

// NewSettlementService 创建结算服务
func NewSettlementService(
	db *gorm.DB,
	permissionService *AccountPermissionService,
	validatorService *ValidatorService,
	cashAccountService *CashAccountService,
) *SettlementService {
	return &SettlementService{
		db:                db,
		permissionService:  permissionService,
		validatorService:   validatorService,
		cashAccountService: cashAccountService,
	}
}

// SettleTaskAfterApproval 任务审核通过后结算
// 流程：
// 1. 从商家账户扣除 campaign_amount（活动总金额）
// 2. 给达人账户增加 creator_amount（达人收入）
// 3. 给员工增加 staff_referral_amount（员工返佣）
// 4. 给服务商增加 provider_amount（服务商分成）
func (s *SettlementService) SettleTaskAfterApproval(task *models.Task, auditorUserID string) error {
	// 获取营销活动信息
	var campaign models.Campaign
	if err := s.db.Where("id = ?", task.CampaignID).First(&campaign).Error; err != nil {
		return fmt.Errorf("获取营销活动失败: %w", err)
	}

	// 验证金额配置
	if campaign.CreatorAmount == nil || *campaign.CreatorAmount <= 0 {
		return errors.New("达人收入金额未配置或为0")
	}

	// 开始事务
	return s.db.Transaction(func(tx *gorm.DB) error {
		transactionGroupID := uuid.New()

		// 1. 从商家冻结账户扣除活动总金额
		merchantAccount, err := s.findOrCreateAccount(tx, campaign.MerchantID, models.OwnerTypeOrgMerchant)
		if err != nil {
			return fmt.Errorf("获取商家账户失败: %w", err)
		}

		if merchantAccount.FrozenBalance < campaign.CampaignAmount {
			return errors.New("商家冻结积分不足")
		}

		// 从冻结余额扣除
		merchantAccount.FrozenBalance -= campaign.CampaignAmount
		if err := tx.Save(&merchantAccount).Error; err != nil {
			return fmt.Errorf("更新商家冻结余额失败: %w", err)
		}

		// 记录商家扣款流水
		taskPublishTransaction := models.CreditTransaction{
			AccountID:        merchantAccount.ID,
			Type:             models.TransactionTaskPublish,
			Amount:           -campaign.CampaignAmount,
			BalanceBefore:     merchantAccount.FrozenBalance + campaign.CampaignAmount,
			BalanceAfter:      merchantAccount.FrozenBalance,
			RelatedCampaignID: &campaign.ID,
			RelatedTaskID:     &task.ID,
			Description:       fmt.Sprintf("任务结算：%s", campaign.Title),
			TransactionGroupID: &transactionGroupID,
			GroupSequence:     intPtr(1),
		}
		if err := tx.Create(&taskPublishTransaction).Error; err != nil {
			return fmt.Errorf("记录商家流水失败: %w", err)
		}

		// 2. 给达人账户增加收入
		creatorUserID, err := s.getCreatorUserID(tx, task.CreatorID)
		if err != nil {
			return fmt.Errorf("获取达人用户ID失败: %w", err)
		}

		creatorAccount, err := s.findOrCreateAccountByUserID(tx, creatorUserID, models.OwnerTypeUserPersonal)
		if err != nil {
			return fmt.Errorf("获取达人账户失败: %w", err)
		}

		creatorBalanceBefore := creatorAccount.Balance
		creatorAccount.Balance += *campaign.CreatorAmount
		if err := tx.Save(&creatorAccount).Error; err != nil {
			return fmt.Errorf("更新达人余额失败: %w", err)
		}

		// 记录达人收入流水
		taskIncomeTransaction := models.CreditTransaction{
			AccountID:        creatorAccount.ID,
			Type:             models.TransactionTaskIncome,
			Amount:           *campaign.CreatorAmount,
			BalanceBefore:     creatorBalanceBefore,
			BalanceAfter:      creatorAccount.Balance,
			RelatedCampaignID: &campaign.ID,
			RelatedTaskID:     &task.ID,
			Description:       fmt.Sprintf("任务收入：%s", campaign.Title),
			TransactionGroupID: &transactionGroupID,
			GroupSequence:     intPtr(2),
		}
		if err := tx.Create(&taskIncomeTransaction).Error; err != nil {
			return fmt.Errorf("记录达人流水失败: %w", err)
		}

		// 3. 员工返佣（如果配置了且任务有邀请人）
		if campaign.StaffReferralAmount != nil && *campaign.StaffReferralAmount > 0 &&
			task.InviterID != nil && *task.InviterID != "" {
			// 查找员工账户
			inviterAccount, err := s.findInviterAccount(tx, *task.InviterID, task.InviterType)
			if err == nil {
				staffBalanceBefore := inviterAccount.Balance
				inviterAccount.Balance += *campaign.StaffReferralAmount
				if err := tx.Save(&inviterAccount).Error; err != nil {
					return fmt.Errorf("更新员工余额失败: %w", err)
				}

				// 记录员工返佣流水
				staffReferralTransaction := models.CreditTransaction{
					AccountID:        inviterAccount.ID,
					Type:             models.TransactionStaffReferral,
					Amount:           *campaign.StaffReferralAmount,
					BalanceBefore:     staffBalanceBefore,
					BalanceAfter:      inviterAccount.Balance,
					RelatedCampaignID: &campaign.ID,
					RelatedTaskID:     &task.ID,
					Description:       fmt.Sprintf("员工返佣：%s", campaign.Title),
					TransactionGroupID: &transactionGroupID,
					GroupSequence:     intPtr(3),
				}
				if err := tx.Create(&staffReferralTransaction).Error; err != nil {
					return fmt.Errorf("记录员工流水失败: %w", err)
				}
			}
			// 如果找不到员工账户，忽略返佣
		}

		// 4. 服务商分成（如果配置了且有服务商）
		if campaign.ProviderAmount != nil && *campaign.ProviderAmount > 0 &&
			campaign.ProviderID != nil {
			providerAccount, err := s.findOrCreateAccount(tx, *campaign.ProviderID, models.OwnerTypeOrgProvider)
			if err != nil {
				return fmt.Errorf("获取服务商账户失败: %w", err)
			}

			providerBalanceBefore := providerAccount.Balance
			providerAccount.Balance += *campaign.ProviderAmount
			if err := tx.Save(&providerAccount).Error; err != nil {
				return fmt.Errorf("更新服务商余额失败: %w", err)
			}

			// 记录服务商收入流水
			providerIncomeTransaction := models.CreditTransaction{
				AccountID:        providerAccount.ID,
				Type:             models.TransactionProviderIncome,
				Amount:           *campaign.ProviderAmount,
				BalanceBefore:     providerBalanceBefore,
				BalanceAfter:      providerAccount.Balance,
				RelatedCampaignID: &campaign.ID,
				RelatedTaskID:     &task.ID,
				Description:       fmt.Sprintf("服务商分成：%s", campaign.Title),
				TransactionGroupID: &transactionGroupID,
				GroupSequence:     intPtr(4),
			}
			if err := tx.Create(&providerIncomeTransaction).Error; err != nil {
				return fmt.Errorf("记录服务商流水失败: %w", err)
			}
		}

		return nil
	})
}

// findOrCreateAccount 查找或创建组织账户
func (s *SettlementService) findOrCreateAccount(tx *gorm.DB, ownerID uuid.UUID, ownerType models.OwnerType) (*models.CreditAccount, error) {
	var account models.CreditAccount
	err := tx.Where("owner_id = ? AND owner_type = ?", ownerID, ownerType).First(&account).Error
	if err == nil {
		return &account, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		account := models.CreditAccount{
			OwnerID:       ownerID,
			OwnerType:     ownerType,
			Balance:       0,
			FrozenBalance: 0,
		}
		if err := tx.Create(&account).Error; err != nil {
			return nil, err
		}
		return &account, nil
	}

	return nil, err
}

// findOrCreateAccountByUserID 通过用户ID查找或创建个人账户
func (s *SettlementService) findOrCreateAccountByUserID(tx *gorm.DB, userID string, ownerType models.OwnerType) (*models.CreditAccount, error) {
	// 将用户ID转换为UUID
	parsedID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("用户ID格式错误: %w", err)
	}

	return s.findOrCreateAccount(tx, parsedID, ownerType)
}

// getCreatorUserID 获取达人对应的用户ID
func (s *SettlementService) getCreatorUserID(tx *gorm.DB, creatorID *uuid.UUID) (string, error) {
	if creatorID == nil {
		return "", errors.New("达人ID为空")
	}

	var creator models.Creator
	if err := tx.Where("id = ?", creatorID).First(&creator).Error; err != nil {
		return "", err
	}

	return creator.UserID, nil
}

// findInviterAccount 查找邀请人账户
func (s *SettlementService) findInviterAccount(tx *gorm.DB, inviterID string, inviterType string) (*models.CreditAccount, error) {
	// 根据邀请人类型确定账户
	switch inviterType {
	case "SERVICE_PROVIDER_STAFF", "SERVICE_PROVIDER_ADMIN":
		// 服务商员工/管理员：查找服务商账户
		var providerStaff models.ServiceProviderStaff
		if err := tx.Where("user_id = ?", inviterID).First(&providerStaff).Error; err != nil {
			return nil, fmt.Errorf("找不到服务商员工记录: %w", err)
		}
		return s.findOrCreateAccount(tx, providerStaff.ProviderID, models.OwnerTypeOrgProvider)

	case "MERCHANT_STAFF", "MERCHANT_ADMIN":
		// 商家员工/管理员：查找商家账户
		var merchantStaff models.MerchantStaff
		if err := tx.Where("user_id = ?", inviterID).First(&merchantStaff).Error; err != nil {
			return nil, fmt.Errorf("找不到商家员工记录: %w", err)
		}
		return s.findOrCreateAccount(tx, merchantStaff.MerchantID, models.OwnerTypeOrgMerchant)

	default:
		// 其他：返回个人账户
		return s.findOrCreateAccountByUserID(tx, inviterID, models.OwnerTypeUserPersonal)
	}
}

// intPtr 返回int指针
func intPtr(i int) *int {
	return &i
}

// SettleCampaignAfterClose 活动关闭后结算，解冻未完成任务的积分
// 流程：
// 1. 检查活动状态是否为 CLOSED
// 2. 统计已完成和未完成的任务数量
// 3. 计算应退还的积分 = 未完成任务数量 * 任务金额
// 4. 将商家冻结余额转回可用余额
func (s *SettlementService) SettleCampaignAfterClose(campaign *models.Campaign) error {
	// 开始事务
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 统计任务完成情况
		var completedTasksCount int64
		if err := tx.Model(&models.Task{}).
			Where("campaign_id = ? AND status = ?", campaign.ID, models.TaskStatusApproved).
			Count(&completedTasksCount).Error; err != nil {
			return fmt.Errorf("统计已完成任务失败: %w", err)
		}

		// 计算未完成任务数量
		totalTasks := int64(campaign.Quota)
		uncompletedTasks := totalTasks - completedTasksCount

		// 如果没有未完成的任务，无需处理
		if uncompletedTasks <= 0 {
			return nil
		}

		// 计算应退还的积分
		refundAmount := uncompletedTasks * int64(campaign.TaskAmount)

		// 获取商家积分账户
		merchantAccount, err := s.findOrCreateAccount(tx, campaign.MerchantID, models.OwnerTypeOrgMerchant)
		if err != nil {
			return fmt.Errorf("获取商家账户失败: %w", err)
		}

		// 验证冻结余额足够
		if merchantAccount.FrozenBalance < int(refundAmount) {
			return errors.New("商家冻结积分不足")
		}

		// 解冻积分：从冻结余额转回可用余额
		merchantAccount.FrozenBalance -= int(refundAmount)
		merchantAccount.Balance += int(refundAmount)

		// 记录解冻流水
		refundTransaction := models.CreditTransaction{
			AccountID:      merchantAccount.ID,
			Type:           "CAMPAIGN_REFUND",
			Amount:         int(refundAmount),
			BalanceBefore:  merchantAccount.Balance - int(refundAmount),
			BalanceAfter:   merchantAccount.Balance,
			RelatedCampaignID: &campaign.ID,
			Description:    fmt.Sprintf("活动关闭退还积分：%s（未完成任务 %d/%d）", campaign.Title, uncompletedTasks, totalTasks),
		}
		if err := tx.Create(&refundTransaction).Error; err != nil {
			return fmt.Errorf("记录解冻流水失败: %w", err)
		}

		// 保存积分账户更新
		if err := tx.Save(&merchantAccount).Error; err != nil {
			return fmt.Errorf("更新积分账户失败: %w", err)
		}

		return nil
	})
}
