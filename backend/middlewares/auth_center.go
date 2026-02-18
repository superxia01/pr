package middlewares

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pr-business/config"
	"pr-business/constants"
	"pr-business/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthCenterMiddleware 账号中心认证中间件（按照 V3.1 统一 Token 模式）
func AuthCenterMiddleware(cfg *config.Config, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取 token（优先从 Header，其次从 URL 参数）
		token := c.GetHeader("Authorization")
		if token == "" {
			token = c.Query("token")
		} else {
			// 移除 "Bearer " 前缀
			token = strings.TrimPrefix(token, "Bearer ")
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "未登录",
			})
			c.Abort()
			return
		}

		// 2. 验证 token
		authCenterUserID, err := verifyToken(cfg.AuthCenterURL, token)
		if err != nil {
			fmt.Printf("[AuthCenterMiddleware] VerifyToken failed: %v\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Token 无效或已过期",
			})
			c.Abort()
			return
		}

		// 3. 获取用户信息
		userInfo, err := getUserInfo(cfg.AuthCenterURL, token)
		if err != nil {
			fmt.Printf("[AuthCenterMiddleware] GetUserInfo failed: %v\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "获取用户信息失败",
			})
			c.Abort()
			return
		}

		// 4. 查找本地用户
		var user models.User
		result := db.Where("auth_center_user_id = ?", authCenterUserID).First(&user)

		if result.Error == gorm.ErrRecordNotFound {
			// 5. 新用户：自动创建
			user = models.User{
				AuthCenterUserID: authCenterUserID,
				Nickname:         userInfo.Nickname,
				AvatarURL:        userInfo.AvatarURL,
				Roles:            models.Roles{constants.RoleBasicUser},
				Status:           "active",
			}
			if err := db.Create(&user).Error; err != nil {
				fmt.Printf("[AuthCenterMiddleware] Create user failed: %v\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"error":   "创建本地用户失败",
				})
				c.Abort()
				return
			}
			fmt.Printf("[AuthCenterMiddleware] Created new user: %s\n", user.ID)
		} else if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "数据库错误",
			})
			c.Abort()
			return
		} else {
			// 6. 更新用户信息（如果变化）
			updated := false
			if userInfo.Nickname != "" && user.Nickname != userInfo.Nickname {
				user.Nickname = userInfo.Nickname
				updated = true
			}
			if userInfo.AvatarURL != "" && user.AvatarURL != userInfo.AvatarURL {
				user.AvatarURL = userInfo.AvatarURL
				updated = true
			}
			if updated {
				db.Save(&user)
			}
		}

		// 7. 存入上下文
		c.Set("user", &user)
		c.Set("userId", user.ID)
		c.Set("roles", user.Roles)

		c.Next()
	}
}

// UserInfo auth-center 用户信息
type UserInfo struct {
	UserID      string `json:"userId"`
	UnionID     string `json:"unionId"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
	Nickname    string `json:"nickname"`
	AvatarURL   string `json:"avatarUrl"`
}

// verifyToken 验证 auth-center token
func verifyToken(authCenterURL, token string) (string, error) {
	url := fmt.Sprintf("%s/api/auth/verify-token", authCenterURL)

	reqBody := map[string]string{"token": token}
	jsonBody, _ := json.Marshal(reqBody)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(url, "application/json", strings.NewReader(string(jsonBody)))
	if err != nil {
		return "", fmt.Errorf("请求账号中心失败: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		Success bool `json:"success"`
		Data    struct {
			UserID  string `json:"userId"`
			UnionID string `json:"unionId"`
		} `json:"data"`
		Error string `json:"error"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if !result.Success || result.Data.UserID == "" {
		return "", fmt.Errorf("token 无效")
	}

	return result.Data.UserID, nil
}

// getUserInfo 获取 auth-center 用户信息
func getUserInfo(authCenterURL, token string) (*UserInfo, error) {
	url := fmt.Sprintf("%s/api/auth/user-info", authCenterURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求账号中心失败: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		Success bool `json:"success"`
		Data    struct {
			UserID      string `json:"userId"`
			UnionID     string `json:"unionId"`
			PhoneNumber string `json:"phoneNumber"`
			Email       string `json:"email"`
			Profile     struct {
				Nickname  string `json:"nickname"`
				AvatarURL string `json:"avatarUrl"`
			} `json:"profile"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("获取用户信息失败")
	}

	return &UserInfo{
		UserID:      result.Data.UserID,
		UnionID:     result.Data.UnionID,
		PhoneNumber: result.Data.PhoneNumber,
		Email:       result.Data.Email,
		Nickname:    result.Data.Profile.Nickname,
		AvatarURL:   result.Data.Profile.AvatarURL,
	}, nil
}
