package controllers

import (
	"errors"
	"net/http"
	"pr-business/constants"
	"pr-business/models"
	"pr-business/services"
	"pr-business/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// WithdrawalEnhancedController 增强提现控制器
type WithdrawalEnhancedController struct {
	db                *gorm.DB
	withdrawalService *services.WithdrawalEnhancedService
	auditService      *services.AuditService
}

func NewWithdrawalEnhancedController(
	db *gorm.DB,
	withdrawalService *services.WithdrawalEnhancedService,
	auditService *services.AuditService,
) *WithdrawalEnhancedController {
	return &WithdrawalEnhancedController{
		db:                db,
		withdrawalService: withdrawalService,
		auditService:      auditService,
	}
}

// CreateWithdrawalRequestRequest 创建提现申请请求
type CreateWithdrawalRequestRequest struct {
	Amount      int     `json:"amount" binding:"required"`       // 单位：积分
	YuanAmount  float64 `json:"yuanAmount" binding:"required"`   // 单位：元
	Description string  `json:"description"`
}

// CreateWithdrawalRequest 创建提现申请
// @Summary 创建提现申请
// @Description 用户创建提现申请，系统自动冻结相应积分
// @Tags 提现管理
// @Accept json
// @Produce json
// @Param request body CreateWithdrawalRequestRequest true "提现申请信息"
// @Success 200 {object} models.WithdrawalRequest
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/withdrawals/enhanced [post]
func (c *WithdrawalEnhancedController) CreateWithdrawalRequest(ctx *gin.Context) {
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

	// 2. 绑定请求参数
	var req CreateWithdrawalRequestRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 3. 验证金额
	if req.Amount <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "提现积分必须大于0"})
		return
	}

	if req.YuanAmount <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "提现金额必须大于0"})
		return
	}

	// 5. 调用服务层创建提现申请
	withdrawalReq := &services.CreateWithdrawalRequest{
		UserID:     userObj.AuthCenterUserID,
		Amount:      req.Amount,
		YuanAmount:  req.YuanAmount,
		Description: req.Description,
	}

	result, err := c.withdrawalService.CreateWithdrawalRequest(withdrawalReq, userObj.AuthCenterUserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建提现申请失败: " + err.Error()})
		return
	}

	// 6. 记录审计日志
	ipAddress := ctx.ClientIP()
	userAgent := ctx.GetHeader("User-Agent")
	_ = c.auditService.LogFinancialOperation(
		userObj.AuthCenterUserID,
		constants.AuditActionWithdrawalRequest,
		constants.AuditResourceWithdrawalRequest,
		result.ID.String(),
		map[string]interface{}{
			"action":       "create",
			"amount":       req.Amount,
			"yuan_amount":  req.YuanAmount,
		},
		ipAddress,
		userAgent,
	)

	ctx.JSON(http.StatusOK, result)
}

// ApproveWithdrawalRequestRequest 审核通过请求
type ApproveWithdrawalRequestRequest struct {
	CashAccountType string `json:"cashAccountType" binding:"required"` // 现金账户类型
}

// ApproveWithdrawalRequest 审核通过提现申请
// @Summary 审核通过提现申请
// @Description 管理员审核通过提现申请，系统解冻积分并扣除现金
// @Tags 提现管理
// @Accept json
// @Produce json
// @Param id path string true "提现申请ID"
// @Param request body ApproveWithdrawalRequestRequest true "审核信息"
// @Success 200 {object} models.WithdrawalRequest
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/withdrawals/enhanced/{id}/approve [post]
func (c *WithdrawalEnhancedController) ApproveWithdrawalRequest(ctx *gin.Context) {
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

	// 2. 权限检查（只有超级管理员和客服管理员可以审核）
	if !utils.IsSuperAdmin(userObj) && !utils.IsServiceProviderAdmin(userObj) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "没有权限执行此操作"})
		return
	}

	// 3. 获取提现申请ID
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "提现申请ID不能为空"})
		return
	}

	// 4. 绑定请求参数
	var req ApproveWithdrawalRequestRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 5. 调用服务层审核通过
	result, err := c.withdrawalService.ApproveWithdrawalRequest(id, userObj.AuthCenterUserID, req.CashAccountType)
	if err != nil {
		if errors.Is(err, services.ErrWithdrawalNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "提现申请不存在"})
			return
		}
		if errors.Is(err, services.ErrInvalidWithdrawalStatus) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "提现申请状态不正确"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "审核通过失败: " + err.Error()})
		return
	}

	// 6. 记录审计日志
	ipAddress := ctx.ClientIP()
	userAgent := ctx.GetHeader("User-Agent")
	_ = c.auditService.LogFinancialOperation(
		userObj.AuthCenterUserID,
		constants.AuditActionWithdrawalApprove,
		constants.AuditResourceWithdrawalRequest,
		result.ID.String(),
		map[string]interface{}{
			"action":           "approve",
			"cash_account_type": req.CashAccountType,
			"amount":           result.Amount,
		},
		ipAddress,
		userAgent,
	)

	ctx.JSON(http.StatusOK, result)
}

// RejectWithdrawalRequestRequest 审核拒绝请求
type RejectWithdrawalRequestRequest struct {
	RejectReason string `json:"rejectReason" binding:"required"` // 拒绝原因
}

// RejectWithdrawalRequest 审核拒绝提现申请
// @Summary 审核拒绝提现申请
// @Description 管理员审核拒绝提现申请，系统退还冻结积分
// @Tags 提现管理
// @Accept json
// @Produce json
// @Param id path string true "提现申请ID"
// @Param request body RejectWithdrawalRequestRequest true "拒绝原因"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/withdrawals/enhanced/{id}/reject [post]
func (c *WithdrawalEnhancedController) RejectWithdrawalRequest(ctx *gin.Context) {
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

	// 2. 权限检查（只有超级管理员和客服管理员可以审核）
	if !utils.IsSuperAdmin(userObj) && !utils.IsServiceProviderAdmin(userObj) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "没有权限执行此操作"})
		return
	}

	// 3. 获取提现申请ID
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "提现申请ID不能为空"})
		return
	}

	// 4. 绑定请求参数
	var req RejectWithdrawalRequestRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 5. 调用服务层审核拒绝
	err := c.withdrawalService.RejectWithdrawalRequest(id, req.RejectReason, userObj.AuthCenterUserID)
	if err != nil {
		if errors.Is(err, services.ErrWithdrawalNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "提现申请不存在"})
			return
		}
		if errors.Is(err, services.ErrInvalidWithdrawalStatus) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "提现申请状态不正确"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "审核拒绝失败: " + err.Error()})
		return
	}

	// 6. 记录审计日志
	ipAddress := ctx.ClientIP()
	userAgent := ctx.GetHeader("User-Agent")
	_ = c.auditService.LogFinancialOperation(
		userObj.AuthCenterUserID,
		constants.AuditActionWithdrawalReject,
		constants.AuditResourceWithdrawalRequest,
		id,
		map[string]interface{}{
			"action":        "reject",
			"reject_reason": req.RejectReason,
		},
		ipAddress,
		userAgent,
	)

	ctx.JSON(http.StatusOK, gin.H{"message": "提现申请已拒绝"})
}

// GetWithdrawalRequests 查询提现申请列表
// @Summary 查询提现申请列表
// @Description 查询提现申请列表，支持分页和状态过滤
// @Tags 提现管理
// @Accept json
// @Produce json
// @Param status query string false "状态过滤"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} utils.PageResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/withdrawals/enhanced [get]
func (c *WithdrawalEnhancedController) GetWithdrawalRequests(ctx *gin.Context) {
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

	// 2. 获取查询参数
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "20")
	status := ctx.Query("status")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 3. 构建查询
	query := c.db.Model(&models.WithdrawalRequest{})

	// 状态过滤
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 权限过滤：普通用户只能查看自己的提现申请
	if !utils.IsSuperAdmin(userObj) && !utils.IsServiceProviderAdmin(userObj) {
		// 查找用户的所有积分账户
		var accountIDs []string
		c.db.Model(&models.CreditAccount{}).
			Where("owner_id = ? AND account_type = ?", userObj.AuthCenterUserID, "CASH").
			Pluck("id", &accountIDs)

		if len(accountIDs) > 0 {
			query = query.Where("account_id IN ?", accountIDs)
		} else {
			// 如果没有账户，返回空列表
			ctx.JSON(http.StatusOK, gin.H{
				"list":      []models.WithdrawalRequest{},
				"total":     0,
				"page":      page,
				"page_size": pageSize,
			})
			return
		}
	}

	// 4. 获取总数
	var total int64
	query.Count(&total)

	// 5. 分页查询
	var withdrawals []models.WithdrawalRequest
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&withdrawals).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "查询提现申请列表失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"list":      withdrawals,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetWithdrawalRequest 获取单个提现申请详情
// @Summary 获取提现申请详情
// @Description 获取指定提现申请的详细信息
// @Tags 提现管理
// @Accept json
// @Produce json
// @Param id path string true "提现申请ID"
// @Success 200 {object} models.WithdrawalRequest
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/withdrawals/enhanced/{id} [get]
func (c *WithdrawalEnhancedController) GetWithdrawalRequest(ctx *gin.Context) {
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

	// 2. 获取提现申请ID
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "提现申请ID不能为空"})
		return
	}

	// 3. 查询提现申请
	var withdrawal models.WithdrawalRequest
	if err := c.db.Where("id = ?", id).First(&withdrawal).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "提现申请不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "查询提现申请失败"})
		return
	}

	// 4. 权限检查：普通用户只能查看自己的提现申请
	if !utils.IsSuperAdmin(userObj) && !utils.IsServiceProviderAdmin(userObj) {
		// 检查是否是申请人自己的提现
		// AccountID存储的是申请人的积分账户ID
		// 需要通过账户找到对应的用户
		var account models.CreditAccount
		if err := c.db.Where("id = ?", withdrawal.AccountID).First(&account).Error; err == nil {
			if account.OwnerID.String() != userObj.AuthCenterUserID {
				ctx.JSON(http.StatusForbidden, gin.H{"error": "无权限查看此提现申请"})
				return
			}
		}
	}

	ctx.JSON(http.StatusOK, withdrawal)
}
