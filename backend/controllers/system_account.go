package controllers

import (
	"net/http"
	"pr-business/constants"
	"pr-business/models"
	"pr-business/services"
	"pr-business/utils"

	"github.com/gin-gonic/gin"
)

// SystemAccountController 系统账户控制器
type SystemAccountController struct {
	systemAccountService *services.SystemAccountService
	auditService        *services.AuditService
}

func NewSystemAccountController(
	systemAccountService *services.SystemAccountService,
	auditService *services.AuditService,
) *SystemAccountController {
	return &SystemAccountController{
		systemAccountService: systemAccountService,
		auditService:        auditService,
	}
}

// GetSystemAccounts 获取系统账户列表
// @Summary 获取系统账户列表
// @Description 获取指定类型的系统账户列表
// @Tags 系统账户管理
// @Accept json
// @Produce json
// @Param account_type query string true "账户类型"
// @Success 200 {array} models.SystemAccount
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/system-accounts [get]
func (c *SystemAccountController) GetSystemAccounts(ctx *gin.Context) {
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
	balance, err := c.systemAccountService.GetEscrowBalance(accountType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "查询系统账户失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"account_type": accountType,
		"balance":     balance,
	})
}

// GetFinancialSummary 获取财务汇总
// @Summary 获取财务汇总
// @Description 获取所有系统账户的财务汇总信息
// @Tags 系统账户管理
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/system-accounts/summary [get]
func (c *SystemAccountController) GetFinancialSummary(ctx *gin.Context) {
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

	// 3. 查询各类账户余额
	summary := make(map[string]interface{})

	// 票务托管账户
	ticketEscrow, _ := c.systemAccountService.GetEscrowBalance(constants.SystemAccountTypeTicketEscrow)
	summary["ticket_escrow"] = ticketEscrow

	// 任务托管账户
	taskEscrow, _ := c.systemAccountService.GetEscrowBalance(constants.SystemAccountTypeTaskEscrow)
	summary["task_escrow"] = taskEscrow

	// 平台收益（统一账户）
	platformRevenue, _ := c.systemAccountService.GetEscrowBalance(constants.SystemAccountTypePlatformRevenue)
	summary["platform_revenue"] = platformRevenue

	// 总托管余额
	totalEscrow := ticketEscrow + taskEscrow
	summary["total_escrow"] = totalEscrow

	ctx.JSON(http.StatusOK, summary)
}
