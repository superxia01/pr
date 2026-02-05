package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"pr-business/models"
	"pr-business/services"
	"pr-business/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreditController struct {
	db                   *gorm.DB
	permissionService    *services.AccountPermissionService
}

func NewCreditController(db *gorm.DB) *CreditController {
	return &CreditController{
		db:                db,
		permissionService: services.NewAccountPermissionService(db),
	}
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

	// 超级管理员使用个人账户
	if utils.IsSuperAdmin(user) || utils.HasRole(user, "admin") {
		parsedID, err := uuid.Parse(user.AuthCenterUserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "用户ID格式错误"})
			return
		}

		account, err := ctrl.findOrCreateAccount(parsedID, models.OwnerTypeUserPersonal, parsedID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询或创建账户失败"})
			return
		}

		c.JSON(http.StatusOK, account)
		return
	}

	// 确定账户类型
	var ownerType models.OwnerType
	var ownerID uuid.UUID

	if utils.IsMerchantAdmin(user) {
		// 商家管理员：查找商家账户
		var merchant models.Merchant
		if err := ctrl.db.Where("admin_id::text = ?", user.AuthCenterUserID).First(&merchant).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "商家信息不存在"})
			return
		}
		ownerType = models.OwnerTypeOrgMerchant
		ownerID = merchant.ID
	} else if utils.IsServiceProviderAdmin(user) || utils.IsServiceProviderStaff(user) {
		// 服务商：查找服务商账户
		var provider models.ServiceProvider
		if utils.IsServiceProviderAdmin(user) {
			if err := ctrl.db.Where("admin_id::text = ?", user.AuthCenterUserID).First(&provider).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "服务商信息不存在"})
				return
			}
		} else {
			// 服务商员工
			var providerStaff models.ServiceProviderStaff
			if err := ctrl.db.Where("user_id::text = ?", user.AuthCenterUserID).First(&providerStaff).Error; err != nil {
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
	} else if utils.IsCreator(user) {
		// 达人：个人账户
		ownerType = models.OwnerTypeUserPersonal
		// 使用user.AuthCenterUserID作为ownerID
		parsedID, err := uuid.Parse(user.AuthCenterUserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "用户ID格式错误"})
			return
		}
		ownerID = parsedID
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "无法确定账户类型"})
		return
	}

	// 获取用户ID用于权限授予
	userID, err := uuid.Parse(user.AuthCenterUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户ID格式错误"})
		return
	}

	// 查找或创建账户（自动授予权限）
	account, err := ctrl.findOrCreateAccount(ownerID, ownerType, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询或创建账户失败"})
		return
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

	// 超级管理员使用个人账户
	if utils.IsSuperAdmin(user) || utils.HasRole(user, "admin") {
		ownerType = models.OwnerTypeUserPersonal
		parsedID, err := uuid.Parse(user.AuthCenterUserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "用户ID格式错误"})
			return
		}
		ownerID = parsedID
	} else if utils.IsMerchantAdmin(user) {
		var merchant models.Merchant
		if err := ctrl.db.Where("admin_id::text = ?", user.AuthCenterUserID).First(&merchant).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "商家信息不存在"})
			return
		}
		ownerType = models.OwnerTypeOrgMerchant
		ownerID = merchant.ID
	} else if utils.IsServiceProviderAdmin(user) {
		var provider models.ServiceProvider
		if err := ctrl.db.Where("admin_id::text = ?", user.AuthCenterUserID).First(&provider).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "服务商信息不存在"})
			return
		}
		ownerType = models.OwnerTypeOrgProvider
		ownerID = provider.ID
	} else if utils.IsCreator(user) {
		ownerType = models.OwnerTypeUserPersonal
		parsedID, err := uuid.Parse(user.AuthCenterUserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "用户ID格式错误"})
			return
		}
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
	if !utils.HasRole(user, "merchant_admin") {
		c.JSON(http.StatusForbidden, gin.H{"error": "只有商家可以充值"})
		return
	}

	// 获取商家信息
	var merchant models.Merchant
	if err := ctrl.db.Where("admin_id::text = ?", user.AuthCenterUserID).First(&merchant).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "商家信息不存在"})
		return
	}

	// 获取用户ID用于权限授予
	userID, err := uuid.Parse(user.AuthCenterUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户ID格式错误"})
		return
	}

	// 获取或创建账户（自动授予权限）
	account, err := ctrl.findOrCreateAccount(merchant.ID, models.OwnerTypeOrgMerchant, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "账户不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询或创建账户失败"})
		}
		return
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

// GetUserAccounts 获取用户可访问的所有账户
// @Summary 获取用户可访问的所有账户
// @Description 根据用户当前角色和权限，返回可访问的账户列表
// @Tags 积分管理
// @Accept json
// @Produce json
// @Success 200 {array} models.CreditAccount
// @Router /api/v1/credit/accounts [get]
func (ctrl *CreditController) GetUserAccounts(c *gin.Context) {
	// 获取当前用户
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	user := currentUser.(*models.User)

	// 解析用户ID
	parsedID, err := uuid.Parse(user.AuthCenterUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户ID格式错误"})
		return
	}

	// 从 account_permissions 表获取用户有权限访问的账户
	var permissions []models.AccountPermission
	if err := ctrl.db.Where("user_id = ?", parsedID).Find(&permissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取账户权限失败"})
		return
	}

	// 如果没有任何权限，返回空列表
	if len(permissions) == 0 {
		c.JSON(http.StatusOK, []models.CreditAccount{})
		return
	}

	// 获取账户ID列表
	accountIDs := make([]uuid.UUID, len(permissions))
	for i, perm := range permissions {
		accountIDs[i] = perm.AccountID
	}

	// 查询所有可访问的账户
	var accounts []models.CreditAccount
	if err := ctrl.db.Where("id IN ?", accountIDs).Find(&accounts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取账户列表失败"})
		return
	}

	// 为每个账户添加权限信息
	type AccountWithPermission struct {
		models.CreditAccount
		CanView    bool `json:"canView"`
		CanOperate bool `json:"canOperate"`
	}

	result := make([]AccountWithPermission, len(accounts))
	permMap := make(map[uuid.UUID]models.AccountPermission)
	for _, perm := range permissions {
		permMap[perm.AccountID] = perm
	}

	for i, account := range accounts {
		perm := permMap[account.ID]
		result[i] = AccountWithPermission{
			CreditAccount: account,
			CanView:       perm.CanView,
			CanOperate:    perm.CanOperate,
		}
	}

	c.JSON(http.StatusOK, result)
}

// 辅助函数：解析整型参数
func parseIntParam(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}

// createAccountWithPermission 创建账户并自动授予权限
func (ctrl *CreditController) createAccountWithPermission(ownerID uuid.UUID, ownerType models.OwnerType, userID uuid.UUID) (*models.CreditAccount, error) {
	account := &models.CreditAccount{
		OwnerID:       ownerID,
		OwnerType:     ownerType,
		UserID:        &userID,
		Balance:       0,
		FrozenBalance: 0,
	}

	if err := ctrl.db.Create(account).Error; err != nil {
		return nil, err
	}

	// 根据账户类型授予权限
	if ownerType == models.OwnerTypeUserPersonal {
		// 个人账户：授予用户完整权限
		if err := ctrl.permissionService.GrantPersonalAccountPermission(userID, account.ID); err != nil {
			return nil, err
		}
	} else if ownerType == models.OwnerTypeOrgMerchant {
		// 商家账户：授予管理员完整权限
		if err := ctrl.permissionService.GrantMerchantAccountPermission(userID, account.ID, true); err != nil {
			return nil, err
		}
	} else if ownerType == models.OwnerTypeOrgProvider {
		// 服务商账户：授予管理员完整权限
		if err := ctrl.permissionService.GrantProviderAccountPermission(userID, account.ID, true); err != nil {
			return nil, err
		}
	}

	return account, nil
}

// findOrCreateAccount 查找或创建账户（带权限）
func (ctrl *CreditController) findOrCreateAccount(ownerID uuid.UUID, ownerType models.OwnerType, userID uuid.UUID) (*models.CreditAccount, error) {
	var account models.CreditAccount
	err := ctrl.db.Where("owner_id = ? AND owner_type = ?", ownerID, ownerType).First(&account).Error
	if err == nil {
		return &account, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ctrl.createAccountWithPermission(ownerID, ownerType, userID)
	}

	return nil, err
}
