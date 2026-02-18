package controllers

import (
	"net/http"
	"pr-business/constants"
	"pr-business/models"
	"pr-business/services"
	"pr-business/utils"

	"github.com/gin-gonic/gin"
)

// CashAccountController 现金账户控制器
type CashAccountController struct {
	cashAccountService *services.CashAccountService
	auditService      *services.AuditService
}

func NewCashAccountController(
	cashAccountService *services.CashAccountService,
	auditService *services.AuditService,
) *CashAccountController {
	return &CashAccountController{
		cashAccountService: cashAccountService,
		auditService:      auditService,
	}
}

// GetCashAccounts 获取现金账户列表
// @Summary 获取现金账户列表
// @Description 获取指定类型的现金账户列表
// @Tags 现金账户管理
// @Accept json
// @Produce json
// @Param account_type query string true "账户类型"
// @Success 200 {array} models.CashAccount
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/cash-accounts [get]
func (c *CashAccountController) GetCashAccounts(ctx *gin.Context) {
	// 1. 获取当前用户
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	userObj, ok := user.(*models.User)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户信息格式错误"})
		return
	}

	// 2. 权限检查（只有管理员可以查看）
	if !utils.IsSuperAdmin(userObj) && !utils.IsServiceProviderAdmin(userObj) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "没有权限执行此操作"})
		return
	}

	// 3. 获取账户类型
	accountType := ctx.Query("account_type")
	if accountType == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "账户类型不能为空"})
		return
	}

	// 4. 调用服务层查询
	accounts, err := c.cashAccountService.GetActiveCashAccounts(accountType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "查询现金账户失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

// CreateCashAccountRequest 创建现金账户请求
type CreateCashAccountRequest struct {
	AccountType    string `json:"account_type" binding:"required"`
	Description    string `json:"description"`
	InitialBalance int    `json:"initial_balance"` // 单位：分
}

// CreateCashAccount 创建现金账户
// @Summary 创建现金账户
// @Description 创建新的现金账户
// @Tags 现金账户管理
// @Accept json
// @Produce json
// @Param request body CreateCashAccountRequest true "账户信息"
// @Success 200 {object} models.CashAccount
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/cash-accounts [post]
func (c *CashAccountController) CreateCashAccount(ctx *gin.Context) {
	// 1. 获取当前用户
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	userObj, ok := user.(*models.User)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户信息格式错误"})
		return
	}

	// 2. 权限检查（只有超级管理员可以创建）
	if !utils.IsSuperAdmin(userObj) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "没有权限执行此操作"})
		return
	}

	// 3. 绑定请求参数
	var req CreateCashAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 4. 验证账户类型
	validTypes := map[string]bool{
		constants.CashAccountTypeWeChat:       true,
		constants.CashAccountTypeAlipay:      true,
		constants.CashAccountTypeBankTransfer: true,
		constants.CashAccountTypeMarketing:    true,
		constants.CashAccountTypeOperations:  true,
	}
	if !validTypes[req.AccountType] {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的账户类型"})
		return
	}

	// 5. 调用服务层创建
	account, err := c.cashAccountService.CreateCashAccount(
		req.AccountType,
		req.Description,
		userObj.AuthCenterUserID,
		req.InitialBalance,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建现金账户失败: " + err.Error()})
		return
	}

	// 6. 记录审计日志
	ipAddress := ctx.ClientIP()
	userAgent := ctx.GetHeader("User-Agent")
	_ = c.auditService.LogFinancialOperation(
		userObj.AuthCenterUserID,
		constants.AuditActionSystemAdjust,
		constants.AuditResourceCashAccount,
		account.ID.String(),
		map[string]interface{}{
			"action":         "create",
			"account_type":   req.AccountType,
			"initial_balance": req.InitialBalance,
		},
		ipAddress,
		userAgent,
	)

	ctx.JSON(http.StatusOK, account)
}

// UpdateCashAccountBalanceRequest 更新现金账户余额请求
type UpdateCashAccountBalanceRequest struct {
	Amount      int    `json:"amount" binding:"required"`      // 变动金额（单位：分），正数为增加，负数为减少
	Description string `json:"description" binding:"required"` // 变动说明
}

// UpdateCashAccountBalance 更新现金账户余额
// @Summary 更新现金账户余额
// @Description 手动调整现金账户余额
// @Tags 现金账户管理
// @Accept json
// @Produce json
// @Param id path string true "账户ID"
// @Param request body UpdateCashAccountBalanceRequest true "余额变动信息"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/cash-accounts/{id}/balance [post]
func (c *CashAccountController) UpdateCashAccountBalance(ctx *gin.Context) {
	// 1. 获取当前用户
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	userObj, ok := user.(*models.User)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户信息格式错误"})
		return
	}

	// 2. 权限检查（只有超级管理员可以调整）
	if !utils.IsSuperAdmin(userObj) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "没有权限执行此操作"})
		return
	}

	// 3. 获取账户ID
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "账户ID不能为空"})
		return
	}

	// 4. 绑定请求参数
	var req UpdateCashAccountBalanceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 5. 调用服务层更新
	err := c.cashAccountService.UpdateBalance(id, req.Amount, req.Description, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新账户余额失败: " + err.Error()})
		return
	}

	// 6. 记录审计日志
	ipAddress := ctx.ClientIP()
	userAgent := ctx.GetHeader("User-Agent")
	_ = c.auditService.LogFinancialOperation(
		userObj.AuthCenterUserID,
		constants.AuditActionSystemAdjust,
		constants.AuditResourceCashAccount,
		id,
		map[string]interface{}{
			"action":      "update_balance",
			"amount":      req.Amount,
			"description": req.Description,
		},
		ipAddress,
		userAgent,
	)

	ctx.JSON(http.StatusOK, gin.H{"message": "余额更新成功"})
}
