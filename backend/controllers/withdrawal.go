package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"pr-business/models"
	"pr-business/utils"
)

type WithdrawalController struct {
	DB *gorm.DB
}

// CreateWithdrawalRequest 申请提现请求
type CreateWithdrawalRequest struct {
	Amount      int                    `json:"amount" binding:"required,min=1"`
	Method      models.WithdrawalMethod `json:"method" binding:"required,oneof=ALIPAY WECHAT BANK"`
	AccountInfo map[string]interface{} `json:"accountInfo" binding:"required"`
}

// AuditWithdrawalRequest 审核提现请求
type AuditWithdrawalRequest struct {
	Approved  bool   `json:"approved" binding:"required"`
	AuditNote string `json:"auditNote"`
}

// CreateWithdrawal 申请提现
func (ctrl *WithdrawalController) CreateWithdrawal(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	var req CreateWithdrawalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// 1. 确定用户的积分账户类型
	var accountType models.OwnerType
	if utils.HasRole(user, "SUPER_ADMIN") {
		accountType = models.OwnerTypeOrgProvider // 超管默认使用服务商账户
	} else if utils.HasRole(user, "SP_ADMIN") {
		accountType = models.OwnerTypeOrgProvider
	} else if utils.HasRole(user, "MERCHANT_ADMIN") || utils.HasRole(user, "MERCHANT_STAFF") {
		accountType = models.OwnerTypeOrgMerchant
	} else {
		accountType = models.OwnerTypeUserPersonal
	}

	// 2. 获取或创建积分账户
	var account models.CreditAccount
	err := ctrl.DB.Where("owner_id = ? AND owner_type = ?", user.ID, accountType).First(&account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "积分账户不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get account"})
		return
	}

	// 3. 检查余额是否足够
	if account.Balance < req.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "积分余额不足"})
		return
	}

	// 4. 计算手续费（这里简单处理：0手续费，实际可以根据金额和方式计算）
	fee := 0
	actualAmount := req.Amount - fee

	// 5. 序列化账户信息并计算哈希
	accountInfoBytes, err := json.Marshal(req.AccountInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account info"})
		return
	}
	accountInfoStr := string(accountInfoBytes)
	hash := sha256.Sum256([]byte(accountInfoStr))
	accountInfoHash := hex.EncodeToString(hash[:])

	// 6. 创建提现记录
	withdrawal := models.Withdrawal{
		AccountID:       account.ID,
		Amount:          req.Amount,
		Fee:             fee,
		ActualAmount:    actualAmount,
		Method:          req.Method,
		AccountInfo:     accountInfoStr,
		AccountInfoHash: accountInfoHash,
		Status:          models.WithdrawalStatusPending,
	}

	if err := ctrl.DB.Create(&withdrawal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create withdrawal"})
		return
	}

	// 7. 冻结积分（从balance转入frozen_balance）
	err = ctrl.DB.Transaction(func(tx *gorm.DB) error {
		// 扣除可用余额
		if err := tx.Model(&account).Update("balance", gorm.Expr("balance - ?", req.Amount)).Error; err != nil {
			return err
		}
		// 增加冻结余额
		if err := tx.Model(&account).Update("frozen_balance", gorm.Expr("frozen_balance + ?", req.Amount)).Error; err != nil {
			return err
		}

		// 记录冻结流水
		transaction := models.CreditTransaction{
			AccountID:       account.ID,
			Type:            "WITHDRAW_FREEZE",
			Amount:          -req.Amount,
			BalanceBefore:   account.Balance,
			BalanceAfter:    account.Balance - req.Amount,
			Description:     fmt.Sprintf("提现申请 %d 积分", req.Amount),
			TransactionGroupID: func() *uuid.UUID { id := uuid.New(); return &id }(),
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to freeze balance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "提现申请已提交，等待审核",
		"withdrawal": gin.H{
			"id":            withdrawal.ID,
			"amount":        withdrawal.Amount,
			"actualAmount":  withdrawal.ActualAmount,
			"status":        withdrawal.Status,
			"createdAt":     withdrawal.CreatedAt,
		},
	})
}

// GetWithdrawals 获取提现记录列表
func (ctrl *WithdrawalController) GetWithdrawals(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	// 分页参数
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")
	status := c.Query("status")

	var withdrawals []models.Withdrawal
	var total int64

	query := ctrl.DB.Model(&models.Withdrawal{})

	// 超管可以查看所有记录，其他用户只能查看自己的记录
	if !utils.HasRole(user, "SUPER_ADMIN") {
		// 需要先获取用户的积分账户ID
		var accountType models.OwnerType
		if utils.HasRole(user, "SP_ADMIN") {
			accountType = models.OwnerTypeOrgProvider
		} else if utils.HasRole(user, "MERCHANT_ADMIN") || utils.HasRole(user, "MERCHANT_STAFF") {
			accountType = models.OwnerTypeOrgMerchant
		} else {
			accountType = models.OwnerTypeUserPersonal
		}

		var account models.CreditAccount
		if err := ctrl.DB.Where("owner_id = ? AND owner_type = ?", user.ID, accountType).First(&account).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{
				"withdrawals": []models.Withdrawal{},
				"total":       0,
			})
			return
		}
		query = query.Where("account_id = ?", account.ID)
	}

	// 状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	query.Count(&total)

	// 分页查询
	offset := (parseInt(page) - 1) * parseInt(pageSize)
	if err := query.Order("created_at DESC").Offset(offset).Limit(parseInt(pageSize)).Find(&withdrawals).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get withdrawals"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"withdrawals": withdrawals,
		"total":       total,
	})
}

// GetWithdrawal 获取提现详情
func (ctrl *WithdrawalController) GetWithdrawal(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	id := c.Param("id")

	var withdrawal models.Withdrawal
	if err := ctrl.DB.Where("id = ?", id).First(&withdrawal).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "提现记录不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get withdrawal"})
		return
	}

	// 权限检查：超管可以查看所有，其他用户只能查看自己的
	if !utils.HasRole(user, "SUPER_ADMIN") {
		var accountType models.OwnerType
		if utils.HasRole(user, "SP_ADMIN") {
			accountType = models.OwnerTypeOrgProvider
		} else if utils.HasRole(user, "MERCHANT_ADMIN") || utils.HasRole(user, "MERCHANT_STAFF") {
			accountType = models.OwnerTypeOrgMerchant
		} else {
			accountType = models.OwnerTypeUserPersonal
		}

		var account models.CreditAccount
		if err := ctrl.DB.Where("owner_id = ? AND owner_type = ?", user.ID, accountType).First(&account).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权访问此记录"})
			return
		}

		if withdrawal.AccountID != account.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权访问此记录"})
			return
		}
	}

	c.JSON(http.StatusOK, withdrawal)
}

// AuditWithdrawal 审核提现（仅超管）
func (ctrl *WithdrawalController) AuditWithdrawal(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	// 只有超管可以审核
	if !utils.HasRole(user, "SUPER_ADMIN") {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权审核提现申请"})
		return
	}

	id := c.Param("id")

	var req AuditWithdrawalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	var withdrawal models.Withdrawal
	if err := ctrl.DB.Where("id = ?", id).First(&withdrawal).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "提现记录不存在"})
		return
	}

	// 检查状态
	if !withdrawal.CanAudit() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该提现申请已被处理"})
		return
	}

	// 获取关联账户
	var account models.CreditAccount
	if err := ctrl.DB.Where("id = ?", withdrawal.AccountID).First(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get account"})
		return
	}

	err := ctrl.DB.Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		if req.Approved {
			// 通过：更新状态为approved
			newStatus := models.WithdrawalStatusApproved
			auditNote := req.AuditNote
			if err := tx.Model(&withdrawal).Updates(map[string]interface{}{
				"status":     newStatus,
				"audit_note": auditNote,
				"audited_by": user.ID,
				"audited_at": &now,
			}).Error; err != nil {
				return err
			}
		} else {
			// 拒绝：解冻积分（从frozen_balance转回balance）
			newStatus := models.WithdrawalStatusRejected
			auditNote := req.AuditNote

			// 更新提现状态
			if err := tx.Model(&withdrawal).Updates(map[string]interface{}{
				"status":     newStatus,
				"audit_note": auditNote,
				"audited_by": user.ID,
				"audited_at": &now,
			}).Error; err != nil {
				return err
			}

			// 解冻积分
			if err := tx.Model(&account).Update("frozen_balance", gorm.Expr("frozen_balance - ?", withdrawal.Amount)).Error; err != nil {
				return err
			}
			if err := tx.Model(&account).Update("balance", gorm.Expr("balance + ?", withdrawal.Amount)).Error; err != nil {
				return err
			}

			// 记录退款流水
			transaction := models.CreditTransaction{
				AccountID:       account.ID,
				Type:            "WITHDRAW_REFUND",
				Amount:          withdrawal.Amount,
				BalanceBefore:   account.FrozenBalance,
				BalanceAfter:    account.FrozenBalance - withdrawal.Amount,
				Description:     fmt.Sprintf("提现拒绝退款 %d 积分", withdrawal.Amount),
				TransactionGroupID: func() *uuid.UUID { id := uuid.New(); return &id }(),
			}
			if err := tx.Create(&transaction).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to audit withdrawal"})
		return
	}

	statusText := "已通过"
	if !req.Approved {
		statusText = "已拒绝"
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("提现申请%s", statusText),
		"withdrawal": gin.H{
			"id":        withdrawal.ID,
			"status":    withdrawal.Status,
			"auditNote": req.AuditNote,
		},
	})
}

// ProcessWithdrawal 处理打款（仅超管）
func (ctrl *WithdrawalController) ProcessWithdrawal(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	// 只有超管可以打款
	if !utils.HasRole(user, "SUPER_ADMIN") {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权处理打款"})
		return
	}

	id := c.Param("id")

	var withdrawal models.Withdrawal
	if err := ctrl.DB.Where("id = ?", id).First(&withdrawal).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "提现记录不存在"})
		return
	}

	// 检查状态
	if !withdrawal.CanProcess() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该提现申请不能打款"})
		return
	}

	// 获取关联账户
	var account models.CreditAccount
	if err := ctrl.DB.Where("id = ?", withdrawal.AccountID).First(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get account"})
		return
	}

	err := ctrl.DB.Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		// 更新状态为completed
		if err := tx.Model(&withdrawal).Updates(map[string]interface{}{
			"status":       models.WithdrawalStatusCompleted,
			"completed_at": &now,
		}).Error; err != nil {
			return err
		}

		// 扣除冻结余额
		if err := tx.Model(&account).Update("frozen_balance", gorm.Expr("frozen_balance - ?", withdrawal.Amount)).Error; err != nil {
			return err
		}

		// 记录提现流水
		transaction := models.CreditTransaction{
			AccountID:       account.ID,
			Type:            "WITHDRAW",
			Amount:          -withdrawal.Amount,
			BalanceBefore:   account.FrozenBalance,
			BalanceAfter:    account.FrozenBalance - withdrawal.Amount,
			Description:     fmt.Sprintf("提现成功 %d 积分", withdrawal.Amount),
			TransactionGroupID: func() *uuid.UUID { id := uuid.New(); return &id }(),
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		// TODO: 实际环境中，这里应该调用第三方支付接口进行打款
		// 例如：支付宝、微信支付、银行转账等

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process withdrawal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "打款成功",
		"withdrawal": gin.H{
			"id":           withdrawal.ID,
			"status":       withdrawal.Status,
			"actualAmount": withdrawal.ActualAmount,
			"completedAt":  withdrawal.CompletedAt,
		},
	})
}

// parseInt 辅助函数：字符串转int
func parseInt(s string) int {
	var result int
	fmt.Sscanf(s, "%d", &result)
	return result
}
