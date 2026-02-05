package middlewares

import (
	"net/http"
	"pr-business/models"
	"pr-business/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Authorization header获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Missing authorization header",
			})
			c.Abort()
			return
		}

		// 解析Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization format",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 解析token
		claims, err := utils.ParseToken(tokenString, jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// 查询用户对象
		db := c.MustGet("db").(*gorm.DB)
		var user models.User
		if err := db.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not found",
			})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("userId", claims.UserID)
		c.Set("user", &user)
		c.Set("currentRole", claims.Role)
		c.Set("roles", claims.Roles)

		c.Next()
	}
}

// RequireRole 角色权限中间件
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文获取当前角色
		currentRole, exists := c.Get("currentRole")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "No role found in context",
			})
			c.Abort()
			return
		}

		// 检查角色是否在允许列表中
		roleStr := currentRole.(string)
		allowed := false
		for _, r := range allowedRoles {
			if r == roleStr {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
