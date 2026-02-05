package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"pr-business/config"
	"pr-business/models"
	"pr-business/utils"
	"strings"
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
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string   `json:"accessToken"`
	RefreshToken string   `json:"refreshToken"`
	ExpiresIn    int64    `json:"expiresIn"`
	UserID       string   `json:"userId"`
	Nickname     string   `json:"nickname"`
	AvatarURL    string   `json:"avatarUrl"`
	Roles        []string `json:"roles"`
	CurrentRole  string   `json:"currentRole"`
}

// AuthCenterResponse auth-center API 响应
type AuthCenterResponse struct {
	Success bool `json:"success"`
	Data    struct {
		UserID      string `json:"userId"`
		Token       string `json:"token"`
		UnionID     string `json:"unionId"`
		PhoneNumber string `json:"phoneNumber"`
		Email       string `json:"email"`
		CreatedAt   string `json:"createdAt"`
		LastLoginAt string `json:"lastLoginAt"`
		Profile     struct {
			Nickname  string `json:"nickname"`
			AvatarURL string `json:"avatarUrl"`
		} `json:"profile"`
	} `json:"data"`
	Error string `json:"error,omitempty"`
}

// AuthCenterUserInfo auth-center 用户信息响应
type AuthCenterUserInfo struct {
	Success bool `json:"success"`
	Data    struct {
		UserID      string `json:"userId"`
		UnionID     string `json:"unionId"`
		PhoneNumber string `json:"phoneNumber"`
		Profile     struct {
			Nickname  string `json:"nickname"`
			AvatarURL string `json:"avatarUrl"`
		} `json:"profile"`
	} `json:"data"`
}

// ============================================
// 微信登录流程
// ============================================

// WeChatLoginRedirect 发起微信登录（重定向到auth-center）
// GET /api/v1/auth/wechat/login
func (ctrl *AuthController) WeChatLoginRedirect(c *gin.Context) {
	// ✅ 重定向到前端登录页（不是后端回调）
	// auth-center 会带着 token 和 userId 回调到前端
	authCenterURL := fmt.Sprintf(
		"%s/api/auth/wechat/login?callbackUrl=%s",
		ctrl.cfg.AuthCenterURL,
		url.QueryEscape(ctrl.cfg.FrontendURL+"/login"),
	)

	c.Redirect(http.StatusFound, authCenterURL)
}

// WeChatLogin 直接用code登录（适用于前端获取code后调用）
// POST /api/v1/auth/wechat
func (ctrl *AuthController) WeChatLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// 调用auth-center API验证微信授权码
	authCenterResp, err := ctrl.callAuthCenterWechatLogin(req.AuthCode, "open")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "认证服务异常",
		})
		return
	}

	// ✅ 输出 auth-center 登录响应用于调试
	responseJSON, _ := json.Marshal(authCenterResp)
	fmt.Fprintf(os.Stderr, "✅ [DEBUG] auth-center登录响应: %s\n", string(responseJSON))

	if !authCenterResp.Success {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "微信授权失败",
		})
		return
	}

	// 创建或获取本地用户
	user, err := ctrl.findOrCreateUser(authCenterResp.Data.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "用户创建失败",
		})
		return
	}

	// ✅ 直接从登录响应中获取用户资料（头像和昵称）
	updated := false
	// 更新昵称
	if authCenterResp.Data.Profile.Nickname != "" && user.Nickname != authCenterResp.Data.Profile.Nickname {
		user.Nickname = authCenterResp.Data.Profile.Nickname
		updated = true
	}
	// 更新头像
	if authCenterResp.Data.Profile.AvatarURL != "" && user.AvatarURL != authCenterResp.Data.Profile.AvatarURL {
		user.AvatarURL = authCenterResp.Data.Profile.AvatarURL
		updated = true
	}

	// 如果有更新，保存到数据库
	if updated {
		ctrl.db.Save(&user)
		fmt.Fprintf(os.Stderr, "✅ [DEBUG] 从登录响应同步用户信息: nickname=%s, avatar=%s\n",
			user.Nickname, user.AvatarURL)
	} else {
		fmt.Fprintf(os.Stderr, "⚠️  [DEBUG] 登录响应中无用户信息或无需更新\n")
	}

	// 更新最后登录时间
	now := time.Now()
	user.LastLoginAt = &now
	ctrl.db.Save(&user)

	// 生成并返回token
	accessToken, refreshToken := ctrl.generateTokens(user)

	c.JSON(http.StatusOK, LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(ctrl.cfg.JWTAccessTokenExpire.Seconds()),
		UserID:       user.ID,
		Nickname:     user.Nickname,
		AvatarURL:    user.AvatarURL,
		Roles:        convertRolesToUpperCase(user.Roles),
		CurrentRole:  user.ActiveRole,
	})
}

// ============================================
// 密码登录
// ============================================

// PasswordLogin 密码登录
// POST /api/v1/auth/password
func (ctrl *AuthController) PasswordLogin(c *gin.Context) {
	var req PasswordLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// 调用auth-center API验证用户名密码
	authCenterResp, err := ctrl.callAuthCenterPasswordLogin(req.PhoneNumber, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "认证服务异常",
		})
		return
	}

	if !authCenterResp.Success {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "手机号或密码错误",
		})
		return
	}

	// 创建或获取本地用户
	user, err := ctrl.findOrCreateUser(authCenterResp.Data.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "用户创建失败",
		})
		return
	}

	// 检查用户状态
	if user.Status != "active" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "用户账号已被禁用",
		})
		return
	}

	// ✅ 直接从登录响应中获取用户资料（头像和昵称）
	updated := false
	// 更新昵称
	if authCenterResp.Data.Profile.Nickname != "" && user.Nickname != authCenterResp.Data.Profile.Nickname {
		user.Nickname = authCenterResp.Data.Profile.Nickname
		updated = true
	}
	// 更新头像
	if authCenterResp.Data.Profile.AvatarURL != "" && user.AvatarURL != authCenterResp.Data.Profile.AvatarURL {
		user.AvatarURL = authCenterResp.Data.Profile.AvatarURL
		updated = true
	}

	// 如果有更新，保存到数据库
	if updated {
		ctrl.db.Save(&user)
		fmt.Printf("✅ 密码登录时从登录响应同步用户信息: nickname=%s, avatar=%s\n", user.Nickname, user.AvatarURL)
	}

	// 更新最后登录时间
	now := time.Now()
	user.LastLoginAt = &now
	ctrl.db.Save(&user)

	// 生成并返回token
	accessToken, refreshToken := ctrl.generateTokens(user)

	c.JSON(http.StatusOK, LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(ctrl.cfg.JWTAccessTokenExpire.Seconds()),
		UserID:       user.ID,
		Nickname:     user.Nickname,
		AvatarURL:    user.AvatarURL,
		Roles:        convertRolesToUpperCase(user.Roles),
		CurrentRole:  user.ActiveRole,
	})
}

// ============================================
// 辅助方法
// ============================================

// callAuthCenterWechatLogin 调用auth-center微信登录API
func (ctrl *AuthController) callAuthCenterWechatLogin(code, loginType string) (*AuthCenterResponse, error) {
	reqBody := map[string]string{
		"code": code,
		"type": loginType,
	}
	jsonBody, _ := json.Marshal(reqBody)

	resp, err := http.Post(
		fmt.Sprintf("%s/api/auth/wechat/login", ctrl.cfg.AuthCenterURL),
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result AuthCenterResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// callAuthCenterPasswordLogin 调用auth-center密码登录API
func (ctrl *AuthController) callAuthCenterPasswordLogin(phoneNumber, password string) (*AuthCenterResponse, error) {
	reqBody := map[string]string{
		"phoneNumber": phoneNumber,
		"password":    password,
	}
	jsonBody, _ := json.Marshal(reqBody)

	resp, err := http.Post(
		fmt.Sprintf("%s/api/auth/password/login", ctrl.cfg.AuthCenterURL),
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result AuthCenterResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// findOrCreateUser 查找或创建本地用户
func (ctrl *AuthController) findOrCreateUser(authCenterUserID string) (*models.User, error) {
	var user models.User
	result := ctrl.db.Where("auth_center_user_id = ?", authCenterUserID).First(&user)

	fmt.Printf("🔍 findOrCreateUser: authCenterUserID=%s, 是否找到用户=%v\n",
		authCenterUserID, result.Error == nil)

	if result.Error == gorm.ErrRecordNotFound {
		// 创建新用户
		var roles models.Roles
		var currentRole string

		// 检查是否是第一个用户（系统没有其他用户）
		var userCount int64
		ctrl.db.Model(&models.User{}).Count(&userCount)

		if userCount == 0 {
			// 第一个用户自动成为超级管理员
			roles = models.Roles{"super_admin", "merchant_admin", "provider_admin", "creator"}
			currentRole = "super_admin"
		} else {
			// 普通用户：创建时无角色，需要通过邀请码获取
			roles = models.Roles{}
			currentRole = ""
		}

		user = models.User{
			AuthCenterUserID: authCenterUserID,
			Nickname:         "新用户",
			Profile:          models.Profile{},
			Roles:            roles,
			ActiveRole:       currentRole,
			Status:           "active",
		}
		if err := ctrl.db.Create(&user).Error; err != nil {
			return nil, err
		}
		fmt.Printf("✅ 创建新用户: ID=%s\n", user.ID)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

// generateTokens 生成访问令牌和刷新令牌
func (ctrl *AuthController) generateTokens(user *models.User) (string, string) {
	accessToken, err := utils.GenerateAccessToken(
		user.ID,
		user.ActiveRole,
		user.Roles,
		ctrl.cfg.JWTSecret,
		ctrl.cfg.JWTAccessTokenExpire,
	)
	if err != nil {
		return "", ""
	}

	refreshToken, err := utils.GenerateRefreshToken(
		user.ID,
		ctrl.cfg.JWTSecret,
		ctrl.cfg.JWTRefreshTokenExpire,
	)
	if err != nil {
		return accessToken, ""
	}

	return accessToken, refreshToken
}

// convertRolesToUpperCase 将小写角色转换为大写格式（用于前端）
func convertRolesToUpperCase(roles models.Roles) []string {
	result := make([]string, len(roles))
	for i, role := range roles {
		switch role {
		case "super_admin":
			result[i] = "SUPER_ADMIN"
		case "merchant_admin":
			result[i] = "MERCHANT_ADMIN"
		case "merchant_staff":
			result[i] = "MERCHANT_STAFF"
		case "service_provider_admin", "provider_admin":
			result[i] = "SP_ADMIN"
		case "service_provider_staff", "provider_staff":
			result[i] = "SP_STAFF"
		case "creator":
			result[i] = "CREATOR"
		default:
			result[i] = strings.ToUpper(role)
		}
	}
	return result
}

// ============================================
// Token刷新和角色切换
// ============================================

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
		user.ActiveRole,
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

	// 将 interface{} 转换为 models.Roles 类型
	rolesBytes, err := json.Marshal(rolesInterface)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to process roles",
		})
		return
	}

	var roles models.Roles
	if err := json.Unmarshal(rolesBytes, &roles); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to process roles",
		})
		return
	}

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

	user.LastUsedRole = user.ActiveRole
	user.ActiveRole = req.NewRole
	if err := ctrl.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to switch role",
		})
		return
	}

	// 生成新token
	accessToken, err := utils.GenerateAccessToken(
		user.ID,
		user.ActiveRole,
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
		"currentRole": user.ActiveRole,
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

	// ✅ 从请求头获取 auth-center token 并同步用户信息
	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		authCenterToken := strings.TrimPrefix(authHeader, "Bearer ")

		// 调用 auth-center 获取用户信息（头像和昵称）
		userInfo, err := ctrl.callAuthCenterGetUserInfo(authCenterToken)
		if err == nil && userInfo.Success {
			updated := false
			// 更新昵称
			if userInfo.Data.Profile.Nickname != "" && user.Nickname != userInfo.Data.Profile.Nickname {
				user.Nickname = userInfo.Data.Profile.Nickname
				updated = true
			}
			// 更新头像
			if userInfo.Data.Profile.AvatarURL != "" && user.AvatarURL != userInfo.Data.Profile.AvatarURL {
				user.AvatarURL = userInfo.Data.Profile.AvatarURL
				updated = true
			}

			// 如果有更新，保存到数据库
			if updated {
				ctrl.db.Save(&user)
				fmt.Printf("✅ 同步用户信息: %s, 头像=%s\n", user.Nickname, user.AvatarURL)
			}
		} else {
			fmt.Printf("⚠️  获取auth-center用户信息失败: %v\n", err)
		}
	}

	// 转换角色为大写格式给前端
	c.JSON(http.StatusOK, gin.H{
		"id":                 user.ID,
		"authCenterUserId":   user.AuthCenterUserID,
		"nickname":           user.Nickname,
		"avatarUrl":          user.AvatarURL,
		"profile":            user.Profile,
		"roles":              convertRolesToUpperCase(user.Roles),
		"currentRole":        user.ActiveRole,
		"lastUsedRole":       user.LastUsedRole,
		"status":             user.Status,
		"lastLoginAt":        user.LastLoginAt,
		"lastLoginIp":        user.LastLoginIP,
		"createdAt":          user.CreatedAt,
		"updatedAt":          user.UpdatedAt,
	})
}

// callAuthCenterGetUserInfo 调用 auth-center 获取用户信息
func (ctrl *AuthController) callAuthCenterGetUserInfo(token string) (*AuthCenterUserInfo, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", ctrl.cfg.AuthCenterURL+"/api/auth/user-info", nil)
	if err != nil {
		return nil, err
	}

	// 设置 Authorization header
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("auth-center 返回错误: %d", resp.StatusCode)
	}

	// 解析响应
	var userInfoResp AuthCenterUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfoResp); err != nil {
		return nil, err
	}

	return &userInfoResp, nil
}

