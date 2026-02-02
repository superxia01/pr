package controllers

import (
	"net/http"
	"pr-business/config"
	"pr-business/models"
	"pr-business/utils"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthController struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewAuthController(db *gorm.DB, cfg *config.Config) *AuthController {
	return &AuthController{
		db:  db,
		cfg: cfg,
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	AuthCode string `json:"authCode" binding:"required"` // 微信授权码
}

// PasswordLoginRequest 密码登录请求
type PasswordLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string   `json:"accessToken"`
	RefreshToken string   `json:"refreshToken"`
	ExpiresIn    int64    `json:"expiresIn"`
	UserID       string   `json:"userId"`
	Roles        []string `json:"roles"`
	CurrentRole  string   `json:"currentRole"`
}

// WeChatLogin 微信登录
func (ctrl *AuthController) WeChatLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// TODO: 调用auth-center验证微信授权码
	// 这里暂时模拟创建/获取用户

	// 模拟从auth-center获取的用户信息
	authCenterUserID := "mock-auth-center-id" // 应从auth-center获取

	// 查找或创建用户
	var user models.User
	result := ctrl.db.Where("auth_center_user_id = ?", authCenterUserID).First(&user)

	if result.Error == gorm.ErrRecordNotFound {
		// 创建新用户
		user = models.User{
			AuthCenterUserID: authCenterUserID,
			Nickname:         "新用户",
			Roles:            []string{},
			Status:           "active",
		}
		if err := ctrl.db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create user",
			})
			return
		}
	} else if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	// 更新最后登录时间
	now := time.Now()
	user.LastLoginAt = &now
	ctrl.db.Save(&user)

	// 生成JWT token
	accessToken, err := utils.GenerateAccessToken(
		user.ID,
		user.CurrentRole,
		user.Roles,
		ctrl.cfg.JWTSecret,
		ctrl.cfg.JWTAccessTokenExpire,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate access token",
		})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(
		user.ID,
		ctrl.cfg.JWTSecret,
		ctrl.cfg.JWTRefreshTokenExpire,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate refresh token",
		})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(ctrl.cfg.JWTAccessTokenExpire.Seconds()),
		UserID:       user.ID,
		Roles:        user.Roles,
		CurrentRole:  user.CurrentRole,
	})
}

// PasswordLogin 密码登录
func (ctrl *AuthController) PasswordLogin(c *gin.Context) {
	var req PasswordLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// TODO: 调用auth-center验证用户名密码
	// 这里暂时模拟登录流程

	// 模拟从auth-center获取的用户信息
	authCenterUserID := "mock-auth-center-id" // 应从auth-center获取

	// 查找用户
	var user models.User
	result := ctrl.db.Where("auth_center_user_id = ?", authCenterUserID).First(&user)

	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username or password",
		})
		return
	} else if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	// 检查用户状态
	if user.Status != "active" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "User account is inactive or banned",
		})
		return
	}

	// 更新最后登录时间
	now := time.Now()
	user.LastLoginAt = &now
	ctrl.db.Save(&user)

	// 生成JWT token
	accessToken, err := utils.GenerateAccessToken(
		user.ID,
		user.CurrentRole,
		user.Roles,
		ctrl.cfg.JWTSecret,
		ctrl.cfg.JWTAccessTokenExpire,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate access token",
		})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(
		user.ID,
		ctrl.cfg.JWTSecret,
		ctrl.cfg.JWTRefreshTokenExpire,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate refresh token",
		})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(ctrl.cfg.JWTAccessTokenExpire.Seconds()),
		UserID:       user.ID,
		Roles:        user.Roles,
		CurrentRole:  user.CurrentRole,
	})
}

// RefreshToken 刷新令牌
func (ctrl *AuthController) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// 验证refresh token
	claims, err := utils.ParseToken(req.RefreshToken, ctrl.cfg.JWTSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid refresh token",
		})
		return
	}

	// 查找用户
	var user models.User
	result := ctrl.db.Where("id = ?", claims.UserID).First(&user)

	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found",
		})
		return
	}

	// 生成新的access token
	accessToken, err := utils.GenerateAccessToken(
		user.ID,
		user.CurrentRole,
		user.Roles,
		ctrl.cfg.JWTSecret,
		ctrl.cfg.JWTAccessTokenExpire,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate access token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken": accessToken,
		"expiresIn":   int64(ctrl.cfg.JWTAccessTokenExpire.Seconds()),
	})
}

// SwitchRole 切换角色
func (ctrl *AuthController) SwitchRole(c *gin.Context) {
	userID := c.GetString("userId")
	rolesInterface, _ := c.Get("roles")
	roles := rolesInterface.([]string)

	var req struct {
		NewRole string `json:"newRole" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// 检查用户是否拥有该角色
	hasRole := false
	for _, role := range roles {
		if role == req.NewRole {
			hasRole = true
			break
		}
	}

	if !hasRole {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "User does not have this role",
		})
		return
	}

	// 更新当前角色
	var user models.User
	if err := ctrl.db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	user.LastUsedRole = user.CurrentRole
	user.CurrentRole = req.NewRole
	if err := ctrl.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to switch role",
		})
		return
	}

	// �新token
	accessToken, err := utils.GenerateAccessToken(
		user.ID,
		user.CurrentRole,
		user.Roles,
		ctrl.cfg.JWTSecret,
		ctrl.cfg.JWTAccessTokenExpire,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate access token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken": accessToken,
		"currentRole": user.CurrentRole,
		"lastUsedRole": user.LastUsedRole,
	})
}

// GetCurrentUser 获取当前用户信息
func (ctrl *AuthController) GetCurrentUser(c *gin.Context) {
	userID := c.GetString("userId")

	var user models.User
	if err := ctrl.db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}
