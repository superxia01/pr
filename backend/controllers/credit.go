package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"pr-business/models"
	"pr-business/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreditController struct {
	db *gorm.DB
}

func NewCreditController(db *gorm.DB) *CreditController {
	return &CreditController{db: db}
}

// RechargeRequest 充值请求
type RechargeRequest struct {
	Amount int `json:"amount" binding:"required,min=1"`
}

// GetAccountBalance 获取积分余额
// @Summary 获取积分余额
// @Description 获取当前用户的积分余额
// @Tags 积分管理
// @Accept json
// @Produce json
// @Success 200 {object} models.CreditAccount
// @Router /api/v1/credit/balance [get]
func (ctrl *CreditController) GetAccountBalance(c *gin.Context) {
	// 获取当前用户
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	user := currentUser.(*models.User)

	// 确定账户类型
	var ownerType models.OwnerType
	var ownerID uuid.UUID

	if utils.HasRole(user, "MERCHANT_ADMIN") {
		// 商家管理员：查找商家账户
		var merchant models.Merchant
		if err := ctrl.db.Where("admin_id = ?", user.ID).First(&merchant).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "商家信息不存在"})
			return
		}
		ownerType = models.OwnerTypeOrgMerchant
		ownerID = merchant.ID
	} else if utils.HasRole(user, "SP_ADMIN") || utils.HasRole(user, "SP_STAFF") {
		// 服务商：查找服务商账户
		var provider models.ServiceProvider
		if utils.HasRole(user, "SP_ADMIN") {
			if err := ctrl.db.Where("admin_id = ?", user.ID).First(&provider).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "服务商信息不存在"})
				return
			}
		} else {
			// 服务商员工
			var providerStaff models.ServiceProviderStaff
			if err := ctrl.db.Where("user_id = ?", user.ID).First(&providerStaff).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "员工信息不存在"})
				return
			}
			if err := ctrl.db.Where("id = ?", providerStaff.ProviderID).First(&provider).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "服务商信息不存在"})
				return
			}
		}
		ownerType = models.OwnerTypeOrgProvider
		ownerID = provider.ID
	} else if utils.HasRole(user, "CREATOR") {
		// 达人：个人账户
		ownerType = models.OwnerTypeUserPersonal
		// 使用user.ID作为ownerID
		parsedID, err := uuid.Parse(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "用户ID格式错误"})
			return
		}
		ownerID = parsedID
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "无法确定账户类型"})
		return
	}

	// 查找或创建账户
	var account models.CreditAccount
	if err := ctrl.db.Where("owner_id = ? AND owner_type = ?", ownerID, ownerType).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 创建新账户
			account = models.CreditAccount{
				OwnerID:       ownerID,
				OwnerType:     ownerType,
				Balance:       0,
				FrozenBalance: 0,
			}
			if err := ctrl.db.Create(&account).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "创建账户失败"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询账户失败"})
			return
		}
	}

	c.JSON(http.StatusOK, account)
}

// GetTransactions 获取积分流水
// @Summary 获取积分流水
// @Description 获取当前用户的积分流水记录
// @Tags 积分管理
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/credit/transactions [get]
func (ctrl *CreditController) GetTransactions(c *gin.Context) {
	// 获取当前用户
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	user := currentUser.(*models.User)

	// 确定账户ID（复用上面的逻辑）
	var ownerType models.OwnerType
	var ownerID uuid.UUID

	if utils.HasRole(user, "MERCHANT_ADMIN") {
		var merchant models.Merchant
		ctrl.db.Where("admin_id = ?", user.ID).First(&merchant)
		ownerType = models.OwnerTypeOrgMerchant
		ownerID = merchant.ID
	} else if utils.HasRole(user, "SP_ADMIN") {
		var provider models.ServiceProvider
		ctrl.db.Where("admin_id = ?", user.ID).First(&provider)
		ownerType = models.OwnerTypeOrgProvider
		ownerID = provider.ID
	} else if utils.HasRole(user, "CREATOR") {
		ownerType = models.OwnerTypeUserPersonal
		parsedID, _ := uuid.Parse(user.ID)
		ownerID = parsedID
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "无法确定账户类型"})
		return
	}

	// 获取账户
	var account models.CreditAccount
	if err := ctrl.db.Where("owner_id = ? AND owner_type = ?", ownerID, ownerType).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "账户不存在"})
		return
	}

	// 分页
	page := 1
	pageSize := 20
	if p := c.Query("page"); p != "" {
		if num, err := parseIntParam(p); err == nil {
			page = num
		}
	}
	if ps := c.Query("page_size"); ps != "" {
		if num, err := parseIntParam(ps); err == nil {
			pageSize = num
		}
	}

	offset := (page - 1) * pageSize

	var transactions []models.CreditTransaction
	var total int64

	query := ctrl.db.Model(&models.CreditTransaction{}).Where("account_id = ?", account.ID)
	query.Count(&total)

	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取流水失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      transactions,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// Recharge 充值（模拟支付）
// @Summary 充值
// @Description 商家充值积分（模拟支付）
// @Tags 积分管理
// @Accept json
// @Produce json
// @Param request body RechargeRequest true "充值请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/credit/recharge [post]
func (ctrl *CreditController) Recharge(c *gin.Context) {
	var req RechargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取当前用户
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	user := currentUser.(*models.User)

	// 权限检查：只有商家可以充值
	if !utils.HasRole(user, "MERCHANT_ADMIN") {
		c.JSON(http.StatusForbidden, gin.H{"error": "只有商家可以充值"})
		return
	}

	// 获取商家信息
	var merchant models.Merchant
	if err := ctrl.db.Where("admin_id = ?", user.ID).First(&merchant).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "商家信息不存在"})
		return
	}

	// 获取账户
	var account models.CreditAccount
	if err := ctrl.db.Where("owner_id = ? AND owner_type = ?", merchant.ID, models.OwnerTypeOrgMerchant).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 创建新账户
			account = models.CreditAccount{
				OwnerID:       merchant.ID,
				OwnerType:     models.OwnerTypeOrgMerchant,
				Balance:       0,
				FrozenBalance: 0,
			}
			if err := ctrl.db.Create(&account).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "创建账户失败"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询账户失败"})
			return
		}
	}

	// 开始事务
	tx := ctrl.db.Begin()

	// 更新余额
	balanceBefore := account.Balance
	account.Balance += req.Amount

	if err := tx.Save(&account).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "充值失败"})
		return
	}

	// 记录流水
	transaction := models.CreditTransaction{
		AccountID:     account.ID,
		Type:          models.TransactionRecharge,
		Amount:        req.Amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  account.Balance,
		Description:   "商家充值",
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "记录流水失败"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "充值成功",
		"account": account,
	})
}

// 辅助函数：解析整型参数
func parseIntParam(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}
