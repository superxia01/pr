package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	"pr-business/config"
	"pr-business/models"
	"pr-business/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthController struct {
	cfg *config.Config
	db  *gorm.DB
}

func NewAuthController(cfg *config.Config, db *gorm.DB) *AuthController {
	return &AuthController{
		cfg: cfg,
		db:  db,
	}
}

// WeChatLoginRedirect 发起微信登录
// GET /api/v1/auth/wechat/login
func (ctrl *AuthController) WeChatLoginRedirect(c *gin.Context) {
	// 重定向到 auth-center，回调到前端登录页
	authCenterURL := fmt.Sprintf(
		"%s/api/auth/wechat/login?callbackUrl=%s",
		ctrl.cfg.AuthCenterURL,
		url.QueryEscape(ctrl.cfg.FrontendURL+"/login"),
	)

	c.Redirect(http.StatusFound, authCenterURL)
}

// GetCurrentUser 获取当前用户信息
// GET /api/v1/user/me
func (ctrl *AuthController) GetCurrentUser(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "未认证",
		})
		return
	}

	u := user.(*models.User)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"id":               u.ID,
			"authCenterUserId": u.AuthCenterUserID,
			"unionId":          u.UnionID,
			"nickname":         u.Nickname,
			"avatarUrl":        u.AvatarURL,
			"roles":            u.Roles,
			"status":           u.Status,
			"createdAt":        u.CreatedAt,
			"updatedAt":        u.UpdatedAt,
		},
	})
}

// GetUsers 获取用户列表（仅超级管理员）
// @Summary 获取用户列表
// @Description 获取系统中所有用户列表，支持分页和过滤（仅超级管理员）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param role query string false "角色过滤"
// @Param status query string false "状态过滤"
// @Success 200 {object} object
// @Router /api/v1/users [get]
func (ctrl *AuthController) GetUsers(c *gin.Context) {
	// 获取当前用户
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	user := currentUser.(*models.User)

	// 权限检查：只有超级管理员可以查看所有用户
	if !utils.IsSuperAdmin(user) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限访问"})
		return
	}

	// 获取分页参数
	page := 1
	pageSize := 20
	if p := c.Query("page"); p != "" {
		if pageNum, err := strconv.Atoi(p); err == nil && pageNum > 0 {
			page = pageNum
		}
	}
	if ps := c.Query("page_size"); ps != "" {
		if size, err := strconv.Atoi(ps); err == nil && size > 0 && size <= 100 {
			pageSize = size
		}
	}

	// 构建查询
	query := ctrl.db.Model(&models.User{})

	// 角色过滤
	if role := c.Query("role"); role != "" {
		query = query.Where("? = ANY(roles)", role)
	}

	// 状态过滤
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// 计算总数
	var total int64
	query.Count(&total)

	// 分页查询
	var users []models.User
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    users,
		"total":   total,
		"page":    page,
		"pageSize": pageSize,
	})
}

