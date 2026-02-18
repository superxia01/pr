package controllers

import (
	"net/http"
	"pr-business/models"
	"pr-business/services"
	"pr-business/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// FinancialAuditController 财务审计日志控制器
type FinancialAuditController struct {
	auditService *services.AuditService
}

func NewFinancialAuditController(auditService *services.AuditService) *FinancialAuditController {
	return &FinancialAuditController{
		auditService: auditService,
	}
}

// GetAuditLogs 查询审计日志列表
// @Summary 查询审计日志列表
// @Description 查询财务审计日志，支持多种过滤条件和分页
// @Tags 审计日志
// @Accept json
// @Produce json
// @Param user_id query string false "用户ID过滤"
// @Param action query string false "操作类型过滤"
// @Param resource_type query string false "资源类型过滤"
// @Param resource_id query string false "资源ID过滤"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} utils.PageResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/financial-audit-logs [get]
func (c *FinancialAuditController) GetAuditLogs(ctx *gin.Context) {
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

	// 3. 获取查询参数
	userID := ctx.Query("user_id")
	action := ctx.Query("action")
	resourceType := ctx.Query("resource_type")
	resourceID := ctx.Query("resource_id")
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	// 4. 调用服务层查询
	logs, err := c.auditService.QueryAuditLogs(
		userID,
		action,
		resourceType,
		resourceID,
		pageSize,
		offset,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "查询审计日志失败: " + err.Error()})
		return
	}

	// 5. 返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"list":      logs,
		"total":     len(logs),
		"page":      page,
		"page_size": pageSize,
	})
}
